package model

import "time"

type Record struct {
	// 不使用Model
	ID                int64  `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id" form:"id"`
	RequestIDPrefix   string `gorm:"index:idx_record_requests_id_prefix;not null"`
	ServerName        string `gorm:"not null"`
	UserID            int64
	UserUID           string  // 修改为你需要的数据类型
	UserToken         string  `gorm:"type:text"`
	IP                string  `gorm:"not null"`
	GeoCode           string  `gorm:"not null"`
	Geo               string  `gorm:"not null"`
	Scheme            string  `gorm:"not null"`
	Method            string  `gorm:"not null"`
	Host              string  `gorm:"not null"`
	Path              string  `gorm:"not null"`
	Query             string  `gorm:"type:json;not null"`
	ContentType       *string `gorm:"not null"`
	RequestsBody      string  `gorm:"type:text;not null"`
	ResponseBody      *string `gorm:"type:text"`
	ResponseBodyError *string `gorm:"size:2000"`
	RequestsHeader    string  `gorm:"type:json;not null"`
	ResponseHeader    *string `gorm:"type:json"`
	StatusCode        int64
	PanicError        *string `gorm:"size:2000"`
	Message           *string `gorm:"type:json"`
	UseTime           int64
	CreateTime        time.Time `gorm:"index:idx_record_create_time;autoUpdateTime"`
	StartTime         *time.Time
	EndTime           *time.Time
}
