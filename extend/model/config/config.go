package config

import "time"

type Config struct {
	Id         int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Name       string    `orm:"unique;column(name);size(64);type(char);default();description(标识)" json:"name" form:"name" validate:"required" label:"配置标识"`
	Title      string    `orm:"column(title);size(32);type(char);default();description(标题)" json:"title" form:"title" validate:"required" label:"配置标题"`
	Class      string    `orm:"column(class);size(32);type(char);default();description(配置分组)" json:"class" form:"class" validate:"required" label:"配置分组"`
	Form       string    `orm:"column(form);size(32);type(char);default();description(配置类型)" json:"form" form:"form" validate:"required" label:"配置类型"`
	Value      string    `orm:"column(value);type(text);default();description(配置值)" json:"value" form:"value"`
	Option     string    `orm:"column(option);type(text);default();description(配置项)" json:"option" form:"option"`
	Tips       string    `orm:"column(tips);size(150);type(char);default();description(配置提示)" json:"tips" form:"tips" validate:"required" label:"提示"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
	Sort       int       `orm:"column(sort);default(0);description(排序序号)" json:"sort" form:"sort"`
	Status     int8      `orm:"index;column(status);type(int);default(0);description(启用状态)" json:"status" form:"status"`
}

/** 单条记录详细解析 **/
type OptionFormConfig struct {
	Config
	OptionObject map[string]string `json:"option_object"`
	ValueInt     int               `json:"value_int"`
}

/** 解析表单结构 **/
type FormConfig struct {
	Child []*OptionFormConfig `json:"child"`
	Name  string              `json:"name"`
	Title string              `json:"title"`
}

// 设置引擎为 INNODB
func (m *Config) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Config) TableName() string {
	return "admin_config"
}
