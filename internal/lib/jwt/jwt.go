package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rautaruukkipalich/go_auth_grpc/internal/domain/models"
)

const (
	ZeroValue = 0
)

func NewJWTToken(user models.User, app models.App, ttl time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(ttl).Unix()
	claims["app_id"] = app.ID

	return token.SignedString([]byte(app.Secret))
}

func GetAppIDFromJWTToken(token string) (int, error) {
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(
		token,
		claims, 
		func(token *jwt.Token) (any, error) {return []byte{}, nil},
	)
	appID := int(claims["app_id"].(float64))
	if appID == ZeroValue {
		return 0, ErrJWTDecode
	}
	return appID, nil
}

func GetSubFromJWTToken(token string, app models.App) (int, error) {
	secret := []byte(app.Secret)

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(
		token,
		claims, 
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrJWTDecode
			}
			return secret, nil
		},
	)

	if err != nil {
		return 0, err
	}

	sub := int(claims["sub"].(float64))

	return sub, nil
}
