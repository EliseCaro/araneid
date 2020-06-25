package service

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
)

type DefaultDisguiseService struct{}

/** 获取一条用户数据 **/
func (service *DefaultDisguiseService) Find(id int) (spider.Disguise, error) {
	item := spider.Disguise{
		Id: id,
	}
	return item, orm.NewOrm().Read(&item)
}

/** 批量删除结果 **/
func (service *DefaultDisguiseService) DeleteArray(array []int) (message error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, message = orm.NewOrm().Delete(&spider.Disguise{Id: v}); message != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if message == nil {
		_ = orm.NewOrm().Commit()
	}
	return message
}

/** 更新状态 **/
func (service *DefaultDisguiseService) UpdateStatus(id int, status int8, field string) (e error) {
	if field == "keyword" {
		_, e = orm.NewOrm().Update(&spider.Disguise{Id: id, Keyword: status}, "Keyword")
	}
	if field == "description" {
		_, e = orm.NewOrm().Update(&spider.Disguise{Id: id, Description: status}, "Description")
	}
	if field == "context" {
		_, e = orm.NewOrm().Update(&spider.Disguise{Id: id, Context: status}, "Context")
	}
	return e
}

/************************************************表格渲染机制 ************************************************************/

/** 获取需要渲染的Column **/
func (service *DefaultDisguiseService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "机器名称", "name": "name", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "挂载次数", "name": "usage", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "自然单词", "name": "keyword", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "自然描述", "name": "description", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "自然内容", "name": "context", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "接口KEY", "name": "api_key", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "接口SECRET", "name": "api_secret", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultDisguiseService) DataTableButtons() []*table.TableButtons {
	var array []*table.TableButtons
	array = append(array, &table.TableButtons{
		Text:      "创建机器",
		ClassName: "btn btn-sm btn-alt-success mt-1 open_iframe",
		Attribute: map[string]string{
			"href":      beego.URLFor("Disguise.Create", ":popup", 1),
			"data-area": "580px,400px",
		},
	})
	array = append(array, &table.TableButtons{
		Text:      "删除机器",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Disguise.Delete")},
	})
	return array
}

/** 处理分页 **/
func (service *DefaultDisguiseService) PageListItems(length, draw, page int, search string) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(spider.Disguise))
	recordsTotal, _ := qs.Count()
	if search != "" {
		qs = qs.Filter("name__icontains", search)
	}
	_, _ = qs.Limit(length, length*(page-1)).OrderBy("-id").ValuesList(&lists, "id", "name", "usage", "keyword", "description", "context", "api_key", "api_secret")
	for _, v := range lists {
		v[6] = service.Substr2HtmlSpan(v[6].(string), 0, 30)
		v[7] = service.Substr2HtmlSpan(v[7].(string), 0, 30)
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
func (service *DefaultDisguiseService) Substr2HtmlSpan(s string, start, end int) string {
	html := fmt.Sprintf(`<span class="badge badge-primary js-tooltip" data-placement="top" data-toggle="tooltip" data-original-title="%s">%s</span>`, s, beego.Substr(s, start, end))
	return html
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultDisguiseService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "switch", "switch", "switch", "string", "string"},
		"fieldName": {"id", "name", "usage", "keyword", "description", "context", "api_key", "api_secret"},
		"action":    {"", "", "", beego.URLFor("Disguise.Status"), beego.URLFor("Disguise.Status"), beego.URLFor("Disguise.Status"), "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultDisguiseService) TableButtonsType() []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "编辑",
			ClassName: "btn btn-sm btn-alt-primary open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Disguise.Edit", ":id", "__ID__", ":popup", 1),
				"data-area": "600px,400px",
			},
		},
		{
			Text:      "删除",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Disguise.Delete"),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}
