package convert

import (
	"html/template"

	"gitee.com/wuntsong/chuangxinyi-communicate/cache"
	"gitee.com/wuntsong/chuangxinyi-communicate/model"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils/html"
)

func ToArticle(article *model.Article) *model.ArticleResponse {
	if article == nil {
		return nil
	}

	rsp := &model.ArticleResponse{}
	rsp.ArticleId = article.ID
	rsp.Title = article.Title
	rsp.Summary = article.Summary
	rsp.Share = article.Share
	rsp.SourceUrl = article.SourceUrl
	rsp.ViewCount = article.ViewCount
	rsp.CreateTime = article.CreateTime

	rsp.User = ToUserDefaultIfNull(article.UserId)

	tagIds := cache.ArticleTagCache.Get(article.ID)
	tags := cache.TagCache.GetList(tagIds)
	rsp.Tags = ToTags(tags)

	mr := html.NewHtml(html.WithTOC()).Run(article.Content)
	rsp.Content = template.HTML(ToHtmlContent(mr.ContentHtml))
	rsp.Toc = template.HTML(mr.TocHtml)
	if len(rsp.Summary) == 0 {
		rsp.Summary = mr.SummaryText
	}

	return rsp
}

func ToArticles(articles []model.Article) []model.ArticleResponse {
	if articles == nil || len(articles) == 0 {
		return nil
	}
	var responses []model.ArticleResponse
	for _, article := range articles {
		responses = append(responses, *ToArticle(&article))
	}
	return responses
}

func ToSimpleArticle(article *model.Article) *model.ArticleSimpleResponse {
	if article == nil {
		return nil
	}

	rsp := &model.ArticleSimpleResponse{}
	rsp.ArticleId = article.ID
	rsp.Title = article.Title
	rsp.Summary = article.Summary
	rsp.Share = article.Share
	rsp.SourceUrl = article.SourceUrl
	rsp.ViewCount = article.ViewCount
	rsp.CreateTime = article.CreateTime

	rsp.User = ToUserDefaultIfNull(article.UserId)

	tagIds := cache.ArticleTagCache.Get(article.ID)
	tags := cache.TagCache.GetList(tagIds)
	rsp.Tags = ToTags(tags)

	if len(rsp.Summary) == 0 {
		mr := html.NewHtml(html.WithTOC()).Run(article.Content)
		rsp.Summary = mr.SummaryText
	}

	return rsp
}

func ToSimpleArticles(articles []model.Article) []model.ArticleSimpleResponse {
	if articles == nil || len(articles) == 0 {
		return nil
	}
	var responses []model.ArticleSimpleResponse
	for _, article := range articles {
		responses = append(responses, *ToSimpleArticle(&article))
	}
	return responses
}
