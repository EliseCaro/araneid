package movie

import "time"

type Movie struct {
	Id          int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Name        string    `orm:"column(name);size(100);type(char);default();description(作品名称)" json:"name" form:"name" validate:"required" label:"作品名称"`
	Actor       string    `orm:"column(actor);type(text);default();description(作品主演)" json:"actor" form:"actor" validate:"required" label:"作品主演"`
	Director    string    `orm:"column(director);size(30);type(char);default();description(作品导演)" json:"director" form:"director" validate:"required" label:"作品导演"`
	Year        string    `orm:"column(year);size(30);type(char);default();description(作品年份)" json:"year" form:"year" validate:"required" label:"作品年份"`
	District    string    `orm:"column(district);size(150);type(char);default();description(作品地区)" json:"district" form:"district" validate:"required" label:"作品地区"`
	Genre       string    `orm:"column(genre);type(text);default();description(作品类型)" json:"genre" form:"genre" validate:"required" label:"作品类型"`
	Cover       string    `orm:"column(cover);size(255);type(char);default();description(作品封面)" json:"cover" form:"cover"`
	Context     string    `orm:"column(context);type(text);default();description(作品介绍)" json:"context" form:"context" validate:"required" label:"作品介绍"`
	ActorCtx    string    `orm:"column(actor_ctx);type(text);default();description(演员介绍)" json:"actor_ctx" form:"actor_ctx"`
	Title       string    `orm:"column(title);size(150);type(char);default();description(作品标题)" json:"title" form:"title" validate:"required" label:"作品标题"`
	Keywords    string    `orm:"column(keywords);size(200);type(char);default();description(作品关键字)" json:"keywords" form:"keywords" validate:"required" label:"作品关键字"`
	Description string    `orm:"column(description);size(255);type(char);default();description(作品描述)" json:"description" form:"description" validate:"required" label:"作品描述"`
	ClassType   int8      `orm:"index;column(class_type);type(tinyint);default(0);description(内容类型[0:电视 1:电影 2:综艺])" json:"class_type" form:"class_type"`
	Short       string    `orm:"column(short);size(80);type(char);default();description(简短标题)" json:"short" form:"short"`
	Logs        string    `orm:"column(logs);size(255);type(char);default();description(发布日志)" json:"logs" form:"logs"`
	Source      string    `orm:"column(source);size(150);type(char);default();description(任务地址)" json:"source" form:"source" validate:"required" label:"任务地址"`
	Status      int8      `orm:"index;column(status);type(tinyint);default(0);description(发布状态[-1:失败 0:等待 1:成功])" json:"status" form:"status"`
	CreateTime  time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime  time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *Movie) TableEngine() string {
	return "MyISAM"
}

/** 设置表名 **/
func (m *Movie) TableName() string {
	return "movie_index"
}

/*** 配置管理 **/
type ConfigMovie struct {
	Id    int    `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Name  string `orm:"unique;index;column(name);size(15);type(char);default();description(标识)" json:"name" form:"name"`
	Value string `orm:"column(value);size(255);type(char);default();description(配置值)" json:"value" form:"value"`
}

/** 设置引擎 **/
func (m *ConfigMovie) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *ConfigMovie) TableName() string {
	return "movie_config"
}

/*** 剧集管理 **/
type EpisodeMovie struct {
	Id          int       `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Pid         int       `orm:"index;column(pid);type(int);default(0);description(所属对象)" json:"pid" form:"pid"`
	Name        string    `orm:"column(name);size(10);type(char);default();description(剧集标识)" json:"name" form:"name" validate:"required" label:"剧集标识"`
	Short       string    `orm:"column(short);size(80);type(char);default();description(扩展标题)" json:"short" form:"short" validate:"required" label:"扩展标题"`
	Title       string    `orm:"column(title);size(150);type(char);default();description(剧集标题)" json:"title" form:"title" validate:"required" label:"剧集标题"`
	Keywords    string    `orm:"column(keywords);size(200);type(char);default();description(剧集关键字)" json:"keywords" form:"keywords" validate:"required" label:"剧集关键字"`
	Description string    `orm:"column(description);size(255);type(char);default();description(剧集描述)" json:"description" form:"description" validate:"required" label:"剧集描述"`
	Context     string    `orm:"column(context);type(text);default();description(剧集详情)" json:"context" form:"context" validate:"required" label:"剧集详情"`
	Logs        string    `orm:"column(logs);size(255);type(char);default();description(发布日志)" json:"logs" form:"logs"`
	Source      string    `orm:"column(source);size(150);type(char);default();description(任务地址)" json:"source" form:"source" validate:"required" label:"任务地址"`
	Status      int8      `orm:"index;column(status);type(tinyint);default(0);description(发布状态[-1:失败 0:等待 1:成功])" json:"status" form:"status"`
	CreateTime  time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime  time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *EpisodeMovie) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *EpisodeMovie) TableName() string {
	return "movie_episode"
}
