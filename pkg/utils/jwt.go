package utils

import (
	"errors"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	Id       uint   `json:"id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

const sign_key = "hello jwt"

// 随机字符串
var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(str_len int) string {
	rand_bytes := make([]rune, str_len)
	for i := range rand_bytes {
		rand_bytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(rand_bytes)
}

func GenerateToken(id uint, username string) (string, error) {
	claim := MyCustomClaims{
		Id:       id,
		UserName: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Auth_Server",                                   // 签发者
			Subject:   username,                                        // 签发对象
			Audience:  jwt.ClaimStrings{"ALL"},                         //签发受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),   //过期时间
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)), //最早使用时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                  //签发时间
			ID:        randStr(10),                                     // wt ID, 类似于盐值
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(sign_key))
	return token, err
}

func ParseToken(token_string string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(token_string, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(sign_key), nil //返回签名密钥
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("claim invalid")
	}

	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok {
		return nil, errors.New("invalid claim type")
	}

	return claims, nil
}
