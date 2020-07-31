package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/service"
)

/** 剧情采集 **/
type Movie struct{ Main }

// @router /movie/index [get,post]
func (c *Movie) Index() {
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.movieService.PageListItems(length, draw, c.PageNumber(), search)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.movieService.TableColumnsType(), c.movieService.TableButtonsType())
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.movieService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.movieService.DataTableButtons() {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
}

/** 爬虫启动与停止 **/
// @router /movie/status [post]
func (c *Movie) Status() {
	status, _ := c.GetInt("status", 0)
	field := c.GetMustString("field", "非法操作～")
	service.SocketInstanceGet().InstructHandle(map[string]interface{}{
		"status": status, "command": "movie",
		"uid": c.UserInfo.Id, "field": field,
	})
	c.Succeed(&controllers.ResultJson{Message: "指令状态已经更改！", Url: beego.URLFor("Movie.Index")})
}
