package service

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/users"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type DefaultUsersService struct{}

/** 获取是否记住密码状态 **/
func (service *DefaultUsersService) IsRemember(str string) (remember bool) {
	if str == "on" {
		remember = true
	}
	return remember
}

/** 获取一条用户数据 **/
func (service *DefaultUsersService) Find(id int) (user users.Users, err error) {
	user.Id = id
	return user, orm.NewOrm().Read(&user)
}

/** 登录程序 **/
func (service *DefaultUsersService) Login(username, password string, IP uint32) (int, string) {
	var user users.Users
	qs := service.buildUsernameType(orm.NewOrm().QueryTable(new(users.Users)), username)
	if err := qs.Filter("status", 1).One(&user); err == nil && user.Id > 0 {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return 0, "账号或密码错误！"
		}
		return service.AutoLogin(user.Id, IP)
	} else {
		return 0, "未查询到用户或者用户被禁用！"
	}
}

/** 自动登录程序 **/
func (service *DefaultUsersService) AutoLogin(id int, IP uint32) (uid int, err string) {
	user, e := service.Find(id)
	if e != nil || user.Status <= 0 {
		return 0, "用户不存在或者被禁用！"
	}
	if id, message := service.adminCheck(user); id == 0 {
		return id, message
	}
	user.LastLoginIp = int64(IP)
	if _, err := orm.NewOrm().Update(&user); err != nil {
		return 0, "更新IP失败，登陆失败！"
	}
	return user.Id, "登陆成功!"
}

/** 角色鉴权 **/
func (service *DefaultUsersService) adminCheck(user users.Users) (uid int, err string) {
	if user.Id != 1 {
		RolesService := DefaultRolesService{}
		if user.Role == 0 || RolesService.IsAccess(user.Role) == false {
			return 0, "未分配角色或该角色禁止访问！"
		}
	}
	return user.Id, "success"
}

/** 解析登录类型 **/
func (service *DefaultUsersService) buildUsernameType(qs orm.QuerySeter, username string) orm.QuerySeter {
	if match, _ := regexp.MatchString(`/^([a-zA-Z0-9_\.\-])+\@(([a-zA-Z0-9\-])+\.)+([a-zA-Z0-9]{2,4})+$/`, username); match == true {
		return qs.Filter("email", username)
	} else if match, _ := regexp.MatchString(`/^1\d{10}$/`, username); match == true {
		return qs.Filter("mobile", username)
	} else {
		return qs.Filter("username", username)
	}
}

/** 该用户名是否存在 **/
func (service *DefaultUsersService) IsExtendUsername(username string, id int) bool {
	var one users.Users
	_ = orm.NewOrm().QueryTable(new(users.Users)).Filter("username", username).One(&one)
	if one.Id > 0 && one.Id != id {
		return true
	} else {
		return false
	}
}

/** 更新状态 **/
func (service *DefaultUsersService) StatusArray(array []int, status int8) (e error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if v > 1 {
			if _, e = orm.NewOrm().Update(&users.Users{Id: v, Status: status}, "Status"); e != nil {
				_ = orm.NewOrm().Rollback()
				break
			}
		} else {
			e = errors.New("超级管理员不允许被操作！")
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if e == nil {
		_ = orm.NewOrm().Commit()
	}
	return e
}

/** 批量删除 **/
func (service *DefaultUsersService) DeleteArray(array []int) (e error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if v > 1 {
			one, _ := service.Find(v)
			attachment := DefaultAdjunctService{}
			if _, e = orm.NewOrm().Delete(&users.Users{Id: v}); e != nil {
				_ = orm.NewOrm().Rollback()
				break
			}
			_ = attachment.Dec(one.Avatar)
		} else {
			e = errors.New("超级管理员不允许删除！")
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if e == nil {
		_ = orm.NewOrm().Commit()
	}
	return e
}

/************************************************表格渲染机制 ************************************************************/

/** 获取需要渲染的Column **/
func (service *DefaultUsersService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "账号", "name": "username", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "昵称", "name": "nickname", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "邮箱", "name": "email", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "手机", "name": "mobile", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "状态", "name": "status", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "最后登录IP", "name": "last_login_ip", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultUsersService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string", "switch", "longIp"},
		"fieldName": {"id", "username", "nickname", "email", "mobile", "status", "last_login_ip"},
		"action":    {"", "", "", "", "", beego.URLFor("Users.Status"), ""},
	}
	return result
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultUsersService) DataTableButtons() []*table.TableButtons {
	var array []*table.TableButtons
	array = append(array, &table.TableButtons{
		Text:      "添加用户",
		ClassName: "btn btn-sm btn-alt-success mt-1 open_iframe",
		Attribute: map[string]string{
			"href":      beego.URLFor("Users.Create", ":popup", 1),
			"data-area": "580px,370px",
		},
	})
	array = append(array, &table.TableButtons{
		Text:      "启用选中",
		ClassName: "btn btn-sm btn-alt-primary mt-1 ids_enables",
		Attribute: map[string]string{"data-action": beego.URLFor("Users.Status"), "data-field": "status"},
	})
	array = append(array, &table.TableButtons{
		Text:      "禁用选中",
		ClassName: "btn btn-sm btn-alt-warning mt-1 ids_disables",
		Attribute: map[string]string{"data-action": beego.URLFor("Users.Status"), "data-field": "status"},
	})
	array = append(array, &table.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Users.Delete")},
	})
	return array
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultUsersService) TableButtonsType() []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "编辑",
			ClassName: "btn btn-sm btn-alt-primary open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Users.Edit", ":id", "__ID__", ":popup", 1),
				"data-area": "580px,370px",
			},
		},
		{
			Text:      "删除",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Users.Delete"),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}

/** 处理分页 **/
func (service *DefaultUsersService) PageListItems(length, draw, page int, search string) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(users.Users))
	recordsTotal, _ := qs.Count()
	if search != "" {
		qs = qs.Filter("nickname__icontains", search)
	}
	_, _ = qs.Limit(length, length*(page-1)).ValuesList(&lists, "id", "username", "nickname", "email", "mobile", "status", "last_login_ip")
	data := map[string]interface{}{
		"draw":            draw,         // 请求次数
		"recordsFiltered": recordsTotal, // 从多少条里面筛选
		"recordsTotal":    recordsTotal, // 总条数
		"data":            lists,        // 筛选结果
	}
	return data
}
