package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenClaims struct {
	ProjectName string `json:"project_name"`
	jwt.RegisteredClaims
}

func RotateJWT(data *Data, project string) {
	mySigningKey := []byte(data.Config.Internal.TG.BotToken)

	jwtExp := time.Duration(data.Config.Internal.JWTExpHours) * time.Hour
	expTime := time.Now().Add(jwtExp)
	claims := TokenClaims{
		project,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			Issuer:    "cto",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(mySigningKey)
	log.Println("Rotated auth token:", ss)

	data.SetObj(fmt.Sprintf("authToken-%s", project), ss, jwtExp)
}
