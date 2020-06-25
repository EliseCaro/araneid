package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/go-playground/validator"
)

/** 自然语言处理 **/
type Disguise struct {
	Main
}

// @router /disguise/index [get,post]
func (c *Disguise) Index() {
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.disguiseService.PageListItems(length, draw, c.PageNumber(), search)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.disguiseService.TableColumnsType(), c.disguiseService.TableButtonsType())
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.disguiseService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.disguiseService.DataTableButtons() {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
}

// @router /disguise/create [get,post]
func (c *Disguise) Create() {
	if c.IsAjax() {
		item := spider.Disguise{}
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Insert(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "创建机器成功", Url: beego.URLFor("Disguise.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "创建机器失败，请稍后再试！"})
		}
	}
}

// @router /disguise/edit [get,post]
func (c *Disguise) Edit() {
	id := c.GetMustInt(":id", "非法请求！")
	if c.IsAjax() {
		item, _ := c.disguiseService.Find(id)
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Update(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "更新机器成功", Url: beego.URLFor("Disguise.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "更新机器失败，请稍后再试！"})
		}
	}
	c.Data["info"], _ = c.disguiseService.Find(id)
}

// @router /disguise/delete [post]
func (c *Disguise) Delete() {
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.disguiseService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除成功！马上返回中。。。",
			Url:     beego.URLFor("Disguise.Index"),
		})
	}
}

// @router /disguise/status [post]
func (c *Disguise) Status() {
	status, _ := c.GetInt8(":status", 0)
	array := c.checkBoxIds(":ids[]", ":ids")
	field := c.GetString(":field", "status")
	if errorMessage := c.disguiseService.UpdateStatus(array[0], status, field); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: "状态更新失败，请稍后再试！",
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "状态更新成功！下次使用生效；马上返回中。。。",
			Url:     beego.URLFor("Disguise.Index"),
		})
	}
}
