package auth

import (
	"errors"
	"fmt"
	"golangoauth2example/internal/model"
	"log"

	"github.com/golang-jwt/jwt/v4"
)

type userClaim struct {
	jwt.RegisteredClaims
	User model.User
}

const key = "my secret jwt key"

func CreateJWTTokenFromUserData(user model.User) (string, error) {

	claim := userClaim{
		RegisteredClaims: jwt.RegisteredClaims{},
		User: model.User{
			ID:              user.ID,
			Name:            user.Name,
			Picture:         user.Picture,
			Email:           user.Email,
			IsEmailVerified: user.IsEmailVerified,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedString, err := token.SignedString([]byte(key))

	if err != nil {
		return "", fmt.Errorf("error creating signed string: %v", err)
	}

	return signedString, nil
}

func ParseJWTToken(jwtToken string) (*model.User, error) {
	user := &userClaim{}

	token, err := jwt.ParseWithClaims(jwtToken, user, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if !token.Valid {
		return nil, errors.New("token not valid")
	}

	result := &model.User{
		ID:              user.User.ID,
		Name:            user.User.Name,
		Picture:         user.User.Picture,
		Email:           user.User.Email,
		IsEmailVerified: user.User.IsEmailVerified,
	}

	return result, nil

}
