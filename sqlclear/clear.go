package sqlclear

import (
	"gitee.com/wuntsong/chuangxinyi-communicate/dao"
	"gitee.com/wuntsong/chuangxinyi-communicate/logger"
	"gitee.com/wuntsong/chuangxinyi-communicate/model"
	errors "github.com/wuntsong-org/wterrors"
	"reflect"
	"time"
)

var ClearList = []interface{}{
	&model.User{}, &model.Tag{}, &model.Article{}, &model.ArticleTag{}, &model.Comment{}, &model.Favorite{},
	&model.Topic{}, &model.Node{}, &model.TopicTag{}, &model.TopicLike{}, &model.Setting{}, &model.Link{},
	&model.UserWatch{}, &model.UserScore{}, &model.UserScoreLog{},
}

func StartClear() errors.WTError {
	tableCount := len(ClearList)

	logger.Logger.WXInfo("开始清理数据库，计划清理 %d 张表", tableCount)

	errorTables := 0
	processTables := 0
	processRows := int64(0)

	deleteTime := time.Now().Add(-24 * time.Hour * 180)

	for _, c := range ClearList {
		res := dao.DB().Unscoped().Delete(c, "delete_time < ?", deleteTime)
		if res.Error != nil {
			errorTables += 1
			logger.Logger.WXInfo("清理数据库表（%s）时出错：%s", reflect.TypeOf(c).Name(), res.Error.Error())
			continue
		}

		if res.RowsAffected > 0 {
			processTables += 1
			processRows += res.RowsAffected
			logger.Logger.WXInfo("清理数据库表（%s）成功，清理了%d条数据", reflect.TypeOf(c).Name(), res.RowsAffected)
		}
	}

	logger.Logger.WXInfo("清理数据库结束，清理出错 %d 张表，清理成功 %d 张表，清理成功删除 %d 行", errorTables, processTables, processRows)

	return nil
}
