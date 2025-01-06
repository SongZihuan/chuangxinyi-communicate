package sqlclear

import (
	"github.com/SongZihuan/chuangxinyi-communicate/dao"
	"github.com/SongZihuan/chuangxinyi-communicate/logger"
	errors "github.com/wuntsong-org/wterrors"
	"time"
)

type RecordClearable func(time time.Time) (int64, errors.WTError)

type RecordClearData struct {
	Name      string
	Clearable RecordClearable
	Day       int64
}

var RecordClearList = []RecordClearData{
	{
		Name:      "record",
		Clearable: dao.RecordDao.DeleteOld,
		Day:       180,
	},
}

func StartRecordClear() errors.WTError {
	tableCount := len(RecordClearList)

	logger.Logger.WXInfo("开始清理数据库，计划清理 %d 张表", tableCount)

	errorTables := 0
	processTables := 0
	processRows := int64(0)

	for _, c := range RecordClearList {
		rows, err := c.Clearable(time.Now().Add(-time.Hour * 24 * time.Duration(c.Day)))
		if err != nil {
			errorTables += 1
			logger.Logger.WXInfo("清理数据库表（%s）时出错：%s", c.Name, err.Error())
			continue
		}

		if rows > 0 {
			processTables += 1
			processRows += rows
			logger.Logger.WXInfo("清理数据库表（%s）成功，清理了%d条数据", c.Name, rows)
		}
	}

	logger.Logger.WXInfo("清理数据库结束，清理出错 %d 张表，清理成功 %d 张表，清理成功删除 %d 行", errorTables, processTables, processRows)

	return nil
}
