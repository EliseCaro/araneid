package admin

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/users"
	"github.com/beatrice950201/araneid/extend/service"
	"strconv"
)

type Main struct {
	controllers.Base
	UserInfo            users.Users
	controllerName      string
	actionName          string
	tableBuilder        table.BuilderTable
	verifyBase          service.DefaultBaseVerify
	menusService        service.DefaultMenusService
	usersService        service.DefaultUsersService
	rolesService        service.DefaultRolesService
	adjunctService      service.DefaultAdjunctService
	collectService      service.DefaultCollectService
	collectorService    service.DefaultCollectorService
	configService       service.DefaultConfigService
	informService       service.DefaultInformService
	dictionariesService service.DefaultDictionariesService
}

/** 准备下一级构造函数 **/
type NextPreparer interface {
	NextPrepare()
}

/** 实现上一级构造 **/
func (c *Main) NestPrepare() {
	if app, ok := c.AppController.(NextPreparer); ok {
		app.NextPrepare()
	}
	c.controllerName, c.actionName = c.GetControllerAndAction()
	thisMenuRole := c.menusService.ControllerAndActionMenu(c.controllerName, c.actionName)
	c.isLogin()
	c.checkRoleMenu()
	c.themeBegin()
	c.setLayout()
	c.Data["popup"], _ = c.GetInt(":popup", 0)
	if c.Data["popup"] == 0 && c.IsAjax() != true {
		c.Data["users"] = c.UserInfo
		c.Data["breadcrumb_menus"] = c.menusService.Breadcrumb(thisMenuRole)
		c.Data["header_menus"] = c.menusService.HeaderMenus(c.UserInfo)
		c.Data["side_bar_menus"] = c.menusService.SideBarMenus(c.UserInfo, thisMenuRole.Id)
		c.Data["breadcrumb_map"] = c.menusService.BreadcrumbMenu(thisMenuRole.Id)
		c.Data["inform_items"] = c.informService.HeaderItems()
	}
}

/** 提交表格构建渲染 **/
func (c *Main) TableColumnsRender(clo []*table.ColumnsItems, order [][]string, btn []*table.TableButtons, length int) {
	tableBuilderColumns, _ := json.Marshal(clo)
	tableBuilderButtons, _ := json.Marshal(btn)
	tableBuilderOrder, _ := json.Marshal(order)
	c.Data["table_builder_columns"] = string(tableBuilderColumns)
	c.Data["table_builder_check"] = string(tableBuilderOrder)
	c.Data["table_builder_buttons"] = string(tableBuilderButtons)
	c.Data["table_page_size_length"] = length
}

/** 主题初始化 **/
func (c *Main) themeBegin() {
	style := c.Ctx.GetCookie("theme_style")
	layout := c.Ctx.GetCookie("layout_style")
	c.Data["theme_style"] = style
	c.Data["layout_style"] = layout
}

/** 扩展模板处理 **/
func (c *Main) setLayout() {
	c.Layout = beego.AppConfig.String("admin_base_html")
	cn, an := c.GetControllerAndAction()
	r := _func.LayoutSections(cn, an, c.Ctx.Request)
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["header"] = r["header"]
	c.LayoutSections["footer"] = r["footer"]
}

/** 是否登陆处理 **/
func (c *Main) isLogin() {
	uid := c.GetSession("uid")
	if uid == nil || uid.(int) <= 0 {
		if c.autoLogin() {
			c.Ctx.Redirect(302, beego.URLFor("Admin.Index"))
		} else {
			c.Ctx.Redirect(302, beego.URLFor("Sign.Login"))
		}
	} else {
		c.UserInfo, _ = c.usersService.Find(uid.(int))
	}
}

/** 自动登录程序 **/
func (c *Main) autoLogin() bool {
	uid := c.Ctx.GetCookie("remember")
	status := false
	if intUID, err := strconv.Atoi(uid); err == nil {
		IP := _func.Ip2long(c.Ctx.Input.IP())
		if intUID > 0 {
			if AutoStatus, _ := c.usersService.AutoLogin(intUID, IP); AutoStatus > 0 {
				c.SetSession("uid", AutoStatus)
				status = true
			}
		}
	}
	return status
}

/** 检测节点是否允许访问 **/
func (c *Main) checkRoleMenu() {
	menu := c.menusService.ControllerAndActionMenu(c.controllerName, c.actionName)
	if menu.Id <= 0 {
		c.Abort500("500", "未检测到节点[ "+c.controllerName+"."+c.actionName+" ]")
	}
	isBool := c.menusService.CheckAuthExport(menu.Id, c.UserInfo.Id, c.UserInfo.Role)
	if isBool == false {
		c.Abort500("500", "您的权限无法访问当前节点")
	}
}

/** 解析ID为数组 */
func (c *Main) checkBoxIds(arrayKey, key string) []int {
	array := c.GetStrings(arrayKey, []string{})
	if len(array) <= 0 {
		id := c.GetMustInt(key, "非法操作，请稍后再试...")
		array = append(array, strconv.Itoa(id))
	}
	var result []int
	for _, v := range array {
		int64String, _ := strconv.Atoi(v)
		result = append(result, int64String)
	}
	return result
}

/** 获得当前页码 **/
func (c *Main) PageNumber() int {
	page, _ := c.GetInt("page", 0)
	if page <= 0 {
		start, _ := c.GetInt("start", 0)
		length, _ := c.GetInt("length", 0)
		page = (start / length) + 1
	}
	return page
}