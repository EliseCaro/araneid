package menus

import (
	"time"
)

type Menus struct {
	Id         int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Pid        int       `orm:"index;column(pid);type(int);default(0);description(上级ID)"  json:"pid" form:"pid"`
	Title      string    `orm:"column(title);size(32);type(char);default();description(节点标题)" json:"title" form:"title" validate:"required" label:"节点标题"`
	Icon       string    `orm:"column(icon);size(64);type(char);default();description(图标)" json:"icon" form:"icon" validate:"required" label:"节点图标"`
	Sort       int       `orm:"column(sort);default(0);description(排序序号)" json:"sort" form:"sort"`
	Status     int8      `orm:"index;column(status);size(1);type(int);default(0);description(启用状态)" json:"status" form:"status"`
	Params     string    `orm:"column(params);size(255);type(char);default();description(附加参数)" json:"params" form:"params"`
	IsMenu     int8      `orm:"column(is_menu);type(int);default(1);description(是否左侧菜单)" json:"is_menu" form:"is_menu"`
	UrlType    string    `orm:"column(url_type);size(16);type(char);default();description(连接类型)" json:"url_type" form:"url_type"`
	UrlValue   string    `orm:"column(url_value);size(255);type(char);default();description(连接值)" json:"url_value" form:"url_value"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *Menus) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Menus) TableName() string {
	return "admin_menu"
}

/** 树结构 **/
type LayerMenus struct {
	Menus
	Child []*LayerMenus `json:"child" form:"children"`
}

/** 高亮结构 **/
type BreadcrumbMenus struct {
	Menus
	Vertex int `json:"vertex" form:"vertex"`
}
