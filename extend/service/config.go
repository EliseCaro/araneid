package service

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_func "github.com/beatrice950201/araneid/extend/func"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/config"
	"strconv"
)

type DefaultConfigService struct{}

/** 根据标识获取详细 **/
func (service *DefaultConfigService) findOneName(name string) config.Config {
	var item config.Config
	_ = orm.NewOrm().QueryTable(new(config.Config)).Filter("name", name).One(&item)
	return item
}

/** 根据ID获取详细 **/
func (service *DefaultConfigService) FindOneId(id int) config.Config {
	var item config.Config
	_ = orm.NewOrm().QueryTable(new(config.Config)).Filter("id", id).One(&item)
	return item
}

/** 解析分组数据 **/
func (service *DefaultConfigService) ClassBuild() []*config.FormConfig {
	var res []*config.FormConfig
	class := _func.ParseAttrConfigMap(beego.AppConfig.String("system_system_class"))
	for k, v := range class {
		res = append(res, &config.FormConfig{Child: service.ClassItems(k), Name: k, Title: v})
	}
	return res
}

/** 根据分组名称获取列表 **/
func (service *DefaultConfigService) ClassItems(class string) []*config.OptionFormConfig {
	var (
		list   []*config.Config
		result []*config.OptionFormConfig
	)
	_, _ = orm.NewOrm().QueryTable(new(config.Config)).Filter("status", 1).OrderBy("sort").Filter("class", class).All(&list)
	if j, _ := json.Marshal(list); len(j) > 0 {
		_ = json.Unmarshal(j, &result)
	}
	for _, v := range result {
		if v.Form == "radio" {
			v.OptionObject = _func.ParseAttrConfigMap(v.Option)
		}
		if v.Form == "image" {
			v.ValueInt, _ = strconv.Atoi(v.Value)
		}
	}
	return result
}

/** 保存配置 **/
func (service *DefaultConfigService) SaveConfig(items map[string]string) (err error) {
	_ = orm.NewOrm().Begin()
	for k, v := range items {
		one := service.findOneName(k)
		if one.Form == "image" { // 挂载图片
			oldImage, _ := strconv.Atoi(one.Value)
			newImage, _ := strconv.Atoi(v)
			adjunctService := DefaultAdjunctService{}
			_ = adjunctService.Inc(newImage, oldImage)
		}
		if _, err = orm.NewOrm().Update(&config.Config{Id: one.Id, Value: v}, "Value"); err != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if err == nil {
		_ = orm.NewOrm().Commit()
	}
	return err
}

/** 该标识是否存在 **/
func (service *DefaultConfigService) IsExtendName(name string, id int) bool {
	var one config.Config
	_ = orm.NewOrm().QueryTable(new(config.Config)).Filter("name", name).One(&one)
	if one.Id > 0 && one.Id != id {
		return true
	} else {
		return false
	}
}

/** 更新状态 **/
func (service *DefaultConfigService) StatusArray(array []int, status int8) (e error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, e = orm.NewOrm().Update(&config.Config{Id: v, Status: status}, "Status"); e != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if e == nil {
		_ = orm.NewOrm().Commit()
	}
	return e
}

/** 批量删除 **/
func (service *DefaultConfigService) DeleteArray(array []int) (e error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, e = orm.NewOrm().Delete(&config.Config{Id: v}); e != nil {
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
func (service *DefaultConfigService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "标识", "name": "id", "className": "text-center", "order": true})
	maps = append(maps, map[string]interface{}{"title": "名称", "name": "name", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "标题", "name": "title", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "分组", "name": "class", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "类型", "name": "form", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "状态", "name": "status", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultConfigService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string", "switch"},
		"fieldName": {"id", "name", "title", "class", "form", "status"},
		"action":    {"", "", "", "", "", beego.URLFor("Config.Status"), ""},
	}
	return result
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultConfigService) DataTableButtons() []*table.TableButtons {
	var array []*table.TableButtons
	array = append(array, &table.TableButtons{
		Text:      "添加配置",
		ClassName: "btn btn-sm btn-alt-success mt-1 open_iframe",
		Attribute: map[string]string{
			"href":      beego.URLFor("Config.Create", ":popup", 1),
			"data-area": "580px,370px",
		},
	})
	array = append(array, &table.TableButtons{
		Text:      "启用选中",
		ClassName: "btn btn-sm btn-alt-primary mt-1 ids_enables",
		Attribute: map[string]string{"data-action": beego.URLFor("Config.Status"), "data-field": "status"},
	})
	array = append(array, &table.TableButtons{
		Text:      "禁用选中",
		ClassName: "btn btn-sm btn-alt-warning mt-1 ids_disables",
		Attribute: map[string]string{"data-action": beego.URLFor("Config.Status"), "data-field": "status"},
	})
	array = append(array, &table.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Config.Delete")},
	})
	return array
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultConfigService) TableButtonsType() []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "编辑",
			ClassName: "btn btn-sm btn-alt-primary open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Config.Edit", ":id", "__ID__", ":popup", 1),
				"data-area": "580px,370px",
			},
		},
		{
			Text:      "删除",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Config.Delete"),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}

/** 处理分页 **/
func (service *DefaultConfigService) PageListItems(length, draw, page int, search string) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(config.Config))
	recordsTotal, _ := qs.Count()
	if search != "" {
		qs = qs.Filter("title__icontains", search)
	}
	_, _ = qs.Limit(length, length*(page-1)).ValuesList(&lists, "id", "name", "title", "class", "form", "status")
	for _, v := range lists {
		v[3] = service.returnClass(v[3].(string))
		v[4] = service.returnType(v[4].(string))
	}
	data := map[string]interface{}{
		"draw":            draw,         // 请求次数
		"recordsFiltered": recordsTotal, // 从多少条里面筛选
		"recordsTotal":    recordsTotal, // 总条数
		"data":            lists,        // 筛选结果
	}
	return data
}

/** 解析分组 **/
func (service *DefaultConfigService) returnClass(t string) string {
	class := _func.ParseAttrConfigMap(beego.AppConfig.String("system_system_class"))
	return class[t]
}

/** 解析类型标题 **/
func (service *DefaultConfigService) returnType(t string) string {
	form := _func.ParseAttrConfigMap(beego.AppConfig.String("system_system_form"))
	return form[t]
}
