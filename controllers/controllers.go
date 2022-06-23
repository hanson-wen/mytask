// Package controllers .
package controllers

import "github.com/astaxie/beego"

type BaseController struct {
	beego.Controller
}

func (c *BaseController) errorShow(errMsg string) {
	c.Data["errMsg"] = errMsg
	c.TplName = "error.html"
	c.Render()
}

func (c *BaseController) sysErrorShow() {
	c.errorShow("服务器繁忙")
}

func init() {
}
