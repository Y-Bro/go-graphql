package jwt

import (
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	SecretKey = []byte("secret")
)

func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["username"] = username

	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(SecretKey)

	if err != nil {
		log.Fatalln("Error while siging token")
		return "", err
	}

	return tokenString, nil

}

func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	fmt.Println("Here2", tokenStr)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)

		return username, nil
	} else {
		fmt.Println("Here")
		return "", err
	}
}
