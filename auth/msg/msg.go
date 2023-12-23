package msg

import (
	"context"
	"fmt"
	"gitee.com/wuntsong/chuangxinyi-communicate/auth"
	"gitee.com/wuntsong/chuangxinyi-communicate/dao"
	"gitee.com/wuntsong/chuangxinyi-communicate/model"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func SendMsgByUserUID(userID string, title string, content string) (bool, error) {
	MsgCaller := viper.GetString("auth.msgCaller")

	var respData auth.SendMsgResp
	_, err := auth.SendRequests(auth.SendMsgReq{
		UserID:  userID,
		Title:   title,
		Content: content,
	}, MsgCaller, &respData)
	if err != nil {
		return false, err
	}

	if !respData.Data.Success {
		return false, errors.Errorf(respData.Msg)
	}

	return true, nil
}

func SendMsgByUserID(userID int64, title string, content string) (bool, error) {
	user := dao.UserDao.Get(userID) // 根据整型主键查找
	if user == nil {
		return false, fmt.Errorf("user not found")
	}

	return SendMsgByUserUID(user.Uid, title, content)
}

func SendMsg(user *model.User, title string, content string) (bool, error) {
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
