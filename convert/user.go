package convert

import (
	"strconv"

	"gitee.com/wuntsong/chuangxinyi-communicate/cache"
	"gitee.com/wuntsong/chuangxinyi-communicate/model"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
)

func ToUserDefaultIfNull(id int64) *model.UserInfo {
	user := cache.UserCache.Get(id)
	if user == nil {
		user = &model.User{}
		user.ID = id
		user.Phone = "18888888888"
		user.Email = utils.SqlNullString("hello@world.com")
		user.Username = utils.SqlNullString(strconv.FormatInt(id, 10))
		user.CreateTime = utils.NowTimestamp()
	}
	return ToUser(user)
}

func ToUserById(id int64) *model.UserInfo {
	user := cache.UserCache.Get(id)
	return ToUser(user)
}

func ToUser(user *model.User) *model.UserInfo {
	if user == nil {
		return nil
	}
	levelName := "普通用户"
	if user.Level == model.UserLevelAdmin {
		levelName = "管理员"
	}
	ret := &model.UserInfo{
		Id:           user.ID,
		Uid:          user.Uid,
		Phone:        user.Phone,
		Username:     user.Username.String,
		Nickname:     user.Nickname.String,
		Header:       user.Header.String,
		Email:        user.Email.String,
		Level:        user.Level,
		LevelName:    levelName,
		TopicCount:   user.TopicCount,
		CommentCount: user.CommentCount,
		Status:       user.Status,
		CreateTime:   user.CreateTime,
	}
	if user.Status == model.StatusDeleted {
		ret.Username = "blacklist"
		ret.Nickname = "黑名单用户"
		ret.Header = ""
		ret.Email = ""
	} else {
		ret.Score = cache.UserCache.GetScore(user.ID)
	}
	return ret
}

func ToUsers(users []model.User) []model.UserInfo {
	if len(users) == 0 {
		return nil
	}
	var responses []model.UserInfo
	for _, user := range users {
		item := ToUser(&user)
		if item != nil {
			responses = append(responses, *item)
		}
	}
	return responses
}
