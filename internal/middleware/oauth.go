package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/ShotaKitazawa/kube-portal/backend/view"
	"github.com/ShotaKitazawa/kube-portal/internal/util"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/jmespath/go-jmespath"
	"github.com/labstack/echo/v4"
)

type AuthConfig struct {
	OIDCVerifier      *oidc.IDTokenVerifier
	RoleAttributePath string
}

func AuthWithConfig(conf AuthConfig, l *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (returnErr error) {
			authHeader := ctx.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return next(ctx)
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			idToken, err := conf.OIDCVerifier.Verify(ctx.Request().Context(), tokenStr)
			if err != nil {
				return view.ResponseJSON(ctx, http.StatusForbidden, "Provided token is invalid")
			}

			var claims any
			if err := idToken.Claims(&claims); err != nil {
				return view.ResponseJSON(ctx, http.StatusInternalServerError, "failed to parse claims in IDToken")
			}
			if ok := authorizeBasedJMESPath(ctx, claims, conf.RoleAttributePath); !ok {
				// response has already set in the above function
				return nil
			}

			util.SetContextIsLogin(ctx, true)
			return next(ctx)
		}
	}
}

func authorizeBasedJMESPath(ctx echo.Context, claims any, roleAttributePath string) bool {
	res, err := jmespath.Search(roleAttributePath, claims)
	if err != nil {
		_ = view.ResponseJSON(ctx, http.StatusInternalServerError, "failed to search claims by JMESPath")
		return false
	}
	if typedRes, ok := res.(bool); !ok {
		_ = view.ResponseJSON(ctx, http.StatusInternalServerError, "the result of JMESPath is not bool")
		return false
	} else if !typedRes {
		_ = view.ResponseJSON(ctx, http.StatusForbidden, "the result of JMESPath is not true")
		return false
	}
	return true
}
