package service

import (
	model "androidProject2/model/user"
	"errors"
	"log"
	"sync"
)

type LevelAction struct {
	UserId       uint
	ModifyUserId uint
}

func LevelActionService(userId uint, modifyUserId uint) error {
	return NewLevelActionService(userId, modifyUserId).Do()
}

func NewLevelActionService(userId uint, modifyUserId uint) *LevelAction {
	return &LevelAction{UserId: userId, ModifyUserId: modifyUserId}
}

func (q *LevelAction) Do() error {
	if err := q.Parameters(); err != nil {
		return err
	}
	if err := q.prepareData(); err != nil {
		return err
	}
	return nil
}

func (q *LevelAction) Parameters() error {
	//查询UserId和ModifyUserId是否存在
	wg := sync.WaitGroup{}
	wg.Add(2)

	errChan := make(chan error, 2)
	defer close(errChan)

	go func() {
		defer wg.Done()
		//根据UserId查询用户是否存在
		var UserDao = model.NewUserDao()
		if !UserDao.QueryUserExistByUserId(q.UserId) {
			errStr := "User Not Exists"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
	}()

	go func() {
		defer wg.Done()
		//根据ModifyUserId是否存在
		var UserDao = model.NewUserDao()
		if !UserDao.QueryUserExistByUserId(q.ModifyUserId) {
			errStr := "ModifyUser Not Exists"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
	}()

	wg.Wait()

	if len(errChan) > 0 {
		return <-errChan
	}

	//判断token中的UserId是否是管理员
	var UserDao = model.NewUserDao()
	UserIdLevel, err := UserDao.QueryLevelByUserId(q.UserId)
	if err != nil {
		return err
	}
	if UserIdLevel != 2 {
		return errors.New("token用户Id不是管理员")
	}
	return nil
}

func (q *LevelAction) prepareData() error {
	var UserDao = model.NewUserDao()
	//查询modifyUser_level
	bef_level, err := UserDao.QueryLevelByUserId(q.ModifyUserId)
	if err != nil {
		return err
	}
	if bef_level != 0 && bef_level != 1 {
		return errors.New("修改的用户id不是普通用户也不是封号用户")
	}
	//修改modifyUser_level
	if err := UserDao.SetUserLevelByUserId(q.ModifyUserId, bef_level^1); err != nil {
		return err
	}
	after_level, err := UserDao.QueryLevelByUserId(q.ModifyUserId)
	if err != nil {
		return err
	}
	if after_level != bef_level^1 {
		return errors.New("level没有修改成功")
	}
	return nil
}
