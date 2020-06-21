package admin

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	"github.com/beatrice950201/araneid/extend/model/menus"
	"github.com/go-playground/validator"
)

type Menu struct{ Main }

// @router /menu/index [get]
func (c *Menu) Index() {
	current, _ := c.GetInt(":id", 1)
	c.Data["current"] = current
	c.Data["groupMenu"] = c.menusService.PidAllMenus(0)
	c.Data["htmlMenus"] = c.menusService.MenuIndexTreeHtml(c.menusService.GroupAllMenu(), current, current)
}

// @router /menu/edit [get,post]
func (c *Menu) Edit() {
	id := c.GetMustInt(":id", "非法访问；已被拒绝！")
	root, _ := c.GetInt(":root", 1)
	if c.IsAjax() {
		var menu menus.Menus
		if err := c.ParseForm(&menu); err != nil {
			c.Fail(&controllers.ResultJson{Message: "操作错误:" + error.Error(err)})
		}
		message := c.verifyBase.Begin().Struct(menu)
		if message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Update(&menu); err == nil {
			c.Succeed(&controllers.ResultJson{
				Message: "更新节点成功",
				Url:     beego.URLFor("Menu.Index", ":id", root),
			})
		} else {
			c.Fail(&controllers.ResultJson{
				Message: "更新节点失败，请稍后再试！",
			})
		}
	}
	c.Data["info"], _ = c.menusService.Find(id)
	c.Data["root"] = root
	c.Data["groupMenu"] = c.menusService.PidAllMenus(0)
	c.Data["urlsType"] = c.menusService.UrlsType()
}

// @router /menu/create [get,post]
func (c *Menu) Create() {
	root, _ := c.GetInt(":root", 1)
	if c.IsAjax() {
		var menu menus.Menus
		if err := c.ParseForm(&menu); err != nil {
			c.Fail(&controllers.ResultJson{Message: "操作错误:" + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(menu); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Insert(&menu); err == nil {
			c.Succeed(&controllers.ResultJson{
				Message: "创建节点成功",
				Url:     beego.URLFor("Menu.Index", ":id", root),
			})
		} else {
			c.Fail(&controllers.ResultJson{
				Message: "创建节点失败，请稍后再试！",
			})
		}
	}
	c.Data["root"] = root
	c.Data["pid"], _ = c.GetInt(":id", 0)
	c.Data["groupMenu"] = c.menusService.PidAllMenus(0)
	c.Data["urlsType"] = c.menusService.UrlsType()
}

// @router /menu/delete [post]
func (c *Menu) Delete() {
	id := c.GetMustInt(":ids", "非法操作，请稍后再试...")
	root, _ := c.GetInt(":root", 1)
	if extend, _ := orm.NewOrm().QueryTable(new(menus.Menus)).Filter("pid", id).Count(); extend > 0 {
		c.Fail(&controllers.ResultJson{
			Message: "该节点存在子级；请先删除所有子节点！",
		})
	}
	if num, _ := orm.NewOrm().Delete(&menus.Menus{Id: id}); num > 0 {
		c.Succeed(&controllers.ResultJson{
			Message: "节点已经删除！马上返回中。。。",
			Url:     beego.URLFor("Menu.Index", ":id", root),
		})
	} else {
		c.Fail(&controllers.ResultJson{Message: "删除节点失败，请稍后再试！"})
	}
}

// @router /menu/status [post]
func (c *Menu) Status() {
	root, _ := c.GetInt(":root", 1)
	ids := c.GetMustInt(":ids", "非法操作，请稍后再试...")
	menu, _ := c.menusService.Find(ids)
	menu.Status, _ = c.GetInt8(":status", 0)
	_, err := orm.NewOrm().Update(&menu)
	if err == nil {
		c.Succeed(&controllers.ResultJson{
			Message: "状态更新成功！马上返回中。。。",
			Url:     beego.URLFor("Menu.Index", ":id", root),
		})
	} else {
		c.Fail(&controllers.ResultJson{Message: "状态更新失败，请稍后再试！"})
	}
}

// @router /menu/tree [post]
func (c *Menu) Tree() {
	id := c.GetMustInt("module", "非法请求,请稍后再试...")
	c.Succeed(&controllers.ResultJson{
		Data:    c.menusService.TreeRender(id),
		Message: "success!!",
	})
}

// @router /menu/sort [get,post]
func (c *Menu) Sort() {
	var sortMenus []*menus.Menus
	root, _ := c.GetInt(":root", 1)
	_ = json.Unmarshal([]byte(c.GetString(":menus")), &sortMenus)
	for _, value := range sortMenus {
		menuOne, _ := c.menusService.Find(value.Id)
		menuOne.Sort = value.Sort
		menuOne.Pid = value.Pid
		_, _ = orm.NewOrm().Update(&menuOne)
	}
	c.Succeed(&controllers.ResultJson{
		Message: "排序更新成功！马上返回中。。。",
		Url:     beego.URLFor("Menu.Index", ":id", root),
	})
}
