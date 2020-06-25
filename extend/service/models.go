package service

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/collect"
	"github.com/beatrice950201/araneid/extend/model/spider"
)

type DefaultModelsService struct{}

/** 获取爬虫列表 **/
func (service *DefaultModelsService) CollectMapColumn() []*collect.Collect {
	var item []*collect.Collect
	_, _ = orm.NewOrm().QueryTable(new(collect.Collect)).All(&item)
	return item
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
	maps = append(maps, map[string]interface{}{"title": "更新时间", "name": "update_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "分类数", "name": "category", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "内容数", "name": "document", "className": "text-center", "order": false})
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
	_, _ = qs.Limit(length, length*(page-1)).OrderBy("status", "-id").ValuesList(&lists, "id", "name", "collect", "template", "masking", "update_time")
	for _, v := range lists {
		v[6] = 0
		v[7] = 1
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
func (service *DefaultModelsService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string", "date", "string", "string"},
		"fieldName": {"id", "name", "collect", "template", "masking", "update_time", "category", "document"},
		"action":    {"", "", "", "", "", "", "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultModelsService) TableButtonsType() []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "删除模型",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Models.Delete"),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}
