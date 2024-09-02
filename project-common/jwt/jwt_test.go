package jwt

import (
	"fmt"
	"testing"
)

// 测试验证解析函数
func TestParseToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjMzOTYyNzIsInRva2VuIjoiMSJ9.ZGCeMUIZ89EbJMj4FnBTIiMsh0eAvDrgwKV44bs2Lv4"
	str, _ := ParseToken(tokenString, "msproject")
	fmt.Println(tokenString + "==>" + str)
}
