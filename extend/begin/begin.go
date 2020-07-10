package begin

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/attachment"
	"github.com/beatrice950201/araneid/extend/model/automatic"
	"github.com/beatrice950201/araneid/extend/model/collect"
	"github.com/beatrice950201/araneid/extend/model/config"
	"github.com/beatrice950201/araneid/extend/model/dictionaries"
	"github.com/beatrice950201/araneid/extend/model/index"
	"github.com/beatrice950201/araneid/extend/model/inform"
	"github.com/beatrice950201/araneid/extend/model/menus"
	"github.com/beatrice950201/araneid/extend/model/roles"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/beatrice950201/araneid/extend/model/users"
	"github.com/beatrice950201/araneid/extend/service"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strconv"
	"time"
)

/** 初始化 **/
func init() {
	coreDbBegin()
	logsDbBegin()
	configBegin()
	socketBegin()
	viewItBegin()
	logsItBegin()
}

/** 初始化日志数据库 **/
func logsDbBegin() {}

/** 初始化数据库 **/
func coreDbBegin() {
	orm.Debug = _func.AnalysisDebug()
	host := beego.AppConfig.String("db_host")
	username := beego.AppConfig.String("db_username")
	password := beego.AppConfig.String("db_password")
	database := beego.AppConfig.String("db_database")
	port := beego.AppConfig.String("db_port")
	charset := beego.AppConfig.String("db_charset")
	prefix := beego.AppConfig.String("db_prefix")
	_ = orm.RegisterDataBase("default", "mysql", fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=%s&loc=%s`, username, password, host, port, database, charset, `Asia%2FShanghai`), 30)
	orm.RegisterModelWithPrefix(
		prefix,
		new(users.Users), new(roles.Roles), new(menus.Menus), new(spider.Match), new(spider.Robot), new(spider.Class),
		new(spider.Prefix), new(spider.Models), new(inform.Inform), new(spider.Detail), new(spider.Domain), new(config.Config),
		new(spider.Keyword), new(spider.Article), new(inform.Context), new(collect.Result), new(spider.Indexes),
		new(collect.Collect), new(spider.Disguise), new(spider.Template), new(spider.Arachnid), new(spider.Category),
		new(automatic.Automatic), new(attachment.Attachment), new(dictionaries.DictConfig), new(dictionaries.Dictionaries),
		new(spider.Journal), new(index.Advantage), new(index.News), new(index.Team),
	)
	_ = orm.RunSyncdb("default", false, _func.AnalysisDebug())
}

/** 覆盖配置 **/
func configBegin() {
	for k, v := range _func.CacheConfig() {
		_ = beego.AppConfig.Set(k, v)
	}
}

/** 启动socket服务 **/
func socketBegin() {
	go service.SocketInstanceGet().ServiceRoom()
}

/** 注册模板函数 **/
func viewItBegin() {
	_ = beego.AddFuncMap("breadcrumbMapTitle", breadcrumbMapTitle)
	_ = beego.AddFuncMap("fileBuyName", fileBuyName)
	_ = beego.AddFuncMap("fileBuySize", fileBuySize)
	_ = beego.AddFuncMap("fileBuyPath", fileBuyPath)
	_ = beego.AddFuncMap("spiderArticleLimit", spiderArticleLimit)
	_ = beego.AddFuncMap("spiderArticleView", spiderArticleView)
}

/** 注册日志 **/
func logsItBegin() {
	logsPath := beego.AppConfig.String("logs_path") + time.Now().Format("20060102") + ".log"
	option := fmt.Sprintf(`{"filename":"%s","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`, logsPath)
	_ = logs.SetLogger(logs.AdapterFile, option)
}

/****************************以下全部为模板函数 **********************************/

/** 获得面包屑主标题 **/
func breadcrumbMapTitle(menus []menus.Menus) string {
	return menus[len(menus)-1].Title
}

/** 获取附件名称 **/
func fileBuyName(id int) string {
	var item attachment.Attachment
	item.Id = id
	_ = orm.NewOrm().Read(&item)
	return item.Name
}

/** 获取附件大小 **/
func fileBuySize(id int) int64 {
	var item attachment.Attachment
	item.Id = id
	_ = orm.NewOrm().Read(&item)
	return item.Size
}

/** 获取附件路径 **/
func fileBuyPath(id int) string {
	var item attachment.Attachment
	item.Id = id
	_ = orm.NewOrm().Read(&item)
	return _func.DomainStatic(item.Driver) + item.Path
}

/**
	根据条件获取文章 todo 考虑关联文章表获得最新的标题组合
    sort [RAND(),FIELD [ASC|DESC]]
**/
func spiderArticleLimit(cid interface{}, limit int, where, sort string) []spider.Article {
	var prefix = beego.AppConfig.String("db_prefix")
	var maps []spider.Article
	var value = cid
	if reflect.TypeOf(cid).String() == "int" {
		value = strconv.Itoa(cid.(int))
	}
	sql := fmt.Sprintf(`SELECT id,title,description FROM %sspider_article WHERE %s%s ORDER BY %s LIMIT %d`, prefix, where, value, sort, limit)
	_, _ = orm.NewOrm().Raw(sql).QueryRows(&maps)
	return maps
}

/** 获取文章阅读量 **/
func spiderArticleView(oid int) int {
	var maps spider.Article
	_ = orm.NewOrm().QueryTable(new(spider.Article)).Filter("id", oid).One(&maps)
	return maps.View
}
