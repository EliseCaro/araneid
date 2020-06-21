package roles

import (
	ms "github.com/beatrice950201/araneid/extend/model/menus"
	"time"
)

type Roles struct {
	Id          int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Pid         int       `orm:"index;column(pid);type(int);default(0);description(上级ID)"  json:"pid" form:"pid"`
	Name        string    `orm:"column(name);size(32);type(char);default();description(分组名称)"  json:"name" form:"name" validate:"required" label:"角色名称"`
	Access      int8      `orm:"index;column(access);type(tinyint);default(0);description(是否可以登陆后台)" json:"access" form:"access"`
	Status      int8      `orm:"index;column(status);type(tinyint);default(0);description(启用状态)" json:"status" form:"status"`
	MenuAuth    string    `orm:"column(menu_auth);type(text);default();description(节点字符串)"  json:"menu_auth" form:"menu_auth" validate:"required" label:"节点授权"`
	Description string    `orm:"column(description);size(150);type(char);default();description(分组描述)"  json:"description" form:"description" validate:"required" label:"角色描述"`
	CreateTime  time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime  time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *Roles) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Roles) TableName() string {
	return "admin_role"
}

/** 创建编辑节点授权渲染 **/
type JsTree struct {
	ms.Menus
	HtmlTree string `json:"html_tree"`
}
