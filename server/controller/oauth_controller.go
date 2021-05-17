package controller

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/sirupsen/logrus"

	"github.com/ShotaKitazawa/kube-portal/server/view"
)

var (
	jwtIssuer               = "kube-portal"
	jwtAndCookieExpiredTime = time.Hour
)

func init() {
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return "github", nil
	}
}

type OAuthController struct {
	logger *logrus.Logger
	view   view.ViewPort

	jwtSecret          string
	allowGitHubUserIDs []string
}

func NewOAuthController(l *logrus.Logger, githubKey, githubSecret, callbackUrl string, jwtSecret string, allowGitHubUserName ...string) (*OAuthController, error) {
	goth.UseProviders(
		github.New(githubKey, githubSecret, callbackUrl),
	)
	// TODO: get GitHub UserIDs

	return &OAuthController{l, &view.JsonView{}, jwtSecret, []string{}}, nil
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

	// TODO: authz
	if user.UserID == "19530785" {
		now := time.Now()
		expiresAt := now.Add(jwtAndCookieExpiredTime)

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
	}

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
