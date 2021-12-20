package porter

import (
	"time"

	"github.com/barklan/cto/pkg/storage/vars"
	log "github.com/sirupsen/logrus"

	"github.com/golang-jwt/jwt/v4"
)

type TokenClaims struct {
	ProjectName string `json:"project_name"`
	jwt.RegisteredClaims
}

func RotateJWT(base *Base, project string) {
	mySigningKey := []byte(base.Config.TG.BotToken)

	jwtExp := time.Duration(base.Config.JWTExpHours) * time.Hour
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

	err := base.Cache.SetVar(project, vars.AuthToken, ss, jwtExp)
	if err != nil {
		log.Panicln("failed to set new jwt token to cache")
	}
}
