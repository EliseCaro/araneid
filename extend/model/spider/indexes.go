package spider

import "time"

/** 站群索引池 **/
type Indexes struct {
	Id         int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Urls       string    `orm:"column(urls);size(255);type(char);default();description(链接地址)" json:"urls" form:"urls" validate:"required" label:"链接地址"`
	Title      string    `orm:"column(title);size(150);type(char);default();description(链接标题)" json:"title" form:"title" validate:"required" label:"链接标题"`
	Arachnid   int       `orm:"column(arachnid);type(int);default(0);description(所属项目)" json:"arachnid" form:"arachnid" validate:"required" label:"所属项目"`
	Sort       int       `orm:"column(sort);default(0);description(链接权重)" json:"sort" form:"sort"`
	Usage      int       `orm:"column(usage);type(int);default(0);description(挂载次数)" json:"usage" form:"usage"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *Indexes) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Indexes) TableName() string {
	return "spider_indexes"
}
