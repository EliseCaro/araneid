package collect

import "time"

type Collect struct {
	Id         int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Name       string    `orm:"column(name);size(32);type(char);default();description(爬虫名字)" json:"name" form:"name" validate:"required" label:"爬虫名称"`
	Status     int8      `orm:"index;column(status);type(tinyint);default(0);description(采集器状态)" json:"status" form:"status"`
	Domain     string    `orm:"column(domain);size(150);type(char);default();description(发布接口地址)" json:"domain" form:"domain" validate:"required" label:"发布接口"`
	Remark     string    `orm:"column(remark);size(200);type(char);default();description(爬虫备注)" json:"remark" form:"remark" validate:"required" label:"爬虫备注"`
	Source     string    `orm:"column(source);type(text);default();description(采集源队列)" json:"source" form:"source" validate:"required" label:"采集源队列"`
	Interval   int8      `orm:"column(interval);type(tinyint);default(5);description(最大并发次数)" json:"interval" form:"interval" validate:"required,min=1" label:"最大并发次数"`
	Translate  int8      `orm:"column(translate);type(tinyint);default(0);description(简繁体转换)" json:"translate" form:"translate"`
	Download   int8      `orm:"column(download);type(tinyint);default(0);description(开启附件下载)" json:"download" form:"download"`
	Matching   string    `orm:"column(matching);type(text);default();description(字段匹配)"  json:"matching" form:"matching" validate:"required" label:"字段匹配"`
	PushTime   int       `orm:"column(push_time);type(int);default(0);description(发布间隔时间)" json:"push_time" form:"push_time" validate:"required,min=1" label:"发布间隔时间"`
	SourceRule string    `orm:"column(source_rule);size(100);type(char);default();description(采集规则)" json:"source_rule" form:"source_rule" validate:"required" label:"采集规则"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
	PushStatus int8      `orm:"index;column(push_status);type(tinyint);default(0);description(发布器状态)" json:"push_status" form:"push_status"`
}

/** 设置引擎 **/
func (m *Collect) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Collect) TableName() string {
	return "admin_collect"
}

/** 爬虫字段解析结构 **/
type Matching struct {
	Form       int    `json:"form" form:"form"`
	Field      string `json:"field" form:"field" validate:"required" label:"匹配字段对应字段"`
	AttrName   string `json:"attr_name" form:"attr_name"`
	Selector   string `json:"selector" form:"selector" validate:"required" label:"选择对象"`
	Filtration int    `json:"filtration" form:"filtration"`
}

/** 爬虫详细 **/
type DetailCollect struct {
	Collect
	MatchingJson  []*Matching `json:"matching_json"`
	MatchingCount int         `json:"matching_count"`
}

/****************************以下为采集结果内容 ****************************/

type Result struct {
	Id         int       `orm:"pk;auto;column(id);type(int);description(主键,自增)" json:"id" form:"id"`
	Logs       string    `orm:"column(logs);size(150);type(char);default();description(发布日志)" json:"logs" form:"logs"`
	Title      string    `orm:"column(title);size(80);type(char);default();description(任务标题)" json:"title" form:"title" validate:"required" label:"任务标题"`
	Source     string    `orm:"column(source);size(150);type(char);default();description(任务地址)" json:"source" form:"source" validate:"required" label:"任务地址"`
	Result     string    `orm:"column(result);type(text);default();description(采集结果)" json:"result" form:"result" validate:"required" label:"采集结果"`
	Status     int8      `orm:"index;column(status);type(tinyint);default(0);description(发布状态[-1:失败 0:等待 1:成功])" json:"status" form:"status"`
	Collect    int       `orm:"index;column(collect);default(0);description(所属爬虫)" json:"collect" form:"collect" validate:"required" label:"所属爬虫"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *Result) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Result) TableName() string {
	return "admin_collect_result"
}
