package task

import (
	"github.com/robfig/cron/v3"
)

var taskCron *cron.Cron

func Init() {
	nodeInit()
	go assignedTasks()
	taskCron = cron.New(cron.WithSeconds())

	/**** crontab task add ****/
	systemCronTask(taskCron)
	go userCronTask(taskCron)
	go userOnceTask()

	//start
	taskCron.Start()
}
