package auth

type UserEasy struct {
	UID            string `json:"id"`
	Phone          string `json:"phone"`
	RoleID         int64  `json:"roleID,omitempty"`
	RoleName       string `json:"roleName"`
	UserName       string `json:"userName,omitempty"`
	NickName       string `json:"nickname,omitempty"`
	Header         string `json:"header,omitempty"`
	Email          string `json:"email,omitempty"`
	UserRealName   string `json:"userRealName,omitempty"`
	CompanyName    string `json:"companyName,omitempty"`
	WeChatNickName string `json:"wechatNickName,omitempty"`
	WeChatHeader   string `json:"wechatHeader,omitempty"`
	UnionID        string `json:"unionID,omitempty"`
	Signin         bool   `json:"signin"`
	Status         string `json:"status"`
}

type UserInfoEsay struct {
	HasVerified bool   `json:"hasVerified"`
	UserName    string `json:"userName"`

	IsCompany       bool   `json:"isCompany"`
	LegalPersonName string `json:"legalPersonName,omitempty"`
	CompanyName     string `json:"companyName,omitempty"`
}

type UserData struct {
	HasPassword             bool   `json:"hasPassword"`
	HasEmail                bool   `json:"hasEmail"`
	Has2FA                  bool   `json:"has2FA"`
	HasWeChat               bool   `json:"hasWeChat"`
	HasUnionID              bool   `json:"hasUnionId"`
	HasVerified             bool   `json:"hasVerified"`
	IsCompany               bool   `json:"isCompany"`
	HasUserOriginal         bool   `json:"hasUserOriginal"`
	HasUserFaceCheck        bool   `json:"hasUserFaceCheck"`
	HasCompanyOriginal      bool   `json:"hasCompanyOriginal"`
	HasLegalPersonFaceCheck bool   `json:"hasLegalPersonFaceCheck"`
	VerifiedPhone           string `json:"verifiedPhone,omitempty"`
}

type SendWorkOrderFile struct {
	FileName string `json:"fileName"`
	File     string `json:"file"`
}

type SendWorkOrder struct {
	UserID  string              `json:"userID"`
	Title   string              `json:"title"`
	Content string              `json:"content"`
	File    []SendWorkOrderFile `json:"file"`
}

type SendMsgData struct {
	Have    bool `json:"have"`
	Success bool `json:"success"`
}

type SendAuditReq struct {
	UserID  string `json:"userID"`
	Content string `json:"content"`
}

type SendMsgResp struct {
	Resp
	Data SendMsgData `json:"data"`
}

type CheckPhoneTokenReq struct {
	Token string `json:"token"`
	Phone string `json:"phone"`
}

type CheckPhoneTokenData struct {
	IsOK bool `json:"isOK"`
}

type CheckPhoneTokenResp struct {
	Resp
	Data CheckPhoneTokenData `json:"data"`
}

type CheckEmailTokenReq struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

type CheckEmailTokenData struct {
	IsOK bool `json:"isOK"`
}

type CheckEmailTokenResp struct {
	Resp
	Data CheckEmailTokenData `json:"data"`
}

type CheckSecondFATokenReq struct {
	Token  string `json:"token"`
	UserID string `json:"userID"`
}

type CheckSecondFATokenData struct {
	IsOK bool `json:"isOK"`
}

type CheckSecondFATokenResp struct {
	Resp
	Data CheckSecondFATokenData `json:"data"`
}

type SendMsgReq struct {
	UserID  string `json:"userID"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
