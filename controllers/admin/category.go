package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/go-playground/validator"
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
		result := c.categoryService.PageListItems(length, draw, c.PageNumber(), search, parent)
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

// @router /category/edit [get,post]
func (c *Category) Edit() {
	id := c.GetMustInt(":id", "非法请求！")
	cate, _ := c.categoryService.Find(id)
	if c.IsAjax() {
		if err := c.ParseForm(&cate); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(cate); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Update(&cate); err == nil {
			_ = c.categoryService.InitializeDomainCate(cate.Domain) // 将域名数据中旧分类更新至新分类
			c.Succeed(&controllers.ResultJson{Message: "更新分类成功", Url: beego.URLFor("Category.Index", ":parent", cate.Domain)})
		} else {
			c.Fail(&controllers.ResultJson{Message: "更新分类失败，请稍后再试！"})
		}
	}
	c.Data["info"] = cate
}

// @router /category/delete [post]
func (c *Category) Delete() {
	parent, _ := c.GetInt(":parent", 0)
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.categoryService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{Message: error.Error(errorMessage)})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除成功！马上返回中。。。",
			Url:     beego.URLFor("Category.Index", ":parent", parent),
		})
	}
}
