package service

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	bmCache "github.com/beatrice950201/araneid/extend/cache"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"net/http"
	"strings"
	"time"
)

type DefaultJournalService struct{}

/** 根据URL获取一条数据 **/
func (service *DefaultJournalService) One(urls, name string) spider.Journal {
	var maps spider.Journal
	_ = orm.NewOrm().QueryTable(new(spider.Journal)).Filter("urls", urls).Filter("spider_name", name).One(&maps)
	return maps
}

/** 写入缓存【用来做七天数据分析跟一条数据分析】 **/
func (service *DefaultJournalService) cachedHandleSet(index spider.Journal) {
	var items []*spider.Journal
	var tags = fmt.Sprintf(`journal_logs_%s`, time.Now().Format("20060102"))
	var cache = _func.GetCache(tags)
	if cache == "" {
		index.Id = 1
	} else {
		items = cache.([]*spider.Journal)
		index.Id = len(items)
	}
	items = append(items, &index)
	_ = bmCache.Bm.Put(tags, items, (86400*10)*time.Second)
}

/** 获取今日缓存 **/
func (service *DefaultJournalService) CachedHandleGetDya() []*spider.Journal {
	var items []*spider.Journal
	var tags = fmt.Sprintf(`journal_logs_%s`, time.Now().Format("20060102"))
	if cache := _func.GetCache(tags); cache != "" {
		items = cache.([]*spider.Journal)
	}
	return items
}

/** 获取一周缓存 **/
func (service *DefaultJournalService) CachedHandleGetWeek() []*spider.Journal {
	var items []*spider.Journal
	var current = time.Now().Unix()
	for i := 0; i <= 6; i++ {
		date := time.Unix(current-(int64(i)*86400), 0).Format("20060102")
		tags := fmt.Sprintf(`journal_logs_%s`, date)
		if cache := _func.GetCache(tags); cache != "" {
			items = append(items, cache.([]*spider.Journal)...)
		}
	}
	return items
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
			Usage:       1,
		}
		var message error
		if one := service.One(item.Urls, item.SpiderName); one.Id > 0 {
			one.Usage += 1
			_, message = orm.NewOrm().Update(&one)
		} else {
			_, message = orm.NewOrm().Insert(&item)
		}
		if message == nil {
			service.cachedHandleSet(item)
		} else {
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
