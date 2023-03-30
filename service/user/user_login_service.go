package service

import (
	"androidProject2/middleware/JWT"
	model2 "androidProject2/model/db"
	model "androidProject2/model/user"
	"errors"
)

var (
	MaxUserNameLen = 100
)

type LoginResponse struct {
	userId int64  `json:"user_id"`
	token  string `json:"token"`
}

type QueryUserLoginFlow struct {
	username string
	password string
	data     *LoginResponse
	userid   int64
	token    string
}

func QueryUserLogin(username, password string) (*LoginResponse, error) {
	return NewQueryUserLoginFlow(username, password).Do() //创建一个结构，然后返回结构的方法
}

func NewQueryUserLoginFlow(username, password string) *QueryUserLoginFlow {
	return &QueryUserLoginFlow{username: username, password: password}
}

func (q *QueryUserLoginFlow) Do() (*LoginResponse, error) {
	//检测参数合法性
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	//准备数据
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	//打包数据
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.data, nil
}

func (q *QueryUserLoginFlow) checkNum() error {
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

func (q *QueryUserLoginFlow) prepareData() error {
	//调用models查数据库
	userLoginDAO := model.NewUserDao()
	var login model2.User
	if err := userLoginDAO.QueryUserLogin(q.username, q.password, &login); err != nil {
		return err
	}
	q.userid = login.UserId

	//登录成功，颁发token
	token, err := JWT.GetToken(uint(q.userid))
	if err != nil {
		return err
	}
	q.token = token
	return nil
}

func (q *QueryUserLoginFlow) packData() error {
	q.data = &LoginResponse{
		userId: q.userid,
		token:  q.token,
	}
	return nil
}
