package begin

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/attachment"
	"github.com/beatrice950201/araneid/extend/model/collect"
	"github.com/beatrice950201/araneid/extend/model/config"
	"github.com/beatrice950201/araneid/extend/model/dictionaries"
	"github.com/beatrice950201/araneid/extend/model/inform"
	"github.com/beatrice950201/araneid/extend/model/menus"
	"github.com/beatrice950201/araneid/extend/model/roles"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/beatrice950201/araneid/extend/model/users"
	"github.com/beatrice950201/araneid/extend/service"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

/** 初始化 **/
func init() {
	databasesBegin()
	configBegin()
	socketBegin()
	viewBegin()
	logsBegin()
}

/** 初始化数据库 **/
func databasesBegin() {
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
		new(users.Users), new(roles.Roles), new(menus.Menus), new(attachment.Attachment),
		new(collect.Collect), new(collect.Result), new(config.Config), new(inform.Inform),
		new(inform.Context), new(dictionaries.DictConfig), new(dictionaries.Dictionaries),
		new(spider.Disguise), new(spider.Template),
	)
	_ = orm.RunSyncdb("default", false, true)
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
func viewBegin() {
	_ = beego.AddFuncMap("breadcrumbMapTitle", breadcrumbMapTitle)
	_ = beego.AddFuncMap("fileBuyName", fileBuyName)
	_ = beego.AddFuncMap("fileBuySize", fileBuySize)
	_ = beego.AddFuncMap("fileBuyPath", fileBuyPath)
}

/** 注册日志 **/
func logsBegin() {
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
