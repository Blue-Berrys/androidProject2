package model

import "gorm.io/gorm"

type User struct {
	gorm.Model             //自带的id是自增id
	UserName        string //用户昵称
	Password        string
	Signature       string //个性签名
	Avatar          string //头像
	BackgroundImage string //背景图片
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
	UserId    int64
	ContextId int64
}

type FriendsChat struct { //这一条朋友圈
	gorm.Model
	Content  string //内容
	ImageUrl string //图片地址

	//Likes    []Like
	//Comments []Comment
}
