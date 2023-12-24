package ip

import (
	"github.com/spf13/viper"
	errors "github.com/wuntsong-org/wterrors"
)

func InitYunIP() errors.WTError {
	if len(viper.GetString("aliyun.appCode")) == 0 {
		return errors.Errorf("aliyun ip app code must be given")
	}

	if len(viper.GetString("aliyun.appKey")) == 0 {
		return errors.Errorf("aliyun ip app secret must be given")
	}

	if len(viper.GetString("aliyun.appSecret")) == 0 {
		return errors.Errorf("aliyun ip app key must be given")
	}
	return nil
}
