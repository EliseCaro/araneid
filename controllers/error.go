package controllers

import (
	"github.com/astaxie/beego"
)

type Error struct {
	Base
}

/** 404 错误处理 **/
func (c *Error) Error404() {
	if c.IsAjax() {
		c.Ctx.Output.Status = 200
		c.Fail(&ResultJson{Message: "啊～哦～你要查看的页面不存在或者已经删除！"})
	} else {
		c.Data["content"] = "啊～哦～你要查看的页面不存在或者已经删除！"
	}
}

/** 500 处理功能 **/
func (c *Error) Error500() {
	if c.IsAjax() {
		c.Ctx.Output.Status = 200
		c.Fail(&ResultJson{Message: error.Error(c.Data["content"].(error))})
	} else {
		c.Data["content"] = c.Data["content"].(error)
		c.Data["title"] = c.Data["title"].(error)
	}
}

/** 数据库检测功能 **/
func (c *Error) ErrorDb() {
	beego.Error("数据库走丢了哦～")
}
