package main

import (
	"gitee.com/wuntsong/chuangxinyi-communicate/accessrecord"
	"net/http"
)

func NewServer(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessrecord.AccessRecordHandle(w, r, func(w http.ResponseWriter, r *http.Request) {
			accessrecord.Options(w, r, handler.ServeHTTP)
		}, []string{"/"}) // / 相当于 /api/v1/ping
	})
}
