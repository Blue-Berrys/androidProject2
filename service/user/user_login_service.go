package service

import (
	"androidProject2/middleware/Bcrypt"
	"androidProject2/middleware/JWT"
	model2 "androidProject2/model/db"
	model "androidProject2/model/user"
	"errors"
	"fmt"
	"log"
	"sync"
)

var (
	MaxUserNameLen = 100
)

type LoginResponse struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
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

func (q *QueryUserLoginFlow) prepareData() error {
	//调用models查数据库
	userLoginDAO := model.NewUserDao()
	var login model2.User
	fmt.Println("名字：", q.username)
	if err := userLoginDAO.QueryUserLogin(q.username, &login); err != nil {
		return errors.New("用户名错误")
	}
	if ok := Bcrypt.QueryEqualEncryptAndPassword(login.Password, q.password); !ok {
		return errors.New("密码错误")
	}

	//登录成功，颁发token
	token, err := JWT.GetToken(login.ID, q.password)
	if err != nil {
		return err
	}

	q.token = token
	q.userid = int64(login.ID)
	return nil
}

func (q *QueryUserLoginFlow) packData() error {
	q.data = &LoginResponse{
		UserId: q.userid,
		Token:  q.token,
	}
	return nil
}
