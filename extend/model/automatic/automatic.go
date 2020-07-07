package automatic

import "time"

/** 百度自动推送记录 **/
type Automatic struct {
	Id         int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Urls       string    `orm:"column(urls);size(200);type(char);default();description(链接)" json:"urls" form:"urls" validate:"required" label:"推送链接"`
	Domain     string    `orm:"column(domain);size(66);type(char);default();description(域名)" json:"domain" form:"domain" validate:"required" label:"项目域名"`
	Status     int8      `orm:"index;column(status);size(1);type(int);default(0);description(推送结果)" json:"status" form:"status"`
	Remark     string    `orm:"column(remark);size(200);type(char);default();description(返回消息)" json:"remark" form:"remark" validate:"required" label:"返回消息"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *Automatic) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Automatic) TableName() string {
	return "spider_automatic"
}
