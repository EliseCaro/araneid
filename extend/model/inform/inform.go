package inform

import "time"

type Inform struct {
	Id        int  `orm:"pk;auto;column(id);type(int);description(主键,自增)" json:"id" form:"id"`
	Form      int8 `orm:"column(form);default(1);description(站内信类型[1:大众爬虫日志])" json:"form" form:"form"`
	Statue    int8 `orm:"column(statue);default(0);description(是否查看)" json:"statue" form:"statue"`
	Receiver  int  `orm:"column(receiver);default(0);description(接受者ID)" json:"receiver" form:"receiver" validate:"required" label:"接受者ID"`
	ObjectId  int  `orm:"column(object_id);default(0);description(类型对象ID)" json:"object_id" form:"object_id"`
	ContextId int  `orm:"column(context_id);default(0);description(消息内容ID)" json:"context_id" form:"context_id" validate:"required" label:"消息内容ID"`
	Expediter int  `orm:"column(expediter);default(0);description(发送者ID)" json:"expediter" form:"expediter" validate:"required" label:"发送者ID"`
}

/** 设置引擎 **/
func (m *Inform) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Inform) TableName() string {
	return "admin_inform"
}

/** socket 消息结构 **/
type Message struct {
	Id         int       `json:"id"`
	Statue     int8      `json:"statue"`
	Receiver   int       `json:"receiver"`
	Context    string    `json:"context"`
	StringTime string    `json:"string_time"`
	CreateTime time.Time `json:"create_time"`
}

/****************************以下为内容表*******************************************/

type Context struct {
	Id         int       `orm:"pk;auto;column(id);type(int);description(主键,自增)" json:"id" form:"id"`
	Context    string    `orm:"column(context);type(text);default();description(站内信内容)" json:"context" form:"context" validate:"required" label:"站内信内容"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);description(发送时间)" json:"create_time"`
}

/** 设置引擎 **/
func (m *Context) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Context) TableName() string {
	return "admin_inform_context"
}
