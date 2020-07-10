package index

import "time"

type Advantage struct {
	Id          int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Icon        string    `orm:"column(icon);type(text);default();description(介绍图标)" json:"icon" form:"icon" validate:"required" label:"介绍图标"`
	Title       string    `orm:"column(title);size(32);type(char);default();description(介绍标题)" json:"title" form:"title" validate:"required" label:"介绍标题"`
	Status      int8      `orm:"column(status);type(tinyint);default(0);description(是否展示)" json:"status" form:"status"`
	Description string    `orm:"column(description);type(text);default();description(介绍描述)" json:"description" form:"description" validate:"required" label:"介绍描述"`
	CreateTime  time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime  time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *Advantage) TableEngine() string {
	return "MyISAM"
}

/** 设置表名 **/
func (m *Advantage) TableName() string {
	return "index_advantage"
}
