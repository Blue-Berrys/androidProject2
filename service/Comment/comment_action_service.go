package service

import (
	model4 "androidProject2/model/Comment"
	model2 "androidProject2/model/db"
	model3 "androidProject2/model/friendschat"
	model "androidProject2/model/user"
	"androidProject2/util"
	"errors"
	"log"
	"sync"
)

const (
	ADDCOMMENT = 1
	DELCOMMENT = 2
)

type CommentActionResponse struct {
	Comment *util.Comment `json:"comment"`
}

type CommentActionFlow struct {
	UserId        uint
	FriendsChatId uint
	ActionType    int
	CommentText   string
	CommentId     uint

	data *util.Comment
	*CommentActionResponse
}

func CommentAction(userId uint, friendschatId uint, actionType int, commentText string, commentId uint) (*CommentActionResponse, error) {
	return NewCommentAction(userId, friendschatId, actionType, commentText, commentId).Do()
}

func NewCommentAction(userId uint, friendschatId uint, actionType int, commentText string, commentId uint) *CommentActionFlow {
	return &CommentActionFlow{UserId: userId, FriendsChatId: friendschatId, ActionType: actionType, CommentText: commentText, CommentId: commentId}
}

func (q *CommentActionFlow) Do() (*CommentActionResponse, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.CommentActionResponse, nil
}

func (q *CommentActionFlow) checkNum() error {
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

func (q *CommentActionFlow) prepareData() error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	errChan := make(chan error, 4)
	defer close(errChan)

	var UserDao = model.NewUserDao()
	var dbUser = model2.User{}
	var CreateData string
	//查用户信息
	go func() {
		defer wg.Done()
		if err := UserDao.QueryUserInfoById(q.UserId, &dbUser); err != nil {
			log.Println(err)
			errChan <- err
		}
	}()
	//增加评论
	if q.ActionType == ADDCOMMENT {
		go func() {
			defer wg.Done()
			var CommentDao = model4.NewCommentDao()
			var dbComment = model2.Comment{UserId: int64(q.UserId), FriendChatId: int64(q.FriendsChatId), CommentText: q.CommentText}
			q.CommentId = dbComment.ID
			if err := CommentDao.AddCommentInfo(&dbComment); err != nil {
				log.Println(err)
				errChan <- err
			}
			CreateData = dbComment.CreatedAt.Format("01-02 15:04:05")
		}()
	} else { //删除评论
		go func() {
			defer wg.Done()
			var CommentDao = model4.NewCommentDao()
			//需要判断是否存在这条评论ID是否存在
			if !CommentDao.IsExistsCommentByCommentId(q.CommentId) {
				errStr := "评论的Id不存在"
				log.Println(errStr)
				errChan <- errors.New(errStr)
			}
			//判断删除的评论ID确实是这个朋友圈ID的
			if !CommentDao.IsCorrectCommentIdAndCommentId(q.FriendsChatId, q.CommentId) {
				errStr := "评论的Id和朋友圈Id不对应"
				log.Println(errStr)
				errChan <- errors.New(errStr)
			}
			if err := CommentDao.DeleteCommentByCommentId(q.CommentId); err != nil {
				log.Println(err)
				errChan <- err
			}
		}()
	}
	wg.Wait()

	//构造util.User
	modelUser := &util.User{
		Id:              dbUser.ID,
		Name:            dbUser.UserName,
		NickName:        dbUser.NickName,
		WorkCount:       dbUser.WorkCount,
		BackGroundImage: dbUser.BackgroundImage,
		Avatar:          dbUser.Avatar,
		Level:           dbUser.Level,
	}

	//构造util.Comment
	Comment := &util.Comment{
		Id:         q.CommentId,
		User:       *modelUser,
		Content:    q.CommentText,
		CreateDate: CreateData,
	}
	q.data = Comment

	if len(errChan) > 0 {
		return <-errChan
	}
	return nil
}

func (q *CommentActionFlow) packData() error {
	q.CommentActionResponse = &CommentActionResponse{Comment: q.data}
	return nil
}
