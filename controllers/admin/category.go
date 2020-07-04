package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
)

/** 蜘蛛池缓存池分类管理 **/
type Category struct{ Main }

// @router /category/index [get,post]
func (c *Category) Index() {
	parent, _ := c.GetInt(":parent", 0)
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.domainService.PageListItems(length, draw, c.PageNumber(), search, parent)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.categoryService.TableColumnsType(), c.categoryService.TableButtonsType(parent))
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	one, _ := c.domainService.Find(parent)
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.categoryService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.categoryService.DataTableButtons(parent, one.Arachnid) {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
	c.Data["info"] = one
	c.Data["action"] = beego.URLFor("Category.Index", ":parent", parent)
}
