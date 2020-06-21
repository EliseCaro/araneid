package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
)

type Collector struct {
	Main
}

// @router /collector/index [get,post]
func (c *Collector) Index() {
	id := c.GetMustInt(":id", "非法请求！")
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.collectorService.PageListItems(length, draw, c.PageNumber(), search, id)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.collectorService.TableColumnsType(), c.collectorService.TableButtonsType(id))
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.collectorService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.collectorService.DataTableButtons(id) {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
	c.Data["info"] = c.collectService.One(id)
}

// @router /collector/delete [post]
func (c *Collector) Delete() {
	parent := c.GetMustInt(":parent", "爬虫ID不合法...")
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.collectorService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除结果成功！马上返回中。。。",
			Url:     beego.URLFor("Collector.Index", ":id", parent),
		})
	}
}

// @router /collector/detail [get]
func (c *Collector) Detail() {
	id := c.GetMustInt(":id", "结果ID不合法...")
	c.Data["info"] = c.collectorService.One(id)
}

// @router /collector/status [post]
func (c *Collector) Status() {
	array := c.checkBoxIds(":ids[]", ":ids")
	parent := c.GetMustInt(":parent", "非法操作；已被拦截")
	for _, v := range array {
		one := c.collectorService.One(v)
		c.collectService.PushDetailAPI(one)
	}
	c.Succeed(&controllers.ResultJson{Message: "提交发布成功！马上返回中。。。", Url: beego.URLFor("Collector.Index", ":id", parent)})
}
