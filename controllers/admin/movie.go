package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/service"
)

/** 剧情采集 **/
type Movie struct {
	Main
	configPrefix map[string]string
}

/** 构造函数  **/
func (c *Movie) NextPrepare() {
	c.configPrefix = map[string]string{
		"interval":    "最大并发",
		"push_time":   "发布间隔",
		"send_domain": "推送接口",
	}
}

// @router /movie/config [get,post]
func (c *Movie) Config() {
	if c.IsAjax() {
		var massage error
		_ = orm.NewOrm().Begin()
		for k, v := range c.configPrefix {
			value := c.GetMustString(k, v+"不能为空！")
			field := c.movieService.ConfigByName(k)
			field.Value = value
			if _, massage = orm.NewOrm().Update(&field); massage != nil {
				_ = orm.NewOrm().Rollback()
				break
			}
		}
		if massage == nil {
			_ = orm.NewOrm().Commit()
			c.Succeed(&controllers.ResultJson{Message: "修改配置成功，重启爬虫生效！", Url: beego.URLFor("Movie.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "修改失败，错误原因: " + error.Error(massage)})
		}
	}
	c.Data["config"] = c.movieService.ConfigMaps()
}

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
