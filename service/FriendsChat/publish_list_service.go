package service

import (
	"androidProject2/config"
	model2 "androidProject2/model/db"
	model3 "androidProject2/model/friendschat"
	model "androidProject2/model/user"
	"androidProject2/util"
	"errors"
	"log"
	"strings"
	"sync"
)

type PublishListResponse struct {
	FriendsChatList []*util.FriendsChat `json:"friendschat_list"`
}

type PublishListFlow struct {
	UserId uint
	SeenId uint
	data   []*util.FriendsChat
	*PublishListResponse
}

func PublishList(userId uint, seenId uint) (*PublishListResponse, error) {
	return NewPublishList(userId, seenId).Do()
}

func NewPublishList(userId uint, seenId uint) *PublishListFlow {
	return &PublishListFlow{UserId: userId, SeenId: seenId}
}

func (q *PublishListFlow) Do() (*PublishListResponse, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.PublishListResponse, nil
}

func (q *PublishListFlow) checkNum() error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	errChan := make(chan error, 2)
	defer close(errChan)

	var UserDao = model.NewUserDao()

	go func() {
		defer wg.Done()
		//判断UserId是否合法
		if !UserDao.QueryUserExistByUserId(q.UserId) {
			errStr := "token用户不存在"
			log.Println(errStr)
			errChan <- errors.New(errStr)
		}
	}()

	go func() {
		defer wg.Done()
		//判断SeenId是否合法
		if q.SeenId != 0 && !UserDao.QueryUserExistByUserId(q.SeenId) {
			errStr := "传入的user_id被看的人id不存在"
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

func (q *PublishListFlow) prepareData() error {
	//先查friends_chat表
	var FriendsChatDao = model3.NewFriendsChatDao()
	var FriendsChats = []*model2.FriendsChat{}
	if q.SeenId == 0 { //查询所有用户
		if err := FriendsChatDao.QueryAllFriendsChat(&FriendsChats); err != nil {
			return err
		}
	} else {
		if err := FriendsChatDao.QueryFriendsChatByUserId(q.SeenId, &FriendsChats); err != nil {
			return err
		}
	}
	for _, FriendsChat := range FriendsChats {
		//查用户信息
		var UserDao = model.NewUserDao()
		var dbUser = model2.User{}
		if err := UserDao.QueryUserInfoById(FriendsChat.UserId, &dbUser); err != nil {
			return err
		}
		//构造util.user
		modelUser := &util.User{
			Id:              dbUser.ID,
			Name:            dbUser.UserName,
			NickName:        dbUser.NickName,
			WorkCount:       dbUser.WorkCount,
			BackGroundImage: config.Miniourl + dbUser.BackgroundImage,
			Avatar:          config.Miniourl + dbUser.Avatar,
		}
		Fields := strings.Fields(FriendsChat.ImageUrl)
		ImageUrl := ""
		for i, field := range Fields {
			if i == 0 {
				ImageUrl = config.Miniourl + field
			} else {
				ImageUrl += " " + config.Miniourl + field
			}
		}
		oneFriendsChat := &util.FriendsChat{
			Id:         FriendsChat.ID,
			User:       *modelUser,
			ImageUrl:   ImageUrl,
			Content:    FriendsChat.Content,
			CreateDate: FriendsChat.CreatedAt.Format("01-02 15:04:05"),
		}
		q.data = append(q.data, oneFriendsChat)
	}

	return nil
}

func (q *PublishListFlow) packData() error {
	q.PublishListResponse = &PublishListResponse{FriendsChatList: q.data}
	return nil
}
