package login

import (
	"context"
	"database/sql"
	"gitee.com/wuntsong/chuangxinyi-communicate/auth"
	"gitee.com/wuntsong/chuangxinyi-communicate/dao"
	"gitee.com/wuntsong/chuangxinyi-communicate/model"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
	"github.com/spf13/viper"
	errors "github.com/wuntsong-org/wterrors"
)

var BannedStatus = errors.NewClass("banned_status")
var RoleNotFound = errors.NewClass("role_not_found")

func CheckLogin(ctx context.Context, token string) (*model.User, errors.WTError) {
	LoginCaller := viper.GetString("auth.loginCaller")

	var respData auth.CheckLoginTokenResp
	_, err := auth.SendRequests(auth.CheckLoginTokenReq{
		Token: token,
	}, LoginCaller, &respData)
	if err != nil {
		return nil, errors.WarpQuick(err)
	} else if !respData.Data.IsLogin {
		return nil, errors.Errorf(respData.Msg)
	}

	user, err := UpdateUserInfo(ctx, respData.Data.User, respData.Data.Info, respData.Data.Data)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	err = StartGetUserInfo(token, user.Uid)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	return user, nil
}

func UpdateUserInfo(ctx context.Context, easy auth.UserEasy, info auth.UserInfoEsay, data auth.UserData) (*model.User, errors.WTError) {
	userNotFound := false

	user := dao.UserDao.GetByUid(easy.UID)
	if user == nil {
		userNotFound = true
		user = new(model.User)
	} else {
		userNotFound = false
	}

	user.Uid = easy.UID
	user.Phone = easy.Phone
	user.Email = sql.NullString{
		Valid:  len(easy.Email) != 0,
		String: easy.Email,
	}
	user.Username = sql.NullString{
		Valid:  len(easy.UserName) != 0,
		String: easy.UserName,
	}
	user.Nickname = sql.NullString{
		Valid:  len(easy.NickName) != 0,
		String: easy.NickName,
	}
	user.Header = sql.NullString{
		Valid:  len(easy.Header) != 0,
		String: easy.Header,
	}

	adminPhone := viper.GetString("admin.phone")

	if userNotFound {
		user.Status = model.StatusOk
		user.CreateTime = utils.NowTimestamp()
		user.UpdateTime = user.CreateTime

		if user.Phone == adminPhone {
			user.Level = model.UserLevelAdmin
		} else {
			user.Level = model.UserLevelGeneral
		}

		err := dao.UserDao.Create(user)
		if err != nil {
			return nil, errors.WarpQuick(err)
		}
	} else {
		if user.Phone == adminPhone {
			user.Level = model.UserLevelAdmin
		}

		user.UpdateTime = utils.NowTimestamp()

		err := dao.UserDao.Update(user)
		if err != nil {
			return nil, errors.WarpQuick(err)
		}
	}

	return user, nil
}
