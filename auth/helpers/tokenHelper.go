package helper

import (
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type SignedDetails struct {
	ID string
	Email    string
	Name     string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

// user.Email, *user.Name, *user.ID
func GenerateAllTokens(ID string, email string, name string) (signedToken string, err error) {
	claims := &SignedDetails{
		Email:    email,
		Name:     name,
		ID: ID,

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24*7)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		return
	}
	return claims, msg
}
