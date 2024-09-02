package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

type JwtToken struct {
	AccessToken  string
	RefreshToken string
	AccessExp    int64
	RefreshExp   int64
}

func CreateToken(val string, exp time.Duration, secret string, refreshSecret string, refreshExp time.Duration) *JwtToken {
	aExpire := time.Now().Add(exp).Unix()
	rExpire := time.Now().Add(refreshExp).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": val,
		"exp":   aExpire,
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": val,
		"exp":   rExpire,
	})
	aToken, _ := accessToken.SignedString([]byte(secret))
	rToken, _ := refreshToken.SignedString([]byte(refreshSecret))
	return &JwtToken{
		AccessToken:  aToken,
		AccessExp:    aExpire,
		RefreshToken: rToken,
		RefreshExp:   rExpire,
	}
}
func ParseToken(tokenStr string, secret string) (string, error) {
	// 去除可能的 "Bearer " 前缀
	tokenStr = strings.TrimPrefix(tokenStr, "bearer ")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v \n", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		val := claims["token"].(string)
		exp := claims["exp"].(float64)
		expTime := int64(exp)
		if expTime <= time.Now().Unix() {
			return "", errors.New("token过期了")
		}
		return val, nil
	} else {
		return "", err
	}
}
