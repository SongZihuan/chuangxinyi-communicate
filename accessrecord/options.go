package accessrecord

import (
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

func Options(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(viper.GetStringSlice("cors.allow_headers"), ", "))
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(viper.GetStringSlice("cors.allow_methods"), ", "))
	w.Header().Set("Access-Control-Allow-Credentials", "true") // 可将将 * 替换为指定的域名

	origin := r.Header.Get("Origin")
	if IsConfigOrigin(origin) { // 同源和配置包内的请求体已经被添加
		WriteOrigin(origin, w)
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	next(w, r)
}

func IsConfigOrigin(origin string) bool {
	if len(origin) == 0 { // 同源请求
		return true
	}

	for _, o := range viper.GetStringSlice("cors.allow_origins") {
		if origin == o {
			return true
		}
	}

	return false
}

func WriteOrigin(origin string, w http.ResponseWriter) {
	if len(origin) == 0 {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	} else {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
}
