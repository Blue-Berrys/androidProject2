package service

import (
	model3 "androidProject2/model/friendschat"
	model "androidProject2/model/user"
	"errors"
	"log"
	"sync"
)

type DeleteResponse struct {
	UserId        uint
	FriendChatsId uint
}

func DeleteAction(userId uint, friendsChatId uint) error {
	return NewDeleteResponse(userId, friendsChatId).Do()
}

func NewDeleteResponse(userId uint, friendsChatId uint) *DeleteResponse {
	return &DeleteResponse{UserId: userId, FriendChatsId: friendsChatId}
}

func (q *DeleteResponse) Do() error {
	if err := q.checkNum(); err != nil {
		return err
	}
	if err := q.prepareData(); err != nil {
		return err
	}
	return nil
}

func (q *DeleteResponse) checkNum() error {
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
		//根据FriendsId查询这条朋友圈是否存在
		var FriendsChatDao = model3.NewFriendsChatDao()
		if !FriendsChatDao.ExistsFriendsChatById(q.FriendChatsId) {
			errStr := "FriendsChat Not Exists"
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

func (q *DeleteResponse) prepareData() error {
	//判断这一条UserId是否对应FriendsId
	var FriendsChatDao = model3.NewFriendsChatDao()
	if !FriendsChatDao.ExistsFriendsChatIdAndUserId(q.UserId, q.FriendChatsId) {
		return errors.New("你不能删除别人的朋友圈")
	}

	//删除这条朋友圈
	if err := FriendsChatDao.DeleteOneFriendsChatById(q.FriendChatsId); err != nil {
		return err
	}
	return nil
}
