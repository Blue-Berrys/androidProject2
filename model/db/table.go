package model

import "gorm.io/gorm"

type User struct {
	gorm.Model             //自带的id是自增id
	UserName        string //用户注册name
	Password        string
	NickName        string //昵称
	Avatar          string //头像
	BackgroundImage string //背景图片
	Level           int    //等级，0-封号，1-普通用户，2-管理员
	WorkCount       int64  //用户发布的朋友圈总数

	//FriendsChats []FriendsChat
	//Likes        []Like
}

type Comment struct {
	gorm.Model
	UserId       int64
	FriendChatId int64  //这一条朋友圈id
	CommentText  string //评论信息

	//User User
}

type Like struct {
	gorm.Model
	UserId        int64
	FriendsChatId int64
}

type FriendsChat struct { //这一条朋友圈
	gorm.Model
	UserId   uint
	Content  string //内容
	ImageUrl string //图片地址

	//Likes    []Like
	//Comments []Comment
}
