package porter

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

func authorize(base *Base, tokenQ string) (string, string, int, bool) {
	if tokenQ == "" {
		base.Log.Warn("no token provided for query")
		return "", "", http.StatusUnauthorized, false
	}
	tokenParsed, err := jwt.Parse(tokenQ, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(base.Config.TG.BotToken), nil
	})
	if err != nil {
		return "", "", http.StatusUnauthorized, false
	}

	if claims, ok := tokenParsed.Claims.(jwt.MapClaims); ok && tokenParsed.Valid {
		projectName := claims["project_name"].(string)
		name := claims["name"].(string)
		return name, projectName, http.StatusOK, true
	} else {
		base.Log.Warn("token is not ok (returning 403)")
		return "", "", http.StatusForbidden, false
	}
}
