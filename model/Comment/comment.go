package model

import (
	model "androidProject2/model/db"
	"errors"
	"gorm.io/gorm"
	"log"
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

// 根据评论Id查询评论是否存在
func (u *CommentDao) IsExistsCommentByCommentId(CommentId uint) bool {
	var comment *model.Comment
	if err := model.DB.Where("id=?", CommentId).First(&comment).Error; err != nil {
		log.Println(err)
	}
	if comment.ID == 0 {
		return false
	}
	return true
}

// 根据评论Id和朋友圈Id查询是否存在这个记录
func (u *CommentDao) IsCorrectCommentIdAndCommentId(FriendsChatId, CommentId uint) bool {
	var comment *model.Comment
	if err := model.DB.Where("id=?", CommentId).Where("friend_chat_id=?", FriendsChatId).First(&comment).Error; err != nil {
		log.Println(err)
	}
	if comment.ID == 0 {
		return false
	}
	return true
}
