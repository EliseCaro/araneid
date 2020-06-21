package admin

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/collect"
	"github.com/beatrice950201/araneid/extend/service"
	"github.com/go-playground/validator"
)

type Collect struct {
	Main
}

// @router /collect/index [get,post]
func (c *Collect) Index() {
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.collectService.PageListItems(length, draw, c.PageNumber(), search)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.collectService.TableColumnsType(), c.collectService.TableButtonsType())
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.collectService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.collectService.DataTableButtons() {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
}

// @router /collect/create [get,post]
func (c *Collect) Create() {
	if c.IsAjax() {
		item := c.parseFormVerify(collect.Collect{})
		if _, err := orm.NewOrm().Insert(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "创建爬虫成功", Url: beego.URLFor("Collect.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "创建爬虫失败，请稍后再试！"})
		}
	}
}

// @router /collect/edit [get,post]
func (c *Collect) Edit() {
	id := c.GetMustInt(":id", "非法请求！")
	if c.IsAjax() {
		item := c.parseFormVerify(c.collectService.One(id))
		item.Status = 0
		item.PushStatus = 0
		if _, err := orm.NewOrm().Update(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "更新爬虫成功", Url: beego.URLFor("Collect.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "更新爬虫失败，请稍后再试！"})
		}
	}
	c.Data["info"] = c.collectService.FindOneDetail(id)
}

// @router /collect/delete [post]
func (c *Collect) Delete() {
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.collectService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除爬虫成功！马上返回中。。。",
			Url:     beego.URLFor("Collect.Index"),
		})
	}
}

// @router /collect/status [post]
func (c *Collect) Status() {
	var (
		status, _ = c.GetInt8(":status", 0)
		field     = c.GetMustString(":field", "非法操作～")
		array     = c.checkBoxIds(":ids[]", ":ids")
	)
	service.SocketInstanceGet().InstructHandle(map[string]interface{}{
		"status": status, "field": field,
		"command": "collect", "id": array[0], "uid": c.UserInfo.Id,
	})
	c.Succeed(&controllers.ResultJson{Message: c.collectService.InstructTypeMessage(field, status)})
}

// @router /collect/test [post]
func (c *Collect) Test() {
	t := c.GetMustString(":type", "操作类型不合法！")
	switch t {
	case "rule":
		c.testRule()
	case "filed":
		c.testField()
	default:
		c.Fail(&controllers.ResultJson{Message: "测试类型不合法！"})
	}
}

/****************************私有方法**********************************/

/** 公用表单解析验证 **/
func (c *Collect) parseFormVerify(item collect.Collect) collect.Collect {
	if err := c.ParseForm(&item); err != nil {
		c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
	}
	if message := c.verifyBase.Begin().Struct(item); message != nil {
		c.Fail(&controllers.ResultJson{Message: c.verifyBase.Translate(message.(validator.ValidationErrors))})
	}
	var matching []*collect.Matching
	_ = json.Unmarshal([]byte(c.GetString("matching")), &matching)
	for _, v := range matching {
		if message := c.verifyBase.Begin().Struct(v); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
			break
		}
	}
	marshalJson, _ := json.Marshal(matching)
	item.Matching = string(marshalJson)
	return item
}

/** 采集规则检测 **/
func (c *Collect) testRule() {
	source := _func.ParseAttrConfigArray(c.GetMustString("source", "采集连接池不能为空！"))
	roles := c.GetMustString("source_rule", "采集规则不能为空！")
	result := c.collectService.ExtractUrls(roles, source)
	c.Succeed(&controllers.ResultJson{Message: "success!", Data: result})
}

/**  字段值检测 **/
func (c *Collect) testField() {
	matchingJson := c.GetMustString("matching", "检索字段至少有一套！")
	source := _func.ParseAttrConfigArray(c.GetMustString("source", "采集连接池不能为空！"))
	roles := c.GetMustString("source_rule", "采集规则不能为空！")
	var matching []*collect.Matching
	_ = json.Unmarshal([]byte(matchingJson), &matching)
	for _, v := range matching {
		if message := c.verifyBase.Begin().Struct(v); message != nil {
			c.Fail(&controllers.ResultJson{Message: c.verifyBase.Translate(message.(validator.ValidationErrors))})
			break
		}
	}
	urls := c.collectService.ExtractUrls(roles, source)
	if len(urls) > 0 {
		if result, message := c.collectService.ExtractDocumentMatching(urls[0]["href"], matching); message == nil {
			c.Succeed(&controllers.ResultJson{Message: "success!", Data: result})
		} else {
			c.Fail(&controllers.ResultJson{Message: error.Error(message)})
		}
	} else {
		c.Fail(&controllers.ResultJson{Message: "未检索到任何匹配URL!"})
	}
}
