package service

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
)

type DefaultArticleService struct{}

/****************** 以下为表格渲染  ***********************/

/** 获取需要渲染的Column **/
func (service *DefaultArticleService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "结果标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "所属模型", "name": "model", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "挂载次数", "name": "usage", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "栏目标题", "name": "title", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "更新时间", "name": "update_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "数据操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultArticleService) DataTableButtons(id int) []*_func.TableButtons {
	var array []*_func.TableButtons
	if id > 0 {
		array = append(array, &_func.TableButtons{
			Text:      "返回模型",
			ClassName: "btn btn-sm btn-alt-success mt-1 jump_urls",
			Attribute: map[string]string{"data-action": beego.URLFor("Models.Index")},
		})
	}
	array = append(array, &_func.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Article.Delete", ":model", id)},
	})
	return array
}

/** 处理分页 **/
func (service *DefaultArticleService) PageListItems(length, draw, page int, search string, id int) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(spider.Article))
	if id > 0 {
		qs = qs.Filter("model", id)
	}
	if search != "" {
		qs = qs.Filter("title__icontains", search)
	}
	recordsTotal, _ := qs.Count()
	_, _ = qs.Limit(length, length*(page-1)).OrderBy("-id").ValuesList(&lists, "id", "model", "usage", "title", "update_time")
	object := DefaultModelsService{}
	for _, v := range lists {
		v[1] = object.One(int(v[1].(int64))).Name
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
func (service *DefaultArticleService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "date"},
		"fieldName": {"id", "model", "usage", "title", "update_time"},
		"action":    {"", "", "", "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultArticleService) TableButtonsType(id int) []*_func.TableButtons {
	buttons := []*_func.TableButtons{
		{
			Text:      "查看详情",
			ClassName: "btn btn-sm btn-alt-warning open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Article.Edit", ":id", "__ID__", ":popup", 1),
				"data-area": "580px,380px",
			},
		},
		{
			Text:      "删除文档",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Article.Delete", ":model", id),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}
