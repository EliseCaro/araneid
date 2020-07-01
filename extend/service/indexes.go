package service

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
)

type DefaultIndexesService struct{}

/** 获取一条数据 **/
func (service *DefaultIndexesService) Find(id int) (spider.Indexes, error) {
	item := spider.Indexes{Id: id}
	return item, orm.NewOrm().Read(&item)
}

/** 挂载一条索引 **/
func (service *DefaultIndexesService) UsageOneIndexes(arachnid int) spider.Indexes {
	var maps spider.Indexes
	_ = orm.NewOrm().QueryTable(new(spider.Indexes)).Filter("arachnid", arachnid).OrderBy("usage").One(&maps)
	_ = service.Inc(maps.Id)
	return maps
}

/** 计数器++ **/
func (service *DefaultIndexesService) Inc(id int) (errorMsg error) {
	_, errorMsg = orm.NewOrm().QueryTable(new(spider.Indexes)).Filter("id", id).Update(orm.Params{
		"usage": orm.ColValue(orm.ColAdd, 1),
	})
	return errorMsg
}

/** 计数器-- **/
func (service *DefaultIndexesService) Dec(id int) error {
	_, errorMessage := orm.NewOrm().QueryTable(new(spider.Indexes)).Filter("id", id).Update(orm.Params{
		"usage": orm.ColValue(orm.ColMinus, 1),
	})
	return errorMessage
}

/** 批量删除 **/
func (service *DefaultIndexesService) DeleteArray(array []int) (message error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, message = orm.NewOrm().Delete(&spider.Indexes{Id: v}); message != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if message == nil {
		_ = orm.NewOrm().Commit()
	}
	return message
}

/** 清空 **/
func (service *DefaultIndexesService) EmptyDelete(arachnid int) {
	var item []*spider.Indexes
	qs := orm.NewOrm().QueryTable(new(spider.Indexes))
	if arachnid > 0 {
		qs = qs.Filter("arachnid", arachnid)
	}
	_, _ = qs.All(&item)
	for _, v := range item {
		_, _ = orm.NewOrm().Delete(&spider.Indexes{Id: v.Id})
	}
}

/************************************************表格渲染机制 ************************************************************/

/** 获取需要渲染的Column **/
func (service *DefaultIndexesService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "标题", "name": "title", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "链接", "name": "urls", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "所属项目", "name": "arachnid", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "内部权重", "name": "sort", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "挂载次数", "name": "usage", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultIndexesService) DataTableButtons(id int) []*table.TableButtons {
	var array []*table.TableButtons
	if id > 0 {
		array = append(array, &table.TableButtons{Text: "返回上级", ClassName: "btn btn-sm btn-alt-dark mt-1 jump_urls", Attribute: map[string]string{"data-action": beego.URLFor("Arachnid.Index")}})
		array = append(array, &table.TableButtons{Text: "添加索引", ClassName: "btn btn-sm btn-alt-success mt-1 open_iframe", Attribute: map[string]string{"href": beego.URLFor("Indexes.Create", ":arachnid", id, ":popup", 1), "data-area": "580px,266px"}})
		array = append(array, &table.TableButtons{Text: "文件导入", ClassName: "btn btn-sm btn-alt-primary mt-1 open_iframe", Attribute: map[string]string{"href": beego.URLFor("Indexes.Import", ":arachnid", id, ":popup", 1), "data-area": "320px,227px"}})
	}
	array = append(array, &table.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-warning mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Indexes.Delete", ":arachnid", id)},
	})
	array = append(array, &table.TableButtons{
		Text:      "清空所有",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Indexes.Empty", ":arachnid", id)},
	})
	return array
}

/** 处理分页 **/
func (service *DefaultIndexesService) PageListItems(length, draw, page int, search string, id int) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(spider.Indexes))
	if search != "" {
		qs = qs.Filter("title__icontains", search)
	}
	if id > 0 {
		qs = qs.Filter("arachnid", id)
	}
	recordsTotal, _ := qs.Count()
	_, _ = qs.Limit(length, length*(page-1)).OrderBy("-id").ValuesList(&lists, "id", "title", "urls", "arachnid", "sort", "usage")
	arachnid := DefaultArachnidService{}
	for _, v := range lists {
		one, _ := arachnid.Find(int(v[3].(int64)))
		v[3] = one.Name
		v[1] = service.substr2HtmlHref(v[1].(string), v[2].(string), 0, 35)
		v[2] = service.substr2HtmlHref(v[2].(string), v[2].(string), 0, 20)
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
func (service *DefaultIndexesService) substr2HtmlHref(s, href string, start, end int) string {
	html := fmt.Sprintf(`<a href="%s" target="_blank" class="badge badge-primary js-tooltip" data-toggle="tooltip" data-original-title="%s">%s...</a>`, href, s, beego.Substr(s, start, end))
	return html
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultIndexesService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string", "string"},
		"fieldName": {"id", "title", "urls", "arachnid", "sort", "usage"},
		"action":    {"", "", "", "", "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultIndexesService) TableButtonsType(id int) []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "编辑",
			ClassName: "btn btn-sm btn-alt-warning open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Indexes.Edit", ":id", "__ID__", ":popup", 1, ":arachnid", id),
				"data-area": "580px,266px",
			},
		},
		{
			Text:      "删除",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Indexes.Delete", ":arachnid", id),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}
