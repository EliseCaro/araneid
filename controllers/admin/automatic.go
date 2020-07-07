package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
)

/** 自动提交链接到百度的日志 **/
type Automatic struct{ Main }

// @router /automatic/index [get,post]
func (c *Automatic) Index() {
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.automaticService.PageListItems(length, draw, c.PageNumber(), search)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.automaticService.TableColumnsType(), c.automaticService.TableButtonsType())
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.automaticService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.automaticService.DataTableButtons() {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
}

// @router /automatic/status [post]
func (c *Automatic) Status() {
	array := c.checkBoxIds(":ids[]", ":ids")
	c.automaticService.StatusArray(array)
	c.Succeed(&controllers.ResultJson{
		Message: "已经处理完成！马上返回中。。。",
		Url:     beego.URLFor("Automatic.Index"),
	})
}

// @router /automatic/delete [post]
func (c *Automatic) Delete() {
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.automaticService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{Message: error.Error(errorMessage)})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除成功！马上返回中。。。",
			Url:     beego.URLFor("Automatic.Index"),
		})
	}
}

// @router /automatic/empty [post]
func (c *Automatic) Empty() {
	c.automaticService.EmptyDelete()
	c.Succeed(&controllers.ResultJson{
		Message: "推送记录已经清空！马上返回中。。。",
		Url:     beego.URLFor("Automatic.Index"),
	})
}
