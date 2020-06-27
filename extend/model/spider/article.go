package spider

import "time"

type Article struct {
	Id          int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Model       int       `orm:"column(model);type(int);default(0);description(所属模型)" json:"model" form:"model" validate:"required" label:"所属模型"`
	Object      int       `orm:"column(object);type(int);default(0);description(采集结果ID)" json:"object" form:"object" validate:"required" label:"采集结果ID"`
	Usage       int       `orm:"column(usage);type(int);default(0);description(挂载次数)" json:"usage" form:"usage"`
	Class       int       `orm:"column(class);type(int);default(0);description(所属分类)" json:"class" form:"class" validate:"required" label:"所属分类"`
	Title       string    `orm:"column(title);size(32);type(char);default();description(标题)" json:"title" form:"title" validate:"required" label:"标题"`
	Cover       string    `orm:"column(cover);size(200);type(char);default();description(封面图)" json:"cover" form:"cover"`
	Keywords    string    `orm:"column(keywords);size(80);type(char);default();description(关键字)" json:"keywords" form:"keywords" validate:"required" label:"关键字"`
	Description string    `orm:"column(description);size(200);type(char);default();description(描述)" json:"description" form:"description" validate:"required" label:"描述"`
	Context     string    `orm:"column(context);type(text);default();description(内容详情)" json:"context" form:"context" validate:"required" label:"内容详情"`
	CreateTime  time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime  time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *Article) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Article) TableName() string {
	return "spider_article"
}
