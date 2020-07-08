package service

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	bmCache "github.com/beatrice950201/araneid/extend/cache"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"math/rand"
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

/** 测试使用 todo 开发完成应该删除 **/
func (service *DefaultJournalService) CachedHandleSetDebug() {
	var maps []*spider.Journal
	_, _ = orm.NewOrm().QueryTable(new(spider.Journal)).All(&maps)
	for _, item := range maps {
		service.cachedHandleSet(*item)
	}
}

/** 获取折线图随机颜色 **/
func (service *DefaultJournalService) randomColor() string {
	color := []string{
		"rgba(27, 158, 183,0.8)",
		"rgba(34, 184, 207, .5)",
		"rgba(132, 94, 247, .3)",
	}
	return color[rand.Intn(len(color))]
}

/** 解析一组数据中每种蜘蛛有几个**/
func (service *DefaultJournalService) CachedHandleAnalysisClass(items []*spider.Journal) map[string]int64 {
	result := make(map[string]int64)
	for _, item := range items {
		if _, ok := result[item.SpiderTitle]; ok == true {
			result[item.SpiderTitle] += 1
		} else {
			result[item.SpiderTitle] = 1
		}
	}
	return result
}

/** 解析成控制面板折线图 **/
func (service *DefaultJournalService) CachedHandleAnalysisWeek() string {
	var result []map[string]interface{}
	var date []string
	data := service.CachedHandleGetWeek()
	class := _func.ParseAttrConfigMap(beego.AppConfig.String("system_spider_class"))
	for n, k := range class {
		var res []int64
		for _, i := range data {
			res = append(res, service.CachedHandleAnalysisDayClassSpider(i["items"].([]*spider.Journal), n))
		}
		if len(result) < 6 { // 最对只展示5个
			result = append(result, map[string]interface{}{
				"label": k, "data": res,
				"backgroundColor": service.randomColor(),
			})
		}
	}
	for _, i := range data {
		date = append(date, i["date"].(string))
	}
	bytes, _ := json.Marshal(map[string]interface{}{"date": date, "items": result})
	return string(bytes)
}

/** 从数据中提取每种蜘蛛每天有几个 **/
func (service *DefaultJournalService) CachedHandleAnalysisDayClassSpider(items []*spider.Journal, name string) int64 {
	var n int64
	for _, i := range items {
		if i.SpiderName == name {
			n += 1
		}
	}
	return n
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
func (service *DefaultJournalService) CachedHandleGetWeek() []map[string]interface{} {
	var items []map[string]interface{}
	var current = time.Now().Unix()
	for i := 0; i <= 6; i++ {
		date := time.Unix(current-(int64(i)*86400), 0).Format("20060102")
		tags := fmt.Sprintf(`journal_logs_%s`, date)
		item := map[string]interface{}{
			"items": []*spider.Journal{},
			"count": 0,
			"date":  time.Unix(current-(int64(i)*86400), 0).Format("01-02"),
		}
		if cache := _func.GetCache(tags); cache != "" {
			item["items"] = cache.([]*spider.Journal)
			item["count"] = len(cache.([]*spider.Journal))
		}
		items = append(items, item)
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
