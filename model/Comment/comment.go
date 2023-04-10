package model

import (
	model "androidProject2/model/db"
	"errors"
	"gorm.io/gorm"
	"sync"
)

type CommentDao struct {
}

var (
	commentDao  *CommentDao
	commentOnce sync.Once
)

func NewCommentDao() *CommentDao {
	commentOnce.Do(func() {
		commentDao = new(CommentDao)
	})
	return commentDao
}

// 增加一条评论
func (u *CommentDao) AddCommentInfo(info *model.Comment) error {
	if info == nil {
		return errors.New("空指针错误")
	}
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(info).Error; err != nil {
			return err
		}
		return nil
	})
}

// 根据朋友圈Id查询所有评论
func (u *CommentDao) QueryCommentByFriendsChatId(FriendsChatId uint, comments *[]*model.Comment) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("friend_chat_id=?", FriendsChatId).Find(&comments).Error; err != nil {
			return err
		}
		return nil
	})
}

// 根据评论Id删除评论
func (u *CommentDao) DeleteCommentByCommentId(CommentId uint) error {
	var comment *model.Comment
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id=?", CommentId).Delete(&comment).Error; err != nil {
			return err
		}
		return nil
	})
}
