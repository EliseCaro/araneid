package _func

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

type BuilderTable struct {
	ColumnsItemsMaps []*ColumnsItems
	ButtonsItemsMaps []*TableButtons
	OrderItemsMaps   [][]string
}

/** 构建Columns结构  **/
type ColumnsItems struct {
	Title     string `json:"title"`
	Name      string `json:"name"`
	Orderable bool   `json:"orderable"`
	ClassName string `json:"className"`
}

/** 按钮结构组 **/
type TableButtons struct {
	Attribute map[string]string `json:"attr"`
	ClassName string            `json:"className"`
	Text      string            `json:"text"`
}

/** 构建Columns  **/
func (builder *BuilderTable) TableColumns(title, name, className string, order bool) *BuilderTable {
	builder.ColumnsItemsMaps = append(builder.ColumnsItemsMaps, &ColumnsItems{Title: title, Name: name, Orderable: order, ClassName: className})
	return builder
}

/** 排序开始位  **/
func (builder *BuilderTable) TableOrder(a []string) *BuilderTable {
	builder.OrderItemsMaps = append(builder.OrderItemsMaps, a)
	return builder
}

/** 渲染按钮 **/
func (builder *BuilderTable) TableButtons(btn *TableButtons) *BuilderTable {
	builder.ButtonsItemsMaps = append(builder.ButtonsItemsMaps, btn)
	return builder
}

/***************************************** 以上为dataTable构建器 *****************************************/

/** 组装字段为dataTable格式  **/
func (builder *BuilderTable) Field(all bool, maps []orm.ParamsList, columnsType map[string][]string, tableButtonsType []*TableButtons) []orm.ParamsList {
	var result []orm.ParamsList
	for _, v := range maps {
		var res []interface{}
		value := builder.FieldAnalysisValue(v, columnsType["fieldName"])
		if all == true {
			res = append(res, builder.checkAll(value))
		}
		v = builder.FieldAnalysisType(v, columnsType)
		result = append(result, append(res, append(v, builder.handle(value, tableButtonsType))...))
	}
	return result
}

/** 重组字段值 **/
func (builder *BuilderTable) FieldAnalysisValue(value orm.ParamsList, f []string) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range value {
		res[f[k]] = v
	}
	return res
}

/** 解析字段值类型 **/
func (builder *BuilderTable) FieldAnalysisType(maps orm.ParamsList, columnsType map[string][]string) orm.ParamsList {
	var (
		columns     = columnsType["columns"]
		fieldNames  = columnsType["fieldName"]
		actionUrls  = columnsType["action"]
		fieldValues = builder.FieldAnalysisValue(maps, fieldNames)
	)
	for k, fieldName := range columns {
		switch fieldName {
		case "switch":
			maps[k] = builder.AnalysisTypeSwitch(fieldValues, fieldNames[k], actionUrls[k], map[int64]string{})
		case "date":
			maps[k] = beego.Date(maps[k].(time.Time), "Y-m-d H:i:s")
		case "longIp":
			maps[k] = Long2IP(uint32(maps[k].(int64)))
		case "fileSize":
			maps[k] = maps[k].(int64) / 1024
		case "fileDriver":
			driver := map[string]string{"local": "本地驱动", "cloud": "远程驱动"}
			maps[k] = driver[maps[k].(string)]
		default:
			maps[k] = maps[k]
		}
	}
	return maps
}

/** 获取全选状态  **/
func (builder *BuilderTable) checkAll(maps map[string]interface{}) string {
	checkId := `check_` + strconv.FormatInt(maps["id"].(int64), 10)
	html := `<div class="custom-control custom-checkbox d-inline-block">`
	html += `<input data-id="` + strconv.FormatInt(maps["id"].(int64), 10) + `" type="checkbox" class="custom-control-input" id="` + checkId + `" name="check[]">`
	html += `<label class="custom-control-label" for="` + checkId + `"></label>`
	html += `</div>`
	return html
}

/** 解析switch类型到实体 **/
func (builder *BuilderTable) AnalysisTypeSwitch(fieldValues map[string]interface{}, fieldName, actionUrls string, option map[int64]string) string {
	var (
		fieldVl   = fieldValues[fieldName].(int64)
		tableId   = strconv.FormatInt(fieldValues["id"].(int64), 10)
		titleName = "禁用中"
		className = ""
	)
	if len(option) > 0 {
		titleName = option[fieldVl]
	} else {
		if fieldVl > 0 {
			titleName = "启用中"
		}
	}
	if fieldVl <= 0 {
		className = "btn-alt-danger ids_enable"
	} else {
		className = "btn-primary ids_disable"
	}
	return `<button type="button" data-field="` + fieldName + `" data-action="` + actionUrls + `" data-ids="` + tableId + `" class="btn btn-sm ` + className + `">` + titleName + `</button>`
}

/** 获取操作按钮  **/
func (builder *BuilderTable) handle(value map[string]interface{}, tableButtonsType []*TableButtons) (s string) {
	for _, v := range tableButtonsType {
		attribute := builder.attributeString(value["id"].(int64), v.Attribute)
		s += `<button type="button" ` + attribute + ` class="` + v.ClassName + `">` + v.Text + `</button> `
	}
	return
}

/** 结构属性转为字符串 **/
func (builder *BuilderTable) attributeString(id int64, maps map[string]string) string {
	attribute := ""
	for k, v := range maps {
		attribute += k + `=` + `"` + v + `"`
	}
	return strings.Replace(attribute, "__ID__", strconv.FormatInt(id, 10), -1)
}
