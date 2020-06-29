package service

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
)

type DefaultArachnidService struct{}

/** 模型挂载列表 **/
func (service *DefaultArachnidService) ModelSelect() []*spider.Models {
	var item []*spider.Models
	_, _ = orm.NewOrm().QueryTable(new(spider.Models)).All(&item)
	return item
}

/** 模型挂载列表 **/
func (service *DefaultArachnidService) MatchSelect() []*spider.Match {
	var item []*spider.Match
	_, _ = orm.NewOrm().QueryTable(new(spider.Match)).All(&item)
	return item
}

/** 获取一条数据 **/
func (service *DefaultArachnidService) Find(id int) (spider.Arachnid, error) {
	item := spider.Arachnid{Id: id}
	return item, orm.NewOrm().Read(&item)
}

/** 更新状态 **/
func (service *DefaultArachnidService) StatusArray(array []int, status int8) (e error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, e = orm.NewOrm().Update(&spider.Arachnid{Id: v, Status: status}, "Status"); e != nil {
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
func (service *DefaultArachnidService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "项目标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "项目名称", "name": "name", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "友链数量", "name": "link", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "轮链数量", "name": "zoology", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "索引数量", "name": "indexes", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "挂载模型", "name": "models", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "KTDB", "name": "matching", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "启用状态", "name": "status", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultArachnidService) DataTableButtons() []*table.TableButtons {
	var array []*table.TableButtons
	array = append(array, &table.TableButtons{
		Text:      "创建项目",
		ClassName: "btn btn-sm btn-alt-success mt-1 open_iframe",
		Attribute: map[string]string{
			"href":      beego.URLFor("Arachnid.Create", ":popup", 1),
			"data-area": "600px,380px",
		},
	})
	array = append(array, &table.TableButtons{
		Text:      "启用选中",
		ClassName: "btn btn-sm btn-alt-primary mt-1 ids_enables",
		Attribute: map[string]string{"data-action": beego.URLFor("Arachnid.Status"), "data-field": "status"},
	})
	array = append(array, &table.TableButtons{
		Text:      "禁用选中",
		ClassName: "btn btn-sm btn-alt-warning mt-1 ids_disables",
		Attribute: map[string]string{"data-action": beego.URLFor("Arachnid.Status"), "data-field": "status"},
	})
	array = append(array, &table.TableButtons{
		Text:      "删除项目",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Arachnid.Delete")},
	})
	return array
}

/** 处理分页 **/
func (service *DefaultArachnidService) PageListItems(length, draw, page int, search string) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(spider.Arachnid))
	if search != "" {
		qs = qs.Filter("name__icontains", search)
	}
	recordsTotal, _ := qs.Count()
	_, _ = qs.Limit(length, length*(page-1)).OrderBy("-id").ValuesList(&lists, "id", "name", "link", "zoology", "indexes", "models", "matching", "status")
	models := DefaultModelsService{}
	matching := DefaultMatchService{}
	for _, v := range lists {
		v[5] = models.One(int(v[5].(int64))).Name
		matchingOne, _ := matching.Find(int(v[6].(int64)))
		v[6] = matchingOne.Name
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
func (service *DefaultArachnidService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string", "string", "string", "switch"},
		"fieldName": {"id", "name", "link", "zoology", "indexes", "models", "matching", "status"},
		"action":    {"", "", "", "", "", "", "", beego.URLFor("Arachnid.Status")},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultArachnidService) TableButtonsType() []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "索引池库",
			ClassName: "btn btn-sm btn-alt-primary jump_urls",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Indexes.Index", ":id", "__ID__"),
			},
		},
		{
			Text:      "关键词库",
			ClassName: "btn btn-sm btn-alt-success jump_urls",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Keyword.Index", ":id", "__ID__"),
			},
		},
		{
			Text:      "编辑项目",
			ClassName: "btn btn-sm btn-alt-warning js-tooltip open_iframe",
			Attribute: map[string]string{
				"href":                beego.URLFor("Arachnid.Edit", ":id", "__ID__", ":popup", 1),
				"data-area":           "600px,380px",
				"data-toggle":         "tooltip",
				"data-original-title": "编辑结果将不作用于已经生成的链接；除非清空链接库；",
			},
		},
		{
			Text:      "删除项目",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Arachnid.Delete"),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}
