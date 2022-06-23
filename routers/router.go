package routers

import (
	"github.com/astaxie/beego"
	"mytask/controllers"
)

func init() {

	// methods run diff fun, use ';' to sep, eg:"get:fun1;post:fun2"
	// all method run the same fun, user ',' between the methods，eg:"get,post:fun"
	/**
	* get: GET request
	* post: POST request
	* put: PUT request
	* delete: DELETE request
	* patch: PATCH request
	* options: OPTIONS request
	* head: HEAD request
	 */

	// about task meta
	beego.Router("/api/taskMeta", &controllers.TaskMetaController{},
		"post:CreateTaskMeta;put:ModifyTaskMeta;get:DescribeTaskMetas;delete:DeleteTaskMeta")

	// about task
	beego.Router("/api/task", &controllers.TaskController{},
		"post:RunTask")
}
