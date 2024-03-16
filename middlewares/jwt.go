package middlewares

import (
	"ALTA_BE_SOSMED/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(email string) (string, error) {
	var data = jwt.MapClaims{}
	data["email"] = email
	data["iat"] = time.Now().Unix()
	data["exp"] = time.Now().Add(time.Hour * 3).Unix()
	var processToken = jwt.NewWithClaims(jwt.SigningMethodRS256, data)
	result, err := processToken.SignedString([]byte(config.JWTSECRET))
	if err != nil {
		return "", err
	}
	return result, nil
}

func DecodeToken(token *jwt.Token) string {
	var result string
	var claim = token.Claims.(jwt.MapClaims)

	if val, found := claim["email"]; found {
		result = val.(string)
	}
	return result
}
