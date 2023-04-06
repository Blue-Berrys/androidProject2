package JWT

import (
	"fmt"
	"testing"
)

// 测试从username和password变成tokenstr
func TestGetToken(t *testing.T) {
	userid := uint(17)
	tokenStr, err := GetToken(userid)
	if err != nil {
		fmt.Println("Test GetToken Failed:", err)
		return
	}
	fmt.Println(tokenStr)
	//测试从tokenstr变成username和password
	claim, ok := TokenToClaim(tokenStr)
	if !ok {
		fmt.Println("Test TokenToClaim Failed")
		return
	}
	fmt.Println(claim.UserId)
}

// JWTMiddleware用apifox测
func TestJWTMiddleware(t *testing.T) {

}
