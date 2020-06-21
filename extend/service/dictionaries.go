package service

import (
	"github.com/astaxie/beego"
	table "github.com/beatrice950201/araneid/extend/func"
)

type DefaultDictionariesService struct{}

/************************************************分类表格渲染机制 ************************************************************/

/** 获取需要渲染的Column **/
func (service *DefaultDictionariesService) DataTableCateColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "标题", "name": "title", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "首字母", "name": "initial", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "源地址", "name": "source", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "发布状态", "name": "status", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "更新时间", "name": "update_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "爬虫操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultDictionariesService) DataTableCateButtons() []*table.TableButtons {
	var array []*table.TableButtons
	array = append(array, &table.TableButtons{
		Text:      "启动爬虫",
		ClassName: "btn btn-sm btn-alt-success mt-1 jump_urls",
		Attribute: map[string]string{"data-action": beego.URLFor("Collect.Index")},
	})
	array = append(array, &table.TableButtons{
		Text:      "爬虫配置",
		ClassName: "btn btn-sm btn-alt-success mt-1 jump_urls",
		Attribute: map[string]string{"data-action": beego.URLFor("Collect.Index")},
	})
	array = append(array, &table.TableButtons{
		Text:      "发布选中",
		ClassName: "btn btn-sm btn-alt-primary mt-1 ids_enables",
		Attribute: map[string]string{
			"data-action": beego.URLFor("Dictionaries.Status"),
			"data-field":  "status",
		},
	})
	array = append(array, &table.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Dictionaries.Delete")},
	})
	return array
}
