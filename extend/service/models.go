package service

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/collect"
	"github.com/beatrice950201/araneid/extend/model/spider"
)

type DefaultModelsService struct{}

/** 获取列表 **/
func (service *DefaultModelsService) One(id int) spider.Models {
	var item spider.Models
	_ = orm.NewOrm().QueryTable(new(spider.Models)).Filter("id", id).One(&item)
	return item
}

/** 获取列表 **/
func (service *DefaultModelsService) CollectMapColumn() []*collect.Collect {
	var item []*collect.Collect
	_, _ = orm.NewOrm().QueryTable(new(collect.Collect)).All(&item)
	return item
}

/** 批量删除结果 **/
func (service *DefaultModelsService) DeleteArray(array []int) (message error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		find := service.One(v)
		if _, message = orm.NewOrm().Delete(&find); message != nil {
			_ = orm.NewOrm().Rollback()
			break
		} else {
			s := DefaultDisguiseService{}
			_ = s.Dec(find.Disguise)
		}
	}
	if message == nil {
		_ = orm.NewOrm().Commit()
	}
	return message
}

/************************************************表格渲染机制 ************************************************************/

/** 获取需要渲染的Column **/
func (service *DefaultModelsService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "名称", "name": "name", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "挂载爬虫", "name": "collect", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "模版分组", "name": "template", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "伪原创", "name": "masking", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "栏目数", "name": "category", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "文档数", "name": "document", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "前缀数", "name": "prefix", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultModelsService) DataTableButtons() []*table.TableButtons {
	var array []*table.TableButtons
	array = append(array, &table.TableButtons{
		Text:      "创建模型",
		ClassName: "btn btn-sm btn-alt-success mt-1 open_iframe",
		Attribute: map[string]string{
			"href":      beego.URLFor("Models.Create", ":popup", 1),
			"data-area": "600px,400px",
		},
	})
	array = append(array, &table.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Models.Delete")},
	})
	return array
}

/** 处理分页 **/
func (service *DefaultModelsService) PageListItems(length, draw, page int, search string) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(spider.Models))
	recordsTotal, _ := qs.Count()
	if search != "" {
		qs = qs.Filter("name__icontains", search)
	}
	_, _ = qs.Limit(length, length*(page-1)).OrderBy("-id").ValuesList(&lists, "id", "name", "collect", "template", "disguise")
	for k, v := range lists {
		v[2] = service.collectName(v[2].(int64))
		v[3] = service.templateName(v[3].(int64))
		v[4] = service.disguiseName(v[4].(int64))
		v = append(v, service.classCount(v[0].(int64)))
		v = append(v, service.articleCount(v[0].(int64)))
		v = append(v, service.prefixCount(v[0].(int64)))
		lists[k] = v
	}
	data := map[string]interface{}{
		"draw":            draw,         // 请求次数
		"recordsFiltered": recordsTotal, // 从多少条里面筛选
		"recordsTotal":    recordsTotal, // 总条数
		"data":            lists,        // 筛选结果
	}
	return data
}

/** 获取爬虫名称 **/
func (service *DefaultModelsService) classCount(id int64) int64 {
	count, _ := orm.NewOrm().QueryTable(new(spider.Class)).Filter("model", id).Count()
	return count
}

/** 获取爬虫名称 **/
func (service *DefaultModelsService) articleCount(id int64) int64 {
	count, _ := orm.NewOrm().QueryTable(new(spider.Article)).Filter("model", id).Count()
	return count
}

/** 获取域名前缀数量 **/
func (service *DefaultModelsService) prefixCount(id int64) int64 {
	count, _ := orm.NewOrm().QueryTable(new(spider.Prefix)).Filter("model", id).Count()
	return count
}

/** 获取爬虫名称 **/
func (service *DefaultModelsService) collectName(id int64) string {
	s := DefaultCollectService{}
	return s.One(int(id)).Name
}

/** 获取模板分组名称 **/
func (service *DefaultModelsService) templateName(id int64) string {
	s := DefaultTemplateService{}
	return s.One(int(id)).Name
}

/** 获取自然语言名称 **/
func (service *DefaultModelsService) disguiseName(id int64) string {
	s := DefaultDisguiseService{}
	res, _ := s.Find(int(id))
	return res.Name
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultModelsService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string", "string", "string", "string"},
		"fieldName": {"id", "name", "collect", "template", "disguise", "category", "document", "prefix"},
		"action":    {"", "", "", "", "", "", "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultModelsService) TableButtonsType() []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "前缀库",
			ClassName: "btn btn-sm btn-alt-success open_iframe js-tooltip",
			Attribute: map[string]string{
				"href":                beego.URLFor("Prefix.Index", ":model", "__ID__", ":popup", 1),
				"data-area":           "600px,400px",
				"data-toggle":         "tooltip",
				"data-original-title": "域名前缀库",
			},
		},
		{
			Text:      "栏目库",
			ClassName: "btn btn-sm btn-alt-primary jump_urls",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Class.Index", ":model", "__ID__"),
			},
		},
		{
			Text:      "内容库",
			ClassName: "btn btn-sm btn-alt-success jump_urls",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Article.Index", ":model", "__ID__"),
			},
		},
		{
			Text:      "编辑",
			ClassName: "btn btn-sm btn-alt-warning open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Models.Edit", ":id", "__ID__", ":popup", 1),
				"data-area": "600px,400px",
			},
		},
		{
			Text:      "删除",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Models.Delete"),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}
