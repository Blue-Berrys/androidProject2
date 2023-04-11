package service

import (
	"androidProject2/middleware/Bcrypt"
	"androidProject2/middleware/JWT"
	model2 "androidProject2/model/db"
	model "androidProject2/model/user"
	"errors"
	"log"
	"sync"
)

type RegisterResponse struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
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
	wg := sync.WaitGroup{}
	wg.Add(3)

	errChan := make(chan error, 3)
	defer close(errChan)

	go func() {
		defer wg.Done()
		if q.username == "" {
			errStr := "用户名为空"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
	}()

	go func() {
		defer wg.Done()
		if q.password == "" {
			errStr := "密码为空"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
	}()

	go func() {
		defer wg.Done()
		if len(q.username) > MaxUserNameLen {
			errStr := "用户名太长"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
	}()

	wg.Wait()

	if len(errChan) > 0 {
		return <-errChan
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
	befPassword := q.password
	q.password = Bcrypt.EncryptionByPassword(q.password) //生成加盐后密码

	//增加这个用户
	user := model2.User{UserName: q.username, Password: q.password}
	if err := userRegisterDao.AddUserInfo(&user); err != nil {
		return err
	}
	q.userId = int64(user.ID)
	log.Println("q.userId:", q.userId)
	//颁发token
	token, err := JWT.GetToken(uint(q.userId), befPassword)
	if err != nil {
		return err
	}
	q.token = token
	q.userId = int64(user.ID)
	return nil
}

func (q *QueryUserRegisterFlow) packData() error {
	q.data = &RegisterResponse{
		UserId: q.userId,
		Token:  q.token,
	}
	return nil
}
