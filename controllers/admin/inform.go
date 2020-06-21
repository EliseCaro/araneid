package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
)

type Inform struct{ Main }

// @router /inform/index [get,post]
func (c *Inform) Index() {
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		statue, _ := c.GetInt8(":statue", -1)
		result := c.informService.PageListItems(length, draw, c.PageNumber(), search, statue)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.informService.TableColumnsType(), c.informService.TableButtonsType())
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.informService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.informService.DataTableButtons() {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
	c.Data["url"] = c.Ctx.Request.URL.String()
}

// @router /inform/status [get,post]
func (c *Inform) Status() {
	status, _ := c.GetInt8(":status", 0)
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.informService.StatusArray(array, status); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "已全部标记～。。。",
			Url:     beego.URLFor("Inform.Index"),
		})
	}
}

// @router /inform/delete [post]
func (c *Inform) Delete() {
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.informService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除消息成功！马上返回中。。。",
			Url:     beego.URLFor("Inform.Index"),
		})
	}
}

// @router /inform/check [get]
func (c *Inform) Check() {
	id := c.GetMustInt(":id", "非法请求！")
	info := c.informService.DetailInform(id)
	c.Data["info"] = info
	_ = c.informService.StatusArray([]int{id}, 1)
}
