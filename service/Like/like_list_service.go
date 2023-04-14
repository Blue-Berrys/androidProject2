package service

import (
	model3 "androidProject2/model/db"
	model2 "androidProject2/model/friendschat"
	model4 "androidProject2/model/like"
	model "androidProject2/model/user"
	"androidProject2/util"
	"errors"
	"log"
	"sync"
)

type LikeListsResponse struct {
	LikeList []*util.Like `json:"like_list"`
}

type LikeListFlow struct {
	UserId        uint
	FriendsChatId uint
	data          []*util.Like
	*LikeListsResponse
}

func LikeList(userId uint, friendsChatId uint) (*LikeListsResponse, error) {
	return NewLikeList(userId, friendsChatId).Do()
}

func NewLikeList(userId uint, friendsChatId uint) *LikeListFlow {
	return &LikeListFlow{UserId: userId, FriendsChatId: friendsChatId}
}

func (q *LikeListFlow) Do() (*LikeListsResponse, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.LikeListsResponse, nil
}

func (q *LikeListFlow) checkNum() error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	errChan := make(chan error, 2)
	defer close(errChan)

	go func() {
		defer wg.Done()
		//判断UserId是否合法
		var UserDao = model.NewUserDao()
		if !UserDao.QueryUserExistByUserId(q.UserId) {
			errStr := "token用户不存在"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
	}()

	go func() {
		defer wg.Done()
		//判断friendschat_id是否合法
		var FriendsChatDao = model2.NewFriendsChatDao()
		if !FriendsChatDao.ExistsFriendsChatById(q.FriendsChatId) {
			errStr := "填入的FriendsChatId不存在"
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

func (q *LikeListFlow) prepareData() error {
	//先查like表
	dblikelist := []*model3.Like{}
	var LikeDao = model4.NewLikeDao()
	if err := LikeDao.QueryAllLikeInfoByFiendsChatId(q.FriendsChatId, &dblikelist); err != nil {
		return err
	}
	for _, dblike := range dblikelist {
		//查User表
		var UserDao = model.NewUserDao()
		dbUser := model3.User{}
		if err := UserDao.QueryUserInfoById(uint(dblike.UserId), &dbUser); err != nil {
			return err
		}
		//构造util.user
		modelUser := &util.User{
			Id:              dbUser.ID,
			Name:            dbUser.UserName,
			NickName:        dbUser.NickName,
			WorkCount:       dbUser.WorkCount,
			BackGroundImage: dbUser.BackgroundImage,
			Avatar:          dbUser.Avatar,
		}
		//构造util.Like
		onelike := &util.Like{
			Id:   dblike.ID,
			User: *modelUser,
		}
		q.data = append(q.data, onelike)
	}

	return nil
}

func (q *LikeListFlow) packData() error {
	q.LikeListsResponse = &LikeListsResponse{LikeList: q.data}
	return nil
}
