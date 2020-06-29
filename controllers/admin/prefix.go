package admin

import (
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	"github.com/beatrice950201/araneid/extend/model/spider"
)

/** 蜘蛛池资源模型前缀管理 **/
type Prefix struct{ Main }

// @router /prefix/index [get,post]
func (c *Prefix) Index() {
	id := c.GetMustInt(":model", "非法请求,请稍后再试...")
	c.Data["select"] = c.prefixService.Select(id)
	c.Data["model"] = id
}

// @router /prefix/create [get,post]
func (c *Prefix) Create() {
	model := c.GetMustInt("model", "参数错误；创建失败")
	if id, err := orm.NewOrm().Insert(&spider.Prefix{Model: model}); err == nil {
		c.Succeed(&controllers.ResultJson{Message: "success!!", Data: id})
	} else {
		c.Fail(&controllers.ResultJson{Message: err.Error()})
	}
}

// @router /prefix/edit [get,post]
func (c *Prefix) Edit() {
	id := c.GetMustInt("id", "参数错误;ID接受失败")
	field := c.GetMustString("field", "参数错误;字段名接受失败")
	one := c.prefixService.One(id)
	if err := c.ParseForm(&one); err != nil {
		c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
	}
	if _, err := orm.NewOrm().Update(&one, field); err == nil {
		c.Succeed(&controllers.ResultJson{Message: "success!!"})
	} else {
		c.Fail(&controllers.ResultJson{Message: err.Error()})
	}
}

// @router /prefix/delete [post]
func (c *Prefix) Delete() {
	id := c.GetMustInt("id", "参数错误;ID接受失败")
	if _, e := orm.NewOrm().Delete(&spider.Prefix{Id: id}); e != nil {
		c.Fail(&controllers.ResultJson{Message: e.Error()})
	} else {
		c.Succeed(&controllers.ResultJson{Message: "success!!"})
	}
}
