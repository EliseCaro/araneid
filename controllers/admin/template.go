package admin

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/go-playground/validator"
)

/** 蜘蛛池模板管理 **/
type Template struct{ Main }

// @router /template/index [get,post]
func (c *Template) Index() {
	var group []*spider.TemplateClass
	class := c.templateService.Groups()
	if j, _ := json.Marshal(class); len(j) > 0 {
		_ = json.Unmarshal(j, &group)
	}
	for _, v := range group {
		v.Child = c.templateService.Items(v.Id)
	}
	c.Data["class"] = group
}

// @router /template/create [get,post]
func (c *Template) Create() {
	if c.IsAjax() {
		item := spider.Template{}
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Insert(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "创建分组成功", Url: beego.URLFor("Template.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "创建分组失败，请稍后再试！"})
		}
	}
}

// @router /template/edit [get,post]
func (c *Template) Edit() {
	id := c.GetMustInt(":id", "非法请求！")
	if c.IsAjax() {
		item := c.templateService.One(id)
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Update(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "更新分组成功", Url: beego.URLFor("Template.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "更新分组失败，请稍后再试！"})
		}
	}
	c.Data["info"] = c.templateService.One(id)
}

// @router /template/delete [post]
func (c *Template) Delete() {
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.templateService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除成功！马上返回中。。。",
			Url:     beego.URLFor("Template.Index"),
		})
	}
}
