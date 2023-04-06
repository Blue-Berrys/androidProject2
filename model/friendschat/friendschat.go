package friendschat

import (
	model "androidProject2/model/db"
	"errors"
	"gorm.io/gorm"
	"sync"
)

type FriendsChatDao struct {
}

var (
	friendschatDao  *FriendsChatDao
	friendschatOnce sync.Once
)

func NewFriendsChatDao() *FriendsChatDao {
	friendschatOnce.Do(func() {
		friendschatDao = new(FriendsChatDao)
	})
	return friendschatDao
}

// 增加一条朋友圈
func (q *FriendsChatDao) AddFriendsChat(info *model.FriendsChat) error {
	if info == nil {
		return errors.New("空指针错误")
	}
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := model.DB.Create(info).Error; err != nil {
			return err
		}
		return nil
	})
}

// 根据userid查他的朋友圈列表
func (q *FriendsChatDao) QueryFriendsChatByUserId(userId uint, info *[]*model.FriendsChat) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := model.DB.Model(&model.FriendsChat{}).Where("user_id=?", userId).Find(info).Error; err != nil {
			return err
		}
		return nil
	})
}

// 查询整个朋友圈所有内容
func (q *FriendsChatDao) QueryAllFriendsChat(info *[]*model.FriendsChat) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := model.DB.Find(&model.FriendsChat{}).Error; err != nil {
			return err
		}
		return nil
	})
}
