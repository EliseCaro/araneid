package service

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"net/http"
	"strings"
)

type DefaultJournalService struct{}

/** 根据URL获取一条数据 **/
func (service *DefaultJournalService) One(urls string) spider.Journal {
	var maps spider.Journal
	_ = orm.NewOrm().QueryTable(new(spider.Journal)).Filter("urls", urls).One(&maps)
	return maps
}

/** 初始化记录接口 **/
func (service *DefaultJournalService) HandleInstantiation(ctx *context.Context) {
	if n, t := service.UserAgent(ctx.Input.UserAgent()); n != "" {
		item := spider.Journal{
			SpiderName:  n,
			SpiderTitle: t,
			Urls:        service.GetURL(ctx.Request),
			Domain:      ctx.Input.Domain(),
			SpiderIp:    ctx.Input.IP(),
		}
		var message error
		if one := service.One(item.Urls); one.Id > 0 {
			one.Usage += 1
			_, message = orm.NewOrm().Update(&one)
		} else {
			_, message = orm.NewOrm().Insert(&item)
		}
		if message != nil {
			logs.Error("蜘蛛访问记录失败；失败原因：%s", message.Error())
		}
	}
}

/** 获取完整url **/
func (service *DefaultJournalService) GetURL(r *http.Request) (Url string) {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	return strings.Join([]string{scheme, r.Host, r.RequestURI}, "")
}

/** 获取蜘蛛标识 **/
func (service *DefaultJournalService) UserAgent(agent string) (n, t string) {
	class := _func.ParseAttrConfigMap(beego.AppConfig.String("system_spider_class"))
	for index, item := range class {
		if strings.Index(agent, index) > 0 {
			n = index
			t = item
			break
		}
	}
	return n, t
}
