package admin

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/go-playground/validator"
)

/** 蜘蛛池缓存池域名管理 **/
type Domain struct{ Main }

// @router /domain/index [get,post]
func (c *Domain) Index() {
	arachnid, _ := c.GetInt(":arachnid", 0)
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.domainService.PageListItems(length, draw, c.PageNumber(), search, arachnid)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.domainService.TableColumnsType(arachnid), c.domainService.TableButtonsType(arachnid))
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.domainService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.domainService.DataTableButtons(arachnid) {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
	c.Data["info"], _ = c.arachnidService.Find(arachnid)
	c.Data["action"] = beego.URLFor("Domain.Index", ":arachnid", arachnid)
}

// @router /domain/status [post]
func (c *Domain) Status() {
	parent, _ := c.GetInt(":parent", 0)
	status, _ := c.GetInt8(":status", 0)
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.domainService.StatusArray(array, status); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "站点状态更新成功！马上返回中。。。",
			Url:     beego.URLFor("Domain.Index", ":arachnid", parent),
		})
	}
}

// @router /domain/delete [post]
func (c *Domain) Delete() {
	parent, _ := c.GetInt(":parent", 0)
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.domainService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{Message: error.Error(errorMessage)})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除成功！马上返回中。。。",
			Url:     beego.URLFor("Domain.Index", ":id", parent),
		})
	}
}

// @router /domain/edit [get,post]
func (c *Domain) Edit() {
	id := c.GetMustInt(":id", "非法请求！")
	domain, _ := c.domainService.Find(id)
	arachnid, _ := c.arachnidService.Find(domain.Arachnid)
	if c.IsAjax() {
		if err := c.ParseForm(&domain); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(domain); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Update(&domain); err == nil {
			c.automaticService.AutomaticClass(domain.Id)
			c.Succeed(&controllers.ResultJson{Message: "更新站点成功", Url: beego.URLFor("Domain.Index", ":arachnid", domain.Arachnid)})
		} else {
			c.Fail(&controllers.ResultJson{Message: "更新站点失败，请稍后再试！"})
		}
	}
	c.Data["info"] = domain
	c.Data["templates"] = c.templateService.Items(arachnid.Models)
}

// @router /domain/empty [post]
func (c *Domain) Empty() {
	parent, _ := c.GetInt(":parent", 0)
	c.domainService.EmptyDelete(parent)
	c.Succeed(&controllers.ResultJson{
		Message: "缓存已经清空！将重制此项目所有站点数据～",
		Url:     beego.URLFor("Domain.Index", ":parent", parent),
	})
}

// @router /domain/links [post,get]
func (c *Domain) Links() {
	parent, _ := c.GetInt(":parent", 0)
	id := c.GetMustInt(":id", "非法请求！")
	if c.IsAjax() {
		value := c.GetString("inputs", "[]")
		if _, e := orm.NewOrm().Update(&spider.Domain{Id: id, Links: value}, "Links"); e == nil {
			c.Succeed(&controllers.ResultJson{
				Message: "更新友情链接成功", Url: beego.URLFor("Domain.Index", ":arachnid", parent),
			})
		} else {
			c.Fail(&controllers.ResultJson{Message: "更新友情链接失败，请稍后再试！"})
		}
	}
	domain, _ := c.domainService.Find(id)
	var links []map[string]string
	_ = json.Unmarshal([]byte(domain.Links), &links)
	c.Data["links"] = links
	c.Data["parent"] = parent
	c.Data["id"] = domain.Id
}
