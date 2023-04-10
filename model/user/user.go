package model

import (
	"androidProject2/config"
	model "androidProject2/model/db"
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
)

type UserDao struct {
}

var (
	userDao  *UserDao
	userOnce sync.Once
)

func NewUserDao() *UserDao {
	userOnce.Do(func() {
		userDao = new(UserDao)
	})
	return userDao
}

func (u *UserDao) AddUserInfo(userinfo *model.User) error {
	if userinfo == nil {
		return errors.New("空指针错误")
	}
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(userinfo).Error; err != nil {
			return err
		}
		return nil
	})

}

func (u *UserDao) QueryUserLogin(username string, login *model.User) error {
	if login == nil {
		return errors.New("登录信息结构体指针为空")
	}
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_name=?", username).First(&login).Error; err != nil {
			return err
		}
		if login.ID == 0 {
			return errors.New("用户未找到,账号或密码出错")
		}
		return nil
	})
}

func (u *UserDao) QueryUserExistByUserName(username string) bool {
	var user model.User
	if err := model.DB.Where("user_name=?", username).First(&user).Error; err != nil {
		log.Println(err)
	}
	if user.ID == 0 {
		return false
	}
	return true
}

func (u *UserDao) QueryUserExistByUserId(userId uint) bool {
	var user model.User
	if err := model.DB.Where("id=?", userId).First(&user).Error; err != nil {
		log.Println(err)
	}
	if user.ID == 0 {
		return false
	}
	return true
}

func (u *UserDao) IsCorrectUserIdAndPassword(userId uint, password string) bool {
	var user model.User
	if err := model.DB.Where("id=?", userId).Where("password=?", password).First(&user).Error; err != nil {
		log.Println(err)
	}
	if user.ID == 0 {
		return false
	}
	return true
}

func (u *UserDao) QueryUserInfoById(userId uint, user *model.User) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", userId).First(&user).Error; err != nil {
			return err
		}
		return nil
	})
}

func (u *UserDao) UpdateUserInfo(user *model.User, userId uint, signature string, avatar string, background_image string) error {
	avatar = config.PlayUrlPrefix + avatar + ".jpg"
	background_image = config.PlayUrlPrefix + background_image + ".jpg"
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&user).Where("id=?", userId).Updates(map[string]interface{}{"signature": signature, "avatar": avatar, "background_image": background_image}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (u *UserDao) AddOneWorkCountByUserId(user *model.User, userId uint, workcount int64) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&user).Where("id=?", userId).Updates(map[string]interface{}{"work_count": workcount}).Error; err != nil {
			return err
		}
		return nil
	})
}
