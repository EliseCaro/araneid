package spider

import "time"

type Disguise struct {
	Id          int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Name        string    `orm:"column(name);size(50);type(char);default();description(机器名称)" json:"name" form:"name" validate:"required" label:"机器名称"`
	Usage       int       `orm:"column(usage);type(int);default(0);description(挂载次数)" json:"usage" form:"usage"`
	Keyword     int8      `orm:"column(keyword);type(tinyint);default(0);description(自然单词)" json:"keyword" form:"keyword"`
	Description int8      `orm:"column(description);type(tinyint);default(0);description(自然描述)" json:"description" form:"description"`
	Context     int8      `orm:"column(context);type(tinyint);default(0);description(自然内容)" json:"context" form:"context"`
	ApiKey      string    `orm:"column(api_key);size(150);type(char);default();description(处理KEY)" json:"api_key" form:"api_key" validate:"required" label:"机器KEY"`
	ApiSecret   string    `orm:"column(api_secret);size(150);type(char);default();description(处理SECRET)" json:"api_secret" form:"api_secret" validate:"required" label:"机器密钥"`
	CreateTime  time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime  time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *Disguise) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Disguise) TableName() string {
	return "spider_disguise"
}
