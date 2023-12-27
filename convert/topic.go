package convert

import (
	"html/template"

	"gitee.com/wuntsong/chuangxinyi-communicate/cache"
	"gitee.com/wuntsong/chuangxinyi-communicate/logger"
	"gitee.com/wuntsong/chuangxinyi-communicate/model"
	"gitee.com/wuntsong/chuangxinyi-communicate/service"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils/html"
)

func ToTopic(topic *model.Topic) *model.TopicResponse {
	if topic == nil {
		return nil
	}

	rsp := &model.TopicResponse{}

	rsp.TopicId = topic.ID
	rsp.Type = topic.Type
	rsp.Title = topic.Title
	rsp.User = ToUserDefaultIfNull(topic.UserId)
	rsp.LastCommentTime = topic.LastCommentTime
	rsp.CreateTime = topic.CreateTime
	rsp.ViewCount = topic.ViewCount
	rsp.CommentCount = topic.CommentCount
	rsp.LikeCount = topic.LikeCount

	if topic.NodeId > 0 {
		node := service.NodeService.Get(topic.NodeId)
		rsp.Node = ToNode(node)
	}

	tags := service.TopicService.GetTopicTags(topic.ID)
	rsp.Tags = ToTags(tags)

	mr := html.NewHtml(html.WithTOC()).Run(topic.Content)
	rsp.Content = template.HTML(ToHtmlContent(mr.ContentHtml))
	rsp.Toc = template.HTML(mr.TocHtml)

	if len(topic.ImageList) > 0 {
		if err := utils.ParseJson(topic.ImageList, &rsp.ImageList); err != nil {
			logger.Logger.Error(err.Error())
		}
	}

	return rsp
}

func ToSimpleTopic(topic *model.Topic) *model.TopicSimpleResponse {
	if topic == nil {
		return nil
	}

	rsp := &model.TopicSimpleResponse{}

	rsp.TopicId = topic.ID
	rsp.Type = topic.Type
	rsp.Title = topic.Title
	rsp.User = ToUserDefaultIfNull(topic.UserId)
	rsp.LastCommentUser = ToUserDefaultIfNull(topic.LastCommentUserId)
	rsp.LastCommentTime = topic.LastCommentTime
	rsp.CreateTime = topic.CreateTime
	rsp.ViewCount = topic.ViewCount
	rsp.CommentCount = topic.CommentCount
	rsp.LikeCount = topic.LikeCount

	if len(topic.ImageList) > 0 {
		if err := utils.ParseJson(topic.ImageList, &rsp.ImageList); err != nil {
			logger.Logger.Error(err.Error())
		}
	}

	if topic.NodeId > 0 {
		node := cache.NodeCache.Get(topic.NodeId)
		rsp.Node = ToNode(node)
	}

	tags := service.TopicService.GetTopicTags(topic.ID)
	rsp.Tags = ToTags(tags)
	return rsp
}

func ToSimpleTopics(topics []model.Topic) []model.TopicSimpleResponse {
	if topics == nil || len(topics) == 0 {
		return nil
	}
	var responses []model.TopicSimpleResponse
	for _, topic := range topics {
		responses = append(responses, *ToSimpleTopic(&topic))
	}
	return responses
}
