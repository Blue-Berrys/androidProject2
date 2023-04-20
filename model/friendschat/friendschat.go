package model

import (
	model "androidProject2/model/db"
	"errors"
	"gorm.io/gorm"
	"log"
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
		if err := tx.Create(info).Error; err != nil {
			return err
		}
		return nil
	})
}

// 根据userid查他的朋友圈列表
func (q *FriendsChatDao) QueryFriendsChatByUserId(userId uint, info *[]*model.FriendsChat) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.FriendsChat{}).Where("user_id=?", userId).Find(info).Error; err != nil {
			return err
		}
		return nil
	})
}

// 查询整个朋友圈所有内容
func (q *FriendsChatDao) QueryAllFriendsChat(info *[]*model.FriendsChat) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Find(info).Error; err != nil {
			return err
		}
		return nil
	})
}

// 根据朋友圈id查这一条朋友圈
func (q *FriendsChatDao) ExistsFriendsChatById(id uint) bool {
	var friendschat = &model.FriendsChat{}
	if err := model.DB.Where("id=?", id).First(friendschat).Error; err != nil {
		log.Println(err)
	}
	if friendschat.ID == 0 {
		return false
	}
	return true
}

// 根据朋友圈id和用户id查是否存在这条记录
func (q *FriendsChatDao) ExistsFriendsChatIdAndUserId(userId uint, friendsChatId uint) bool {
	var friendschat = &model.FriendsChat{}
	if err := model.DB.Where("id=?", friendsChatId).Where("user_id", userId).First(friendschat).Error; err != nil {
		log.Println(err)
	}
	if friendschat.ID == 0 {
		return false
	}
	return true
}

// 根据朋友圈id删除这条朋友圈
func (q *FriendsChatDao) DeleteOneFriendsChatById(id uint) error {
	var friendschat = &model.FriendsChat{}
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(friendschat).Delete("id=?", id).Error; err != nil {
			return err
		}
		return nil
	})
}
