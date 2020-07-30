package service

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/movie"
)

type DefaultMovieService struct{}

/*********************************以下为表格渲染*******************************************8/

/** 获取需要渲染的Column **/
func (service *DefaultMovieService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "作品标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "作品名称", "name": "name", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "作品年份", "name": "year", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "内容类型", "name": "class_type", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "简短标题", "name": "short", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "对方地址", "name": "source", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "发布状态", "name": "status", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "更新时间", "name": "update_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "爬虫操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultMovieService) DataTableButtons() []*table.TableButtons {
	var array []*table.TableButtons
	array = append(array, &table.TableButtons{
		Text:      "发布选中",
		ClassName: "btn btn-sm btn-alt-primary mt-1 ids_enables",
		Attribute: map[string]string{
			"data-action": beego.URLFor("Movie.Push"),
			"data-field":  "status",
		},
	})
	array = append(array, &table.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Movie.Delete")},
	})
	return array
}

/*** 获取分页数据 **/
func (service *DefaultMovieService) PageListItems(length, draw, page int, search string) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(movie.Movie))
	if search != "" {
		qs = qs.Filter("title__icontains", search)
	}
	recordsTotal, _ := qs.Count()
	_, _ = qs.Limit(length, length*(page-1)).ValuesList(&lists, "id", "name", "year", "class_type", "short", "source", "status", "update_time")
	for _, v := range lists {
		v[4] = service.Int2HtmlStatus(v[4], v[0], beego.URLFor("Movie.Push"))
	}
	data := map[string]interface{}{
		"draw":            draw,         // 请求次数
		"recordsFiltered": recordsTotal, // 从多少条里面筛选
		"recordsTotal":    recordsTotal, // 总条数
		"data":            lists,        // 筛选结果
	}
	return data
}

/** 将状态码转为html **/
func (service *DefaultMovieService) Int2HtmlStatus(val interface{}, id interface{}, url string) string {
	t := table.BuilderTable{}
	return t.AnalysisTypeSwitch(map[string]interface{}{"status": val, "id": id}, "status", url, map[int64]string{-1: "已失败", 0: "待发布", 1: "已发布"})
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultMovieService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string", "string", "string", "date"},
		"fieldName": {"id", "name", "year", "class_type", "short", "source", "status", "update_time"},
		"action":    {"", "", "", "", "", "", "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultMovieService) TableButtonsType() []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "查看",
			ClassName: "btn btn-sm btn-alt-primary open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Movie.Detail", ":id", "__ID__", ":popup", 1),
				"data-area": "620px,400px",
			},
		},
		{
			Text:      "删除",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Movie.Delete"),
				"data-ids":    "__ID__",
			},
		},
		{
			Text:      "剧情",
			ClassName: "btn btn-sm btn-alt-success jump_urls",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Movie.Index", ":id", "__ID__"),
			},
		},
	}
	return buttons
}
