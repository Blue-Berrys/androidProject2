package model

import (
	model "androidProject2/model/db"
	"errors"
	"gorm.io/gorm"
	"sync"
)

type LikeDao struct {
}

var (
	likeDao  *LikeDao
	likeOnce sync.Once
)

func NewLikeDao() *LikeDao {
	likeOnce.Do(func() {
		likeDao = new(LikeDao)
	})
	return likeDao
}

// 判断是否存在userId点赞了FriendschatId
// 找不到就是没点赞
func (u *LikeDao) IsLikeByUserIdAndVideoId(userId uint, FriendschatId uint) bool {
	var like model.Like
	model.DB.Where("user_id = ?", userId).Where("friends_chat_id = ?", FriendschatId).First(&like)
	if like.ID == 0 {
		return false //没找到
	}
	return true
}

// 增加一个赞
func (u *LikeDao) AddOneLikeByFriendschatIdAndUserId(FriendschatId, UserId uint) error {
	var like = &model.Like{FriendsChatId: int64(FriendschatId), UserId: int64(UserId)}
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&like).Error; err != nil {
			return err
		}
		return nil
	})
}

// 减少一个赞
func (u *LikeDao) SubOneLikeByFriendschatIdAndUserId(FriendschatId, UserId uint) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("friends_chat_id=?", FriendschatId).Where("user_id=?", UserId).Delete(&model.Like{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// 查看视频点赞总数
func (u *LikeDao) QueryLenFavorVideoListByVideoId(FriendschatId int64) (int, error) {
	var likeList *[]*model.Like
	err := model.DB.Where("FriendschatId=?", FriendschatId).Find(&likeList).Error
	if err != nil {
		return 0, err
	}
	if len(*likeList) == 0 {
		return 0, errors.New("没有人给这个视频点赞")
	}
	return len(*likeList), nil
}
