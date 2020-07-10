package index

type Team struct {
	Id          int    `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	QQ          string `orm:"column(qq);size(15);type(char);default();description(QQ)" json:"qq" form:"qq"`
	Blog        string `orm:"column(blog);size(100);type(char);default();description(个人网站)" json:"blog" form:"blog"`
	Name        string `orm:"column(name);size(50);type(char);default();description(姓名)" json:"name" form:"name" validate:"required" label:"姓名"`
	Major       string `orm:"column(major);size(50);type(char);default();description(专业)" json:"major" form:"major" validate:"required" label:"专业"`
	Weibo       string `orm:"column(weibo);size(100);type(char);default();description(微博)" json:"weibo" form:"weibo"`
	Github      string `orm:"column(github);size(100);type(char);default();description(Github)" json:"github" form:"github"`
	Avatar      int    `orm:"column(avatar);type(int);default(0);description(头像)" json:"avatar" form:"avatar" validate:"required,min=1" label:"头像"`
	Status      int8   `orm:"column(status);type(tinyint);default(0);description(是否展示)" json:"status" form:"status"`
	Description string `orm:"column(description);type(text);default();description(介绍)" json:"description" form:"description" validate:"required" label:"介绍"`
}

/** 设置引擎 **/
func (m *Team) TableEngine() string {
	return "MyISAM"
}

/** 设置表名 **/
func (m *Team) TableName() string {
	return "index_team"
}
