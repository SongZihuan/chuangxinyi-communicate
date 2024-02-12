package convert

import (
	"gitee.com/wuntsong/chuangxinyi-communicate/urls"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"gitee.com/wuntsong/chuangxinyi-communicate/model"
	"gitee.com/wuntsong/chuangxinyi-communicate/service"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
)

func ToFavorites(favorites []model.Favorite) []model.FavoriteResponse {
	if favorites == nil || len(favorites) == 0 {
		return nil
	}
	var responses []model.FavoriteResponse
	for _, favorite := range favorites {
		responses = append(responses, *ToFavorite(&favorite))
	}
	return responses
}

func ToFavorite(favorite *model.Favorite) *model.FavoriteResponse {
	rsp := &model.FavoriteResponse{}
	rsp.FavoriteId = favorite.ID
	rsp.EntityType = favorite.EntityType
	rsp.EntityId = favorite.EntityId
	rsp.CreateTime = favorite.CreateTime

	if favorite.EntityType == model.EntityTypeArticle {
		article := service.ArticleService.Get(favorite.EntityId)
		if article == nil || article.Status != model.StatusOk {
			rsp.Deleted = true
		} else {
			rsp.User = ToUserById(article.UserId)
			rsp.Title = article.Title
			rsp.Content = utils.GetHtmlSummary(article.Content)
		}
	} else {
		topic := service.TopicService.Get(favorite.EntityId)
		if topic == nil || topic.Status != model.StatusOk {
			rsp.Deleted = true
		} else {
			rsp.User = ToUserById(topic.UserId)
			rsp.Title = topic.Title
			rsp.Content = utils.GetHtmlSummary(topic.Content)
		}
	}
	return rsp
}

func ToHtmlContent(htmlContent string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return htmlContent
	}

	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		href := selection.AttrOr("href", "")

		if len(href) == 0 {
			return
		}

		// 不是内部链接
		if !urls.IsInternalUrl(href) {
			selection.SetAttr("target", "_blank")
			selection.SetAttr("rel", "external nofollow") // 标记站外链接，搜索引擎爬虫不传递权重值
		}

		// 如果是锚链接
		if urls.IsAnchor(href) {
			selection.ReplaceWithHtml(selection.Text())
		}

		// 如果a标签没有title，那么设置title
		title := selection.AttrOr("title", "")
		if len(title) == 0 {
			selection.SetAttr("title", selection.Text())
		}
	})

	// 处理图片
	doc.Find("img").Each(func(i int, selection *goquery.Selection) {
		src := selection.AttrOr("src", "")
		// 处理第三方图片
		if strings.Contains(src, "qpic.cn") {
			src = utils.ParseUrl("/api/img/proxy").AddQuery("url", src).BuildStr()
			// selection.SetAttr("src", src)
		}
		selection.SetAttr("data-src", src)
	})

	html, err := doc.Find("body").Html()
	if err != nil {
		return htmlContent
	}
	return html
}
