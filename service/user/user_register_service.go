package service

import (
	"androidProject2/middleware/JWT"
	model2 "androidProject2/model/db"
	model "androidProject2/model/user"
	"errors"
)

type RegisterResponse struct {
	userId int64  `json:"user_id"`
	token  string `json:"token"`
}

type QueryUserRegisterFlow struct {
	username string
	password string
	data     *RegisterResponse
	userId   int64
	token    string
}

func NewQueryUserRegister(username, password string) (*RegisterResponse, error) {
	return NewQueryUserRegisterFLow(username, password).Do()
}

func NewQueryUserRegisterFLow(username, password string) *QueryUserRegisterFlow {
	return &QueryUserRegisterFlow{username: username, password: password}
}

func (q *QueryUserRegisterFlow) Do() (*RegisterResponse, error) {
	//对参数进行合法性检验
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.data, nil
}

func (q *QueryUserRegisterFlow) checkNum() error {
	if q.username == "" {
		return errors.New("用户名为空")
	}
	if q.password == "" {
		return errors.New("密码为空")
	}
	if len(q.username) > MaxUserNameLen {
		return errors.New("用户名太长")
	}
	return nil
}

func (q *QueryUserRegisterFlow) prepareData() error {
	userRegisterDao := model.NewUserDao()
	//判断用户名是否存在
	exist := userRegisterDao.QueryUserExistByUserName(q.username)
	if exist {
		return errors.New("用户名已存在，请重新输入用户名")
	}

	//增加这个用户
	user := model2.User{UserName: q.username, Password: q.password}
	if err := userRegisterDao.AddUserInfo(&user); err != nil {
		return err
	}

	//颁发token
	token, err := JWT.GetToken(uint(q.userId))
	if err != nil {
		return err
	}
	q.token = token
	q.userId = user.UserId
	return nil
}

func (q *QueryUserRegisterFlow) packData() error {
	q.data = &RegisterResponse{
		userId: q.userId,
		token:  q.token,
	}
	return nil
}
