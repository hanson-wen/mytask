package controllers

import (
	"errors"
	"github.com/robfig/cron/v3"
	"mytask/dao"
	errx "mytask/error"
	"mytask/model"
	"mytask/util"
	"time"
)

// TaskMetaController .
type TaskMetaController struct {
	BaseController
}

// CreateTaskMeta add task meta
func (c *TaskMetaController) CreateTaskMeta() {
	// step1 get param
	taskName := c.GetString("task_name")
	taskDesc := c.GetString("task_desc")
	taskParam := c.GetString("task_param")
	taskType, _ := c.GetInt8("task_type")
	timeRule := c.GetString("time_rule")
	taskKey := util.GenerateUuid()
	user := model.Auth(c.GetString("token"))

	// step2 check param
	if taskParam == "" {
		taskParam = "{}"
	}
	if taskType == model.TaskMetaType.Cron && timeRule == "" {
		c.Ctx.WriteString(errx.ResponseError(errx.DbExecuteErr, errors.New("cron task must need time rule")))
		return
	}

	// step3 for getting the next execute time
	var nextExecutedTime time.Time
	var timeParser = cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	if taskType == model.TaskMetaType.Once {
		nextExecutedTime = time.Now()
	} else if taskType == model.TaskMetaType.Cron {
		schedule, _ := timeParser.Parse(timeRule)
		start := time.Now()
		nextExecutedTime = schedule.Next(start)
	}

	// step4 combine data to insert db
	meta := &model.TaskMeta{
		TaskKey:         taskKey,
		TaskName:        taskName,
		TaskType:        taskType,
		TaskParam:       taskParam,
		TimeRule:        timeRule,
		TaskState:       model.TaskState.Undo,
		ExecutedTimes:   0,
		NextExecuteTime: nextExecutedTime,
		TaskDesc:        taskDesc,
		Creator:         user.Name,
		CreatedTime:     time.Now(),
		ModifiedTime:    time.Now(),
	}
	err := dao.CreateData(meta)
	if err != nil {
		c.Ctx.WriteString(errx.ResponseError(errx.DbExecuteErr, errors.New("system busy, try again later")))
		return
	}
	c.Ctx.WriteString(errx.SuccessWithData(meta))
	return
}

// ModifyTaskMeta modify task meta data
func (c *TaskMetaController) ModifyTaskMeta() {
	// step1 get param, task_type can not be modified
	taskKey := c.GetString("task_key")
	taskName := c.GetString("task_name")
	taskDesc := c.GetString("task_desc")
	taskParam := c.GetString("task_param")
	timeRule := c.GetString("time_rule")
	user := model.Auth(c.GetString("token"))

	// step2 check param
	if taskParam == "" {
		taskParam = "{}"
	}

	// step3 get info and check valid
	task, err := dao.DescribeTaskMetaByKey(taskKey)
	if err != nil {
		c.Ctx.WriteString(errx.ResponseError(errx.DbExecuteErr, errors.New("system busy, try again later")))
		return
	}
	if task.ID == 0 {
		c.Ctx.WriteString(errx.ResponseError(errx.LogicErr, errors.New("record not exists")))
		return
	}
	if task.Creator != user.Name {
		c.Ctx.WriteString(errx.ResponseError(errx.LogicErr, errors.New("only can modify your own task")))
		return
	}

	// todo other check, if task is running,return

	// step 4 for getting the new next execute time
	var nextExecutedTime time.Time
	var timeParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.
		Descriptor)
	if task.TaskType == model.TaskMetaType.Once {
		nextExecutedTime = time.Now()
	} else if task.TaskType == model.TaskMetaType.Cron {
		schedule, _ := timeParser.Parse(timeRule)
		start := time.Now()
		nextExecutedTime = schedule.Next(start)
	}

	// step 5 combine data to insert to db
	meta := &model.TaskMeta{
		ID:              task.ID,
		TaskName:        taskName,
		TaskParam:       taskParam,
		TimeRule:        timeRule,
		NextExecuteTime: nextExecutedTime,
		TaskDesc:        taskDesc,
	}
	err = dao.ModifyData(meta)
	if err != nil {
		c.Ctx.WriteString(errx.ResponseError(errx.DbExecuteErr, errors.New("system busy, try again later")))
		return
	}
	c.Ctx.WriteString(errx.SuccessWithData(meta))
	return
}

// DeleteTaskMeta delete
func (c *TaskMetaController) DeleteTaskMeta() {
	// step1 get param
	id, _ := c.GetInt("id")
	user := model.Auth(c.GetString("token"))

	// step2 check valid
	task, err := dao.DescribeTaskMetaByID(id)
	if err != nil {
		c.Ctx.WriteString(errx.ResponseError(errx.DbExecuteErr, errors.New("system busy, try again later")))
		return
	}
	if task.ID == 0 {
		c.Ctx.WriteString(errx.ResponseError(errx.LogicErr, errors.New("record not exists")))
		return
	}
	if task.Creator != user.Name {
		c.Ctx.WriteString(errx.ResponseError(errx.LogicErr, errors.New("only can delete your own task")))
		return
	}

	// todo other check, eg: the task is running,return

	// step4 delete
	err = dao.DeleteTaskMeta(id)
	if err != nil {
		c.Ctx.WriteString(errx.ResponseError(errx.DbExecuteErr, errors.New("system busy, try again later")))
		return
	}
	c.Ctx.WriteString(errx.SuccessWithoutData())
	return
}

// DescribeTaskMetas describe task meta list
func (c *TaskMetaController) DescribeTaskMetas() {
	// filter param
	limit, _ := c.GetInt("limit")
	offset, _ := c.GetInt("offset")
	id, _ := c.GetInt("id")
	user := model.Auth(c.GetString("token"))
	query := &dao.TaskMetaQuery{
		ID:      id,
		Creator: user.Name,
		Limit:   limit,
		Offset:  offset,
		Order:   "id desc",
	}

	// execute select
	list, total, err := dao.DescribeTaskMetas(query)
	if err != nil {
		c.Ctx.WriteString(errx.ResponseError(errx.DbExecuteErr, errors.New("system busy, try again later")))
		return
	}
	c.Ctx.WriteString(errx.SuccessWithData(struct {
		Total int
		List  []*model.TaskMeta
	}{Total: total, List: list}))
	return
}
