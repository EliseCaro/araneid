package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/roles"
	"github.com/go-playground/validator"
)

type Roles struct {
	Main
	tableIsCheck bool
}

/** 构造函数  **/
func (c *Roles) NextPrepare() {
	c.tableIsCheck = true
}

// @router /roles/index [get,post]
func (c *Roles) Index() {
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.rolesService.PageListItems(length, draw, c.PageNumber(), search)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.rolesService.TableColumnsType(), c.rolesService.TableButtonsType())
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.rolesService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.rolesService.DataTableButtons() {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, 10)
}

// @router /roles/status [post]
func (c *Roles) Status() {
	status, _ := c.GetInt8(":status", 0)
	array := c.checkBoxIds(":ids[]", ":ids")
	field := c.GetString(":field", "status")
	if errorMessage := c.rolesService.UpdateAccessAndStatus(array, status, field); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: "批量状态更新失败，请稍后再试！",
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "状态更新成功！马上返回中。。。",
			Url:     beego.URLFor("Roles.Index"),
		})
	}
}

// @router /roles/create [get,post]
func (c *Roles) Create() {
	if c.IsAjax() {
		item := roles.Roles{}
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Insert(&item); err == nil {
			c.Succeed(&controllers.ResultJson{
				Message: "创建角色成功",
				Url:     beego.URLFor("Roles.Index"),
			})
		} else {
			c.Fail(&controllers.ResultJson{
				Message: "创建角色失败，请稍后再试！",
			})
		}
	}
	c.Data["pid"] = c.rolesService.ToLayerRole(c.rolesService.RoleAll(), 0, 0)
	c.Data["group"] = c.rolesService.BuildJsTree(c.UserInfo, 0)
}

// @router /roles/edit [get,post]
func (c *Roles) Edit() {
	id := c.GetMustInt(":id", "非法访问；已被拒绝！")
	if c.IsAjax() {
		item := roles.Roles{Id: id}
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Update(&item); err == nil {
			c.Succeed(&controllers.ResultJson{
				Message: "修改角色成功",
				Url:     beego.URLFor("Roles.Index"),
			})
		} else {
			c.Fail(&controllers.ResultJson{
				Message: "修改角色失败，请稍后再试！",
			})
		}
	}
	c.Data["pid"] = c.rolesService.ToLayerRole(c.rolesService.RoleAll(), 0, 0)
	c.Data["group"] = c.rolesService.BuildJsTree(c.UserInfo, id)
	c.Data["info"], _ = c.rolesService.Find(id)
}

/** 删除节点 **/
// @router /roles/delete [post]
func (c *Roles) Delete() {
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.rolesService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: "删除角色失败，请稍后再试！",
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除角色成功！马上返回中。。。",
			Url:     beego.URLFor("Roles.Index"),
		})
	}
}
