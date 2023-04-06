package Redis

import (
	"errors"
	"fmt"
	"strconv"
)

// 点赞，更新喜欢列表，喜欢列表是所有点赞的朋友圈
func (u *RedisDao) UpdatePostLike(userId uint, friendschatId uint, state bool) error {
	strKey := fmt.Sprintf("%s-%d", LIKE, userId)
	if state {
		if err := Rdb.SAdd(ctx, strKey, friendschatId).Err(); err != nil {
			fmt.Println("UpdatePostLike SAdd Failed:", err)
			return err
		}
	} else {
		if err := Rdb.SRem(ctx, strKey, friendschatId).Err(); err != nil {
			fmt.Println("UpdatePostLike SRem Failed:", err)
			return err
		}
	}
	return nil
}

// 查询用户对该朋友圈是否点赞
func (u *RedisDao) GetLikeState(userId uint, friendschatId uint) (bool, error) {
	strKey := fmt.Sprintf("%s-%d", LIKE, userId)
	ok, err := Rdb.SIsMember(ctx, strKey, friendschatId).Result()
	if err != nil {
		fmt.Println("RedisContains SIsMember Failed:", err)
		return false, err
	}
	return ok, nil
}

// 根据friendschatId设置该朋友圈被点赞的次数
func (u *RedisDao) SetLikeNumByfriendschatId(friendschatId uint, num int) error {
	strKey := fmt.Sprintf("%s-%d", LIKENUM, friendschatId)
	if err := Rdb.Set(ctx, strKey, num, 0).Err(); err != nil {
		return err
	}
	return nil
}

// 根据friendschatId查询该朋友圈被点赞的次数
func (u *RedisDao) GetLikeNumByfriendschatId(friendschatId uint) (int, error) {
	strKey := fmt.Sprintf("%s-%d", LIKENUM, friendschatId)
	strRes, err := Rdb.Get(ctx, strKey).Result()
	if err != nil {
		return 0, err
	}
	res, _ := strconv.ParseInt(strRes, 10, 64)
	return int(res), nil
}

// 给friendschatId增加一个赞
func (u *RedisDao) AddOneLikeNumByfriendschatId(friendschatId uint) error {
	strKey := fmt.Sprintf("%s-%d", LIKENUM, friendschatId)
	if err := Rdb.Incr(ctx, strKey).Err(); err != nil {
		return err
	}
	return nil
}

// 给friendschatId减少一个赞
func (u *RedisDao) SubOneLikeNumByfriendschatId(friendschatId uint) error {
	strKey := fmt.Sprintf("%s-%d", LIKENUM, friendschatId)
	res, err := Rdb.Decr(ctx, strKey).Result()
	if err != nil {
		return err
	}
	if res <= -1 {
		return errors.New("The Key Not Found, You Can't Decrease it")
	}
	return nil
}
