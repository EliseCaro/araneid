package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/go-playground/validator"
)

/** 蜘蛛池缓存池文档缓存管理 **/
type Detail struct{ Main }

// @router /detail/index [get,post]
func (c *Detail) Index() {
	parent, _ := c.GetInt(":parent", 0)
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.detailService.PageListItems(length, draw, c.PageNumber(), search, parent)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.detailService.TableColumnsType(), c.detailService.TableButtonsType(parent))
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	one, _ := c.domainService.Find(parent)
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.detailService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.detailService.DataTableButtons(parent, one.Arachnid) {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
	c.Data["info"] = one
	c.Data["action"] = beego.URLFor("Detail.Index", ":parent", parent)
}

// @router /detail/delete [post]
func (c *Detail) Delete() {
	parent, _ := c.GetInt(":parent", 0)
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.detailService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{Message: error.Error(errorMessage)})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除成功！马上返回中。。。",
			Url:     beego.URLFor("Detail.Index", ":parent", parent),
		})
	}
}

// @router /detail/empty [post]
func (c *Detail) Empty() {
	parent, _ := c.GetInt(":parent", 0)
	c.detailService.EmptyDelete(parent)
	c.Succeed(&controllers.ResultJson{
		Message: "缓存已经清空！将重新生成文章！",
		Url:     beego.URLFor("Detail.Index", ":parent", parent),
	})
}

// @router /detail/edit [get,post]
func (c *Detail) Edit() {
	id := c.GetMustInt(":id", "非法请求！")
	parent, _ := c.GetInt(":parent", 0)
	detail, _ := c.detailService.Find(id)
	if c.IsAjax() {
		if err := c.ParseForm(&detail); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(detail); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Update(&detail); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "更新文章成功", Url: beego.URLFor("Detail.Index", ":parent", parent)})
		} else {
			c.Fail(&controllers.ResultJson{Message: "更新文章失败，请稍后再试！"})
		}
	}
	c.Data["info"] = detail
	c.Data["parent"] = parent
}
