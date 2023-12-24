package cron

import (
	"github.com/robfig/cron"
	errors "github.com/wuntsong-org/wterrors"

	"gitee.com/wuntsong/chuangxinyi-communicate/service"
)

func Setup() errors.WTError {
	err := startSchedule()
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}

func startSchedule() errors.WTError {
	c := cron.New()

	// Generate RSS
	err := addCronFunc(c, "@every 30m", func() {
		service.ArticleService.GenerateRss()
		service.TopicService.GenerateRss()
	})
	if err != nil {
		return errors.WarpQuick(err)
	}

	c.Start()

	return nil
}

func addCronFunc(c *cron.Cron, sepc string, cmd func()) errors.WTError {
	err := c.AddFunc(sepc, cmd)
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}
