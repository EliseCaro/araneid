package admin

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/service"
	"os"
	"strconv"
	"time"
)

type Sign struct {
	controllers.Base
	usersService service.DefaultUsersService
	adminService service.DefaultAdminService
}

/** 实现上一级构造 **/
func (c *Sign) NestPrepare() {
	c.DomainCheck(c.AdminCheck)
}

/** 检测后台域名 **/
func (c *Sign) AdminCheck(prefix, main string) (bool, string, string) {
	adminDomain := beego.AppConfig.String("system_admin_domain")
	if adminDomain == fmt.Sprintf("%s.%s", prefix, main) {
		return true, "success!!", "success!!"
	} else {
		return false, "域名检测失败", "该域名不是有效合法的访问域名；请检查是否授权该域名访问主系统！"
	}
}

// @router /sign/login [get,post]
func (c *Sign) Login() {
	if c.IsAjax() {
		username := c.GetMustString("username", "账号不能为空;请正确填写账号！")
		password := c.GetMustString("password", "密码不能为空;请正确填写密码！")
		remember := c.usersService.IsRemember(c.GetString("remember", "off"))
		IP := _func.Ip2long(c.Ctx.Input.IP())
		if uid, massage := c.usersService.Login(username, password, IP); uid > 0 {
			c.SetSession("uid", uid)
			if remember {
				c.Ctx.SetCookie("remember", strconv.Itoa(uid), 86400*7, "/")
			}
			c.Succeed(&controllers.ResultJson{Message: "登陆系统成功！正在跳转...", Url: beego.URLFor("Admin.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: massage})
		}
	}
}

// @router /sign/quit [get]
func (c *Sign) Quit() {
	c.DelSession("uid")
	c.Ctx.SetCookie("remember", "0", 100, "/")
	c.Ctx.Redirect(302, beego.URLFor("Admin.Index"))
}

// @router /sign/theme [post]
func (c *Sign) Theme() {
	s := c.GetString("theme_style", "")
	if s != "" {
		c.Ctx.SetCookie("theme_style", s, 86400*30*12, "/")
	}
	c.Succeed(&controllers.ResultJson{Message: "success!!"})
}

// @router /sign/layout [post]
func (c *Sign) Layout() {
	s := c.GetString("layout_style", "")
	if s != "" {
		c.Ctx.SetCookie("layout_style", s, 86400*30*12, "/")
	}
	c.Succeed(&controllers.ResultJson{Message: "success!!"})
}

// @router /sign/chart [post]
func (c *Sign) Chart() {
	init, _ := c.GetInt("init", 1)
	if init == 1 {
		_ = _func.SetCache("network_cache", "")
	}
	c.Succeed(&controllers.ResultJson{
		Message: beego.Date(time.Now(), "H:i:s"),
		Data:    c.adminService.DashboardInitialized(),
	})
}

// @router /sign/clear [post]
func (c *Sign) Clear() {
	path := c.GetMustString("file", "清理类型非法！")
	config := beego.AppConfig.String(path)
	if err := os.RemoveAll(config); err == nil {
		if path == "sessionproviderconfig" {
			c.Succeed(&controllers.ResultJson{Message: "会话以全部清空；请重新登录", Url: beego.URLFor("Admin.Index")})
		} else {
			c.Succeed(&controllers.ResultJson{Message: "文件已经全部删除～"})
		}
	} else {
		c.Fail(&controllers.ResultJson{Message: "清理失败！失败原因:" + err.Error()})
	}
}
