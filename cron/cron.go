package cron

import (
	"github.com/robfig/cron"
	errors "github.com/wuntsong-org/wterrors"
)

func Setup() errors.WTError {
	// 暂时无任务
	//err := startSchedule()
	//if err != nil {
	//	return errors.WarpQuick(err)
	//}
	//
	return nil
}

func startSchedule() errors.WTError {
	c := cron.New()

	//err := addCronFunc(c, "@every 30m", func() {
	//})
	//if err != nil {
	//	return errors.WarpQuick(err)
	//}

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
