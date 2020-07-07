package controllers

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	_func "github.com/beatrice950201/araneid/extend/func"
	"strconv"
	"strings"
	"time"
)

type Base struct {
	beego.Controller
	DomainMain   string
	DomainPrefix string
	Module       string
}

/*  子级构造函数  */
type NestPreparer interface {
	NestPrepare()
}

/** JSON 返回格式 */
type ResultJson struct {
	Code    int         `json:"code"`
	Status  bool        `json:"status"`
	Message string      `json:"message,omitempty"`
	Url     string      `json:"url,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

/** 构造行数 **/
func (c *Base) Prepare() {
	c.setStaticVersions()
	c.SetTemplate(c.extractModuleName())
	if app, ok := c.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
}

/** 域名鉴权 **/
func (c *Base) DomainCheck(callback func(string, string) (bool, string, string)) {
	domain := strings.Split(c.Ctx.Input.Domain(), ".")
	if len(domain) == 2 || len(domain) == 3 {
		var (
			is      bool
			title   string
			context string
		)
		if len(domain) == 2 {
			is, title, context = callback("", fmt.Sprintf("%s.%s", domain[0], domain[1]))
		} else {
			is, title, context = callback(domain[0], fmt.Sprintf("%s.%s", domain[1], domain[2]))
		}
		if is == false {
			c.Abort500(title, context)
		}
	} else {
		c.Abort500("域名解析错误", "域名不合法！请检查您的域名解析配置流程以及域名格式;")
	}
}

/** 自动设置模板路径 **/
func (c *Base) SetTemplate(module string) {
	controller, action := c.GetControllerAndAction()
	c.TplName = module + "/" + strings.ToLower(controller) + "/" + strings.ToLower(action) + ".html"
}

/** 提取模块名称 **/
func (c *Base) extractModuleName() string {
	URI := c.Ctx.Request.RequestURI
	URIS := strings.Split(URI, "/")
	if URIS[1] != beego.AppConfig.String("admin_name") {
		return "index"
	}
	return beego.AppConfig.String("admin_name")
}

/** 设置资源版本号码 **/
func (c *Base) setStaticVersions() {
	debug := _func.AnalysisDebug()
	versions := beego.AppConfig.String("web_versions")
	if debug == true {
		versions = strconv.FormatInt(time.Now().Unix(), 10)
	}
	c.Data["versions"] = versions
}

/** 成功返回函数 **/
func (c *Base) Succeed(o *ResultJson) {
	o.Status = true
	c.Json(o)
}

/** 失败返回函数 **/
func (c *Base) Fail(o *ResultJson) {
	o.Status = false
	c.Json(o)
}

/** 返回 **/
func (c *Base) Json(o *ResultJson) {
	c.Data["json"] = o
	c.ServeJSON()
	c.Finish()
	c.StopRun()
}

/** 检测字符串参数返回 **/
func (c *Base) GetMustString(key string, message string) string {
	str := c.GetString(key, "")
	if len(str) == 0 {
		c.Fail(&ResultJson{Message: message})
	}
	return str
}

/** 检测INT参数返回 **/
func (c *Base) GetMustInt(key string, message string) int {
	str, _ := c.GetInt(key, 0)
	if str == 0 {
		c.Fail(&ResultJson{Message: message})
	}
	return str
}

/** 手动终止并且抛出错误 **/
func (c *Base) Abort500(title, content string) {
	if c.IsAjax() {
		c.Ctx.Output.Status = 200
		c.Fail(&ResultJson{Message: content})
	} else {
		c.Data["content"] = errors.New(content)
		c.Data["title"] = errors.New(title)
		c.Abort("500")
	}
}
