package service

import (
	model2 "androidProject2/model/db"
	model "androidProject2/model/user"
	"androidProject2/util"
	"errors"
)

type UserLevelResponse struct {
	LevelResponse []*util.User `json:"user_list"`
}

type LevelListFlow struct {
	UserId uint
	Level  int
	data   []*util.User
	*UserLevelResponse
}

func QueryLevelList(userId uint, level int) (*UserLevelResponse, error) {
	return NewQueryLevelList(userId, level).Do()
}

func NewQueryLevelList(userId uint, level int) *LevelListFlow {
	return &LevelListFlow{UserId: userId, Level: level}
}

func (q *LevelListFlow) Do() (*UserLevelResponse, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.UserLevelResponse, nil
}

func (q *LevelListFlow) checkNum() error {
	//判断token中的UserId是否存在
	UserDao := model.NewUserDao()
	if !UserDao.QueryUserExistByUserId(q.UserId) {
		return errors.New("token用户Id不存在")
	}

	//判断token中的UserId是否是管理员
	UserIdLevel, err := UserDao.QueryLevelByUserId(q.UserId)
	if err != nil {
		return err
	}
	if UserIdLevel != 2 {
		return errors.New("token用户Id不是管理员")
	}
	return nil
}

func (q *LevelListFlow) prepareData() error {
	//查询level中的用户列表
	UserDao := model.NewUserDao()
	users := []*model2.User{{}}
	if err := UserDao.QueryUserInfoByLevel(q.Level, &users); err != nil {
		return err
	}

	for _, dbUser := range users {
		user := util.User{
			Id:              dbUser.ID,
			Name:            dbUser.UserName,
			NickName:        dbUser.NickName,
			WorkCount:       dbUser.WorkCount,
			BackGroundImage: dbUser.BackgroundImage,
			Level:           dbUser.Level,
			Avatar:          dbUser.Avatar,
		}
		q.data = append(q.data, &user)
	}
	return nil
}

func (q *LevelListFlow) packData() error {
	q.UserLevelResponse = &UserLevelResponse{LevelResponse: q.data}
	return nil
}
