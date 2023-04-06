package service

import (
	"androidProject2/cache/Redis"
	model2 "androidProject2/model/db"
	model3 "androidProject2/model/friendschat"
	model4 "androidProject2/model/like"
	model "androidProject2/model/user"
	"androidProject2/util"
	"errors"
	"fmt"
)

type PublishListResponse struct {
	FriendsChatList []*util.FriendsChat `json:"friendschat"`
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
	//判断UserId是否合法
	var UserDao = model.NewUserDao()
	if !UserDao.QueryUserExistByUserId(q.UserId) {
		return errors.New("token用户不存在")
	}
	//判断SeenId是否合法
	if q.SeenId != 0 && !UserDao.QueryUserExistByUserId(q.SeenId) {
		return errors.New("传入的user_id被看的人id不存在")
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
		//构造model.user
		modelUser := &util.User{
			Id:              dbUser.ID,
			Name:            dbUser.UserName,
			Signature:       dbUser.Signature,
			WorkCount:       dbUser.WorkCount,
			BackGroundImage: dbUser.BackgroundImage,
			Avatar:          dbUser.Avatar,
		}
		//先用Redis查这个视频点赞次数
		var RedisDao = Redis.NewRedisDao()
		num := 0
		Redisnum, err := RedisDao.GetLikeNumByfriendschatId(FriendsChat.ID)
		var LikeDao = model4.NewLikeDao()
		if err != nil {
			//找不到去Like表里找
			Sqlnum, err := LikeDao.QueryLenFavorVideoListByVideoId(int64(FriendsChat.ID))
			if err != nil {
				return err
			}
			num = Sqlnum
		} else {
			num = Redisnum
		}
		//查询本人是否点赞这个朋友圈
		isLike := LikeDao.IsLikeByUserIdAndVideoId(q.UserId, FriendsChat.ID)

		oneFriendsChat := &util.FriendsChat{
			Id:            FriendsChat.ID,
			User:          *modelUser,
			ImageUrl:      FriendsChat.ImageUrl,
			FavoriteCount: int64(num),
			IsFavorite:    isLike,
		}
		fmt.Println(oneFriendsChat)
	}

	return nil
}

func (q *PublishListFlow) packData() error {

	return nil
}
