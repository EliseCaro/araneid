package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/users"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	Main
	tableIsCheck bool
}

/** 构造函数  **/
func (c *Users) NextPrepare() {
	c.tableIsCheck = true
}

// @router /users/index [get,post]
func (c *Users) Index() {
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.usersService.PageListItems(length, draw, c.PageNumber(), search)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.usersService.TableColumnsType(), c.usersService.TableButtonsType())
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.usersService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.usersService.DataTableButtons() {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, 10)
}

// @router /users/create [get,post]
func (c *Users) Create() {
	if c.IsAjax() {
		item := users.Users{}
		verify := c.GetMustString("verify_password", "请正确填写用户密码")
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if verify != item.Password {
			c.Fail(&controllers.ResultJson{Message: "密码双重验证失败～"})
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
		item.Password = string(hash)
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		isExtend := c.usersService.IsExtendUsername(item.Username, 0)
		if isExtend == true {
			c.Fail(&controllers.ResultJson{Message: "该用户名已经存在！！"})
		}
		if _, err := orm.NewOrm().Insert(&item); err == nil {
			_ = c.adjunctService.Inc(item.Avatar, 0)
			c.Succeed(&controllers.ResultJson{
				Message: "创建用户成功",
				Url:     beego.URLFor("Users.Index"),
			})
		} else {
			c.Fail(&controllers.ResultJson{
				Message: "创建用户失败，请稍后再试！",
			})
		}
	}
	c.Data["roles"] = c.rolesService.RoleAllStatus()
}

// @router /users/edit [get,post]
func (c *Users) Edit() {
	id := c.GetMustInt(":id", "非法请求！")
	if c.IsAjax() {
		item, _ := c.usersService.Find(id)
		oldImage := item.Avatar
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if pwd := c.GetString("password", ""); pwd != "" {
			verify := c.GetString("verify_password", "")
			if verify != pwd {
				c.Fail(&controllers.ResultJson{Message: "密码双重验证失败～"})
			}
			hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
			item.Password = string(hash)
		}
		isExtend := c.usersService.IsExtendUsername(item.Username, item.Id)
		if isExtend == true {
			c.Fail(&controllers.ResultJson{Message: "该用户名已经存在！！"})
		}
		if _, err := orm.NewOrm().Update(&item); err == nil {
			_ = c.adjunctService.Inc(item.Avatar, oldImage)
			c.Succeed(&controllers.ResultJson{
				Message: "修改用户资料成功",
				Url:     beego.URLFor("Users.Index"),
			})
		} else {
			c.Fail(&controllers.ResultJson{
				Message: "修改用户资料失败，请稍后再试！",
			})
		}
	}
	_info, _ := c.usersService.Find(id)
	c.Data["roles"] = c.rolesService.RoleAllStatus()
	c.Data["info"] = _info
	c.Data["avatar"] = c.adjunctService.FindId(_info.Avatar)
}

// @router /users/status [post]
func (c *Users) Status() {
	status, _ := c.GetInt8(":status", 0)
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.usersService.StatusArray(array, status); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "用户状态更新成功！马上返回中。。。",
			Url:     beego.URLFor("Users.Index"),
		})
	}
}

// @router /users/delete [post]
func (c *Users) Delete() {
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.usersService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除用户成功！马上返回中。。。",
			Url:     beego.URLFor("Users.Index"),
		})
	}
}
