package admin

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/addons"
	"github.com/go-playground/validator"
)

/** 第三方程序插件 **/
type Addons struct{ Main }

// @router /addons/index [get,post]
func (c *Addons) Index() {
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

// @router /addons/create [get,post]
func (c *Addons) Create() {
	if c.IsAjax() {
		item := addons.Addons{}
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Insert(&item); err == nil {
			_ = c.adjunctService.Inc(item.File, 0)
			c.Succeed(&controllers.ResultJson{Message: "添加成功", Url: beego.URLFor("Addons.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "添加失败，请稍后再试！"})
		}
	}
}

// @router /addons/edit [get,post]
func (c *Addons) Edit() {
	id := c.GetMustInt(":id", "非法请求！")
	item, _ := c.Find(id)
	if c.IsAjax() {
		oldImage := item.File
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Update(&item); err == nil {
			_ = c.adjunctService.Inc(item.File, oldImage)
			c.Succeed(&controllers.ResultJson{Message: "更新成功", Url: beego.URLFor("Addons.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "更新失败，请稍后再试！"})
		}
	}
	c.Data["info"] = item
	c.Data["file"] = c.adjunctService.FindId(item.File)
}

// @router /addons/download [get]
func (c *Addons) Download() {
	id := c.GetMustInt(":id", "非法请求！")
	item, _ := c.Find(id)
	path := table.FilePath(item.File)
	c.Ctx.Redirect(302, path)
}

// @router /addons/delete [post]
func (c *Addons) Delete() {
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除成功！马上返回中。。。",
			Url:     beego.URLFor("Addons.Index"),
		})
	}
}

/************************************************表格渲染机制 ************************************************************/

/** 批量删除 **/
func (c *Addons) DeleteArray(array []int) (e error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		one, _ := c.Find(v)
		if _, e = orm.NewOrm().Delete(&addons.Addons{Id: v}); e != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
		_ = c.adjunctService.Dec(one.File)
	}
	if e == nil {
		_ = orm.NewOrm().Commit()
	}
	return e
}

/** 获取一条数据 **/
func (c *Addons) Find(id int) (addons.Addons, error) {
	item := addons.Addons{
		Id: id,
	}
	return item, orm.NewOrm().Read(&item)
}

/** 获取需要渲染的Column **/
func (c *Addons) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "ID", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "插件名称", "name": "title", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "插件版本", "name": "version", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "插件作者", "name": "author", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "插件功能", "name": "description", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (c *Addons) DataTableButtons() []*table.TableButtons {
	var array []*table.TableButtons
	array = append(array, &table.TableButtons{
		Text:      "上传插件",
		ClassName: "btn btn-sm btn-alt-primary open_iframe",
		Attribute: map[string]string{
			"href":      beego.URLFor("Addons.Create", ":id", "__ID__", ":popup", 1),
			"data-area": "600px,390px",
		},
	})
	array = append(array, &table.TableButtons{
		Text:      "删除插件",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Addons.Delete")},
	})
	return array
}

/** 处理分页 **/
func (c *Addons) PageListItems(length, draw, page int, search string) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(addons.Addons))
	if search != "" {
		qs = qs.Filter("title__icontains", search)
	}
	recordsTotal, _ := qs.Count()
	_, _ = qs.Limit(length, length*(page-1)).ValuesList(&lists, "id", "title", "version", "author", "description")
	for _, v := range lists {
		v[4] = c.substr2HtmlHref("javascript:void(0)", v[4].(string))
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
func (c *Addons) substr2HtmlHref(u, s string) string {
	html := fmt.Sprintf(`<a href="%s" target="_blank" class="badge badge-primary js-tooltip" data-placement="top" data-toggle="tooltip" data-original-title="%s">功能简介</a>`, u, s)
	return html
}

/** 返回表单结构字段如何解析 **/
func (c *Addons) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string"},
		"fieldName": {"id", "title", "version", "author", "description"},
		"action":    {"", "", "", "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (c *Addons) TableButtonsType() []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "下载插件",
			ClassName: "btn btn-sm btn-alt-primary jump_urls",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Addons.Download", ":id", "__ID__"),
			},
		},
		{
			Text:      "编辑插件",
			ClassName: "btn btn-sm btn-alt-primary open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Addons.Edit", ":id", "__ID__", ":popup", 1),
				"data-area": "600px,390px",
			},
		},
		{
			Text:      "删除插件",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Addons.Delete"),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}
