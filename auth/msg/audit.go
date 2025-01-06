package msg

import (
	"fmt"
	"github.com/SongZihuan/chuangxinyi-communicate/auth"
	"github.com/SongZihuan/chuangxinyi-communicate/model"
	"github.com/spf13/viper"
	errors "github.com/wuntsong-org/wterrors"
)

func NewAudit(user *model.User, content string, args ...any) {
	go func() {
		_, _ = SendAudit(user, content, args)
	}()
}

func SendAudit(user *model.User, content string, args ...any) (bool, errors.WTError) {
	AuditCaller := viper.GetString("auth.auditCaller")

	var respData auth.SendMsgResp
	_, err := auth.SendRequests(auth.SendAuditReq{
		UserID:  user.Uid,
		Content: fmt.Sprintf(content, args...),
	}, AuditCaller, &respData)
	if err != nil {
		return false, errors.WarpQuick(err)
	}

	if !respData.Data.Success {
		return false, errors.Errorf(respData.Msg)
	}

	return true, nil
}
