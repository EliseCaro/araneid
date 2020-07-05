package spider

import "time"

/** 域名挂载缓存池 **/
type Category struct {
	Id          int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Cid         int       `orm:"column(cid);type(int);default(0);description(原始ID)" json:"cid" form:"cid" validate:"required" label:"原始ID"`
	Model       int       `orm:"column(model);type(int);default(0);description(所属模型)" json:"model" form:"model" validate:"required" label:"所属模型"`
	Name        string    `orm:"column(name);size(32);type(char);default();description(名称)" json:"name" form:"name" validate:"required" label:"名称"`
	Title       string    `orm:"column(title);size(150);type(char);default();description(标题)" json:"title" form:"title" validate:"required" label:"标题"`
	Domain      int       `orm:"column(domain);type(int);default(0);description(挂载域名)" json:"domain" form:"domain" validate:"required" label:"挂载域名"`
	Arachnid    int       `orm:"column(arachnid);type(int);default(0);description(所属项目)" json:"arachnid" form:"arachnid" validate:"required" label:"所属项目"`
	Keywords    string    `orm:"column(keywords);size(200);type(char);default();description(关键字)" json:"keywords" form:"keywords" validate:"required" label:"关键字"`
	Description string    `orm:"column(description);size(255);type(char);default();description(描述)" json:"description" form:"description" validate:"required" label:"描述"`
	CreateTime  time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime  time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *Category) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Category) TableName() string {
	return "spider_category"
}
