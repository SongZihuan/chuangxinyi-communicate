package auth

type RespInterface interface {
	GetCode() string
	GetSubCode() string
	GetMsg() string
}

type Resp struct {
	Code       string `json:"code"`
	SubCode    string `json:"subCode"`
	NumCode    int64  `json:"_code"`
	NumSubCode int64  `json:"_subCode"`
	Msg        string `json:"msg,omitempty"`
}

type AuthResp struct {
	Code string `json:"code"`
	Msg  string `json:"msg,omitempty"`
}

type CheckLoginTokenData struct {
	IsLogin bool         `json:"isLogin"`
	User    UserEasy     `json:"user"`
	Data    UserData     `json:"data"`
	Info    UserInfoEsay `json:"info"`
}

type CheckLoginTokenResp struct {
	Resp
	Data CheckLoginTokenData `json:"data"`
}

type CreateDefrayData struct {
	Token   string `json:"token"`
	TradeID string `json:"tradeID"`
}

type CreateDefrayResp struct {
	Resp
	Data CreateDefrayData `json:"data"`
}

type QueryDefrayData struct {
	Status       string `json:"status"`
	PayerID      string `json:"payerID"`
	PayAt        int64  `json:"payAt"`
	ReturnAt     int64  `json:"returnAt"`
	ReturnReason string `json:"returnReason"`
}

type QueryDefrayResp struct {
	Resp
	Data QueryDefrayData `json:"data"`
}

type ReturnDefrayReq struct {
	TradeID string `json:"tradeID"`
	Reason  string `json:"reason"`
	Must    bool   `json:"must"`
}

func (r Resp) GetCode() string {
	return r.Code
}

func (r Resp) GetSubCode() string {
	return r.SubCode
}

func (r Resp) GetMsg() string {
	return r.Msg
}
