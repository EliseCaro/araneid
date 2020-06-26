package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
)

/** 蜘蛛池资源文档库管理 **/
type Article struct{ Main }

// @router /article/index [get,post]
func (c *Article) Index() {
	model, _ := c.GetInt(":model", 0)
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.articleService.PageListItems(length, draw, c.PageNumber(), search, model)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.articleService.TableColumnsType(), c.articleService.TableButtonsType(model))
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.articleService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.articleService.DataTableButtons(model) {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
	c.Data["action"] = beego.URLFor("Article.Index", ":model", model)
}
