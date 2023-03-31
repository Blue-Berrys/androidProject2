package service

import (
	model2 "androidProject2/model/db"
	model "androidProject2/model/user"
	"androidProject2/util"
	"errors"
)

type InfoResponse struct {
	MyUser *util.User
}

type QueryUserInfoFlow struct {
	Id   uint
	data *util.User
	*InfoResponse
}

func QueryUserInfo(id uint) (*InfoResponse, error) {
	return NewQueryUserInfo(id).Do()
}

func NewQueryUserInfo(id uint) *QueryUserInfoFlow {
	return &QueryUserInfoFlow{Id: id}
}

func (q *QueryUserInfoFlow) Do() (*InfoResponse, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.InfoResponse, nil
}

func (q *QueryUserInfoFlow) checkNum() error {
	//查询用户id是否在数据库中
	var UserDao = model.NewUserDao()
	if q.Id == 0 || UserDao.QueryUserExistByUserId(q.Id) {
		return errors.New("用户不存在")
	}
	return nil
}

func (q *QueryUserInfoFlow) prepareData() error {
	//根据id查询用户信息
	var UserDao = model.NewUserDao()
	var dbUser = model2.User{}
	if err := UserDao.QueryUserInfoById(q.Id, &dbUser); err != nil {
		return err
	}
	//设置信息
	q.data.Avatar = dbUser.Avatar
	q.data.Signature = dbUser.Signature
	q.data.BackGroundImage = dbUser.BackgroundImage
	q.data.UserId = dbUser.UserId
	q.data.WorkCount = dbUser.WorkCount
	return nil
}

func (q *QueryUserInfoFlow) packData() error {
	q.InfoResponse = &InfoResponse{MyUser: q.data}
	return nil
}
