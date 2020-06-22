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
		"interval":  "最大并发",
		"push_time": "发布间隔",
		"translate": "语言转换",
		"download":  "资源下载",
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
			c.Succeed(&controllers.ResultJson{Message: "修改配置成功，重启爬虫生效！", Url: beego.URLFor("Dictionaries.Cate")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "修改失败，错误原因: " + error.Error(massage)})
		}
	}
	c.Data["config"] = c.dictionariesService.ConfigMaps()
}

/** 爬虫启动与停止 **/
// @router /dictionaries/status_cate [post]
func (c *Dictionaries) StatusCate() {
	status, _ := c.GetInt("status", 0)
	field := c.GetMustString("field", "非法操作～")
	service.SocketInstanceGet().InstructHandle(map[string]interface{}{
		"status": status, "command": "dict_cate",
		"uid": c.UserInfo.Id, "field": field,
	})
	c.Succeed(&controllers.ResultJson{Message: "指令状态已经更改！", Url: beego.URLFor("Dictionaries.Cate")})
}

/** 爬虫启动与停止 **/
// @router /dictionaries/delete_cate [post]
func (c *Dictionaries) DeleteCate() {
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.dictionariesService.DeleteArrayCate(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除结果成功！马上返回中。。。",
			Url:     beego.URLFor("Dictionaries.Cate"),
		})
	}
}

/** 查看结果  **/
// @router /dictionaries/detail_cate [get]
func (c *Dictionaries) DetailCate() {
	id := c.GetMustInt(":id", "结果ID不合法...")
	c.Data["info"] = c.dictionariesService.DetailCateOne(id)
}

/** 发布分类结果  **/
// @router /dictionaries/push_cate [post]
func (c *Dictionaries) PushCate() {
	array := c.checkBoxIds(":ids[]", ":ids")
	for _, v := range array {
		one := c.dictionariesService.OneCate(v)
		c.dictionariesService.PushDetailAPICate(one)
	}
	c.Succeed(&controllers.ResultJson{Message: "提交发布成功！马上返回中。。。", Url: beego.URLFor("Dictionaries.Cate")})

}

// 采集分类
// @router /dictionaries/cate [get,post]
func (c *Dictionaries) Cate() {
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.dictionariesService.PageListItemsCate(length, draw, c.PageNumber(), search)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.dictionariesService.TableColumnsTypeCate(), c.dictionariesService.TableButtonsTypeCate())
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.dictionariesService.DataTableCateColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.dictionariesService.DataTableCateButtons() {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
}
