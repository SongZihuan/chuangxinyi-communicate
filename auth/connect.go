package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
	"github.com/spf13/viper"
)

var PriKey *rsa.PrivateKey
var AuthPubKey *rsa.PublicKey

func InitAuth() error {
	signPriKey := viper.GetString("auth.signPriKey")
	getPubKey := viper.GetString("auth.getPubKey")

	priKeyString, err := base64.StdEncoding.DecodeString(signPriKey)
	if err != nil {
		return err
	}

	PriKey, err = utils.ReadRsaPrivateKey(priKeyString)
	if err != nil {
		return err
	}

	pubKeyString, err := base64.StdEncoding.DecodeString(getPubKey)
	if err != nil {
		return err
	}

	AuthPubKey, err = utils.ReadRsaPublicKey(pubKeyString)
	if err != nil {
		return err
	}

	return nil
}
