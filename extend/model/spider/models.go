package spider

import "time"

type Models struct {
	Id         int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Name       string    `orm:"column(name);size(32);type(char);default();description(模型名称)" json:"name" form:"name" validate:"required" label:"模型名称"`
	Collect    int       `orm:"column(collect);default(0);description(内容爬虫)" json:"collect" form:"collect" validate:"required" label:"分类爬虫"`
	Masking    int       `orm:"column(masking);default(0);description(伪原创分组)" json:"masking" form:"masking" validate:"required" label:"伪原创分组"`
	Template   int       `orm:"column(template);default(0);description(模板分组)" json:"template" form:"template" validate:"required" label:"模板分组"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *Models) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Models) TableName() string {
	return "spider_models"
}
