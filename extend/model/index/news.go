package index

import "time"

type News struct {
	Id          int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Cover       int       `orm:"column(cover);type(int);default(0);description(封面图片)" json:"cover" form:"cover" validate:"required,min=1" label:"封面图片"`
	Title       string    `orm:"column(title);size(100);type(char);default();description(标题)" json:"title" form:"title" validate:"required" label:"标题"`
	Status      int8      `orm:"column(status);type(tinyint);default(0);description(是否展示)" json:"status" form:"status"`
	Context     string    `orm:"column(context);type(text);default();description(内容详情)" json:"context" form:"context" validate:"required" label:"内容详情"`
	SeoTitle    string    `orm:"column(seo_title);size(150);type(char);default();description(seo标题)" json:"seo_title" form:"seo_title" validate:"required" label:"SEO标题"`
	Keywords    string    `orm:"column(keywords);size(80);type(char);default();description(关键字)" json:"keywords" form:"keywords" validate:"required" label:"关键字"`
	Description string    `orm:"column(description);type(text);default();description(介绍描述)" json:"description" form:"description" validate:"required" label:"介绍描述"`
	CreateTime  time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime  time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *News) TableEngine() string {
	return "MyISAM"
}

/** 设置表名 **/
func (m *News) TableName() string {
	return "index_news"
}
