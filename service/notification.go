package service

import (
	"fmt"
	"github.com/SongZihuan/chuangxinyi-communicate/auth/msg"
	errors "github.com/wuntsong-org/wterrors"
	"sync"

	"github.com/SongZihuan/chuangxinyi-communicate/cache"
	"github.com/SongZihuan/chuangxinyi-communicate/dao"
	"github.com/SongZihuan/chuangxinyi-communicate/logger"
	"github.com/SongZihuan/chuangxinyi-communicate/model"
	"github.com/SongZihuan/chuangxinyi-communicate/utils"
)

type Notification struct {
	FromID       int64
	ToID         int64
	Content      string
	QuoteContent string
	Type         int
}

var NotificationService = newNotificationService()

func newNotificationService() *notificationService {
	return &notificationService{
		notificationsChan: make(chan *Notification, 10),
	}
}

type notificationService struct {
	notificationsConsumeOnce sync.Once
	notificationsChan        chan *Notification
}

func (s *notificationService) Create(t *Notification) errors.WTError {
	title := ""

	switch t.Type {
	case model.MsgTypeUserWatch:
		title = "新增关注"
	case model.MsgTypeComment:
		title = "新增评论"
	case model.MsgTypeTopicLike:
		title = "新增喜爱"
	default:
		return errors.Errorf("bad type")
	}

	from := cache.UserCache.Get(t.FromID)
	fromName := ""
	if from == nil {
		fromName = "陌生人"
	} else {
		fromName = utils.GetUserName(from.Uid, from.Username, from.Nickname)
	}

	content := t.Content
	if len(t.QuoteContent) > 0 {
		content += fmt.Sprintf("\n引用消息：%s", t.QuoteContent)
	}
	content += fmt.Sprintf("\n回复人：%s", fromName)

	_, _ = msg.SendMsgByUserID(t.ToID, title, content)

	return nil
}

// 用户关注
func (s *notificationService) SendUserWatchNotification(userWatch *model.UserWatch) {
	var (
		fromId       = userWatch.WatcherID // 消息发送人
		authorId     int64                 // 被关注人
		content      string                // 消息内容
		quoteContent string                // 引用内容
	)

	authorId = userWatch.UserID
	content = "我关注了你"
	quoteContent = ""

	if authorId <= 0 {
		return
	}
	// 给被关注者发消息
	s.Produce(fromId, authorId, content, quoteContent, model.MsgTypeUserWatch)
}

// 内容被点赞
func (s *notificationService) SendTopicLikeNotification(topicLike *model.TopicLike) {
	var (
		fromId       = topicLike.UserId // 消息发送人
		authorId     int64              // 点赞者编号
		content      string             // 消息内容
		quoteContent string             // 引用内容
	)
	topic := dao.TopicDao.Get(topicLike.TopicId)
	if topic != nil {
		authorId = topic.UserId
		content = "我点赞了你的话题：" + topic.Title
		quoteContent = ""
	}

	if authorId <= 0 {
		return
	}
	// 给帖子作者发消息
	s.Produce(fromId, authorId, content, quoteContent, model.MsgTypeTopicLike)
}

// 评论被回复消息
func (s *notificationService) SendCommentNotification(comment *model.Comment) {
	quote := s.getQuoteComment(comment.QuoteId)
	summary := utils.GetHtmlSummary(comment.Content)

	var (
		fromId       = comment.UserId // 消息发送人
		authorId     int64            // 帖子作者编号
		content      string           // 消息内容
		quoteContent string           // 引用内容
	)

	if comment.EntityType == model.EntityTypeArticle { // 文章被评论
		article := dao.ArticleDao.Get(comment.EntityId)
		if article != nil {
			authorId = article.UserId
			content = "我回复了你的文章：" + summary
			quoteContent = "《" + article.Title + "》"
		}
	} else if comment.EntityType == model.EntityTypeTopic { // 话题被评论
		topic := dao.TopicDao.Get(comment.EntityId)
		if topic != nil {
			authorId = topic.UserId
			content = "我回复了你的话题：" + summary
			quoteContent = "《" + topic.Title + "》"
		}
	}

	if authorId <= 0 {
		return
	}

	if quote != nil { // 回复跟帖
		if comment.UserId != authorId && quote.UserId != authorId { // 回复人和帖子作者不是同一个人，并且引用的用户不是帖子作者，需要给帖子作者也发送一下消息
			// 给帖子作者发消息
			s.Produce(fromId, authorId, content, quoteContent, model.MsgTypeComment)
		}

		// 给被引用的人发消息
		s.Produce(fromId, quote.UserId, "我回复了你的评论："+summary, utils.GetHtmlSummary(quote.Content), model.MsgTypeComment)
	} else if comment.UserId != authorId { // 回复主贴，并且不是自己回复自己
		// 给帖子作者发消息
		s.Produce(fromId, authorId, content, quoteContent, model.MsgTypeComment)
	}
}

func (s *notificationService) getQuoteComment(quoteId int64) *model.Comment {
	if quoteId <= 0 {
		return nil
	}
	return dao.CommentDao.Get(quoteId)
}

// 生产，将消息数据放入chan
func (s *notificationService) Produce(fromId, toId int64, content, quoteContent string, msgType int) {
	s.Consume()
	s.notificationsChan <- &Notification{
		FromID:       fromId,
		ToID:         toId,
		Content:      content,
		QuoteContent: quoteContent,
		Type:         msgType,
	}
}

// 消费，消费chan中的消息
func (s *notificationService) Consume() {
	s.notificationsConsumeOnce.Do(func() {
		go func() {
			logger.Logger.Info("开始消费系统消息...")
			for {
				func() {
					defer utils.Recover(logger.Logger, nil, "")

					m := <-s.notificationsChan
					logger.Logger.Info("处理消息：from=%s to=%s", m.FromID, m.ToID)

					if err := s.Create(m); err != nil {
						logger.Logger.Info("创建消息发生异常...")
					}
				}()
			}
		}()
	})
}
