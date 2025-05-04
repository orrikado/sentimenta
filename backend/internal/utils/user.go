package utils

import (
	"errors"

	"github.com/labstack/echo/v4"
)

func GetUserID(c echo.Context) (string, error) {
	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		return "", errors.New("userID not found in context")
	}
	return userID, nil
}
