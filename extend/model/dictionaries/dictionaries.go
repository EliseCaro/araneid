package dictionaries

import "time"

type DictConfig struct {
	Id    int    `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Name  string `orm:"unique;index;column(name);size(15);type(char);default();description(标识)" json:"name" form:"name"`
	Value string `orm:"column(value);size(255);type(char);default();description(配置值)" json:"value" form:"value"`
}

/** 设置引擎 **/
func (m *DictConfig) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *DictConfig) TableName() string {
	return "dictionaries_config"
}

/*************** 爬虫分类结果 **/

type DictCate struct {
	Id          int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Logs        string    `orm:"column(logs);size(150);type(char);default();description(发布日志)" json:"logs" form:"logs"`
	Name        string    `orm:"column(name);size(15);type(char);default();description(名称)" json:"name" form:"name" validate:"required" label:"分类名称"`
	Keywords    string    `orm:"column(keywords);size(80);type(char);default();description(关键字)" json:"keywords" form:"keywords" validate:"required" label:"分类关键字"`
	Description string    `orm:"column(description);size(200);type(char);default();description(描述)" json:"description" form:"description" validate:"required" label:"分类描述"`
	Source      string    `orm:"column(source);size(150);type(char);default();description(任务地址)" json:"source" form:"source" validate:"required" label:"任务地址"`
	Status      int8      `orm:"index;column(status);type(tinyint);default(0);description(发布状态[-1:失败 0:等待 1:成功])" json:"status" form:"status"`
	Initial     string    `orm:"column(initial);size(5);type(char);default();description(首字母)" json:"initial" form:"initial" validate:"required" label:"分类首字母"`
	CreateTime  time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime  time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *DictCate) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *DictCate) TableName() string {
	return "dictionaries_cate"
}
