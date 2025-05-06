package security

import (
	cfg "sentimenta/internal/config"
	errs "sentimenta/internal/errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	config *cfg.Config
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

func (j JWT) ParseJWT(tokenStr string, secretKey string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Проверка метода подписи
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.ErrUnsupportedSignatureMethod
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}

	// Проверка на валидность токена и срок действия
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Проверка срока действия
		if exp, ok := claims["exp"].(float64); ok {
			if exp < float64(time.Now().Unix()) {
				return "", errs.ErrTokenExpired // Если токен истёк
			}
		} else {
			return "", errs.ErrNoExpClaim // Если в токене нет поля exp
		}

		// Извлечение userID из claim
		if uid, ok := claims["sub"].(string); ok {
			return uid, nil
		}
	}

	return "", errs.ErrNotFoundInJWT
}

func NewJWT(cfg *cfg.Config) *JWT {
	return &JWT{config: cfg}
}
