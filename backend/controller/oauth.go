package controller

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"

	"github.com/ShotaKitazawa/kube-portal/backend/models/ports"
	"github.com/ShotaKitazawa/kube-portal/backend/utils"
)

var (
	jwtIssuer = "kube-portal"
)

func init() {
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return "github", nil
	}
}

type OAuthController struct {
	logger *slog.Logger

	jwtSecret          string
	expiredTime        time.Duration
	allowGitHubUserIDs []string
}

func NewOAuthController(l *slog.Logger, githubClient ports.GitHub, githubKey, githubSecret, callbackUrl string, jwtSecret string, expiredTime time.Duration, allowGitHubUserName ...string) (*OAuthController, error) {
	// set OAuth Provider
	goth.UseProviders(
		github.New(githubKey, githubSecret, callbackUrl),
	)

	// get GitHub UserIDs
	var allowIDs []string
	for _, username := range allowGitHubUserName {
		id, err := githubClient.GetUserIDByUserName(context.Background(), username)
		if err != nil {
			return nil, err
		}
		allowIDs = append(allowIDs, id)
	}

	return &OAuthController{l, jwtSecret, expiredTime, allowIDs}, nil
}

type MyCustomClaims struct {
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	jwt.StandardClaims
}

func (c OAuthController) Callback(ctx echo.Context) error {
	user, err := gothic.CompleteUserAuth(ctx.Response(), ctx.Request())
	if err != nil {
		fmt.Fprintln(ctx.Response(), err)
		return err
	}

	// Authorization: check login user contains c.allowGitHubUserIDs
	if !utils.Contains(c.allowGitHubUserIDs, user.UserID) {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}

	// Set Cookie the JWT token
	now := time.Now()
	expiresAt := now.Add(c.expiredTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		user.Name,
		user.AvatarURL,
		jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: expiresAt.Unix(),
			Issuer:    jwtIssuer,
		},
	})
	signedToken, err := token.SignedString([]byte(c.jwtSecret))
	if err != nil {
		return err
	}
	ctx.SetCookie(&http.Cookie{
		Name:    "jwt",
		Value:   signedToken,
		Path:    "/",
		Expires: expiresAt,
	})

	return ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func (c OAuthController) Logout(ctx echo.Context) error {
	if err := gothic.Logout(ctx.Response(), ctx.Request()); err != nil {
		return err
	}
	if cookie, err := ctx.Request().Cookie("jwt"); err == nil {
		cookie.Path = "/"
		cookie.MaxAge = -1
		http.SetCookie(ctx.Response(), cookie)
	}

	return ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func (c OAuthController) Login(ctx echo.Context) error {
	gothic.BeginAuthHandler(ctx.Response(), ctx.Request())
	if status := ctx.Response().Status; status >= 400 {
		c.logger.Debug(fmt.Sprintf("Not Authorized: status code %d", status))
	}
	return nil
}
