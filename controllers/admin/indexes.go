package admin

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/go-playground/validator"
	"strconv"
)

/** 蜘蛛池索引池管理 **/
type Indexes struct{ Main }

// @router /indexes/index [get,post]
func (c *Indexes) Index() {
	id, _ := c.GetInt(":id", 0)
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.indexesService.PageListItems(length, draw, c.PageNumber(), search, id)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.indexesService.TableColumnsType(), c.indexesService.TableButtonsType(id))
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.indexesService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.indexesService.DataTableButtons(id) {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
	c.Data["action"] = beego.URLFor("Indexes.Index", ":id", id)
	c.Data["info"], _ = c.arachnidService.Find(id)
}

// @router /indexes/create [get,post]
func (c *Indexes) Create() {
	arachnid := c.GetMustInt(":arachnid", "非法请求～")
	if c.IsAjax() {
		item := spider.Indexes{Arachnid: arachnid}
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Insert(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "添加索引成功", Url: beego.URLFor("Indexes.Index", ":id", arachnid)})
		} else {
			c.Fail(&controllers.ResultJson{Message: "添加索引失败，请稍后再试！"})
		}
	}
	c.Data["arachnid"] = arachnid
}

// @router /indexes/edit [get,post]
func (c *Indexes) Edit() {
	id := c.GetMustInt(":id", "非法请求！")
	if c.IsAjax() {
		item, _ := c.indexesService.Find(id)
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Update(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "更新索引成功", Url: beego.URLFor("Indexes.Index", ":id", item.Arachnid)})
		} else {
			c.Fail(&controllers.ResultJson{Message: "更新索引失败，请稍后再试！"})
		}
	}
	c.Data["info"], _ = c.indexesService.Find(id)
}

// @router /indexes/delete [post]
func (c *Indexes) Delete() {
	arachnid, _ := c.GetInt(":arachnid", 0)
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.indexesService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{Message: error.Error(errorMessage)})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除成功！马上返回中。。。",
			Url:     beego.URLFor("Indexes.Index", ":id", arachnid),
		})
	}
}

// @router /indexes/empty [post]
func (c *Indexes) Empty() {
	arachnid, _ := c.GetInt(":arachnid", 0)
	c.indexesService.EmptyDelete(arachnid)
	c.Succeed(&controllers.ResultJson{
		Message: "索引已经清空！马上返回中。。。",
		Url:     beego.URLFor("Indexes.Index", ":id", arachnid),
	})
}

// @router /indexes/import [get,post]
func (c *Indexes) Import() {
	arachnid := c.GetMustInt(":arachnid", "非法请求～")
	if c.IsAjax() {
		file := c.GetMustInt("files", "请上传正确格式的xlsx文件！")
		path := c.adjunctService.FindId(file).Path
		if f, err := excelize.OpenFile("." + path); err == nil {
			rows, _ := f.GetRows("Sheet1")
			var items []*spider.Indexes
			for _, row := range rows {
				if len(row) == 2 {
					items = append(items, &spider.Indexes{Arachnid: arachnid, Title: row[0], Urls: row[1]})
				}
			}
			if index, err := orm.NewOrm().InsertMulti(100, items); err == nil {
				c.Succeed(&controllers.ResultJson{Message: "批量导入索引共" + strconv.FormatInt(index, 10) + "个;正在刷新返回！", Url: beego.URLFor("Indexes.Index", ":id", arachnid)})
			} else {
				c.Fail(&controllers.ResultJson{Message: "批量导入索引失败；失败原因:" + error.Error(err)})
			}
		} else {
			c.Fail(&controllers.ResultJson{Message: "打开文件失败；失败原因:" + error.Error(err)})
		}
	}
	c.Data["arachnid"] = arachnid
}
