package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/go-playground/validator"
)

/** 蜘蛛池内容模型 **/
type Models struct {
	Main
}

// @router /models/index [get,post]
func (c *Models) Index() {
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.modelsService.PageListItems(length, draw, c.PageNumber(), search)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.modelsService.TableColumnsType(), c.modelsService.TableButtonsType())
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.modelsService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.modelsService.DataTableButtons() {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
}

// @router /models/create [get,post]
func (c *Models) Create() {
	if c.IsAjax() {
		item := spider.Models{}
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Insert(&item); err == nil {
			_ = c.disguiseService.Inc(item.Disguise, 0) // 挂载自然语言使用次数
			c.Succeed(&controllers.ResultJson{Message: "创建模型成功", Url: beego.URLFor("Models.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "创建模型失败，请稍后再试！"})
		}
	}
	c.Data["collect"] = c.modelsService.CollectMapColumn()
	c.Data["disguise"] = c.disguiseService.Groups()
	c.Data["theme"] = c.templateService.Groups()
}

// @router /models/edit [get,post]
func (c *Models) Edit() {
	id := c.GetMustInt(":id", "非法请求！")
	if c.IsAjax() {
		item := c.modelsService.One(id)
		disguise := item.Disguise
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Update(&item); err == nil {
			_ = c.disguiseService.Inc(item.Disguise, disguise)
			c.Succeed(&controllers.ResultJson{Message: "更新模型成功", Url: beego.URLFor("Models.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "更新模型失败，请稍后再试！"})
		}
	}
	c.Data["collect"] = c.modelsService.CollectMapColumn()
	c.Data["disguise"] = c.disguiseService.Groups()
	c.Data["theme"] = c.templateService.Groups()
	c.Data["info"] = c.modelsService.One(id)
}

// @router /models/delete [post]
func (c *Models) Delete() {
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.modelsService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除成功！马上返回中。。。",
			Url:     beego.URLFor("Disguise.Index"),
		})
	}
}
