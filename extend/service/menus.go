package service

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/menus"
	"github.com/beatrice950201/araneid/extend/model/users"
	"strconv"
)

type DefaultMenusService struct{}

/** 根据控制器和方法名获取对应节点 **/
func (service *DefaultMenusService) ControllerAndActionMenu(c, a string) menus.Menus {
	var menu menus.Menus
	_ = orm.NewOrm().QueryTable(new(menus.Menus)).Filter("url_value", c+"."+a).Filter("status", 1).Filter("pid__gt", 0).One(&menu)
	return menu
}

/** 获取导航小地图 **/
func (service *DefaultMenusService) BreadcrumbMenu(id int) (items []menus.Menus) {
	o, _ := service.Find(id)
	if o.UrlType == "module_admin" && o.UrlValue != "" {
		o.UrlValue = beego.URLFor(o.UrlValue) //todo 处理URL参数问题
	}
	items = append(items, o)
	if o.Pid > 0 {
		items = append(service.BreadcrumbMenu(o.Pid), o)
	}
	return items
}

/** 获取顶部导航 **/
func (service *DefaultMenusService) HeaderMenus(user users.Users) (items []*menus.Menus) {
	debug := _func.AnalysisDebug()
	cacheTag := "role_header_" + strconv.Itoa(user.Role)
	menusCache := _func.GetCache(cacheTag)
	if menusCache != "" {
		items = menusCache.([]*menus.Menus)
	} else {
		temp, _, _ := service.PidAll(0, 1)
		items = service.checkAuthAll(temp, user)
	}
	if debug == false {
		_ = _func.SetCache(cacheTag, items)
	}
	return items
}

/** 根据pid获取所有数据 **/
func (service *DefaultMenusService) PidAll(pid, status int) (items []*menus.Menus, lens int64, err error) {
	lens, err = orm.NewOrm().QueryTable(new(menus.Menus)).Filter("status", status).Filter("pid", pid).All(&items)
	return
}

/** 获取当前导航高亮 **/
func (service *DefaultMenusService) Breadcrumb(m menus.Menus) (menu menus.BreadcrumbMenus) {
	if m.IsMenu == 0 {
		m, _ = service.Find(m.Pid)
	}
	if j, _ := json.Marshal(m); len(j) > 0 {
		_ = json.Unmarshal(j, &menu)
	}
	menu.Vertex = service.recursionCurrentHeaderMenu(m.Id).Id
	return menu
}

/** 获取一条数据 **/
func (service *DefaultMenusService) Find(id int) (m menus.Menus, err error) {
	m.Id = id
	return m, orm.NewOrm().Read(&m)
}

/** 获取所有符合条件的节点 **/
func (service *DefaultMenusService) SideBarAll() (items []*menus.Menus, err error) {
	debug := _func.AnalysisDebug()
	menusCache := _func.GetCache("side_bar_all")
	if menusCache != "" {
		items = menusCache.([]*menus.Menus)
	} else {
		_, err = orm.NewOrm().QueryTable(new(*menus.Menus)).Filter("status", 1).Filter("is_menu", 1).OrderBy("sort").All(&items)
	}
	if debug == false {
		_ = _func.SetCache("side_bar_all", items)
	}
	return items, err
}

/** 获取左侧菜单 **/
func (service *DefaultMenusService) SideBarMenus(user users.Users, tid int) (item []*menus.LayerMenus) {
	currentParent := service.recursionCurrentHeaderMenu(tid).Id
	debug := _func.AnalysisDebug()
	cacheTag := "role_left_" + strconv.Itoa(user.Role) + "_" + strconv.Itoa(currentParent)
	menusCache := _func.GetCache(cacheTag)
	if menusCache != "" {
		item = menusCache.([]*menus.LayerMenus)
	} else {
		items, _ := service.SideBarAll()
		items = service.checkAuthAll(items, user)
		item = service.toLayerMenus(service.marshalMenus(items), currentParent)
	}
	if debug == false {
		_ = _func.SetCache(cacheTag, item)
	}
	return item
}

/** 检测节点权限导出方法 **/
func (service *DefaultMenusService) CheckAuthExport(id, uid, role int) bool {
	return service.checkAuth(id, uid, role)
}

/** 检测节点权限导出方法 **/
func (service *DefaultMenusService) MarshalMenusExport(items []*menus.Menus) (tree []*menus.LayerMenus) {
	return service.marshalMenus(items)
}

func (service *DefaultMenusService) ToLayerMenusExport(list []*menus.LayerMenus, pid int) (tree []*menus.LayerMenus) {
	return service.toLayerMenus(list, pid)
}

/**  获取所有顶级节点 **/
func (service *DefaultMenusService) PidAllMenus(pid int) (items []*menus.Menus) {
	debug := _func.AnalysisDebug()
	cacheTag := "pid_" + strconv.Itoa(pid) + "_menus"
	menusCache := _func.GetCache(cacheTag)
	if menusCache != "" {
		items = menusCache.([]*menus.Menus)
	} else {
		_, _ = orm.NewOrm().QueryTable(new(menus.Menus)).Filter("pid", pid).OrderBy("-sort").All(&items)
	}
	if debug == false {
		_ = _func.SetCache(cacheTag, items)
	}
	return items
}

//// 获取所有节点
func (service *DefaultMenusService) GroupAllMenu() (list []*menus.Menus) {
	debug := _func.AnalysisDebug()
	menusCache := _func.GetCache("group_all_menus")
	if menusCache != "" {
		list = menusCache.([]*menus.Menus)
	} else {
		_, _ = orm.NewOrm().QueryTable(new(menus.Menus)).OrderBy("sort").All(&list)
	}
	if debug == false {
		_ = _func.SetCache("group_all_menus", list)
	}
	return
}

// 获取树结构HTML渲染列表页
func (service *DefaultMenusService) MenuIndexTreeHtml(list []*menus.Menus, pid, root int) string {
	var html string
	for _, v := range list {
		var disable = ""
		var stringId = strconv.Itoa(v.Id)
		if v.Status == 0 {
			disable = "dd-disable"
		}
		if pid == v.Pid {
			html += `<li class="dd-item dd3-item ` + disable + `" data-id="` + stringId + `">`
			html += service.treeHtml(v, root)
			if childHtml := service.MenuIndexTreeHtml(list, v.Id, root); childHtml != "" {
				html += "<ol class='dd-list'>" + childHtml + "</ol>"
			}
			html += "</li>"
		}
	}
	return html
}

///** 返回URL类型MAP **/
func (service *DefaultMenusService) UrlsType() map[string]string {
	return map[string]string{
		"module_admin":  "模块连接",
		"module_others": "外部链接",
	}
}

///** 获取节点增加修改动态树 **/
func (service *DefaultMenusService) TreeRender(id int) (tree []*menus.LayerMenus) {
	groupAll := service.GroupAllMenu()
	marshal := service.marshalMenus(groupAll)
	return service.toLayerMenus(marshal, id)
}

/**************************************************以下为私有方法************************************/

/** 节点鉴权[ 单个 ] **/
func (service *DefaultMenusService) checkAuth(id, uid, role int) (s bool) {
	if uid == 1 {
		s = true
	}
	RolesService := DefaultRolesService{}
	temp := RolesService.IdArray(role)
	if _func.InArray(id, temp) {
		s = true
	}
	return s
}

// 转换结构体
func (service *DefaultMenusService) marshalMenus(items []*menus.Menus) (tree []*menus.LayerMenus) {
	if jsonData, _ := json.Marshal(items); len(jsonData) > 0 {
		_ = json.Unmarshal(jsonData, &tree)
	}
	return
}

/** 转换为树结构 **/
func (service *DefaultMenusService) toLayerMenus(list []*menus.LayerMenus, pid int) (tree []*menus.LayerMenus) {
	for _, v := range list {
		if v.Pid == pid {
			if child := service.toLayerMenus(list, v.Id); len(child) > 0 {
				v.Child = child
			}
			tree = append(tree, v)
		}
	}
	return tree
}

/** 根据 ID 递归找上一层 一直找到pid为0 **/
func (service *DefaultMenusService) recursionCurrentHeaderMenu(id int) menus.Menus {
	var menu menus.Menus
	_ = orm.NewOrm().QueryTable(new(menus.Menus)).Filter("status", 1).Filter("id", id).One(&menu)
	if menu.Pid > 0 {
		return service.recursionCurrentHeaderMenu(menu.Pid)
	}
	return menu
}

/** 节点鉴权[ 列表 ] **/
func (service *DefaultMenusService) checkAuthAll(items []*menus.Menus, user users.Users) (resource []*menus.Menus) {
	for _, v := range items {
		if service.checkAuth(v.Id, user.Id, user.Role) == true {
			if v.UrlType == "module_admin" && v.UrlValue != "" {
				v.UrlValue = beego.URLFor(v.UrlValue) //todo 处理URL参数问题
			}
			resource = append(resource, v)
		}
	}
	return resource
}

/** 获得单个节点树html **/
func (service *DefaultMenusService) treeHtml(m *menus.Menus, root int) (html string) {
	stringId := strconv.Itoa(m.Id)
	createUrls := beego.URLFor("Menu.Create", ":id", m.Id, ":root", root, ":popup", 1)
	updateUrls := beego.URLFor("Menu.Edit", ":id", m.Id, ":root", root, ":popup", 1)
	statusUrls := beego.URLFor("Menu.Status", ":root", root)
	deleteUrls := beego.URLFor("Menu.Delete", ":root", root)
	html += `<div class='dd-handle dd3-handle'>拖拽</div>`
	html += `<div class='dd3-content'>`
	html += `<i class='fa fa-` + m.Icon + `'></i> ` + m.Title
	if len([]byte(m.UrlValue)) > 0 {
		html += `<span class="link"><i class="fa fa-link"></i> ` + beego.URLFor(m.UrlValue) + `</span>`
	}
	html += `<div class='action'>`
	html += `<a href="` + createUrls + `" data-toggle='tooltip' data-original-title='新增子节点' data-area="580px,395px" class='open_iframe'>`
	html += `<i class='list-icon fa fa-plus fa-fw'></i>`
	html += `</a>`
	html += `<a href="` + updateUrls + `" data-toggle='tooltip' data-original-title='编辑节点' data-area="580px,395px" class='open_iframe'>`
	html += `<i class='list-icon fa fa-pencil-alt fa-fw'></i>`
	html += `</a>`
	if m.Status <= 0 {
		html += `<a href='javascript:void(0);' data-action="` + statusUrls + `" data-ids="` + stringId + `" class='ids_enable' data-toggle='tooltip' data-original-title='启用'>`
		html += `<i class='list-icon fa fa-check-circle fa-fw'></i>`
		html += `</a>`
	} else {
		html += `<a href='javascript:void(0);' data-action="` + statusUrls + `" data-ids="` + stringId + `" class='ids_disable' data-toggle='tooltip' data-original-title='禁用'>`
		html += `<i class='list-icon fa fa-ban fa-fw'></i>`
		html += `</a>`
	}
	html += `<a href='javascript:void(0);' data-action="` + deleteUrls + `" data-ids=` + stringId + ` class='ids_delete' data-toggle='tooltip' data-original-title='删除'>`
	html += `<i class='list-icon fa fa-times fa-fw'></i>`
	html += `</a>`
	html += `</div>`
	html += `</div>`
	return html
}
