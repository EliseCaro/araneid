package service

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
)

type DefaultKeywordService struct{}

/** 获取一条数据 **/
func (service *DefaultKeywordService) One(id int) spider.Keyword {
	var item spider.Keyword
	_ = orm.NewOrm().QueryTable(new(spider.Keyword)).Filter("id", id).One(&item)
	return item
}

/** 批量删除 **/
func (service *DefaultKeywordService) DeleteArray(array []int) (message error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, message = orm.NewOrm().Delete(&spider.Keyword{Id: v}); message != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if message == nil {
		_ = orm.NewOrm().Commit()
	}
	return message
}

/** 清空关键词 **/
func (service *DefaultKeywordService) EmptyDelete(arachnid int) {
	var item []*spider.Keyword
	qs := orm.NewOrm().QueryTable(new(spider.Keyword))
	if arachnid > 0 {
		qs = qs.Filter("arachnid", arachnid)
	}
	_, _ = qs.All(&item)
	for _, v := range item {
		_, _ = orm.NewOrm().Delete(&spider.Keyword{Id: v.Id})
	}
}

/************************************************表格渲染机制 ************************************************************/

/** 获取需要渲染的Column **/
func (service *DefaultKeywordService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "所属项目", "name": "arachnid", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "关键词", "name": "title", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "更新时间", "name": "update_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "创建时间", "name": "create_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultKeywordService) DataTableButtons(id int) []*table.TableButtons {
	var array []*table.TableButtons
	if id > 0 {
		array = append(array, &table.TableButtons{Text: "返回上级", ClassName: "btn btn-sm btn-alt-dark mt-1 jump_urls", Attribute: map[string]string{"data-action": beego.URLFor("Arachnid.Index")}})
		array = append(array, &table.TableButtons{Text: "添加关键词", ClassName: "btn btn-sm btn-alt-success mt-1 open_iframe", Attribute: map[string]string{"href": beego.URLFor("Keyword.Create", ":arachnid", id, ":popup", 1), "data-area": "350px,67px", "data-skin": "rounded-lg-custom"}})
		array = append(array, &table.TableButtons{Text: "批量导入", ClassName: "btn btn-sm btn-alt-primary mt-1 open_iframe", Attribute: map[string]string{"href": beego.URLFor("Keyword.Import", ":arachnid", id, ":popup", 1), "data-area": "320px,227px"}})
	}
	array = append(array, &table.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-warning mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Keyword.Delete", ":arachnid", id)},
	})
	array = append(array, &table.TableButtons{
		Text:      "清空所有",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Keyword.Empty", ":arachnid", id)},
	})
	return array
}

/** 处理分页 **/
func (service *DefaultKeywordService) PageListItems(length, draw, page int, search string, id int) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(spider.Keyword))
	if search != "" {
		qs = qs.Filter("title__icontains", search)
	}
	if id > 0 {
		qs = qs.Filter("arachnid", id)
	}
	recordsTotal, _ := qs.Count()
	_, _ = qs.Limit(length, length*(page-1)).OrderBy("-id").ValuesList(&lists, "id", "title", "arachnid", "create_time", "update_time")
	arachnid := DefaultArachnidService{}
	for _, v := range lists {
		one, _ := arachnid.Find(int(v[2].(int64)))
		v[2] = one.Name
	}
	data := map[string]interface{}{
		"draw":            draw,         // 请求次数
		"recordsFiltered": recordsTotal, // 从多少条里面筛选
		"recordsTotal":    recordsTotal, // 总条数
		"data":            lists,        // 筛选结果
	}
	return data
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultKeywordService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "date", "date"},
		"fieldName": {"id", "title", "arachnid", "create_time", "update_time"},
		"action":    {"", "", "", "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultKeywordService) TableButtonsType(id int) []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "编辑",
			ClassName: "btn btn-sm btn-alt-warning open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Keyword.Edit", ":id", "__ID__", ":popup", 1, ":arachnid", id),
				"data-area": "350px,67px",
				"data-skin": "rounded-lg-custom",
			},
		},
		{
			Text:      "删除",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Keyword.Delete", ":arachnid", id),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}
