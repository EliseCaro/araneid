package spider

import "time"

/** 域名前缀挂载到模型管理 **/
/** 域名已文本框输入多个单独存一个表； **/
/** 域名组合结果单独存一个表；里面存入挂载的模板信息 **/
/** 索引池单独存一个表在蜘蛛池自节点里面导入**/
/** 挂载一个ktdb模板 **/

type Arachnid struct {
	Id         int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Name       string    `orm:"column(name);size(32);type(char);default();description(名称)" json:"name" form:"name" validate:"required" label:"项目名称"`
	Link       int       `orm:"column(link);type(int);default(0);description(友链数量)" json:"link" form:"link" validate:"required" label:"友链数量"`
	Zoology    int       `orm:"column(zoology);type(int);default(3);description(轮链最大数)" json:"zoology" form:"zoology" validate:"required" label:"轮链数量"`
	Indexes    int       `orm:"column(indexes);type(int);default(5);description(索引最大数)" json:"indexes" form:"indexes" validate:"required" label:"索引数量"`
	Models     int       `orm:"column(models);type(int);default(0);description(挂载模型)" json:"models" form:"models" validate:"required" label:"挂载模型"`
	Matching   int       `orm:"column(matching);type(int);default(0);description(匹配模板)" json:"matching" form:"matching" validate:"required" label:"匹配模板"`
	Domain     string    `orm:"column(domain);type(text);default();description(项目域名)" json:"domain" form:"domain" validate:"required" label:"项目域名"`
	Statistics string    `orm:"column(statistics);type(text);default();description(统计代码)" json:"statistics" form:"statistics" validate:"required" label:"统计代码"`
	Status     int8      `orm:"column(status);type(tinyint);default(0);description(启用状态)" json:"status" form:"status"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *Arachnid) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Arachnid) TableName() string {
	return "spider_arachnid"
}
