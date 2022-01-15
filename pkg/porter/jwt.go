package porter

import (
	"time"

	"github.com/barklan/cto/pkg/storage/vars"
	log "github.com/sirupsen/logrus"

	"github.com/golang-jwt/jwt/v4"
)

type TokenClaims struct {
	Name        string `json:"name"`
	ProjectName string `json:"project_name"`
	jwt.RegisteredClaims
}

func CreateJWT(base *Base, name, project string) string {
	mySigningKey := []byte(base.Config.TG.BotToken)

	jwtExp := time.Duration(base.Config.JWTExpHours) * time.Hour
	expTime := time.Now().Add(jwtExp)
	claims := TokenClaims{
		name,
		project,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			Issuer:    "cto",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Panicln("failed to create jwt token", err)
	}
	return ss
}

func RotateJWT(base *Base, name, project string) {
	ss := CreateJWT(base, name, project)
	log.Println("Rotated auth token:", ss)

	jwtExp := time.Duration(base.Config.JWTExpHours) * time.Hour
	err := base.Cache.SetVar(project, vars.AuthToken, ss, jwtExp)
	if err != nil {
		log.Panicln("failed to set new jwt token to cache")
	}
}
