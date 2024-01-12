package accessrecord

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"gitee.com/wuntsong/chuangxinyi-communicate/dao"
	"gitee.com/wuntsong/chuangxinyi-communicate/global/peername"
	"gitee.com/wuntsong/chuangxinyi-communicate/ip"
	"gitee.com/wuntsong/chuangxinyi-communicate/logger"
	"gitee.com/wuntsong/chuangxinyi-communicate/model"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	errors "github.com/wuntsong-org/wterrors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

func AccessRecordHandle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc, notRecordPath []string) {
	var err error
	ctx := r.Context()
	notRecord := utils.StringIn(r.URL.Path, notRecordPath)
	contentType := r.Header.Get("Content-Type")
	method := r.Method
	path := r.URL.Path
	queryValues := r.URL.Query()
	host := r.Host
	scheme := r.Header.Get("X-Forwarded-Proto")
	if len(scheme) == 0 {
		scheme = "http"
	}

	if websocket.IsWebSocketUpgrade(r) {
		scheme = fmt.Sprintf("%s/ws", scheme)
	}

	if len(contentType) > 80 {
		contentType = contentType[:80]
	}

	if len(method) > 15 {
		method = method[:15]
	}

	if len(scheme) > 15 {
		scheme = scheme[:15]
	}

	if len(path) > 100 {
		path = path[:100]
	}

	queryMap := make(map[string][]string, len(queryValues))
QUERY:
	for k, v := range queryValues {
		if len(v) == 0 {
			continue
		}

		if len(v) > 5 {
			continue
		}

		for _, vv := range v {
			if len(vv) > 100 {
				continue QUERY
			}
		}

		queryMap[k] = v
	}
	queryByte, err := utils.JsonMarshal(queryMap)
	if err != nil {
		queryByte = []byte("{}")
	}
	query := string(queryByte)
	if len(query) > 6000 {
		query = "{}"
	}

	if len(host) > 500 {
		host = host[:500]
	}

	allowHeader := viper.GetStringSlice("cors.allow_headers")
	headerMap := make(map[string]string, len(allowHeader))
	for _, h := range allowHeader {
		v := r.Header.Get(h)
		if len(v) != 0 {
			headerMap[h] = v
		}
	}
	headerByte, err := utils.JsonMarshal(headerMap)
	if err != nil {
		queryByte = []byte("{}")
	}
	header := string(headerByte)
	if len(header) > 6000 {
		header = "{}"
	}

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		accessRecordError(w, r, errors.WarpQuick(err), "获取请求体错误")
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyByte)) // 塞回去

	var body string
	err = r.ParseMultipartForm(10 << 20) // 限制最大10MB大小的表单数据
	if err == nil {
		type File struct {
			FileName string `json:"fileName"`
			FileSize int64  `json:"fileSize"`
		}

		data := struct {
			Value map[string][]string `json:"value"`
			File  map[string][]File
		}{Value: r.MultipartForm.Value, File: make(map[string][]File, len(r.MultipartForm.File))}

		for n, k := range r.MultipartForm.File {
			FileList := make([]File, 0, len(k))
			for _, h := range k {
				FileList = append(FileList, File{
					FileName: h.Filename,
					FileSize: h.Size,
				})
			}
			data.File[n] = FileList
		}

		bodyByte, err := utils.JsonMarshal(data)
		if err != nil {
			accessRecordError(w, r, err, "打包form请求体错误")
			return
		}

		body = string(bodyByte)
	} else {
		var data map[string]interface{}
		err := utils.JsonUnmarshal(bodyByte, &data)
		if err == nil {
			for k, v := range data {
				vString, ok := v.(string)
				if !ok {
					continue
				}

				if strings.HasPrefix(vString, "base64:") {
					data[k] = ""
				}
			}

			newBody, err := utils.JsonMarshal(data)
			if err != nil {
				newBody = bodyByte
			}

			body = string(newBody)
			if len(body) > 65000 {
				body = body[:65000]
			}
		} else {
			if utf8.Valid(bodyByte) {
				body = string(bodyByte)
				if len(body) > 65000 {
					body = body[:65000]
				}
			} else {
				body = fmt.Sprintf("<media type: %s bytes: %d>", utils.GetMediaType(bodyByte), len(bodyByte))
			}
		}
	}

	realIp := utils.GetTargetIP(r)
	var code, geo string
	code, geo, err = ip.GetGeo(ctx, realIp)
	if err != nil {
		code = ip.UnknownGeoCode
		geo = "未知"
	}

	var requestIDPrefix string
	err = func() (errRes errors.WTError) {
		defer utils.Recover(logger.Logger, &errRes, "create requests id prefix resp")
		unixNano := time.Now().UnixNano()
		text := fmt.Sprintf("%s\n%d\n%s\n%s\n", realIp, unixNano, r.URL.String(), peername.PeerName)
		requestIDPrefix = fmt.Sprintf("%s-%d", utils.HashSHA256WithBase62(text), unixNano)
		return nil
	}()
	if err != nil {
		accessRecordError(w, r, errors.WarpQuick(err), "创建请求前缀错误")
		return
	}

	access := &model.Record{
		RequestIDPrefix: requestIDPrefix,
		ServerName:      peername.PeerName,
		IP:              realIp,
		GeoCode:         code,
		Geo:             geo,
		Method:          method,
		Path:            path,
		Scheme:          scheme,
		Host:            host,
		Query:           query,
		ContentType:     &contentType,
		RequestsBody:    body,
		RequestsHeader:  header,
	}

	if !notRecord && geo != ip.LocalGeo {
		err := dao.RecordDao.Create(access)
		if err != nil {
			logger.Logger.Error("mysql resp: %s", err.Error())
			return
		}
	}

	w.Header().Set("X-Server", peername.PeerName)

	writer := MakeNewWriter(w)

	record := &Record{RequestsID: fmt.Sprintf("%s-%d", access.RequestIDPrefix, access.ID)}

	ctx = context.WithValue(ctx, "X-Record", record)
	ctx = context.WithValue(ctx, "X-Real-IP-Geo-Code", code)
	ctx = context.WithValue(ctx, "X-Real-IP-Geo-Code", code)
	ctx = context.WithValue(ctx, "X-Real-IP-Geo", geo)
	ctx = context.WithValue(ctx, "X-Real-IP", realIp)

	var startTime time.Time
	var endTime time.Time

	nextFuncErr := func() (errRes errors.WTError) {
		defer utils.Recover(logger.Logger, &errRes, "next func resp")
		defer func() {
			endTime = time.Now()
		}()
		startTime = time.Now()
		next(writer, r.WithContext(ctx))
		return nil
	}()

	go func() {
		if record.User != nil {
			access.UserID = record.User.ID
			access.UserUID = record.User.Uid
		}

		if len(record.UserToken) != 0 {
			access.UserToken = record.UserToken
		}

		if !strings.HasPrefix(path, "/api/v1/admin/accessrecord") && writer.Header().Get("X-Not-Record") != "True" {
			body := writer.Body
			if len(body) > 65000 {
				body = body[:65000]
			}
			access.ResponseBody = &body
		}

		headerByte, err := utils.JsonMarshal(writer.Header())
		if err != nil {
			headerByte = []byte("{}")
		}
		header := string(headerByte)
		if len(header) > 6000 {
			header = "{}"
		}
		access.ResponseHeader = &header

		writeError := writer.WriteError
		if len(writeError) > 2000 {
			writeError = writeError[:2000]
		}
		access.ResponseBodyError = &writeError

		access.StatusCode = writer.Status

		if nextFuncErr != nil {
			panicError := nextFuncErr.Error()
			if len(panicError) > 2000 {
				panicError = panicError[:2000]
			}
			access.PanicError = &panicError
		}

		var msgByte []byte
		msgByte, err = utils.JsonMarshal(struct {
			Msg string `json:"msg"`
		}{
			Msg: record.Msg,
		})
		if err != nil {
			msgByte = []byte("{}")
		}
		msg := string(msgByte)
		if len(msg) > 6000 {
			msg = msg[:6000]
		}
		access.Message = &msg

		use := endTime.Sub(startTime).Milliseconds() // 毫秒
		access.UseTime = use

		access.StartAt = &startTime

		access.EndAt = &endTime

		if !notRecord && geo != ip.LocalGeo {
			err := dao.RecordDao.Update(access)
			if err != nil {
				logger.Logger.Error("mysql resp: %s", err)
			}
		}
	}()

}

func accessRecordError(w http.ResponseWriter, r *http.Request, err errors.WTError, reason string) {
	_, _ = w.Write([]byte("error"))
	w.WriteHeader(http.StatusBadRequest)
}

type NewWriter struct {
	http.ResponseWriter
	Status     int64
	Body       string
	WriteError string
}

func MakeNewWriter(w http.ResponseWriter) *NewWriter {
	return &NewWriter{
		ResponseWriter: w,
		Status:         0,
		Body:           "",
		WriteError:     "",
	}
}

func (w *NewWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *NewWriter) WriteHeader(statusCode int) {
	w.Status = int64(statusCode)
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *NewWriter) Write(data []byte) (int, error) {
	res, err := w.ResponseWriter.Write(data)
	w.Body = string(data)
	if err != nil {
		w.WriteError = err.Error()
	}
	return res, err
}

func (w *NewWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		panic("response does not implement http.Hijacker")
	}

	return h.Hijack()
}

func (w *NewWriter) CloseNotify() <-chan bool {
	h, ok := w.ResponseWriter.(http.CloseNotifier)
	if !ok {
		panic("response does not implement http.CloseNotifier")
	}

	return h.CloseNotify()
}

func (w *NewWriter) Flush() {
	h, ok := w.ResponseWriter.(http.Flusher)
	if !ok {
		panic("response does not implement http.Flusher")
	}

	h.Flush()
}

type Record struct {
	RequestsID string
	User       *model.User
	UserToken  string
	Msg        string
}

func GetRecord(ctx context.Context) *Record {
	res, ok := ctx.Value("X-Record").(*Record)
	if !ok {
		logger.Logger.Error("bad X-Record")
		return &Record{RequestsID: "unknown"} // 返回一个无效的，防止报错
	}
	return res
}

func GetGinRecord(ctx *gin.Context) *Record {
	return GetRecord(ctx.Request.Context())
}

func GetRecordIfExists(ctx context.Context) *Record {
	res, ok := ctx.Value("X-Record").(*Record)
	if !ok {
		return &Record{RequestsID: "unknown"} // 返回一个无效的，防止报错
	}
	return res
}

func GetGinRecordIfExists(ctx *gin.Context) *Record {
	return GetRecordIfExists(ctx.Request.Context())
}
