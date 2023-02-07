// @Author Bing 
// @Date 2023/2/6 19:22:00 
// @Desc
package mycron

import (
	"fmt"
	"github.com/qingconglaixueit/abing_logger"
	"github.com/robfig/cron/v3"
)

func StartCron(cronStr string, fn func()) {
	abing_logger.SugarLogger.Info("Starting...")

	// 定义一个cron运行器
	c := cron.New(cron.WithSeconds()) //精确到秒`
	defer c.Stop()

	_, err := c.AddFunc(cronStr, fn)
	if err != nil {
		tErrStr := fmt.Sprintf("cronErr :%+v", err)
		abing_logger.SugarLogger.Errorln(tErrStr)
		panic(tErrStr)
	}

	// 开始
	c.Start()
	select {}
}
