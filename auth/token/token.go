package token

import (
	"gitee.com/wuntsong/chuangxinyi-communicate/auth"
	"github.com/spf13/viper"
	errors "github.com/wuntsong-org/wterrors"
)

func CheckPhoneToken(token string, phone string) (bool, errors.WTError) {
	PhoneCheckCaller := viper.GetString("auth.phoneCheckCaller")

	var respData auth.CheckPhoneTokenResp
	_, err := auth.SendRequests(auth.CheckPhoneTokenReq{
		Phone: phone,
		Token: token,
	}, PhoneCheckCaller, &respData)
	if err != nil {
		return false, errors.WarpQuick(err)
	}

	if !respData.Data.IsOK {
		return false, errors.Errorf(respData.Msg)
	}

	return true, nil
}

func CheckEmailToken(token string, email string) (bool, errors.WTError) {
	EmailCheckCaller := viper.GetString("auth.emailCheckCaller")

	var respData auth.CheckEmailTokenResp
	_, err := auth.SendRequests(auth.CheckEmailTokenReq{
		Email: email,
		Token: token,
	}, EmailCheckCaller, &respData)
	if err != nil {
		return false, errors.WarpQuick(err)
	}

	if !respData.Data.IsOK {
		return false, errors.Errorf(respData.Msg)
	}

	return true, nil
}

func Check2FAToken(token string, userID string) (bool, errors.WTError) {
	SecondFACheckCaller := viper.GetString("auth.secondFACheckCaller")

	var respData auth.CheckSecondFATokenResp
	_, err := auth.SendRequests(auth.CheckSecondFATokenReq{
		UserID: userID,
		Token:  token,
	}, SecondFACheckCaller, &respData)
	if err != nil {
		return false, errors.WarpQuick(err)
	}

	if !respData.Data.IsOK {
		return false, errors.Errorf(respData.Msg)
	}

	return true, nil
}
