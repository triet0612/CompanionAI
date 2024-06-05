package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type CustomJWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type Middleware struct {
	JWT_SECRET      []byte
	JWT_AUTH_METHOD string
}

func NewClaim(c echo.Context) jwt.Claims {
	return new(CustomJWTClaims)
}

func (a *Middleware) AuthMiddleware() echo.MiddlewareFunc {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != a.JWT_AUTH_METHOD {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return []byte(a.JWT_SECRET), nil
	}

	return echojwt.WithConfig(echojwt.Config{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "login") ||
				strings.Contains(c.Path(), "register") ||
				strings.Contains(c.Path(), "docs") ||
				strings.Contains(c.Path(), "assets")
		},
		SigningMethod: a.JWT_AUTH_METHOD,
		SigningKey:    a.JWT_SECRET,
		TokenLookup:   "cookie:jwt,header:Authorization:Bearer ",
		NewClaimsFunc: NewClaim,
		KeyFunc:       keyFunc,
		ErrorHandler: func(c echo.Context, err error) error {
			if c.Request().Header.Get("Content-Type") == "application/json" {
				return err
			}
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		},
	})
}
