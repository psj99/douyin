package utils

import (
	"douyin/service/types/response"

	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 自定义错误类型
var ErrorTokenInvalid = errors.New("token无效")

type CustomClaims struct {
	User_ID  uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

const sign_key = "hello jwt"

// 随机字符串基础
var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(str_len int) string {
	rand_bytes := make([]rune, str_len)
	for i := range rand_bytes {
		rand_bytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(rand_bytes)
}

func GenerateToken(user_id uint, username string) (string, error) {
	claim := CustomClaims{
		User_ID:  user_id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Auth_Server",                                   // 签发者
			Subject:   username,                                        // 签发对象
			Audience:  jwt.ClaimStrings{"ALL"},                         // 签发受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),   // 过期时间
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)), // 最早使用时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                  // 签发时间
			ID:        randStr(10),                                     // wt ID, 类似于盐值
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(sign_key))
	return token, err
}

func ParseToken(token_string string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(token_string, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(sign_key), nil // 返回签名密钥
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrorTokenInvalid
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, ErrorTokenInvalid
	}

	return claims, nil
}

// gin中间件
// jwt鉴权 验证token并提取user_id与username
func MiddlewareAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 尝试从GET中提取token
		tokenStr := ctx.Query("token")
		// 若失败则尝试从POST中提取token
		if tokenStr == "" {
			tokenStr = ctx.PostForm("token")
		}
		// 若无法提取token
		if tokenStr == "" {
			ZapLogger.Warnf("MiddlewareAuth warn: 未授权请求")
			ctx.JSON(http.StatusUnauthorized, response.CommonResp{Status_Code: -1, Status_Msg: "需要token"})
			ctx.Abort()
			return
		}

		// 解析/校验token (自动验证有效期等)
		claims, err := ParseToken(tokenStr)
		if err != nil {
			if err == ErrorTokenInvalid {
				ZapLogger.Warnf("MiddlewareAuth warn: 未授权请求")
				ctx.JSON(http.StatusUnauthorized, response.CommonResp{
					Status_Code: -1,
					Status_Msg:  "token无效",
				})
				ctx.Abort()
				return
			} else {
				ZapLogger.Errorf("MiddlewareAuth err: token解析失败")
				ctx.JSON(http.StatusInternalServerError, response.CommonResp{
					Status_Code: -1,
					Status_Msg:  "token解析失败",
				})
				ctx.Abort()
				return
			}
		}

		// 提取user_id和username
		ctx.Set("user_id", claims.User_ID)
		ctx.Set("username", claims.Username)

		ctx.Next()
	}
}
