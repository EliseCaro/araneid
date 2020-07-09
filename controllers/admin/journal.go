package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"strings"
)

/** 蜘蛛池资源文档库管理 **/
type Journal struct{ Main }

// @router /journal/index [get,post]
func (c *Journal) Index() {
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		var searchMap map[string]string
		if search := c.GetString("search[value]", ""); search != "" {
			searchMap = c.searchMap(search)
		}
		if searchOne := c.GetString(":search", ""); searchOne != "" {
			searchMap = c.searchMap(searchOne)
		}
		result := c.journalService.PageListItems(length, draw, c.PageNumber(), searchMap)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.journalService.TableColumnsType(), c.journalService.TableButtonsType())
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.journalService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.journalService.DataTableButtons() {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
}

/** 解析条件 **/
func (c *Journal) searchMap(str string) map[string]string {
	str = strings.ReplaceAll(str, "|", "\n")
	maps := _func.ParseAttrConfigMap(str)
	return maps
}

// @router /journal/delete [post]
func (c *Journal) Delete() {
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.journalService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{Message: error.Error(errorMessage)})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除成功！马上返回中。。。",
			Url:     beego.URLFor("Journal.Index"),
		})
	}
}

// @router /journal/empty [post]
func (c *Journal) Empty() {
	c.journalService.EmptyDelete()
	c.Succeed(&controllers.ResultJson{
		Message: "蜘蛛记录已经清空！马上返回中。。。",
		Url:     beego.URLFor("Journal.Index"),
	})
}
