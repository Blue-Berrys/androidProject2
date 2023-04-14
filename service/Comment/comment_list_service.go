package service

import (
	model2 "androidProject2/model/Comment"
	model4 "androidProject2/model/db"
	model3 "androidProject2/model/friendschat"
	model "androidProject2/model/user"
	"androidProject2/util"
	"errors"
	"log"
	"sync"
)

type CommentListResponse struct {
	CommentList []*util.Comment `json:"comment_list"`
}

type CommentListFlow struct {
	UserId        uint
	FriendsChatId uint
	data          []*util.Comment
	*CommentListResponse
}

func CommentList(userId uint, friendsChatId uint) (*CommentListResponse, error) {
	return NewCommentList(userId, friendsChatId).Do()
}

func NewCommentList(userId uint, friendsChatId uint) *CommentListFlow {
	return &CommentListFlow{UserId: userId, FriendsChatId: friendsChatId}
}

func (q *CommentListFlow) Do() (*CommentListResponse, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.CommentListResponse, nil
}

func (q *CommentListFlow) checkNum() error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	errChan := make(chan error, 2)
	defer close(errChan)
	go func() {
		wg.Done()
		//判断UserId是否合法
		var UserDao = model.NewUserDao()
		if !UserDao.QueryUserExistByUserId(q.UserId) {
			errStr := "token用户不存在"
			log.Println(errStr)
			err := errors.New(errStr)
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		//根据FriendsId查询这条朋友圈是否存在
		var FriendsChatDao = model3.NewFriendsChatDao()
		if !FriendsChatDao.ExistsFriendsChatById(q.FriendsChatId) {
			errStr := "FriendsChatId不存在"
			log.Println(errStr)
			err := errors.New(errStr)
			errChan <- err
		}
	}()

	wg.Wait()
	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}

func (q *CommentListFlow) prepareData() error {
	//先查Comment表
	var CommentDao = model2.NewCommentDao()
	var Comments = []*model4.Comment{}
	//查询所有FriendsChatId的评论
	if err := CommentDao.QueryCommentByFriendsChatId(q.FriendsChatId, &Comments); err != nil {
		return err
	}

	for _, Comment := range Comments {
		//查用户信息
		var UserDao = model.NewUserDao()
		var dbUser = model4.User{}
		if err := UserDao.QueryUserInfoById(q.UserId, &dbUser); err != nil {
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
		//构造util.Comment
		Comment := &util.Comment{
			Id:         Comment.ID,
			User:       *modelUser,
			Content:    Comment.CommentText,
			CreateDate: Comment.CreatedAt.Format("01-02 15:04:05"),
		}
		q.data = append(q.data, Comment)
	}

	return nil
}

func (q *CommentListFlow) packData() error {
	q.CommentListResponse = &CommentListResponse{CommentList: q.data}
	return nil
}
