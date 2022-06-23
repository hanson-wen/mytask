package model

import "time"

// TaskMeta table model
type TaskMeta struct {
	ID              int       `gorm:"primaryKey;column:id;type:int;not null" json:"id"`
	TaskKey         string    `gorm:"index:task_key;column:task_key;type:varchar(255);not null" json:"task_key"`
	TaskName        string    `gorm:"column:task_name;type:varchar(255)" json:"task_name"`
	TaskParam       string    `gorm:"column:task_param;type:json" json:"task_param"`
	TaskType        int8      `gorm:"column:task_type;type:tinyint" json:"task_type"`
	TimeRule        string    `gorm:"column:time_rule;type:varchar(255)" json:"time_rule"`
	EntryID         int       `gorm:"column:entry_id;type:int" json:"entry_id"`
	NextExecuteTime time.Time `gorm:"column:next_execute_time;type:datetime;default:null" json:"next_execute_time"`
	TaskState       int8      `gorm:"column:task_state;type:tinyint" json:"task_state"`
	ExecutedTimes   int       `gorm:"column:executed_times;type:int" json:"executed_times"`
	TaskDesc        string    `gorm:"column:task_desc;type:varchar(255)" json:"task_desc"`
	Creator         string    `gorm:"column:creator;type:varchar(64)" json:"creator"`
	ExecutorSerial  string    `gorm:"column:executor_serial;type:varchar(255)" json:"executor_serial"`
	CreatedTime     time.Time `gorm:"column:created_time;type:datetime;default:CURRENT_TIMESTAMP" json:"created_time"`
	ModifiedTime    time.Time `gorm:"column:modified_time;type:datetime;default:CURRENT_TIMESTAMP" json:"modified_time"`
}

// TableName get sql table name.xxx
func (m *TaskMeta) TableName() string {
	return "t_task_meta"
}

// TaskMetaTypeDefined all enum of the task type
type TaskMetaTypeDefined struct {
	Once int8
	Cron int8
}

// TaskStateDefined all enum of the task state
type TaskStateDefined struct {
	Cancel int8
	Undo   int8
	Doing  int8
	Done   int8
}

var (
	// TaskMetaType 1 once，2 cron
	TaskMetaType = TaskMetaTypeDefined{
		Once: 1,
		Cron: 2,
	}

	// TaskState -1cancel,1 undo,2 doing,3 done
	TaskState = TaskStateDefined{
		Cancel: -1,
		Undo:   1,
		Doing:  2,
		Done:   3,
	}
)
