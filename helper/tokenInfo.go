package helper

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
	SignKey          = "niChatTest"
)

type CustomClaims struct {
	UserId    string `json:"user_id"`
	NetId     string `json:"net_id"`
	TimeStamp int64  `json:"time_stamp"`
	jwt.StandardClaims
}

type JWTInfo struct {
	SignKey []byte
}

func NewJWT() *JWTInfo {
	return &JWTInfo{[]byte(SignKey)}
}

// 创建token
func (j *JWTInfo) TokenCreate(tokenInfo CustomClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodES256, tokenInfo)
	return t.SignedString(j.SignKey)
}

// 解析token
func (j *JWTInfo) TokenParse(token string) (*CustomClaims, error) {
	t, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return j.SignKey, nil
	})
	if err != nil {
		if ve, ok := err.(jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := t.Claims.(*CustomClaims); ok && t.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// 更新token
func (j *JWTInfo) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Now()
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SignKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.TokenCreate(*claims)
	}
	return "", TokenInvalid
}
