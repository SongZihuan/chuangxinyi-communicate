package login

import (
	"context"
	"fmt"
	"gitee.com/wuntsong/chuangxinyi-communicate/auth"
	"gitee.com/wuntsong/chuangxinyi-communicate/logger"
	"gitee.com/wuntsong/chuangxinyi-communicate/rand"
	"gitee.com/wuntsong/chuangxinyi-communicate/redis"
	"gitee.com/wuntsong/chuangxinyi-communicate/signalexit"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	errors "github.com/wuntsong-org/wterrors"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"
)

var UserInfoChan chan WSWebsiteMessage

const (
	WebsiteClose    = "CLOSE"
	WebsitePong     = "PONG"
	WebsiteBadCode  = "BAD_CODE"
	WebsiteUserInfo = "USER_INFO"
	WebsiteBadToken = "BAD_TOKEN"
)

const (
	WebsitePing        = "PING"
	WebsiteGetUserInfo = "GET_USER_INFO"
)

type WSWebsiteMessage struct {
	Code string `json:"code"`
	Data string `json:"data"`
}

func WriteMessage(c chan WSWebsiteMessage, msg WSWebsiteMessage) {
	go func() {
		defer func() {
			recover()
			// 不需要utils.Recover
		}()
		c <- msg
	}()
}

func WriteJson(data any) string {
	s, err := utils.JsonMarshal(data)
	if err != nil {
		return errors.WarpQuick(err).Error()
	}
	return string(s)
}

func RestartGetUserInfo() errors.WTError {
	res := redis.RedisClient.Keys(context.Background(), "logintoken:getinfo:*:*")
	lst, err := res.Result()
	if err != nil {
		return errors.WarpQuick(err)
	}

	for _, l := range lst {
		tmp := strings.Split(l, ":")
		if len(tmp) != 4 {
			continue
		}

		err = StartGetUserInfo(tmp[2], tmp[3])
		if err != nil {
			logger.Logger.Error("restart get user info resp: %s", err.Error())
		}
	}

	return nil
}

func StartGetUserInfo(token string, uid string) errors.WTError {
	key := fmt.Sprintf("logintoken:getinfo:%s:%s", token, uid)
	_ = redis.RedisClient.Del(context.Background(), key) // 删除这个key，避免重复

	loginTokenExpire := viper.GetInt64("auth.loginTokenExpire")

	go func() {
		defer signalexit.AddExitFuncAsDeferNotDo(func() {
			_ = redis.RedisClient.Set(context.Background(), key, "1", time.Second*time.Duration(loginTokenExpire))
		})()

		for {
			time.Sleep(time.Duration(15+rand.GlobalRander.Intn(30)) * time.Minute) // 加扰动因素
			func() {
				key := fmt.Sprintf("ws:info:%s", uid)
				if !redis.AcquireLock(context.Background(), key, time.Second*30) {
					return
				}
				// 不需要释放锁，等待锁过期即可

				WriteMessage(UserInfoChan, WSWebsiteMessage{
					Code: WebsiteGetUserInfo,
					Data: WriteJson(struct {
						Token string `json:"token"`
					}{Token: token}),
				})
			}()
		}
	}()

	return nil
}

func ConnectWebSocket() (chan WSWebsiteMessage, errors.WTError) {
	UserInfoChan = make(chan WSWebsiteMessage, 10)

	websocketURL := viper.GetString("auth.websocket")

	go func() {
		defer func() {
			close(UserInfoChan)
		}()

		for {
			func() {
				v := url.Values{}
				v.Add("xrunmode", "release")

				u, jsonErr := auth.SendGetRequests(websocketURL, v)
				if jsonErr != nil {
					logger.Logger.Error("websocket resp: %s", jsonErr)
					return
				}

				dialer := websocket.Dialer{}
				conn, _, err := dialer.Dial(u, nil)
				if err != nil {
					return
				}
				defer utils.Close(conn)

				utils.OnceFunc(func() {
					err = RestartGetUserInfo() // 只需要执行一次
					if err != nil {
						logger.Logger.Error("restart get user info resp: %s", err.Error())
					}
				})()

				var deadlineMutex = new(sync.Mutex)
				deadline := time.Now().Add(time.Minute * 5)
				stop := false

				deadlineNow := func() {
					deadlineMutex.Lock()
					defer deadlineMutex.Unlock()
					deadline = time.Now()
					stop = true
				}

				go func() {
					for {
						WriteMessage(UserInfoChan, WSWebsiteMessage{
							Code: WebsitePing,
						})

						time.Sleep(1 * time.Minute)

						if func() bool {
							defer utils.Recover(logger.Logger, nil, "ws check deadline")

							deadlineMutex.Lock()
							defer deadlineMutex.Unlock()

							if stop {
								return true
							}

							if time.Now().After(deadline) {
								func() {
									logger.Logger.Error("ws not pong message: %d, %d", deadline.Unix(), time.Now().Unix())
									WriteMessage(UserInfoChan, WSWebsiteMessage{
										Code: WebsiteClose,
									})
									_ = conn.Close()
								}()
								return true
							}

							return false
						}() {
							return // 退出go
						}
					}
				}()

				go func() {
					defer signalexit.AddExitFuncAsDefer(func() {
						defer utils.Recover(logger.Logger, nil, "ws reader")
						time.Sleep(1 * time.Second)
						WriteMessage(UserInfoChan, WSWebsiteMessage{
							Code: WebsiteClose,
						})
						_ = conn.Close()
						deadlineNow()
					})()

					for {
						var msg WSWebsiteMessage
						err = conn.ReadJSON(&msg)
						if utils.IsNetClose(err) || errors.Is(err, io.ErrUnexpectedEOF) || websocket.IsUnexpectedCloseError(err) {
							return
						} else if err != nil {
							logger.Logger.Error("websocket read resp: %s", err.Error())
							return
						}

						switch msg.Code {
						case WebsitePong:
							func() {
								defer func() {
									recover()
								}()

								deadlineMutex.Lock()
								defer deadlineMutex.Unlock()

								deadline = time.Now().Add(time.Minute * 5)
							}()
						case WebsiteUserInfo, WebsiteBadToken:
							func() {
								data := struct {
									Success bool              `json:"success"`
									Error   bool              `json:"resp"`
									Token   string            `json:"token"`
									Msg     string            `json:"msg"`
									User    auth.UserEasy     `json:"user"`
									Data    auth.UserData     `json:"data"`
									Info    auth.UserInfoEsay `json:"info"`
								}{}

								err = utils.JsonUnmarshal([]byte(msg.Data), &data)
								if err != nil {
									return
								}

								if data.Error {
									logger.Logger.Error("get user info fail: %s", data.Msg)
									return
								}

								if !data.Success {
									return
								}

								_, err = UpdateUserInfo(context.Background(), data.User, data.Info, data.Data)
								if err != nil {
									logger.Logger.Error("websocket resp: %s", err.Error())
								}
							}()
						}
					}
				}()

				func() { // 不使用go协程
					defer signalexit.AddExitFuncAsDefer(func() {
						defer utils.Recover(logger.Logger, nil, "ws writer")
						_ = conn.Close()
						deadlineNow()
					})()

					defer utils.Recover(logger.Logger, nil, "ws writer")

					for {
						msg := <-UserInfoChan

						switch msg.Code {
						case WebsiteClose:
							return
						}

						err = conn.WriteJSON(msg)
						if utils.IsNetClose(err) || websocket.IsUnexpectedCloseError(err) {
							return
						} else if err != nil {
							WriteMessage(UserInfoChan, msg) // 重新把chan放回去
							logger.Logger.Error("websocket write resp: %s", err.Error())
							return
						}
					}
				}()
			}()
			time.Sleep(3 * time.Second)
		}
	}()

	return UserInfoChan, nil
}
