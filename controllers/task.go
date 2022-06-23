package controllers

import (
	"errors"
	"mytask/dao"
	errx "mytask/error"
	"mytask/model"
	"mytask/task"
)

// TaskController .
type TaskController struct {
	BaseController
}

// RunTask manual run task
func (c *TaskController) RunTask() {
	// step1 get param
	id, _ := c.GetInt("id")
	user := model.Auth(c.GetString("token"))

	// step2 check valid
	taskMeta, err := dao.DescribeTaskMetaByID(id)
	if err != nil {
		c.Ctx.WriteString(errx.ResponseError(errx.DbExecuteErr, errors.New("system busy, try again later")))
		return
	}
	if taskMeta.ID == 0 {
		c.Ctx.WriteString(errx.ResponseError(errx.LogicErr, errors.New("record not exists")))
		return

	}
	if taskMeta.Creator != user.Name {
		c.Ctx.WriteString(errx.ResponseError(errx.LogicErr, errors.New("only can delete your own taskMeta")))
		return
	}
	err = task.BeginTask(taskMeta, true)
	if err != nil {
		c.Ctx.WriteString(errx.ResponseError(errx.LogicErr, err))
		return
	}
	c.Ctx.WriteString(errx.SuccessWithoutData())
	return
}
