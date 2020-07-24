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
	"strconv"
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
		service.cachedHandleSetMonitorDomain(*item)
		service.cachedHandleSetMonitorTags(*item)
	}
}

/** 统计总蜘蛛数量 **/
func (service *DefaultJournalService) SumAll() int64 {
	var maps []orm.Params
	var size int64
	var prefix = beego.AppConfig.String("db_prefix")
	_, _ = orm.NewOrm().Raw("SELECT SUM(`usage`) AS size FROM " + prefix + "spider_journal").Values(&maps)
	if len(maps) > 0 && maps[0]["size"] != nil {
		size, _ = strconv.ParseInt(maps[0]["size"].(string), 10, 64)
	}
	return size
}

/** 统计日总蜘蛛数量 **/
func (service *DefaultJournalService) SumDay() int64 {
	return int64(len(service.CachedHandleGetDya()))
}

/** 获取热门网站 **/
func (service *DefaultJournalService) HotDomain() []*map[string]interface{} {
	var result []*map[string]interface{}
	var maps []orm.Params
	var prefix = beego.AppConfig.String("db_prefix")
	sql := fmt.Sprintf(`SELECT domain,COUNT(id) AS num FROM %sspider_journal GROUP BY domain ORDER BY num DESC LIMIT 4`, prefix)
	_, _ = orm.NewOrm().Raw(sql).Values(&maps)
	for _, item := range maps {
		if domain := new(DefaultDomainService).OneDomain(item["domain"].(string)); domain.Name != "" {
			result = append(result, &map[string]interface{}{
				"id":     domain.Id,
				"name":   domain.Name,
				"domain": domain.Domain,
				"count":  item["num"].(string),
			})
		}
	}
	return result
}

/** 获取折线图随机颜色 **/
func (service *DefaultJournalService) randomColor() string {
	color := []string{
		"rgba(27, 158, 183,0.8)",
		"rgba(132, 94, 247, .3)",
		"rgba(233, 236, 239, 1)",
		"rgba(34, 184, 207, .3)",
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var num int
	for i := 0; i < 10; i++ {
		num = r.Intn(len(color))
	}
	return color[num]
}

/** 解析一组数据中每种蜘蛛有几个**/
func (service *DefaultJournalService) CachedHandleAnalysisClass(items []*spider.Journal) map[string]map[string]interface{} {
	var result = make(map[string]map[string]interface{})
	for _, item := range items {
		if _, ok := result[item.SpiderName]; ok == true {
			num := result[item.SpiderName]["count"].(int) + 1
			result[item.SpiderName]["count"] = num
		} else {
			result[item.SpiderName] = map[string]interface{}{
				"count": 1,
				"title": item.SpiderTitle,
				"tags":  item.SpiderName,
				"urls":  beego.URLFor("Journal.Index", ":search", fmt.Sprintf(`date:%s|spider_name:%s`, time.Now().Format("2006-01-02"), item.SpiderName)),
			}
		}
	}
	if _, ok := result["Googlebot"]; ok == true && len(result) >= 5 {
		delete(result, "Googlebot")
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
				"label": k, "data": res, "fill": true,
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

/** 匹配沙盒监听 **/
func (service *DefaultJournalService) SpiderIp(ip string) (n, t string) {
	class := _func.ParseAttrConfigMap(beego.AppConfig.String("system_spider_monitor"))
	for index, item := range class {
		if strings.Index(ip, index) >= 0 {
			n = index
			t = item
			break
		}
	}
	return n, t
}

/** 写入沙盒数据 **/
func (service *DefaultJournalService) cachedHandleSetMonitorDomain(items spider.Journal) {
	if n, _ := service.SpiderIp(items.SpiderIp); n != "" {
		var domainResult []*spider.Journal
		domain := _func.GetCache(items.Domain)
		if domain != "" {
			domainResult = domain.([]*spider.Journal)
			items.Id = len(domainResult)
		} else {
			items.Id = 1
		}
		domainResult = append(domainResult, &items)
		_ = bmCache.Bm.Put(items.Domain, domainResult, (86400*365)*time.Second)
	}
}

/** 写入数量统计缓存 **/
func (service *DefaultJournalService) cachedHandleSetMonitorTags(items spider.Journal) {
	if n, t := service.SpiderIp(items.SpiderIp); n != "" {
		var countResult = make(map[string]map[string]interface{})
		count := _func.GetCache("spider_monitor_count")
		if count != "" {
			countResult = count.(map[string]map[string]interface{})
			if _, ok := countResult[n]; ok == true {
				countResult[n]["count"] = countResult[n]["count"].(int) + 1
			} else {
				countResult[n] = map[string]interface{}{
					"urls":  beego.URLFor("Journal.Index", ":search", fmt.Sprintf(`date:%s|spider_ip:%s`, time.Now().Format("2006-01-02"), n)),
					"count": 1,
					"title": t,
				}
			}
		} else {
			countResult[n] = map[string]interface{}{
				"urls":  beego.URLFor("Journal.Index", ":search", fmt.Sprintf(`date:%s|spider_ip:%s`, time.Now().Format("2006-01-02"), n)),
				"count": 1,
				"title": t,
			}
		}
		_ = bmCache.Bm.Put("spider_monitor_count", countResult, (86400*365)*time.Second)
	}
}

/** 获取数量统计缓存 **/
func (service *DefaultJournalService) CachedHandleGetMonitorTags() map[string]map[string]interface{} {
	count := _func.GetCache("spider_monitor_count")
	var countResult = make(map[string]map[string]interface{})
	if count != "" {
		countResult = count.(map[string]map[string]interface{})
	}
	return countResult
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
			service.cachedHandleSetMonitorDomain(item)
			service.cachedHandleSetMonitorTags(item)
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
		if strings.Index(agent, index) >= 0 {
			n = index
			t = item
			break
		}
	}
	return n, t
}

/****************** 以下为表格渲染  ***********************/

/** 批量删除 **/
func (service *DefaultJournalService) DeleteArray(array []int) (message error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, message = orm.NewOrm().Delete(&spider.Journal{Id: v}); message != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if message == nil {
		_ = orm.NewOrm().Commit()
	}
	return message
}

/** 清空 **/
func (service *DefaultJournalService) EmptyDelete() {
	var item []*spider.Journal
	_, _ = orm.NewOrm().QueryTable(new(spider.Journal)).All(&item)
	for _, v := range item {
		_, _ = orm.NewOrm().Delete(&spider.Journal{Id: v.Id})
	}
}

/** 获取需要渲染的Column **/
func (service *DefaultJournalService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "访问标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "访问地址", "name": "urls", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "访问次数", "name": "usage", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "访问域名", "name": "domain", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "访问IP", "name": "spider_ip", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "蜘蛛归类", "name": "spider_title", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "访问时间", "name": "create_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "数据操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultJournalService) DataTableButtons() []*_func.TableButtons {
	var array []*_func.TableButtons
	array = append(array, &_func.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-warning mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Journal.Delete")},
	})
	array = append(array, &_func.TableButtons{
		Text:      "清空记录",
		ClassName: "btn btn-sm btn-alt-danger mt-1 js-tooltip ids_deletes",
		Attribute: map[string]string{
			"data-action":         beego.URLFor("Journal.Empty"),
			"data-toggle":         "tooltip",
			"data-original-title": "清空后将无法恢复，请谨慎操作！",
		},
	})
	return array
}

/** 处理分页 **/
func (service *DefaultJournalService) PageListItems(length, draw, page int, search map[string]string) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(spider.Journal))
	qs = service.searchAnalysis(search, qs)
	recordsTotal, _ := qs.Count()
	_, _ = qs.Limit(length, length*(page-1)).OrderBy("-id").ValuesList(&lists, "id", "urls", "usage", "domain", "spider_ip", "spider_title", "create_time")
	for _, v := range lists {
		v[1] = service.substr2HtmlHref(v[1].(string), v[1].(string), 0, 20)
	}
	data := map[string]interface{}{
		"draw":            draw,         // 请求次数
		"recordsFiltered": recordsTotal, // 从多少条里面筛选
		"recordsTotal":    recordsTotal, // 总条数
		"data":            lists,        // 筛选结果
	}
	return data
}

/** 解析map写入条件 **/
func (service *DefaultJournalService) searchAnalysis(search map[string]string, qs orm.QuerySeter) orm.QuerySeter {
	if _, ok := search["date"]; ok == true {
		qs = qs.Filter("create_time__gte", fmt.Sprintf(`%s 00:00:00`, search["date"]))
		qs = qs.Filter("create_time__lte", fmt.Sprintf(`%s 23:59:59`, search["date"]))
	}
	if _, ok := search["spider_name"]; ok == true {
		qs = qs.Filter("spider_name__icontains", search["spider_name"])
	}
	if _, ok := search["spider_title"]; ok == true {
		qs = qs.Filter("spider_title__icontains", search["spider_title"])
	}
	if _, ok := search["spider_ip"]; ok == true {
		qs = qs.Filter("spider_ip__icontains", search["spider_ip"])
	}
	if _, ok := search["domain"]; ok == true {
		qs = qs.Filter("domain", search["domain"])
	}
	return qs
}

/**  转为pop提示 **/
func (service *DefaultJournalService) substr2HtmlHref(u, s string, start, end int) string {
	html := fmt.Sprintf(`<a href="%s" target="_blank" class="badge badge-primary js-tooltip" data-placement="top" data-toggle="tooltip" data-original-title="%s">%s...</a>`, u, s, beego.Substr(s, start, end))
	return html
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultJournalService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string", "string", "date"},
		"fieldName": {"id", "urls", "usage", "domain", "spider_ip", "spider_title", "create_time"},
		"action":    {"", "", "", "", "", "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultJournalService) TableButtonsType() []*_func.TableButtons {
	buttons := []*_func.TableButtons{
		{
			Text:      "删除",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Journal.Delete"),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}
