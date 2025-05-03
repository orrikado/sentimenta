package utils

import (
	cfg "sentimenta/internal/config"
	"sentimenta/internal/jwt"

	"github.com/labstack/echo/v4"
)

var config = cfg.NewConfig()

func GetUserID(c echo.Context) (string, error) {
	jwtCookie, err := c.Cookie(config.JWT_COOKIE_NAME)
	if err != nil {
		return "", err
	}

	jwtToken := jwtCookie.Value

	userID, err := jwt.ParseJWT(jwtToken)
	if err != nil {
		return "", err
	}

	return userID, nil
}
