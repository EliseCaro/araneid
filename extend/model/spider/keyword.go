package spider

import "time"

type Keyword struct {
	Id         int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Title      string    `orm:"column(title);size(200);type(char);default();description(关键词)" json:"title" form:"title" validate:"required" label:"关键词"`
	Arachnid   int       `orm:"column(arachnid);type(int);default(0);description(所属项目)" json:"arachnid" form:"arachnid" validate:"required" label:"所属项目"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *Keyword) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Keyword) TableName() string {
	return "spider_keyword"
}
