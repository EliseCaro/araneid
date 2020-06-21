package index

import (
	"github.com/astaxie/beego"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
)

type Index struct {
	Main
}

// @router / [get]
func (c *Index) Index() {
	beego.Info(_func.WebPageSize())
	c.Data["Title"] = "home"
}

// @router /index/test [post]
func (c *Index) Test() {
	maps := make(map[string]interface{})
	if e := c.ParseForm(&maps); e != nil {
		beego.Warn(e)
	} else {
		beego.Info(maps)
	}
	//c.Succeed(&controllers.ResultJson{Message: " 发布数据成功拉！"})
	c.Fail(&controllers.ResultJson{Message: " 发布数据失败拉！"})
}
