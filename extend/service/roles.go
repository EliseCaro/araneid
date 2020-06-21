package service

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_func "github.com/beatrice950201/araneid/extend/func"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/menus"
	"github.com/beatrice950201/araneid/extend/model/roles"
	"github.com/beatrice950201/araneid/extend/model/users"
	"strconv"
	"strings"
)

type DefaultRolesService struct{}

/** 更新状态 **/
func (service *DefaultRolesService) UpdateStatus(id int, status int8) error {
	item, _ := service.Find(id)
	item.Status = status
	_, err := orm.NewOrm().Update(&item)
	return err
}

/** 更新状态 **/
func (service *DefaultRolesService) UpdateAccess(id int, status int8) error {
	item, _ := service.Find(id)
	item.Access = status
	_, err := orm.NewOrm().Update(&item)
	return err
}

/** 更新状态 **/
func (service *DefaultRolesService) UpdateAccessAndStatus(array []int, status int8, field string) (e error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		switch field {
		case "status":
			if e = service.UpdateStatus(v, status); e != nil {
				_ = orm.NewOrm().Rollback()
				break
			}
		case "access":
			if e = service.UpdateAccess(v, status); e != nil {
				_ = orm.NewOrm().Rollback()
				break
			}
		default:
		}
	}
	if e == nil {
		_ = orm.NewOrm().Commit()
	}
	return e
}

/** 删除节点 **/
func (service *DefaultRolesService) DeleteArray(array []int) (e error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, e = orm.NewOrm().Delete(&roles.Roles{Id: v}); e != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if e == nil {
		_ = orm.NewOrm().Commit()
	}
	return e
}

/** 处理分页 **/
func (service *DefaultRolesService) PageListItems(length, draw, page int, search string) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(roles.Roles))
	recordsTotal, _ := qs.Count()
	if search != "" {
		qs = qs.Filter("name__icontains", search)
	}
	_, _ = qs.Limit(length, length*(page-1)).ValuesList(&lists, "id", "name", "access", "status", "update_time")
	data := map[string]interface{}{
		"draw":            draw,         // 请求次数
		"recordsFiltered": recordsTotal, // 从多少条里面筛选
		"recordsTotal":    recordsTotal, // 总条数
		"data":            lists,        // 筛选结果
	}
	return data
}

/** 获取一条数据 **/
func (service *DefaultRolesService) Find(id int) (m roles.Roles, err error) {
	m.Id = id
	return m, orm.NewOrm().Read(&m)
}

/** 是否可以登录后台 **/
func (service *DefaultRolesService) IsAccess(id int) (status bool) {
	role := roles.Roles{Id: id}
	if err := orm.NewOrm().Read(&role); err == nil && role.Status > 0 && role.Access > 0 {
		status = true
	}
	return status
}

/** 获取节点组 **/
func (service *DefaultRolesService) IdArray(id int) (ids []int) {
	role := roles.Roles{Id: id}
	if err := orm.NewOrm().Read(&role); err == nil {
		_ = json.Unmarshal([]byte(role.MenuAuth), &ids)
	}
	return ids
}

/** 准备组装渲染结构 **/
func (service *DefaultRolesService) BuildJsTree(u users.Users, rid int) (menusAuthTree []*roles.JsTree) {
	var Tree []*roles.JsTree
	var menusService DefaultMenusService
	group, _, _ := menusService.PidAll(0, 1)
	if j, _ := json.Marshal(group); len(j) > 0 {
		_ = json.Unmarshal(j, &Tree)
	}
	for _, v := range Tree {
		if menusService.CheckAuthExport(v.Id, u.Id, u.Role) == true {
			groupAllMenu := menusService.GroupAllMenu()
			marshalMenus := menusService.MarshalMenusExport(groupAllMenu)
			layerChildren := menusService.ToLayerMenusExport(marshalMenus, v.Id)
			v.HtmlTree = service.buildTree(layerChildren, service.IdArray(rid), u.Id, u.Role)
			menusAuthTree = append(menusAuthTree, v)
		}
	}
	return menusAuthTree
}

/** 获得节点授权选择层级 **/
func (service *DefaultRolesService) ToLayerRole(items []*roles.Roles, pid, level int) []*roles.Roles {
	var trees []*roles.Roles
	for _, v := range items {
		if pid == v.Pid {
			titlePrefix := strings.Repeat("&nbsp;", level*8) + "┝ "
			if pid > 0 {
				v.Name = titlePrefix + v.Name
			}
			child := service.ToLayerRole(items, v.Id, level+1)
			trees = append(append(trees, v), child...)
		}
	}
	return trees
}

/** 获取所有角色组 **/
func (service *DefaultRolesService) RoleAll() (r []*roles.Roles) {
	debug := _func.AnalysisDebug()
	cacheTag := "roles_all"
	menusCache := _func.GetCache(cacheTag)
	if menusCache != "" {
		r = menusCache.([]*roles.Roles)
	} else {
		_, _ = orm.NewOrm().QueryTable(new(roles.Roles)).All(&r)
	}
	if debug == false {
		_ = _func.SetCache(cacheTag, r)
	}
	return r
}

/** 获取所有可用角色 **/
func (service *DefaultRolesService) RoleAllStatus() (r []*roles.Roles) {
	debug := _func.AnalysisDebug()
	cacheTag := "roles_all_status"
	menusCache := _func.GetCache(cacheTag)
	if menusCache != "" {
		r = menusCache.([]*roles.Roles)
	} else {
		_, _ = orm.NewOrm().QueryTable(new(roles.Roles)).Filter("status", 1).All(&r)
	}
	if debug == false {
		_ = _func.SetCache(cacheTag, r)
	}
	return r
}

/** 获取需要渲染的Column **/
func (service *DefaultRolesService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "唯一标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "角色名称", "name": "name", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "登陆后台", "name": "access", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "启用状态", "name": "status", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "更新时间", "name": "update_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "角色操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultRolesService) DataTableButtons() []*table.TableButtons {
	var array []*table.TableButtons
	array = append(array, &table.TableButtons{
		Text:      "添加角色",
		ClassName: "btn btn-sm btn-alt-success mt-1 jump_urls",
		Attribute: map[string]string{"data-action": beego.URLFor("Roles.Create")},
	})
	array = append(array, &table.TableButtons{
		Text:      "启用选中",
		ClassName: "btn btn-sm btn-alt-primary mt-1 ids_enables",
		Attribute: map[string]string{"data-action": beego.URLFor("Roles.Status"), "data-field": "status"},
	})
	array = append(array, &table.TableButtons{
		Text:      "禁用选中",
		ClassName: "btn btn-sm btn-alt-warning mt-1 ids_disables",
		Attribute: map[string]string{"data-action": beego.URLFor("Roles.Status"), "data-field": "status"},
	})
	array = append(array, &table.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Roles.Delete")},
	})
	return array
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultRolesService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "switch", "switch", "date"},
		"fieldName": {"id", "name", "access", "status", "datetime"},
		"action":    {"", "", beego.URLFor("Roles.Status"), beego.URLFor("Roles.Status"), ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultRolesService) TableButtonsType() []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "编辑",
			ClassName: "btn btn-sm btn-alt-primary jump_urls",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Roles.Edit", ":id", "__ID__"),
			},
		},
		{
			Text:      "删除",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Roles.Delete"),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}

/**************************************************以下为私有方法************************************/

/** 组装一个生成html的配置参数  **/
func (service *DefaultRolesService) buildTreeMapJson(icon string, selected bool) string {
	j := make(map[string]interface{})
	j["opened"] = true
	j["icon"] = "fa fa-fw fa-" + icon
	j["selected"] = selected
	o, _ := json.Marshal(j)
	return string(o)
}

/** 开始组装渲染结构 **/
func (service *DefaultRolesService) buildTree(m []*menus.LayerMenus, auth []int, uid, role int) (h string) {
	var menusService DefaultMenusService
	for _, v := range m {
		if menusService.CheckAuthExport(v.Id, uid, role) == true {
			oString := service.buildTreeMapJson("fa fa-fw fa-"+v.Icon, _func.InArray(v.Id, auth))
			if v.UrlValue != "" {
				v.UrlValue = " ( " + beego.URLFor(v.UrlValue) + " ) "
			}
			title := v.Title + v.UrlValue
			menuSid := strconv.Itoa(v.Id)
			if len(v.Child) > 0 {
				title += service.buildTree(v.Child, auth, uid, role)
			}
			h += "<li id='" + menuSid + "' data-jstree='" + oString + "'>" + title + "</li>"
		}
	}
	return "<ul>" + h + "</ul>"
}
