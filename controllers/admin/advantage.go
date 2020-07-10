package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/index"
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

/************************************************表格渲染机制 ************************************************************/

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
		ClassName: "btn btn-sm btn-alt-success mt-1 jump_urls",
		Attribute: map[string]string{"data-action": beego.URLFor("Advantage.Create")},
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
	data := map[string]interface{}{
		"draw":            draw,         // 请求次数
		"recordsFiltered": recordsTotal, // 从多少条里面筛选
		"recordsTotal":    recordsTotal, // 总条数
		"data":            lists,        // 筛选结果
	}
	return data
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
				"data-area": "580px,370px",
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
