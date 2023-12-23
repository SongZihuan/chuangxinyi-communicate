package cron

import (
	"github.com/robfig/cron"

	"gitee.com/wuntsong/chuangxinyi-communicate/service"
)

func Setup() error {
	err := startSchedule()
	if err != nil {
		return err
	}

	return nil
}

func startSchedule() error {
	c := cron.New()

	// Generate RSS
	err := addCronFunc(c, "@every 30m", func() {
		service.ArticleService.GenerateRss()
		service.TopicService.GenerateRss()
	})
	if err != nil {
		return err
	}

	c.Start()

	return nil
}

func addCronFunc(c *cron.Cron, sepc string, cmd func()) error {
	err := c.AddFunc(sepc, cmd)
	if err != nil {
		return err
	}

	return nil
}
