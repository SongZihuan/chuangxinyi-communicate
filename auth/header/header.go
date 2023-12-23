package header

import (
	"fmt"
	"github.com/spf13/viper"
	"net/url"
)

func GetHeaderUrl(userID string) string {
	Header := viper.GetString("auth.header")

	if len(userID) == 0 {
		return Header
	}

	v := url.Values{}
	v.Add("id", userID)
	return fmt.Sprintf("%s?%s", Header, v.Encode())
}
