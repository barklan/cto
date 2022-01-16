package security

import (
	"fmt"
	"time"

	"github.com/barklan/cto/pkg/storage"
	"github.com/golang-jwt/jwt/v4"
)

type TokenClaims struct {
	Name        string `json:"name"`
	ProjectName string `json:"project_name"`
	jwt.RegisteredClaims
}

func CreateJWT(conf *storage.InternalConfig, email, project string) (string, error) {
	mySigningKey := []byte(conf.TG.BotToken)

	jwtExp := time.Duration(conf.JWTExpHours) * time.Hour
	expTime := time.Now().Add(jwtExp)
	claims := TokenClaims{
		email,
		project,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			Issuer:    "cto",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", fmt.Errorf("failed to create jwt token: %w", err)
	}
	return ss, nil
}
