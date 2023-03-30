package model

import (
	"errors"
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	//需要现在本地建database：androidProject2
	Init()
	//到本地数据库中查看是否建立成功
	existComments := DB.Migrator().HasTable("friends_chats")
	if !existComments {
		errors.New("Create Table friends_chats Failed")
		return
	}
	existLikes := DB.Migrator().HasTable("likes")
	if !existLikes {
		errors.New("Create Table likes Failed")
		return
	}
	existUsers := DB.Migrator().HasTable("users")
	if !existUsers {
		errors.New("Create Table users Failed")
		return
	}
	existVideos := DB.Migrator().HasTable("comments")
	if !existVideos {
		errors.New("Create Table comments Failed")
		return
	}
	fmt.Println("Create Table Successfully")
}
