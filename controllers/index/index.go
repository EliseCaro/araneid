package index

import (
	"github.com/astaxie/beego"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"regexp"
)

type Index struct {
	Main
}

// @router / [get]
func (c *Index) Index() {
	beego.Info(_func.WebPageSize())
	beego.Info(regexp.MustCompile(`https://www.chazidian.com/([a-z]+)_([a-z]+)_([a-zA-Z0-9]{32})/$`).MatchString("https://www.chazidian.com/r_ci_a4d0a723bbf7bdb8e473c2737fc62255/"))
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
