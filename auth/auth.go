package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/zedisdog/armor/config"
	"github.com/zedisdog/armor/model"
	"strconv"
)

var (
	TokenIsInvalid = errors.New("token is invalid")
)

type MyCustomClaims struct {
	DouDouBirthday string `json:"doudou_birthday"`
	jwt.StandardClaims
}

func GenerateToken(account model.HasId) (string, error) {
	claims := MyCustomClaims{
		"19960415",
		jwt.StandardClaims{
			Id: strconv.FormatUint(account.GetId(), 10),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(config.Conf.Bytes("jwt.key"))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func ParseToken(token string) (*MyCustomClaims, error) {
	t, err := jwt.ParseWithClaims(token, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.Conf.Bytes("jwt.key"), nil
	})
	if err != nil {
		return nil, err
	}
	if t.Valid {
		return t.Claims.(*MyCustomClaims), nil
	} else {
		return nil, TokenIsInvalid
	}
}