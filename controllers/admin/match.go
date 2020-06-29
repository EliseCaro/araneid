package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/go-playground/validator"
)

/** 蜘蛛池KTDB调用模板管理 **/
type Match struct{ Main }

// @router /match/index [get,post]
func (c *Match) Index() {
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.matchService.PageListItems(length, draw, c.PageNumber(), search)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.matchService.TableColumnsType(), c.matchService.TableButtonsType())
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.matchService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.matchService.DataTableButtons() {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())

}

// @router /match/create [get,post]
func (c *Match) Create() {
	if c.IsAjax() {
		item := spider.Match{}
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Insert(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "创建模板成功", Url: beego.URLFor("Match.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "创建模板失败，请稍后再试！"})
		}
	}
}

// @router /match/edit [get,post]
func (c *Match) Edit() {
	id := c.GetMustInt(":id", "非法请求！")
	if c.IsAjax() {
		item, _ := c.matchService.Find(id)
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Update(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "更新模板成功;", Url: beego.URLFor("Match.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "更新模板失败，请稍后再试！"})
		}
	}
	c.Data["info"], _ = c.matchService.Find(id)
}

// @router /match/delete [post]
func (c *Match) Delete() {
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.matchService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{Message: error.Error(errorMessage)})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除成功！马上返回中。。。",
			Url:     beego.URLFor("Match.Index"),
		})
	}
}
