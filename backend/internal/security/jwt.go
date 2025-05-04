package security

import (
	"sentimenta/internal/config"
	cfg "sentimenta/internal/config"
	errs "sentimenta/internal/errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	config cfg.Config
}

func (j JWT) GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 720).Unix(),
		"iat": time.Now().Unix(),
		"iss": "my-api",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.JWT_SECRET))
}

func (j JWT) ParseJWT(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Проверка метода подписи
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.ErrUnsupportedSignatureMethod
		}
		return []byte(j.config.JWT_SECRET), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if uid, ok := claims["sub"].(string); ok {
			return uid, nil
		}
	}

	return "", errs.ErrNotFoundInJWT
}

func NewJWT(cfg config.Config) JWT {
	return JWT{config: cfg}
}
