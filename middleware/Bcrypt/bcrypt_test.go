package Bcrypt

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestEncryptionByPassword(t *testing.T) {
	password := "8950po@"
	newPassword1 := EncryptionByPassword(password)
	newPassword2 := EncryptionByPassword(password)

	//fmt.Println(newPassword1)

	nowG := time.Now()
	fmt.Println("加密后", newPassword1, "耗时", time.Now().Sub(nowG))
	nowC := time.Now()
	err := bcrypt.CompareHashAndPassword([]byte(newPassword1), []byte(password))
	fmt.Println("验证耗费时间", time.Now().Sub(nowC))
	fmt.Println(err)

	nowG = time.Now()
	fmt.Println("加密后", newPassword2, "耗时", time.Now().Sub(nowG))
	nowC = time.Now()
	err = bcrypt.CompareHashAndPassword([]byte(newPassword2), []byte(password))
	fmt.Println("验证耗费时间", time.Now().Sub(nowC))
	fmt.Println(err)
}
