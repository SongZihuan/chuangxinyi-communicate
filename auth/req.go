package auth

type CheckLoginTokenReq struct {
	Token string `json:"token"`
}

type CreateDefrayReq struct {
	OwnerID            string `json:"ownerID,optional"`
	Subject            string `json:"subject"`   // 标题
	Price              int64  `json:"price"`     // 价格
	Quantity           int64  `json:"quantity"`  // 数量
	UnitPrice          int64  `json:"unitPrice"` // 单价
	Describe           string `json:"describe"`  // 描述
	InvitePre          int64  `json:"invitePre,optional"`
	DistributionLevel1 int64  `json:"distributionLevel1,optional"`
	DistributionLevel2 int64  `json:"distributionLevel2,optional"`
	DistributionLevel3 int64  `json:"distributionLevel3,optional"`
	CanWithdraw        bool   `json:"canWithdraw"`
	ReturnURL          string `json:"returnURL"`
	MustSelfDefray     bool   `json:"mustSelfDefray"`
}

type QueryDefrayReq struct {
	TradeID string `json:"tradeID"`
}
