package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const hmacSecret string = "sdhfbviuasbfilvabsid"

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 2)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func GenJWT(userId int) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"userId": fmt.Sprint(userId)})

	jwt, err := token.SignedString([]byte(hmacSecret))
	if err != nil {
		fmt.Print(err.Error())
		return "", fmt.Errorf("error generating signed jwt string")
	}
	return jwt, nil
}

func ExtractJwtFromAuthHeader(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	tokenarr := strings.Split(header, " ")
	if len(tokenarr) < 2 {
		return "", fmt.Errorf("pass Bearer token to proceed")
	}
	token := tokenarr[1]
	return token, nil
}

func ValidateJwtAndExtractClaims(tokenString string) (string, error) {
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(hmacSecret), nil
	})
	if err != nil {

		return "", err
	}
	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok {
		userId := claims["userId"].(string)
		return userId, nil
	} else {
		fmt.Println(err)
		return "", err
	}

}
