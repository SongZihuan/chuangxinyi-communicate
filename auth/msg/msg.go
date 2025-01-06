package msg

import (
	"context"
	"fmt"
	"github.com/SongZihuan/chuangxinyi-communicate/auth"
	"github.com/SongZihuan/chuangxinyi-communicate/cache"
	"github.com/SongZihuan/chuangxinyi-communicate/model"
	"github.com/spf13/viper"
	errors "github.com/wuntsong-org/wterrors"
)

func SendMsgByUserUID(userID string, title string, content string) (bool, errors.WTError) {
	MsgCaller := viper.GetString("auth.msgCaller")

	var respData auth.SendMsgResp
	_, err := auth.SendRequests(auth.SendMsgReq{
		UserID:  userID,
		Title:   title,
		Content: content,
	}, MsgCaller, &respData)
	if err != nil {
		return false, errors.WarpQuick(err)
	}

	if !respData.Data.Success {
		return false, errors.Errorf(respData.Msg)
	}

	return true, nil
}

func SendMsgByUserID(userID int64, title string, content string) (bool, errors.WTError) {
	user := cache.UserCache.Get(userID)
	if user == nil {
		return false, errors.Errorf("user not found")
	}

	return SendMsgByUserUID(user.Uid, title, content)
}

func SendMsg(user *model.User, title string, content string) (bool, errors.WTError) {
	return SendMsgByUserUID(user.Uid, title, content)
}

func MessageSendLogin(user *model.User, ctx context.Context) {
	geo, ok := ctx.Value("X-Real-IP-Geo").(string)
	if !ok {
		geo = "未知"
	}

	ip, ok := ctx.Value("X-Real-IP").(string)
	if !ok {
		ip = "未知"
	}

	ReadableName := viper.GetString("readableName")

	go func() {
		_, _ = SendMsg(user, fmt.Sprintf("%s用户登录提醒", ReadableName), fmt.Sprintf("用户成功登录%s，登录IP为：%s，登录地为：%s", ReadableName, ip, geo))
	}()
}
