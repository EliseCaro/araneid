package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/config"
	"github.com/go-playground/validator"
)

type Config struct {
	Main
}

// @router /config/index [get,post]
func (c *Config) Index() {
	if c.IsAjax() {
		items := make(map[string]string)
		for k, v := range c.Ctx.Request.PostForm {
			items[k] = v[0]
		}
		if err := c.configService.SaveConfig(items); err != nil {
			c.Fail(&controllers.ResultJson{Message: error.Error(err)})
		} else {
			for k, v := range _func.CacheConfig() {
				_ = beego.AppConfig.Set(k, v)
			}
			c.Succeed(&controllers.ResultJson{
				Message: "配置已更新成功～",
				Url:     beego.URLFor("Config.Index"),
			})
		}
	}
	c.Data["class"] = c.configService.ClassBuild()
}

// @router /config/table [get,post]
func (c *Config) Table() {
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.configService.PageListItems(length, draw, c.PageNumber(), search)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.configService.TableColumnsType(), c.configService.TableButtonsType())
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.configService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.configService.DataTableButtons() {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
}

// @router /config/create [get,post]
func (c *Config) Create() {
	if c.IsAjax() {
		item := config.Config{}
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		isExtend := c.configService.IsExtendName(item.Name, 0)
		if isExtend == true {
			c.Fail(&controllers.ResultJson{Message: "该配置标识名已经存在！！"})
		}
		if _, err := orm.NewOrm().Insert(&item); err == nil {
			c.Succeed(&controllers.ResultJson{
				Message: "创建配置项成功",
				Url:     beego.URLFor("Config.Table"),
			})
		} else {
			c.Fail(&controllers.ResultJson{
				Message: "创建配置项失败，请稍后再试！",
			})
		}
	}
	c.Data["class"] = _func.ParseAttrConfigMap(beego.AppConfig.String("system_system_class"))
	c.Data["form"] = _func.ParseAttrConfigMap(beego.AppConfig.String("system_system_form"))
}

// @router /config/edit [get,post]
func (c *Config) Edit() {
	id := c.GetMustInt(":id", "非法请求！")
	if c.IsAjax() {
		item := c.configService.FindOneId(id)
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		isExtend := c.configService.IsExtendName(item.Name, id)
		if isExtend == true {
			c.Fail(&controllers.ResultJson{Message: "该配置标识名已经存在！！"})
		}
		if _, err := orm.NewOrm().Update(&item); err == nil {
			c.Succeed(&controllers.ResultJson{
				Message: "修改配置项成功",
				Url:     beego.URLFor("Config.Table"),
			})
		} else {
			c.Fail(&controllers.ResultJson{
				Message: "修改配置项失败，请稍后再试！",
			})
		}
	}
	c.Data["info"] = c.configService.FindOneId(id)
	c.Data["class"] = _func.ParseAttrConfigMap(beego.AppConfig.String("system_system_class"))
	c.Data["form"] = _func.ParseAttrConfigMap(beego.AppConfig.String("system_system_form"))
}

// @router /config/status [post]
func (c *Config) Status() {
	status, _ := c.GetInt8(":status", 0)
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.configService.StatusArray(array, status); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "配置状态更新成功！马上返回中。。。",
			Url:     beego.URLFor("Config.Table"),
		})
	}
}

// @router /config/delete [post]
func (c *Config) Delete() {
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.configService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除配置成功！马上返回中。。。",
			Url:     beego.URLFor("Config.Table"),
		})
	}
}
