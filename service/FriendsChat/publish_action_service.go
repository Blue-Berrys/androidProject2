package service

import (
	"androidProject2/config"
	"androidProject2/middleware/minio"
	model2 "androidProject2/model/db"
	"androidProject2/model/friendschat"
	model "androidProject2/model/user"
	"androidProject2/util"
	"errors"
	"mime/multipart"
)

type PublishActionResponse struct {
	FriendsChat *util.FriendsChat `json:"friendschat"`
}

type PublishActionFlow struct {
	UserId     uint
	Content    string
	ActionType int
	Images     []*multipart.FileHeader
	ImageName  []string

	ImageUrl string

	data *util.FriendsChat
	*PublishActionResponse
}

func PublishAction(userId uint, content string, action_type int, imageName []string, images []*multipart.FileHeader) (*PublishActionResponse, error) {
	return NewPublishAction(userId, content, action_type, imageName, images).Do()
}

func NewPublishAction(userId uint, context string, action_type int, imageName []string, images []*multipart.FileHeader) *PublishActionFlow {
	return &PublishActionFlow{UserId: userId, Content: context, ActionType: action_type, ImageName: imageName, Images: images}
}

func (q *PublishActionFlow) Do() (*PublishActionResponse, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.PublishActionResponse, nil
}

func (q *PublishActionFlow) checkNum() error {
	//查action_type是否只有1或2
	if q.ActionType != 1 && q.ActionType != 2 {
		return errors.New("action_type不是1和2，输入错误")
	}
	//查数据库，这个id是否存在
	var UserDao = model.NewUserDao()
	if q.UserId == 0 || !UserDao.QueryUserExistByUserId(q.UserId) {
		return errors.New("UserId用户不存在")
	}
	return nil
}

func (q *PublishActionFlow) prepareData() error {
	//根据id查询用户信息
	var UserDao = model.NewUserDao()
	dbUser := model2.User{}
	if err := UserDao.QueryUserInfoById(q.UserId, &dbUser); err != nil {
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
	//增加朋友圈记录
	var ResImageUrl string
	if q.ActionType == 1 { //有图片
		for i, image := range q.Images {
			if err := minio.ImageToMinio(image, q.ImageName[i]); err != nil {
				return err
			}
			if i == 0 {
				q.ImageUrl = config.PlayUrlPrefix + q.ImageName[i] + ".jpg"
				ResImageUrl = config.Miniourl + config.PlayUrlPrefix + q.ImageName[i] + ".jpg"
			} else {
				q.ImageUrl += " " + config.PlayUrlPrefix + q.ImageName[i] + ".jpg"
				ResImageUrl += " " + config.Miniourl + config.PlayUrlPrefix + q.ImageName[i] + ".jpg"
			}
		}
	} else {
		q.ImageUrl = ""
		ResImageUrl = ""
	}
	//插入FriendsChat
	var FriendsChatDao = friendschat.NewFriendsChatDao()
	var friendschat = model2.FriendsChat{
		UserId:   q.UserId,
		Content:  q.Content,
		ImageUrl: q.ImageUrl,
	}
	if err := FriendsChatDao.AddFriendsChat(&friendschat); err != nil {
		return err
	}

	//构造返回Json
	q.data = &util.FriendsChat{
		Id:            friendschat.ID,
		User:          *modelUser,
		ImageUrl:      ResImageUrl,
		FavoriteCount: 0,
		IsFavorite:    false,
		Content:       q.Content,
		CreateDate:    friendschat.CreatedAt.Format("01-02 15:04:05"),
	}
	return nil
}

func (q *PublishActionFlow) packData() error {
	q.PublishActionResponse = &PublishActionResponse{FriendsChat: q.data}
	return nil
}
