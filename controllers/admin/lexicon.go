package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/dictionaries"
)

/** 词典采集列表 **/
type Lexicon struct {
	Main
}

// @router /lexicon/index [get,post]
func (c *Lexicon) Index() {
	id := c.GetMustInt(":id", "非法请求！")
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.pageListItems(length, draw, c.PageNumber(), search, id)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.tableColumnsType(), c.tableButtonsType(id))
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.dataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.dataTableButtons(id) {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
	c.Data["info"] = c.dictionariesService.One(id)
}

// @router /lexicon/delete [post]
func (c *Lexicon) Delete() {
	parent := c.GetMustInt(":parent", "爬虫ID不合法...")
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.dictionariesService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除结果成功！马上返回中。。。",
			Url:     beego.URLFor("Lexicon.Index", ":id", parent),
		})
	}
}

// @router /lexicon/status [post]
func (c *Lexicon) Status() {
	array := c.checkBoxIds(":ids[]", ":ids")
	parent := c.GetMustInt(":parent", "非法操作；已被拦截")
	for _, v := range array {
		c.dictionariesService.PushDetailAPI(c.dictionariesService.One(v))
	}
	c.Succeed(&controllers.ResultJson{Message: "提交发布成功！马上返回中。。。", Url: beego.URLFor("Lexicon.Index", ":id", parent)})
}

/******************表格渲染  ***********************/

/** 获取需要渲染的Column **/
func (c *Lexicon) dataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "结果标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "采集标题", "name": "title", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "采集地址", "name": "source", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "发布状态", "name": "status", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "更新时间", "name": "update_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "数据操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (c *Lexicon) dataTableButtons(id int) []*_func.TableButtons {
	var array []*_func.TableButtons
	array = append(array, &_func.TableButtons{
		Text:      "返回上级",
		ClassName: "btn btn-sm btn-alt-success mt-1 jump_urls",
		Attribute: map[string]string{"data-action": beego.URLFor("Dictionaries.Index")},
	})
	array = append(array, &_func.TableButtons{
		Text:      "发布选中",
		ClassName: "btn btn-sm btn-alt-primary mt-1 ids_enables",
		Attribute: map[string]string{
			"data-action": beego.URLFor("Lexicon.Status", ":parent", id),
			"data-field":  "status",
		},
	})
	array = append(array, &_func.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Lexicon.Delete", ":parent", id)},
	})
	return array
}

/** 处理分页 **/
func (c *Lexicon) pageListItems(length, draw, page int, search string, id int) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(dictionaries.DictCate)).Filter("pid", id)
	recordsTotal, _ := qs.Count()
	if search != "" {
		qs = qs.Filter("title__icontains", search)
	}
	_, _ = qs.Limit(length, length*(page-1)).OrderBy("status", "-id").ValuesList(&lists, "id", "title", "source", "status", "update_time")
	for _, v := range lists {
		v[1] = c.dictionariesService.Substr2HtmlSpan(v[1].(string), 0, 25)
		v[2] = c.dictionariesService.Substr2HtmlHref(v[2].(string), 0, 30)
		v[3] = c.dictionariesService.Int2HtmlStatus(v[3], v[0], beego.URLFor("Lexicon.Status", ":parent", id))
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
func (c *Lexicon) tableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "date"},
		"fieldName": {"id", "title", "source", "status", "update_time"},
		"action":    {"", "", "", "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (c *Lexicon) tableButtonsType(id int) []*_func.TableButtons {
	buttons := []*_func.TableButtons{
		{
			Text:      "查看结果",
			ClassName: "btn btn-sm btn-alt-primary open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Dictionaries.Detail", ":id", "__ID__", ":popup", 1),
				"data-area": "620px,400px",
			},
		},
		{
			Text:      "删除结果",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Lexicon.Delete", ":parent", id),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}

/* ********************表格渲染完成*********************** */
