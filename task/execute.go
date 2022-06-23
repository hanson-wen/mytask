package task

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/robfig/cron/v3"
	"mytask/dao"
	"mytask/model"
	"mytask/util"
	"time"
)

// BeginTask pre logic
func BeginTask(taskMeta model.TaskMeta, manual bool) (err error) {
	logs.Debug("============== task name:%s", taskMeta.TaskName)
	taskMeta, _ = dao.DescribeTaskMetaByID(taskMeta.ID) // get the newest

	// step1 re check the state, if cancel or done, stop
	if manual == false &&
		(taskMeta.TaskState == model.TaskState.Cancel || taskMeta.TaskState == model.TaskState.Done) {
		taskCron.Remove(cron.EntryID(taskMeta.EntryID))
		return nil
	}
	// step2 insert task record
	taskRecord := &model.TaskRecord{
		TaskKey:        taskMeta.TaskKey,
		TaskName:       taskMeta.TaskName,
		TaskParam:      taskMeta.TaskParam,
		ExecutorSerial: taskMeta.ExecutorSerial,
		TaskProgress:   0,
		TaskResult:     "{}",
		CreatedTime:    time.Now(),
		ModifiedTime:   time.Now(),
	}
	_ = dao.CreateData(taskRecord)

	// step3 update task state
	newMeta := &model.TaskMeta{
		ID:            taskMeta.ID,
		ExecutedTimes: taskMeta.ExecutedTimes + 1,
		ModifiedTime:  time.Now(),
	}
	if taskMeta.TaskState == model.TaskState.Undo {
		newMeta.TaskState = model.TaskState.Doing
	}
	_ = dao.ModifyData(newMeta)

	// step4 execute task
	go ExecuteTask(taskMeta, *taskRecord)
	return err
}

// ExecuteTask run the task
func ExecuteTask(taskMeta model.TaskMeta, taskRecord model.TaskRecord) {
	// step1 parse task_param to do different task logic
	err := taskLogic()
	if err != nil {
		util.Alarm(fmt.Sprintf("[task_name:%s] run err: %s", taskMeta.TaskName, err.Error()))
	}

	// step3 if success, update newMeta info
	newMeta := &model.TaskMeta{
		ID:           taskMeta.ID,
		ModifiedTime: time.Now(),
	}
	if taskMeta.TaskType == model.TaskMetaType.Cron {
		var nextExecutedTime time.Time
		var timeParser = cron.NewParser(
			cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
		schedule, _ := timeParser.Parse(taskMeta.TimeRule)
		nextExecutedTime = schedule.Next(time.Now())
		newMeta.NextExecuteTime = nextExecutedTime
	} else if taskMeta.TaskType == model.TaskMetaType.Once {
		if taskMeta.TaskState == model.TaskState.Doing {
			newMeta.TaskState = model.TaskState.Done
		}
	}

	_ = dao.ModifyData(newMeta)

	// step4 update task record
	updateTaskRecord := &model.TaskRecord{
		ID:           taskRecord.ID,
		TaskProgress: 100,
		TaskResult:   "{\"res\":\"success\"}",
		ModifiedTime: time.Now(),
	}
	_ = dao.ModifyData(updateTaskRecord)
}

// systemCronTask .
func systemCronTask(task *cron.Cron) {
	// node heartbeat check
	var err error
	_, err = task.AddFunc("*/10 * * * * ?", func() {
		nodePing(false)
	})
	if err != nil {
		logs.Info("add crontab fun: nodePing, err: %s", err.Error())
	}

	// elect master node and set
	_, err = task.AddFunc("*/10 * * * * ?", func() {
		node, err := electMaster()
		if err != nil {
			logs.Error("elect master err: %s", err.Error())
		} else {
			logs.Info("master ip is : %s", node.Host)
		}
	})
	if err != nil {
		logs.Error("add crontab fun: electMaster, err: %s", err.Error())
	}

	// monitor the broken node and alarm
	_, err = task.AddFunc("*/20 * * * * ?", monitorNode)
	if err != nil {
		logs.Info("add crontab fun: monitorNode, err: %s", err.Error())
	}
}

// userCronTask .
func userCronTask(task *cron.Cron) {
	ip, err := util.GetIPV4()
	if err != nil {
		util.Alarm("get ip err:" + err.Error())
		return
	}
	// step1 get the tasks which undo assigned to this node
	query := &dao.TaskMetaQuery{
		TaskType:        model.TaskMetaType.Cron,
		TaskStateUndone: true,
		NoEntryID:       true,
		ExecutorSerial:  ip,
		Limit:           10,
	}
	for {
		tasks, err := dao.DescribeAllTaskMeta(query)
		if err != nil {
			time.Sleep(time.Second * 10)
			continue
		}
		if len(tasks) == 0 {
			time.Sleep(time.Second * 10)
			continue
		}

		// step2 add cron
		for _, t := range tasks {
			temp := *t
			entryID, err := task.AddFunc(t.TimeRule, func() {
				_ = BeginTask(temp, false)
			})
			logs.Info("==========entryID:%v", entryID)
			if err != nil {
				logs.Info("add crontab task: %s, err: %s", t.TaskName, err.Error())
			} else {
				// set meta running
				_ = dao.ModifyData(&model.TaskMeta{ID: temp.ID, TaskState: model.TaskState.Doing, EntryID: int(entryID)})
			}
		}
		time.Sleep(time.Second * 10)
	}
}

// userOnceTask
func userOnceTask() {
	ip, err := util.GetIPV4()
	if err != nil {
		util.Alarm("get ip err:" + err.Error())
		return
	}
	// step1 get the tasks which undo assigned to this node
	query := &dao.TaskMetaQuery{
		TaskType:       model.TaskMetaType.Once,
		TaskState:      model.TaskState.Undo,
		ExecutorSerial: ip,
		Limit:          10,
	}
	for {
		tasks, err := dao.DescribeAllTaskMeta(query)
		if err != nil {
			time.Sleep(time.Second * 10)
			continue
		}
		if len(tasks) == 0 {
			time.Sleep(time.Second * 10)
			continue
		}

		for _, t := range tasks {
			temp := *t
			newMeta := &model.TaskMeta{
				ID:            temp.ID,
				TaskState:     model.TaskState.Doing,
				ExecutedTimes: 1,
				ModifiedTime:  time.Now(),
			}
			_ = dao.ModifyData(newMeta)
			BeginTask(temp, false)
		}
	}
}

// taskLogic .
func taskLogic() (err error) {
	return err
}
