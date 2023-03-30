package model

import (
	model "androidProject2/model/db"
	"errors"
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
	return model.DB.Create(userinfo).Error
}

func (u *UserDao) QueryUserLogin(username, password string, login *model.User) error {
	if login == nil {
		return errors.New("登录信息结构体指针为空")
	}
	model.DB.Where("username=? and password=?", username, password).First(&login)
	if login.ID == 0 {
		return errors.New("用户未找到,账号或密码出错")
	}
	return nil
}

func (u *UserDao) QueryUserExistByUserName(username string) bool {
	var userLogin model.User
	model.DB.Where("username=?", username).First(&userLogin)
	if userLogin.ID == 0 {
		return false
	}
	return true
}
