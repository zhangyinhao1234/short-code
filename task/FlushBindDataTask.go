package task

import (
	"github.com/robfig/cron/v3"
	"log"
	"short-code/service/domain/binding"
	"short-code/service/domain/unuse"
)

var (
	bindingService   = binding.BindingService{}
	unUseCodeService = unuse.UnUseCodeService{}
)

type FlushBindDataTask struct {
}

func (e *FlushBindDataTask) Run() {
	c := cron.New(cron.WithSeconds())
	spec := "*/20 * * * * *" // 每隔20s执行一次，cron格式（秒，分，时，天，月，周）
	_, err := c.AddFunc(spec, func() {
		//global.LOG.Info("定时刷新缓存数据")
		bindingService.Flush()
	})
	if err != nil {
		log.Fatal("定时任务初始化异常", err)
	}
	c.Start()
}
