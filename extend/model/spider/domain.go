package spider

import "time"

/** 域名挂载缓存池 **/
type Domain struct {
	Id          int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Cate        string    `orm:"column(cate);type(text);default();description(挂载分类)" json:"cate" form:"cate" validate:"required" label:"挂载分类"`
	Name        string    `orm:"column(name);size(32);type(char);default();description(站点名)" json:"name" form:"name" validate:"required" label:"站点名"`
	Title       string    `orm:"column(title);size(150);type(char);default();description(标题)" json:"title" form:"title" validate:"required" label:"标题"`
	Submit      string    `orm:"column(submit);size(80);type(char);default();description(百度自动提交token)" json:"submit" form:"submit"`
	Links       string    `orm:"column(links);type(text);default();description(底部友联)" json:"links" form:"links" form:"description" validate:"required" label:"底部友链"`
	Domain      string    `orm:"column(domain);size(32);type(char);default();description(域名)" json:"domain" form:"domain" validate:"required" label:"域名"`
	Arachnid    int       `orm:"column(arachnid);type(int);default(0);description(所属项目)" json:"arachnid" form:"arachnid" validate:"required" label:"所属项目"`
	Keywords    string    `orm:"column(keywords);size(200);type(char);default();description(关键字)" json:"keywords" form:"keywords" validate:"required" label:"关键字"`
	Description string    `orm:"column(description);size(255);type(char);default();description(描述)" json:"description" form:"description" validate:"required" label:"描述"`
	Template    string    `orm:"column(template);size(32);type(char);default();description(挂载模板)" json:"template" form:"template" validate:"required" label:"挂载模板"`
	Status      int8      `orm:"column(status);type(tinyint);default(0);description(启用状态)" json:"status" form:"status"`
	CreateTime  time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime  time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *Domain) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Domain) TableName() string {
	return "spider_domain"
}
