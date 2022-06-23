package util

import "github.com/astaxie/beego/logs"

// Alarm send sth to notice developers
func Alarm(s interface{}) {
	logs.Error("=======: %v", s)
}
