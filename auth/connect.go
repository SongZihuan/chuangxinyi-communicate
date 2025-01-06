package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"github.com/SongZihuan/chuangxinyi-communicate/utils"
	"github.com/spf13/viper"
	errors "github.com/wuntsong-org/wterrors"
)

var PriKey *rsa.PrivateKey
var AuthPubKey *rsa.PublicKey

func InitAuth() errors.WTError {
	signPriKey := viper.GetString("auth.signPriKey")
	getPubKey := viper.GetString("auth.getPubKey")

	priKeyString, err := base64.StdEncoding.DecodeString(signPriKey)
	if err != nil {
		return errors.WarpQuick(err)
	}

	PriKey, err = utils.ReadRsaPrivateKey(priKeyString)
	if err != nil {
		return errors.WarpQuick(err)
	}

	pubKeyString, err := base64.StdEncoding.DecodeString(getPubKey)
	if err != nil {
		return errors.WarpQuick(err)
	}

	AuthPubKey, err = utils.ReadRsaPublicKey(pubKeyString)
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}
