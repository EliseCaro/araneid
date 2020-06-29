package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/go-playground/validator"
)

/** 蜘蛛池管理 **/
type Arachnid struct{ Main }

// @router /arachnid/index [get,post]
func (c *Arachnid) Index() {
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.arachnidService.PageListItems(length, draw, c.PageNumber(), search)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.arachnidService.TableColumnsType(), c.arachnidService.TableButtonsType())
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.arachnidService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.arachnidService.DataTableButtons() {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
}

// @router /arachnid/create [get,post]
func (c *Arachnid) Create() {
	if c.IsAjax() {
		item := spider.Arachnid{}
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Insert(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "创建项目成功", Url: beego.URLFor("Arachnid.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "创建项目失败，请稍后再试！"})
		}
	}
	c.Data["model"] = c.arachnidService.ModelSelect() // 获取挂载模型记录
	c.Data["match"] = c.arachnidService.MatchSelect() // 获取挂载匹配模板记录
}

// @router /arachnid/edit [get,post]
func (c *Arachnid) Edit() {
	id := c.GetMustInt(":id", "非法请求！")
	if c.IsAjax() {
		item, _ := c.arachnidService.Find(id)
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Update(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "更新项目成功;", Url: beego.URLFor("Arachnid.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "更新项目失败，请稍后再试！"})
		}
	}
	c.Data["info"], _ = c.arachnidService.Find(id)
	c.Data["model"] = c.arachnidService.ModelSelect() // 获取挂载模型记录
	c.Data["match"] = c.arachnidService.MatchSelect() // 获取挂载匹配模板记录
}

// @router /arachnid/status [post]
func (c *Arachnid) Status() {
	status, _ := c.GetInt8(":status", 0)
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.arachnidService.StatusArray(array, status); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "状态更新成功！马上返回中。。。",
			Url:     beego.URLFor("Arachnid.Index"),
		})
	}
}

// todo 待后边来做；可能涉及删除的地方很多
// @router /arachnid/delete [post]
func (c *Arachnid) Delete() {

}
