package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type (
	TokenOptions struct {
		AccessSecret string
		AccessExpire int64
		Fileds       map[string]interface{}
	}
	Token struct {
		AccessToken  string `json:"access_token"`
		AccessExipre int64  `json:"access_expire"`
	}
)

func BuildToken(opt TokenOptions) (Token, error) {
	var token Token
	now := time.Now().Add(-time.Minute).Unix()
	accessToken, err := genToken(now, opt.AccessSecret, opt.Fileds, opt.AccessExpire)
	if err != nil {
		return token, err
	}
	token.AccessToken = accessToken
	token.AccessExipre = now + opt.AccessExpire

	return token, nil
}

func genToken(iat int64, secretKey string, payLoad map[string]interface{}, seconds int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	for k, v := range payLoad {
		claims[k] = v
	}
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
