package index

import "time"

type Message struct {
	Id         int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Name       string    `orm:"column(name);size(20);type(char);default();description(姓名)" json:"name" form:"name" validate:"required" label:"姓名"`
	Email      string    `orm:"column(email);size(64);type(char);default();description(邮箱)" json:"email" form:"email" validate:"required,email" label:"邮箱"`
	Title      string    `orm:"column(title);size(100);type(char);default();description(标题)" json:"title" form:"title" validate:"required" label:"标题"`
	Reply      string    `orm:"column(reply);type(text);default();description(问题回复)" json:"reply" form:"reply"`
	Status     int8      `orm:"column(status);type(tinyint);default(0);description(是否处理)" json:"status" form:"status"`
	Context    string    `orm:"column(context);type(text);default();description(问题详情)" json:"context" form:"context" validate:"required" label:"内容详情"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *Message) TableEngine() string {
	return "MyISAM"
}

/** 设置表名 **/
func (m *Message) TableName() string {
	return "index_message"
}
