package service

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
)

type DefaultMatchService struct{}

/** 获取一条数据 **/
func (service *DefaultMatchService) Find(id int) (spider.Match, error) {
	item := spider.Match{Id: id}
	return item, orm.NewOrm().Read(&item)
}

/** 获取一条数据 **/
func (service *DefaultMatchService) IsExtend(id int) (is bool) {
	var item spider.Arachnid
	if err := orm.NewOrm().QueryTable(new(spider.Arachnid)).Filter("matching", id).One(&item); err == nil && item.Id > 0 {
		is = true
	}
	return is
}

/** 批量删除结果 **/
func (service *DefaultMatchService) DeleteArray(array []int) (message error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if service.IsExtend(v) == false {
			if _, message = orm.NewOrm().Delete(&spider.Match{Id: v}); message != nil {
				_ = orm.NewOrm().Rollback()
				break
			}
		} else {
			message = errors.New("模板已被蜘蛛池挂载；不能删除使用中的模板～")
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if message == nil {
		_ = orm.NewOrm().Commit()
	}
	return message
}

/************************************************表格渲染机制 ************************************************************/

/** 获取需要渲染的Column **/
func (service *DefaultMatchService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "模板标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "模板名称", "name": "name", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "首页标题", "name": "index_title", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "列表页标题", "name": "cate_title", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "详情页标题", "name": "detail_title", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultMatchService) DataTableButtons() []*table.TableButtons {
	var array []*table.TableButtons
	array = append(array, &table.TableButtons{
		Text:      "创建模板",
		ClassName: "btn btn-sm btn-alt-success mt-1 open_iframe",
		Attribute: map[string]string{
			"href":      beego.URLFor("Match.Create", ":popup", 1),
			"data-area": "600px,410px",
		},
	})
	array = append(array, &table.TableButtons{
		Text:      "删除模板",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Match.Delete")},
	})
	return array
}

/** 处理分页 **/
func (service *DefaultMatchService) PageListItems(length, draw, page int, search string) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(spider.Match))
	if search != "" {
		qs = qs.Filter("name__icontains", search)
	}
	recordsTotal, _ := qs.Count()
	_, _ = qs.Limit(length, length*(page-1)).OrderBy("-id").ValuesList(&lists, "id", "name", "index_title", "cate_title", "detail_title")
	for _, v := range lists {
		v[2] = service.substr2HtmlSpan(v[2].(string), 0, 20)
		v[3] = service.substr2HtmlSpan(v[3].(string), 0, 20)
		v[4] = service.substr2HtmlSpan(v[4].(string), 0, 20)
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
func (service *DefaultMatchService) substr2HtmlSpan(s string, start, end int) string {
	html := fmt.Sprintf(`<span class="badge badge-primary js-tooltip" data-placement="top" data-toggle="tooltip" data-original-title="%s">%s...</span>`, s, beego.Substr(s, start, end))
	return html
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultMatchService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string"},
		"fieldName": {"id", "name", "index_title", "cate_title", "detail_title"},
		"action":    {"", "", "", "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultMatchService) TableButtonsType() []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "编辑模板",
			ClassName: "btn btn-sm btn-alt-primary js-tooltip open_iframe",
			Attribute: map[string]string{
				"href":                beego.URLFor("Match.Edit", ":id", "__ID__", ":popup", 1),
				"data-area":           "600px,410px",
				"data-toggle":         "tooltip",
				"data-original-title": "编辑模板结果将不作用于已经生成的链接；除非清空链接库；",
			},
		},
		{
			Text:      "删除模板",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Match.Delete"),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}
