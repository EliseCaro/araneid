package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/service"
)

/** 查词典采集器 **/
type Dictionaries struct {
	Main
	configPrefix map[string]string
}

/** 构造函数  **/
func (c *Dictionaries) NextPrepare() {
	c.configPrefix = map[string]string{
		"interval":    "最大并发",
		"push_time":   "发布间隔",
		"translate":   "语言转换",
		"download":    "资源下载",
		"send_domain": "推送接口",
	}
}

// @router /dictionaries/config [get,post]
func (c *Dictionaries) Config() {
	if c.IsAjax() {
		var massage error
		_ = orm.NewOrm().Begin()
		for k, v := range c.configPrefix {
			value := c.GetMustString(k, v+"不能为空！")
			field := c.dictionariesService.ConfigByName(k)
			field.Value = value
			if _, massage = orm.NewOrm().Update(&field); massage != nil {
				_ = orm.NewOrm().Rollback()
				break
			}
		}
		if massage == nil {
			_ = orm.NewOrm().Commit()
			c.Succeed(&controllers.ResultJson{Message: "修改配置成功，重启爬虫生效！", Url: beego.URLFor("Dictionaries.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "修改失败，错误原因: " + error.Error(massage)})
		}
	}
	c.Data["config"] = c.dictionariesService.ConfigMaps()
}

// @router /dictionaries/index [get,post]
func (c *Dictionaries) Index() {
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.dictionariesService.PageListItems(length, draw, c.PageNumber(), search)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.dictionariesService.TableColumnsType(), c.dictionariesService.TableButtonsType())
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.dictionariesService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.dictionariesService.DataTableButtons() {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
}

/** 查看结果  **/
// @router /dictionaries/detail [get]
func (c *Dictionaries) Detail() {
	id := c.GetMustInt(":id", "结果ID不合法...")
	c.Data["info"] = c.dictionariesService.DetailOne(id)
}

/** 爬虫启动与停止 **/
// @router /dictionaries/delete [post]
func (c *Dictionaries) Delete() {
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.dictionariesService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除结果成功！马上返回中。。。",
			Url:     beego.URLFor("Dictionaries.Index"),
		})
	}
}

/** 爬虫启动与停止 **/
// @router /dictionaries/status [post]
func (c *Dictionaries) Status() {
	status, _ := c.GetInt("status", 0)
	field := c.GetMustString("field", "非法操作～")
	service.SocketInstanceGet().InstructHandle(map[string]interface{}{
		"status": status, "command": "dict",
		"uid": c.UserInfo.Id, "field": field,
	})
	c.Succeed(&controllers.ResultJson{Message: "指令状态已经更改！", Url: beego.URLFor("Dictionaries.Index")})
}

/** 发布结果  **/
// @router /dictionaries/push [post]
func (c *Dictionaries) Push() {
	array := c.checkBoxIds(":ids[]", ":ids")
	for _, v := range array {
		one := c.dictionariesService.One(v)
		c.dictionariesService.PushDetailAPI(one)
	}
	c.Succeed(&controllers.ResultJson{Message: "提交发布成功！马上返回中。。。", Url: beego.URLFor("Dictionaries.Index")})
}
