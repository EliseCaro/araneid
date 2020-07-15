package addons

type Addons struct {
	Id          int     `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	File        int     `orm:"column(file);type(int);default(0);description(文件对象)" json:"file" form:"file" validate:"required,min=1" label:"文件对象"`
	Title       string  `orm:"column(title);size(150);type(char);default();description(标题)" json:"title" form:"title" validate:"required" label:"标题"`
	Author      string  `orm:"column(author);size(150);type(char);default();description(作者)" json:"author" form:"author" validate:"required" label:"插件作者"`
	Version     float64 `orm:"column(version);default(0.01);description(版本号码)" json:"version" form:"version" validate:"required" label:"版本号码"`
	Description string  `orm:"column(description);size(255);type(char);default();description(描述)" json:"description" form:"description" validate:"required" label:"描述"`
}

/** 设置引擎 **/
func (m *Addons) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Addons) TableName() string {
	return "admin_addons"
}
