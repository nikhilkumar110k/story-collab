package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const Secretkey = "totallsecretkeylol"

func Generatetoken(email string, userid int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Email":  email,
		"userid": userid,
		"exp":    time.Now().Add(time.Hour * 20).Unix(),
	})
	return token.SignedString([]byte(Secretkey))
}

func VerifyToken(token string) (error, int64) {
	tokenparsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Wrong token")
		}
		return []byte(Secretkey), nil
	})
	if err != nil {
		return errors.New("couldn't parse"), 0
	}
	validornot := tokenparsed.Valid
	if !validornot {
		return errors.New("the token is not valid or just expired"), 0
	}
	claims, ok := tokenparsed.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("Invalid Claims"), 0
	}
	useridRaw, exists := claims["userid"]
	if !exists {
		return errors.New("userid not found in claims"), 0
	}
	var Userid int64
	switch v := useridRaw.(type) {
	case float64:
		Userid = int64(v)
	case int64:
		Userid = v
	default:
		return errors.New("unexpected type for userid in claims"), 0
	}

	return nil, Userid
}
