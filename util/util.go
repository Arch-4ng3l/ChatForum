package util

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
)

func CreateHash(s string) string {
	hash := sha256.New()
	return hex.EncodeToString(hash.Sum([]byte(s)))
}

func DecodeJson(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func CreateJWT(name string) (string, error) {
	claims := &jwt.MapClaims{
		"ExpiresAt": 15000,
		"Name":      name,
	}
	secret := os.Getenv("JWT-SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return []byte(secret), nil
	})
}

func AuthJWT(uname, tokenString string) bool {

	token, err := ValidateJWT(tokenString)

	if err != nil {
		return false
	}

	claims := token.Claims.(jwt.MapClaims)

	name := claims["Name"]

	if name != uname {
		return false
	}

	return true
}
