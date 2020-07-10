package admin

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/index"
	"github.com/go-playground/validator"
)

/** 功能与优势 **/
type Advantage struct{ Main }

// @router /advantage/index [get,post]
func (c *Advantage) Index() {
	if c.IsAjax() {
		length, _ := c.GetInt("length", table.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.PageListItems(length, draw, c.PageNumber(), search)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.TableColumnsType(), c.TableButtonsType())
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.DataTableButtons() {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, table.WebPageSize())
}

// @router /advantage/create [get,post]
func (c *Advantage) Create() {
	if c.IsAjax() {
		item := index.Advantage{}
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Insert(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "添加成功", Url: beego.URLFor("Advantage.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "添加失败，请稍后再试！"})
		}
	}
}

// @router /advantage/edit [get,post]
func (c *Advantage) Edit() {
	id := c.GetMustInt(":id", "非法请求！")
	if c.IsAjax() {
		item, _ := c.Find(id)
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Update(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "更新成功", Url: beego.URLFor("Advantage.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "更新失败，请稍后再试！"})
		}
	}
	c.Data["info"], _ = c.Find(id)
}

// @router /advantage/status [post]
func (c *Advantage) Status() {
	status, _ := c.GetInt8(":status", 0)
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.StatusArray(array, status); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "状态更新成功！马上返回中。。。",
			Url:     beego.URLFor("Advantage.Index"),
		})
	}
}

// @router /advantage/delete [post]
func (c *Advantage) Delete() {
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除成功！马上返回中。。。",
			Url:     beego.URLFor("Advantage.Index"),
		})
	}
}

/************************************************表格渲染机制 ************************************************************/

/** 批量删除 **/
func (c *Advantage) DeleteArray(array []int) (e error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, e = orm.NewOrm().Delete(&index.Advantage{Id: v}); e != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if e == nil {
		_ = orm.NewOrm().Commit()
	}
	return e
}

/** 更新状态 **/
func (c *Advantage) StatusArray(array []int, status int8) (e error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, e = orm.NewOrm().Update(&index.Advantage{Id: v, Status: status}, "Status"); e != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if e == nil {
		_ = orm.NewOrm().Commit()
	}
	return e
}

/** 获取一条数据 **/
func (c *Advantage) Find(id int) (index.Advantage, error) {
	item := index.Advantage{
		Id: id,
	}
	return item, orm.NewOrm().Read(&item)
}

/** 获取需要渲染的Column **/
func (c *Advantage) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "ID", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "标题", "name": "title", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "描述", "name": "description", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "状态", "name": "status", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "更新时间", "name": "update_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "创建时间", "name": "create_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "爬虫操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (c *Advantage) DataTableButtons() []*table.TableButtons {
	var array []*table.TableButtons
	array = append(array, &table.TableButtons{
		Text:      "添加",
		ClassName: "btn btn-sm btn-alt-primary open_iframe",
		Attribute: map[string]string{
			"href":      beego.URLFor("Advantage.Create", ":id", "__ID__", ":popup", 1),
			"data-area": "600px,400px",
		},
	})
	array = append(array, &table.TableButtons{
		Text:      "删除",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Advantage.Delete")},
	})
	array = append(array, &table.TableButtons{
		Text:      "启用",
		ClassName: "btn btn-sm btn-alt-primary mt-1 ids_enables",
		Attribute: map[string]string{"data-action": beego.URLFor("Advantage.Status"), "data-field": "status"},
	})
	array = append(array, &table.TableButtons{
		Text:      "禁用",
		ClassName: "btn btn-sm btn-alt-warning mt-1 ids_disables",
		Attribute: map[string]string{"data-action": beego.URLFor("Advantage.Status"), "data-field": "status"},
	})
	return array
}

/** 处理分页 **/
func (c *Advantage) PageListItems(length, draw, page int, search string) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(index.Advantage))
	if search != "" {
		qs = qs.Filter("title__icontains", search)
	}
	recordsTotal, _ := qs.Count()
	_, _ = qs.Limit(length, length*(page-1)).ValuesList(&lists, "id", "title", "description", "status", "create_time", "update_time")
	for _, v := range lists {
		v[2] = c.substr2HtmlHref("javascript:void(0)", v[2].(string), 0, 20)
	}
	data := map[string]interface{}{
		"draw":            draw,         // 请求次数
		"recordsFiltered": recordsTotal, // 从多少条里面筛选
		"recordsTotal":    recordsTotal, // 总条数
		"data":            lists,        // 筛选结果
	}
	return data
}

/**  转为pop提示 **/
func (c *Advantage) substr2HtmlHref(u, s string, start, end int) string {
	html := fmt.Sprintf(`<a href="%s" target="_blank" class="badge badge-primary js-tooltip" data-placement="top" data-toggle="tooltip" data-original-title="%s">%s...</a>`, u, s, beego.Substr(s, start, end))
	return html
}

/** 返回表单结构字段如何解析 **/
func (c *Advantage) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "switch", "date", "date"},
		"fieldName": {"id", "title", "description", "status", "create_time", "update_time"},
		"action":    {"", "", "", beego.URLFor("Advantage.Status"), "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (c *Advantage) TableButtonsType() []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "编辑",
			ClassName: "btn btn-sm btn-alt-primary open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Advantage.Edit", ":id", "__ID__", ":popup", 1),
				"data-area": "600px,400px",
			},
		},
		{
			Text:      "删除",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Advantage.Delete"),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}
