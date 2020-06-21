package service

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/collect"
)

type DefaultCollectorService struct{}

/** 获取一条数据 **/
func (service *DefaultCollectorService) One(id int) collect.Result {
	var item collect.Result
	_ = orm.NewOrm().QueryTable(new(collect.Result)).Filter("id", id).One(&item)
	return item
}

/** 批量删除 **/
func (service *DefaultCollectorService) DeleteArray(array []int) (e error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, e = orm.NewOrm().Delete(&collect.Result{Id: v}); e != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if e == nil {
		_ = orm.NewOrm().Commit()
	}
	return e
}

/************************************************表格渲染机制 ************************************************************/

/** 获取需要渲染的Column **/
func (service *DefaultCollectorService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "结果标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "采集标题", "name": "title", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "采集地址", "name": "source", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "发布状态", "name": "status", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "更新时间", "name": "update_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "爬虫操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultCollectorService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "date"},
		"fieldName": {"id", "title", "source", "status", "update_time"},
		"action":    {"", "", "", "", ""},
	}
	return result
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultCollectorService) DataTableButtons(id int) []*table.TableButtons {
	var array []*table.TableButtons
	array = append(array, &table.TableButtons{
		Text:      "返回上级",
		ClassName: "btn btn-sm btn-alt-success mt-1 jump_urls",
		Attribute: map[string]string{"data-action": beego.URLFor("Collect.Index")},
	})
	array = append(array, &table.TableButtons{
		Text:      "发布选中",
		ClassName: "btn btn-sm btn-alt-primary mt-1 ids_enables",
		Attribute: map[string]string{
			"data-action": beego.URLFor("Collector.Status"),
			"data-field":  "status",
		},
	})
	array = append(array, &table.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Collector.Delete", ":parent", id)},
	})
	return array
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultCollectorService) TableButtonsType(id int) []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "查看结果",
			ClassName: "btn btn-sm btn-alt-primary open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Collector.Detail", ":id", "__ID__", ":popup", 1),
				"data-area": "620px,400px",
			},
		},
		{
			Text:      "删除结果",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Collector.Delete", ":parent", id),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}

/** 处理分页 **/
func (service *DefaultCollectorService) PageListItems(length, draw, page int, search string, id int) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(collect.Result)).Filter("collect", id)
	recordsTotal, _ := qs.Count()
	if search != "" {
		qs = qs.Filter("title__icontains", search)
	}
	_, _ = qs.Limit(length, length*(page-1)).OrderBy("status", "-id").ValuesList(&lists, "id", "title", "source", "status", "update_time")
	for _, v := range lists {
		v[1] = service.substr2HtmlSpan(v[1].(string), 0, 25)
		v[2] = service.substr2HtmlHref(v[2].(string), 0, 30)
		v[3] = service.int2HtmlStatus(v[3], v[0], id)
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
func (service *DefaultCollectorService) substr2HtmlSpan(s string, start, end int) string {
	html := fmt.Sprintf(`<span class="badge badge-primary js-tooltip" data-placement="top" data-toggle="tooltip" data-original-title="%s">%s...</span>`, s, beego.Substr(s, start, end))
	return html
}

/**  转为pop提示 **/
func (service *DefaultCollectorService) substr2HtmlHref(s string, start, end int) string {
	html := fmt.Sprintf(`<a href="%s" target="_blank" class="badge badge-primary js-tooltip" data-placement="top" data-toggle="tooltip" data-original-title="%s">%s...</a>`, s, s, beego.Substr(s, start, end))
	return html
}

/** 将状态码转为html **/
func (service *DefaultCollectorService) int2HtmlStatus(val interface{}, id interface{}, collect int) string {
	t := table.BuilderTable{}
	return t.AnalysisTypeSwitch(map[string]interface{}{"status": val, "id": id}, "status", beego.URLFor("Collector.Status", ":parent", collect), map[int64]string{-1: "已失败", 0: "待发布", 1: "已发布"})
}
