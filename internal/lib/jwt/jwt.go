package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rautaruukkipalich/go_auth_grpc/internal/domain/models"
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

