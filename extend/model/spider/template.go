package spider

import "time"

type Template struct {
	Id         int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Name       string    `orm:"column(name);size(50);type(char);default();description(分组名称)" json:"name" form:"name" validate:"required" label:"分组名称"`
	Remark     string    `orm:"column(remark);size(200);type(char);default();description(分组备注)" json:"remark" form:"remark" validate:"required" label:"分组备注"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *Template) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Template) TableName() string {
	return "spider_template"
}

/*** 模板配置结构 ***/
type TemplateConfig struct {
	Name   string `json:"name"`
	Class  int    `json:"class"`
	Title  string `json:"title"`
	Cover  string `json:"cover"`
	Remark string `json:"remark"`
}

/*** 模板分组结构 ***/
type TemplateClass struct {
	Template
	Child []*TemplateConfig `json:"child"`
}
