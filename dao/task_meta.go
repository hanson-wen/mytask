package dao

import (
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	"mytask/model"
)

// TaskMetaQuery the filter
type TaskMetaQuery struct {
	ID                 int
	TaskKey            string
	Creator            string
	ExecutorSerial     string
	IsNoExecutorSerial bool
	TaskState          int8
	TaskStateUndone    bool
	TaskType           int8
	NoEntryID          bool
	Limit              int
	Offset             int
	Order              string
}

// TaskMetaCondition all condition maybe used
func TaskMetaCondition(db *gorm.DB, query *TaskMetaQuery, total *int) *gorm.DB {
	if query.ID > 0 {
		db = db.Where("id = ?", query.ID)
	}
	if query.TaskKey != "" {
		db = db.Where("task_key = ?", query.TaskKey)
	}
	if query.Creator != "" {
		db = db.Where("creator = ?", query.Creator)
	}
	if query.TaskState > 0 {
		db = db.Where("task_state = ?", query.TaskState)
	}
	if query.TaskType > 0 {
		db = db.Where("task_type = ?", query.TaskType)
	}
	if query.ExecutorSerial != "" {
		db = db.Where("executor_serial = ?", query.ExecutorSerial)
	}
	if query.IsNoExecutorSerial == true {
		db = db.Where("executor_serial = ?", "")
	}
	if query.TaskStateUndone == true {
		db = db.Where("task_state != ? and task_state < ?", model.TaskState.Cancel, model.TaskState.Done)
	}
	if query.NoEntryID == true {
		db = db.Where("entry_id = ?", 0)
	}
	if total != nil {
		db = db.Count(total)
	}
	if query.Limit > 0 {
		db = db.Limit(query.Limit)
	}
	if query.Offset > 0 {
		db = db.Offset(query.Offset)
	}
	if query.Order != "" {
		db = db.Order(query.Order)
	}
	return db
}

// DescribeTaskMetaByKey .
func DescribeTaskMetaByKey(taskKey string) (task model.TaskMeta, err error) {
	if taskKey == "" {
		return
	}
	list, err := DescribeAllTaskMeta(&TaskMetaQuery{TaskKey: taskKey})
	if err != nil {
		return
	}
	if len(list) > 0 {
		task = *list[0]
	}
	return
}

// DescribeTaskMetaByID .
func DescribeTaskMetaByID(id int) (task model.TaskMeta, err error) {
	if id == 0 {
		return
	}
	list, err := DescribeAllTaskMeta(&TaskMetaQuery{ID: id})
	if err != nil {
		return
	}
	if len(list) > 0 {
		task = *list[0]
	}
	return
}

// DeleteTaskMeta .
func DeleteTaskMeta(id int) (err error) {
	err = db.Delete(&model.TaskMeta{}, "id = ?", id).Error
	return
}

// DescribeAllTaskMeta all query list, will not count total
func DescribeAllTaskMeta(query *TaskMetaQuery) (list []*model.TaskMeta, err error) {
	err = TaskMetaCondition(db.Model(&model.TaskMeta{}), query, nil).Find(&list).Error
	if err != nil {
		logs.Error("describe task meta err:%s", err.Error())
	}
	return
}

// AssignTaskNoSerial maybe some node was down ,need to assign the task to no serial
// for redo assign later
// pay attention to task state(only undo and doing)
func AssignTaskNoSerial(serial string) (err error) {
	err = db.Model(&model.TaskMeta{}).Where("executor_serial = ? and task_state in(1,2)", serial).
		Update(map[string]interface{}{"executor_serial": "", "entry_id": 0}).Error
	return
}

// DescribeTaskMetas .
func DescribeTaskMetas(query *TaskMetaQuery) (list []*model.TaskMeta, total int, err error) {
	err = TaskMetaCondition(db.Model(&model.TaskMeta{}), query, &total).Find(&list).Error
	if err != nil {
		logs.Error("describe task meta err:%s", err.Error())
	}
	return
}
