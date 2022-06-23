package main

import (
	_ "mytask/routers"
	"mytask/task"

	"github.com/astaxie/beego"
)

func main() {
	// don't render automatically, there are apis
	beego.BConfig.WebConfig.AutoRender = false
	task.Init()
	beego.Run()
}
