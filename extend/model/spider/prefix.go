package spider

import "time"

type Prefix struct {
	Id         int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Title      string    `orm:"column(title);size(32);type(char);default();description(前缀备注)" json:"title" form:"title"`
	Tags       string    `orm:"column(tags);size(32);type(char);default();description(前缀标签)" json:"tags" form:"tags"`
	Model      int       `orm:"column(model);type(int);default(0);description(挂载模型)" json:"model" form:"model" validate:"required" label:"挂载模型"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *Prefix) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Prefix) TableName() string {
	return "spider_prefix"
}
