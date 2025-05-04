package middlewares

import (
	"net/http"
	"sentimenta/internal/config"
	"sentimenta/internal/security"

	"github.com/labstack/echo/v4"
)

func NewJWTMiddleware(cfg config.Config, JWT security.JWT) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(cfg.JWT_COOKIE_NAME)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "требуется аутентификация"})
			}

			userID, err := JWT.ParseJWT(cookie.Value)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "невалидный токен"})
			}

			c.Set("userID", userID)
			return next(c)
		}
	}
}
