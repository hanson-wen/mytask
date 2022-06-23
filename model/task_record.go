package model

import "time"

// TaskRecord task runned record
type TaskRecord struct {
	ID             int       `gorm:"primaryKey;column:id;type:int;not null" json:"-"`
	TaskKey        string    `gorm:"index:task_key;column:task_key;type:varchar(255);not null" json:"task_key"`
	TaskName       string    `gorm:"column:task_name;type:varchar(255)" json:"task_name"`
	TaskParam      string    `gorm:"column:task_param;type:json" json:"task_param"`
	TaskProgress   int8      `gorm:"column:task_progress;type:tinyint" json:"task_progress"`
	TaskResult     string    `gorm:"column:task_result;type:json" json:"task_result"`
	ExecutorSerial string    `gorm:"column:executor_serial;type:varchar(255)" json:"executor_serial"`
	CreatedTime    time.Time `gorm:"column:created_time;type:datetime;default:CURRENT_TIMESTAMP" json:"created_time"`
	ModifiedTime   time.Time `gorm:"column:modified_time;type:datetime;default:CURRENT_TIMESTAMP" json:"modified_time"`
}

// TableName get sql table name
func (m *TaskRecord) TableName() string {
	return "t_task_record"
}
