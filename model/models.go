package model

var Models = []interface{}{
	&User{}, &Tag{}, &Article{}, &ArticleTag{}, &Comment{}, &Favorite{},
	&Topic{}, &Node{}, &TopicTag{}, &TopicLike{}, &Setting{}, &Link{},
	&UserWatch{}, &UserScore{}, &UserScoreLog{}, Record{},
}

type Model struct {
	ID int64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id" form:"id"`
}

const (
	StatusOk      = 0 // 正常
	StatusDeleted = 1 // 删除
	StatusPending = 2 // 待审核

	UserLevelGeneral = 0  // 普通用户
	UserLevelAdmin   = 10 // 管理员

	ContentTypeHtml = "html"

	EntityTypeArticle = "article"
	EntityTypeTopic   = "topic"
	EntityTypeComment = "comment"
	EntityTypeUser    = "user"

	MsgTypeComment   = 0 // 回复消息
	MsgTypeTopicLike = 1 // 话题点赞
	MsgTypeUserWatch = 2 // 用户关注

	ScoreTypeIncr = 0 // 积分+
	ScoreTypeDecr = 1 // 积分-

	TopicTypeNormal  = 0 // 普通帖子
	TopicTypeTwitter = 1 // 推文
)
