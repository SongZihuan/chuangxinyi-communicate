package auth

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func SendGetRequests(u string, values url.Values) (string, error) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	n := utils.GenerateUniqueNumber(18)

	domainUID := viper.GetString("auth.domainUID")

	values.Add("xdomain", domainUID)
	values.Add("xtimestamp", timestamp)
	values.Add("xn", n)

	up, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	signText := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n", http.MethodGet, domainUID, timestamp, n, up.Path, values.Encode())
	sign, err := utils.SignRsaHash256Sign(signText, PriKey)
	if err != nil {
		return "", err
	}

	values.Add("xsign", base64.StdEncoding.EncodeToString(sign))
	return fmt.Sprintf("%s?%s", u, values.Encode()), nil
}

func SendRequests(data any, u string, r RespInterface) (*http.Response, error) {
	domainUID := viper.GetString("auth.domainUID")
	xRunMode := viper.GetString("auth.xRunMode")

	dataByte, jsonErr := utils.JsonMarshal(data)
	if jsonErr != nil {
		return nil, jsonErr
	}

	uu, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(dataByte))
	if err != nil {
		return nil, err
	}

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	n := utils.GenerateUniqueNumber(18)

	signText := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n", http.MethodPost, domainUID, timestamp, n, uu.Path, string(dataByte))
	sign, err := utils.SignRsaHash256Sign(signText, PriKey)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Domain", domainUID)
	req.Header.Set("X-Timestamp", timestamp)
	req.Header.Set("X-N", n)
	req.Header.Set("X-Sign", base64.StdEncoding.EncodeToString(sign))
	req.Header.Set("X-RunMode", xRunMode)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("get bad status code: %s", body)
	}

	err = utils.JsonUnmarshal(body, r)
	if err != nil {
		return nil, err
	}

	// TODO 设置Code
	//code := r.GetCode()
	//if code == "SUCCESS" {
	//	return resp, nil
	//} else if code == "WEBSITE_DENY" {
	//	return nil, fmt.Errorf("get website deny: %s", r.GetMsg()).SetCode(fmt.Sprintf(r.GetSubCode()))
	//} else if code == "LOGIC_DENY" {
	//	return nil, fmt.Errorf("get logic resp: %s", r.GetMsg()).SetCode(fmt.Sprintf(r.GetSubCode()))
	//}
	//
	//return nil, fmt.Errorf("requests fail: %s", r.GetMsg()).SetCode(fmt.Sprintf(r.GetCode()))

	code := r.GetCode()
	if code == "SUCCESS" {
		return resp, nil
	} else if code == "WEBSITE_DENY" {
		return nil, fmt.Errorf("get website deny: %s", r.GetMsg())
	} else if code == "LOGIC_DENY" {
		return nil, fmt.Errorf("get logic resp: %s", r.GetMsg())
	}

	return nil, fmt.Errorf("requests fail: %s", r.GetMsg())
}

func WriteRespFail(w http.ResponseWriter) {
	WriteResp(w, "FAIL")
}

func WriteRespSuccess(w http.ResponseWriter) {
	WriteResp(w, "SUCCESS")
}

func WriteResp(w http.ResponseWriter, status string) {
	data, err := utils.JsonMarshal(struct {
		AuthResp
		Data struct {
			Status string `json:"status"`
		}
	}{
		AuthResp: AuthResp{
			Code: "SUCCESS",
			Msg:  "成功",
		},
		Data: struct {
			Status string `json:"status"`
		}{
			Status: status,
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}

func Verify(w http.ResponseWriter, r *http.Request) (res bool) {
	defer func() {
		if !res {
			WriteResp(w, "FAIL")
		}
	}()

	selfDomainUID := viper.GetString("auth.domainUID")

	domainUID := r.Header.Get("X-Domain")
	if len(domainUID) == 0 {
		domainUID = r.URL.Query().Get("xdomain")
	}

	n := r.Header.Get("X-N")
	if len(n) == 0 {
		n = r.URL.Query().Get("xn")
	}

	timestamp := r.Header.Get("X-Timestamp")
	if len(timestamp) == 0 {
		timestamp = r.URL.Query().Get("xtimestamp")
	}

	signBase64 := strings.TrimSpace(r.Header.Get("X-Sign"))
	if len(signBase64) == 0 {
		signBase64 = r.URL.Query().Get("xsign")
	}

	if domainUID != selfDomainUID {
		return false
	}

	if len(timestamp) == 0 || len(n) != 18 || len(signBase64) == 0 || len(domainUID) == 0 {
		return false
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}

	t := time.Unix(ts, 0)
	if time.Now().After(t.Add(time.Minute * 5)) {
		return false
	}

	var signText string
	if r.Method == http.MethodGet || r.Method == http.MethodHead || websocket.IsWebSocketUpgrade(r) {
		q := r.URL.Query()
		q.Del("xsign")
		signText = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n", r.Method, domainUID, timestamp, n, r.URL.Path, q.Encode())
	} else {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return false
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body)) // 塞回去
		signText = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n", r.Method, domainUID, timestamp, n, string(body))
	}

	sign, err := base64.StdEncoding.DecodeString(signBase64)
	if err != nil {
		return false
	}

	if utils.VerifyRsaHash256Sign(signText, sign, AuthPubKey) != nil {
		return false
	}

	return true
}
