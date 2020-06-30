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

/** 蜘蛛池关键词管理 **/
type Keyword struct{ Main }

// @router /keyword/index [get,post]
func (c *Keyword) Index() {
	id, _ := c.GetInt(":id", 0)
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.keywordService.PageListItems(length, draw, c.PageNumber(), search, id)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.keywordService.TableColumnsType(), c.keywordService.TableButtonsType(id))
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.keywordService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.keywordService.DataTableButtons(id) {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.Data["action"] = beego.URLFor("Keyword.Index", ":id", id)
	c.Data["info"], _ = c.arachnidService.Find(id)
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
}

// @router /keyword/create [get,post]
func (c *Keyword) Create() {
	arachnid := c.GetMustInt(":arachnid", "非法请求～")
	if c.IsAjax() {
		item := spider.Keyword{Arachnid: arachnid}
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Insert(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "添加关键词成功", Url: beego.URLFor("Keyword.Index", ":id", arachnid)})
		} else {
			c.Fail(&controllers.ResultJson{Message: "添加关键词失败，请稍后再试！"})
		}
	}
	c.Data["arachnid"] = arachnid
}

// @router /keyword/import [get,post]
func (c *Keyword) Import() {
	arachnid := c.GetMustInt(":arachnid", "非法请求～")
	if c.IsAjax() {
		file := c.GetMustInt("files", "请上传正确格式的xlsx文件！")
		path := c.adjunctService.FindId(file).Path
		if f, err := excelize.OpenFile("." + path); err == nil {
			rows, _ := f.GetRows("Sheet1")
			var items []*spider.Keyword
			for _, row := range rows {
				items = append(items, &spider.Keyword{Arachnid: arachnid, Title: row[0]})
			}
			if index, err := orm.NewOrm().InsertMulti(100, items); err == nil {
				c.Succeed(&controllers.ResultJson{Message: "批量导入关键词共" + strconv.FormatInt(index, 10) + "个;正在刷新返回！", Url: beego.URLFor("Keyword.Index", ":id", arachnid)})
			} else {
				c.Fail(&controllers.ResultJson{Message: "批量导入关键词失败；失败原因:" + error.Error(err)})
			}
		} else {
			c.Fail(&controllers.ResultJson{Message: "打开文件失败；失败原因:" + error.Error(err)})
		}
	}
	c.Data["arachnid"] = arachnid
}

// @router /keyword/edit [get,post]
func (c *Keyword) Edit() {
	id := c.GetMustInt(":id", "非法请求！")
	if c.IsAjax() {
		item := c.keywordService.One(id)
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Update(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "更新关键词成功", Url: beego.URLFor("Keyword.Index", ":id", item.Arachnid)})
		} else {
			c.Fail(&controllers.ResultJson{Message: "更新关键词失败，请稍后再试！"})
		}
	}
	c.Data["info"] = c.keywordService.One(id)
}

// @router /keyword/delete [post]
func (c *Keyword) Delete() {
	arachnid, _ := c.GetInt(":arachnid", 0)
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.keywordService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{Message: error.Error(errorMessage)})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除成功！马上返回中。。。",
			Url:     beego.URLFor("Keyword.Index", ":id", arachnid),
		})
	}
}

// @router /keyword/empty [post]
func (c *Keyword) Empty() {
	arachnid, _ := c.GetInt(":arachnid", 0)
	c.keywordService.EmptyDelete(arachnid)
	c.Succeed(&controllers.ResultJson{
		Message: "关键词已经清空！马上返回中。。。",
		Url:     beego.URLFor("Keyword.Index", ":id", arachnid),
	})
}
