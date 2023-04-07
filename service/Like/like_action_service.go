package service

import (
	"androidProject2/cache/Redis"
	model3 "androidProject2/model/friendschat"
	model2 "androidProject2/model/like"
	model "androidProject2/model/user"
	"errors"
)

const (
	LIKE    = 1
	DISLIKE = 2
)

type LikeState struct {
	UserId        uint
	FriendsChatId uint
	actionType    int
}

func LikeAction(userid uint, friendschatid uint, actiontype int) error {
	return NewLikeState(userid, friendschatid, actiontype).Do()
}

func NewLikeState(userid uint, friendschatid uint, actiontype int) *LikeState {
	return &LikeState{
		UserId:        userid,
		FriendsChatId: friendschatid,
		actionType:    actiontype,
	}
}

func (q *LikeState) Do() error {
	if err := q.Parameters(); err != nil {
		return err
	}
	//因为前面已经判断了,只能是LIKE or UNLIKE
	if q.actionType == LIKE {
		if err := q.LikeFriendsChat(); err != nil {
			return err
		}
	} else {
		if err := q.UnLikeFriendsChat(); err != nil {
			return err
		}
	}
	return nil
}

func (q *LikeState) Parameters() error {
	//根据UserId查询用户是否存在
	var UserDao = model.NewUserDao()
	if !UserDao.QueryUserExistByUserId(q.UserId) {
		return errors.New("User Not Exists")
	}
	var FriendsChatDao = model3.NewFriendsChatDao()
	if !FriendsChatDao.ExistsFriendsChatById(q.FriendsChatId) {
		return errors.New("FriendsChat Not Exists")
	}
	//判断actionType是否合法
	if q.actionType != LIKE && q.actionType != DISLIKE {
		return errors.New("actionType illegal")
	}
	return nil
}

func (q *LikeState) LikeFriendsChat() error {
	//先判断这个记录存不存在
	var RedisDao = Redis.NewRedisDao()
	var LikeDao = model2.NewLikeDao()
	ok, err := RedisDao.GetLikeState(q.UserId, q.FriendsChatId)
	if err != nil {
		//在Redis中没找到
		//去数据库中找like表
		ok = LikeDao.IsLikeByUserIdAndFriendschatId(q.UserId, q.FriendsChatId)
	}
	if ok {
		return errors.New("you can't like again after you've already liked it")
	}
	//在Mysql里点赞，增加一条记录
	if err := LikeDao.AddOneLikeByFriendschatIdAndUserId(q.FriendsChatId, q.UserId); err != nil {
		return err
	}
	//在redis里置为喜欢
	if err := RedisDao.UpdatePostLike(q.UserId, q.FriendsChatId, true); err != nil {
		return err
	}
	//在redis里给这条朋友圈点赞的人数加1
	if err := RedisDao.AddOneLikeNumByfriendschatId(q.FriendsChatId); err != nil {
		return err
	}
	return nil
}
func (q *LikeState) UnLikeFriendsChat() error {
	//先判断这个记录存不存在
	var RedisDao = Redis.NewRedisDao()
	var LikeDao = model2.NewLikeDao()
	ok, err := RedisDao.GetLikeState(q.UserId, q.FriendsChatId)
	if err != nil {
		//在Redis中没找到
		//去数据库中找like表
		ok = LikeDao.IsLikeByUserIdAndFriendschatId(q.UserId, q.FriendsChatId)
	}
	if !ok {
		return errors.New("you can't cancel like again after you've already disliked it")
	}
	//在Mysql里取消点赞，减少一条记录
	if err := LikeDao.SubOneLikeByFriendschatIdAndUserId(q.FriendsChatId, q.UserId); err != nil {
		return err
	}
	//在redis里置为喜欢
	if err := RedisDao.UpdatePostLike(q.UserId, q.FriendsChatId, false); err != nil {
		return err
	}
	//在redis里给这条朋友圈点赞的人数减1
	if err := RedisDao.SubOneLikeNumByfriendschatId(q.FriendsChatId); err != nil {
		return err
	}
	return nil
}
