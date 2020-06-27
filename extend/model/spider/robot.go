package spider

import "time"

type Robot struct {
	Id          int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Name        string    `orm:"column(name);size(80);type(char);default();description(初始词)" json:"name" form:"name"`
	Resemblance string    `orm:"column(resemblance);type(text);default();description(近义词列表)" json:"resemblance" form:"resemblance"`
	CreateTime  time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime  time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *Robot) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Robot) TableName() string {
	return "robot_keyword"
}
