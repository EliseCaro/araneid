package spider

import "time"

type Journal struct {
	Id          int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Urls        string    `orm:"column(urls);size(200);type(char);default();description(地址)" json:"urls" form:"urls"`
	Usage       int       `orm:"column(usage);type(int);default(0);description(访问次数)" json:"usage" form:"usage"`
	Domain      string    `orm:"column(domain);size(32);type(char);default();description(域名)" json:"domain" form:"domain"`
	SpiderIp    string    `orm:"column(spider_ip);size(32);type(char);default();description(访问IP)" json:"spider_ip" form:"spider_ip"`
	SpiderName  string    `orm:"column(spider_name);size(32);type(char);default();description(蜘蛛标识)" json:"spider_name" form:"spider_name"`
	SpiderTitle string    `orm:"column(spider_title);size(32);type(char);default();description(蜘蛛标题)" json:"spider_title" form:"spider_title"`
	CreateTime  time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
}

/** 设置引擎 **/
func (m *Journal) TableEngine() string {
	return "MyISAM"
}

/** 设置表名 **/
func (m *Journal) TableName() string {
	return "spider_journal"
}
