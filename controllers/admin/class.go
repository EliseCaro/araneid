package admin

import (
	"errors"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/go-playground/validator"
	"strconv"
)

/** 蜘蛛池资源栏目库管理 **/
type Class struct{ Main }

// @router /class/index [get,post]
func (c *Class) Index() {
	model, _ := c.GetInt(":model", 0)
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.pageListItems(length, draw, c.PageNumber(), search, model)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.tableColumnsType(), c.tableButtonsType(model))
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.dataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.dataTableButtons(model) {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
	c.Data["action"] = beego.URLFor("Class.Index", ":model", model)
}

// @router /class/import [get,post]
func (c *Class) Import() {
	model, _ := c.GetInt(":model", 0)
	if c.IsAjax() {
		file := c.GetMustInt("files", "请上传正确格式的xlsx文件！")
		path := c.adjunctService.FindId(file).Path
		if f, err := excelize.OpenFile("." + path); err == nil {
			var items []*spider.Class
			for _, sheet := range f.GetSheetList() {
				rows, _ := f.GetRows(sheet)
				for _, row := range rows {
					if len(row) >= 3 && c.classService.OneExtends(row[0]).Id == 0 {
						items = append(items, &spider.Class{Model: model, Title: row[0], Keywords: row[1], Description: row[2]})
					}
				}
			}
			if len(items) > 0 {
				if index, err := orm.NewOrm().InsertMulti(100, items); err == nil {
					c.Succeed(&controllers.ResultJson{Message: "批量导入栏目共" + strconv.FormatInt(index, 10) + "个;正在刷新返回！", Url: beego.URLFor("Class.Index", ":model", model)})
				} else {
					c.Fail(&controllers.ResultJson{Message: "批量导入栏目失败；失败原因:" + error.Error(err)})
				}
			} else {
				c.Fail(&controllers.ResultJson{Message: "导入完成；未添加任何数据；已被去重过滤或者格式非法！"})
			}
		} else {
			c.Fail(&controllers.ResultJson{Message: "打开文件失败；失败原因:" + error.Error(err)})
		}
	}
	c.Data["model"] = model
}

// @router /class/create [get,post]
func (c *Class) Create() {
	model, _ := c.GetInt(":model", 0)
	if c.IsAjax() {
		item := spider.Class{}
		item.Model = model
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Insert(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "创建栏目成功", Url: beego.URLFor("Class.Index", ":model", model)})
		} else {
			c.Fail(&controllers.ResultJson{Message: "创建栏目失败，请稍后再试！"})
		}
	}
	c.Data["model"] = model
}

// @router /class/edit [get,post]
func (c *Class) Edit() {
	id := c.GetMustInt(":id", "非法请求！")
	if c.IsAjax() {
		item := c.one(id)
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Update(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "更新栏目成功", Url: beego.URLFor("Class.Index", ":model", item.Model)})
		} else {
			c.Fail(&controllers.ResultJson{Message: "更新栏目失败，请稍后再试！"})
		}
	}
	c.Data["info"] = c.one(id)
}

// @router /class/delete [post]
func (c *Class) Delete() {
	model, _ := c.GetInt(":model", 0)
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.deleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除成功！马上返回中。。。",
			Url:     beego.URLFor("Class.Index", ":model", model),
		})
	}
}

/**************以下为service层内容***********/

/** 获取一条数据 **/
func (c *Class) one(id int) spider.Class {
	var item spider.Class
	_ = orm.NewOrm().QueryTable(new(spider.Class)).Filter("id", id).One(&item)
	return item
}

/** 检测分类是否存在文章 **/
func (c *Class) extArticle(id int) spider.Article {
	var item spider.Article
	_ = orm.NewOrm().QueryTable(new(spider.Article)).Filter("class", id).One(&item)
	return item
}

/** 批量删除结果 **/
func (c *Class) deleteArray(array []int) (message error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if c.one(v).Usage == 0 && c.extArticle(v).Id == 0 {
			if _, message = orm.NewOrm().Delete(&spider.Class{Id: v}); message != nil {
				_ = orm.NewOrm().Rollback()
				break
			}
		} else {
			message = errors.New("被挂载或者存在文档数据，不允许删除使用中的栏目！")
			break
		}
	}
	if message == nil {
		_ = orm.NewOrm().Commit()
	}
	return message
}

/****************** 以下为表格渲染  ***********************/
/** 获取需要渲染的Column **/
func (c *Class) dataTableColumns() []map[string]interface{} {
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
func (c *Class) dataTableButtons(id int) []*_func.TableButtons {
	var array []*_func.TableButtons
	if id > 0 {
		array = append(array, &_func.TableButtons{
			Text:      "返回模型",
			ClassName: "btn btn-sm btn-alt-success mt-1 jump_urls",
			Attribute: map[string]string{"data-action": beego.URLFor("Models.Index")},
		})
		array = append(array, &_func.TableButtons{
			Text:      "添加栏目",
			ClassName: "btn btn-sm btn-alt- mt-1 btn-alt-warning open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Class.Create", ":model", id, ":popup", 1),
				"data-area": "580px,380px",
			},
		})
		array = append(array, &_func.TableButtons{
			Text:      "批量导入",
			ClassName: "btn btn-sm btn-alt-primary mt-1 open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Class.Import", ":model", id, ":popup", 1),
				"data-area": "320px,227px",
			},
		})
	}
	array = append(array, &_func.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Class.Delete", ":model", id)},
	})
	return array
}

/** 处理分页 **/
func (c *Class) pageListItems(length, draw, page int, search string, id int) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(spider.Class))
	if id > 0 {
		qs = qs.Filter("model", id)
	}
	if search != "" {
		qs = qs.Filter("title__icontains", search)
	}
	recordsTotal, _ := qs.Count()
	_, _ = qs.Limit(length, length*(page-1)).OrderBy("-id").ValuesList(&lists, "id", "model", "usage", "title", "update_time")
	for _, v := range lists {
		v[1] = c.modelsService.One(int(v[1].(int64))).Name
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
func (c *Class) tableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "date"},
		"fieldName": {"id", "model", "usage", "title", "update_time"},
		"action":    {"", "", "", "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (c *Class) tableButtonsType(id int) []*_func.TableButtons {
	buttons := []*_func.TableButtons{
		{
			Text:      "编辑栏目",
			ClassName: "btn btn-sm btn-alt-warning open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Class.Edit", ":id", "__ID__", ":popup", 1),
				"data-area": "580px,380px",
			},
		},
		{
			Text:      "删除栏目",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Class.Delete", ":model", id),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}
