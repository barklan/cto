package querying

import (
	"fmt"
	"log"
	"net/http"

	"github.com/barklan/cto/pkg/storage"
	"github.com/golang-jwt/jwt/v4"
)

func authorize(data *storage.Data, tokenQ string) (string, int, bool) {
	if tokenQ == "" {
		log.Println("No token provided for query.")
		return "", http.StatusUnauthorized, false
	}
	if tokenQ == data.Config.Internal.MagicJWTToken {
		return "nftg", http.StatusOK, true // HACK
	}
	tokenParsed, err := jwt.Parse(tokenQ, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(data.Config.Internal.TG.BotToken), nil
	})
	if err != nil {
		return "", http.StatusUnauthorized, false
	}

	if claims, ok := tokenParsed.Claims.(jwt.MapClaims); ok && tokenParsed.Valid {
		log.Println("token is valid, claims:", claims)
		projectName := claims["project_name"].(string)
		return projectName, http.StatusOK, true
	} else {
		log.Println("token is not ok (returning 403):", err)
		return "", http.StatusForbidden, false
	}
}
