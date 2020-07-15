package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/collect"
	ccst "github.com/go-cc/cc-jianfan"
	"github.com/go-playground/validator"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type DefaultCollectService struct{}

/** 获取远行中的数量 **/
func (service *DefaultCollectService) aliveNum() int64 {
	index, _ := orm.NewOrm().QueryTable(new(collect.Collect)).Filter("status", 1).Count()
	return index
}

/** 获取远行中的发布器数量 **/
func (service *DefaultCollectService) alivePushNum() int64 {
	index, _ := orm.NewOrm().QueryTable(new(collect.Collect)).Filter("push_status", 1).Count()
	return index
}

/** 发布一条数据;返回格式：{status:[false|true],message:[message]} **/
func (service *DefaultCollectService) PushDetailAPI(item collect.Result) {
	collectDetail := service.One(item.Collect)
	if resp, err := http.PostForm(collectDetail.Domain, url.Values{
		"password": {beego.AppConfig.String("collect_collect_password")},
		"object":   {strconv.Itoa(item.Id)},
		"title":    {item.Title}, "source": {item.Source}, "result": {item.Result},
		"collect": {strconv.Itoa(collectDetail.Id)},
		"update":  {strconv.FormatInt(item.UpdateTime.Unix(), 10)},
		"create":  {strconv.FormatInt(item.CreateTime.Unix(), 10)},
	}); err == nil {
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			service.pushDetailAPIResult(item.Id, -1, err)
		} else {
			result := make(map[string]interface{})
			if err := json.Unmarshal(body, &result); err == nil && len(result) > 0 {
				if result["status"].(bool) == true {
					service.pushDetailAPIResult(item.Id, 1, errors.New(result["message"].(string)))
				} else {
					service.pushDetailAPIResult(item.Id, -1, errors.New(result["message"].(string)))
				}
			} else {
				service.pushDetailAPIResult(item.Id, -1, errors.New("接口返回结果解析失败！请检查返回格式！"))
			}
		}
		defer resp.Body.Close()
	} else {
		service.pushDetailAPIResult(item.Id, -1, err)
	}
}

/** 获取一条爬虫详细数据 **/
func (service *DefaultCollectService) FindOneDetail(id int) collect.DetailCollect {
	var (
		item collect.Collect
		rest collect.DetailCollect
	)
	_ = orm.NewOrm().QueryTable(new(collect.Collect)).Filter("id", id).One(&item)
	if j, _ := json.Marshal(item); len(j) > 0 {
		_ = json.Unmarshal(j, &rest)
	}
	_ = json.Unmarshal([]byte(rest.Matching), &rest.MatchingJson)
	rest.MatchingCount = len(rest.MatchingJson)
	return rest
}

/** 获取一条数据 **/
func (service *DefaultCollectService) One(id int) collect.Collect {
	var item collect.Collect
	_ = orm.NewOrm().QueryTable(new(collect.Collect)).Filter("id", id).One(&item)
	return item
}

/** 控制器获取指令类型返回message **/
func (service *DefaultCollectService) InstructTypeMessage(field string, status int8) string {
	mapsString := map[string][]string{
		"push_status": {
			"发布器已经停止发布！",
			"发布器已经启动；发布数据将推送到指定接口！",
		},
		"status": {
			"采集器已经停止采集！",
			"采集器已经启动；您可以查看采集结果！",
		},
	}
	return mapsString[field][status]
}

/** 批量删除 **/
func (service *DefaultCollectService) DeleteArray(array []int) (e error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, e = orm.NewOrm().Delete(&collect.Collect{Id: v}); e != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if e == nil {
		_ = orm.NewOrm().Commit()
	}
	return e
}

/** 根据采集源跟匹配规则提取全部URL **/
func (service *DefaultCollectService) ExtractUrls(role string, source []string) []map[string]string {
	var result []map[string]string
	domain := new(DefaultCollectService).queueUrlDomain(source[0])
	collector := service.collectInstance(5, 2, domain, true)
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		matchText := service.eliminateTrim(e.Text, []string{" ", "\n"})
		if service.checkSourceRule(role, e.Attr("href")) && matchText != "" && service.inMapHref(e.Attr("href"), result) == false {
			result = append(result, map[string]string{
				"title": matchText,
				"href":  service.completionSrc(source[0], e.Attr("href")),
			})
		}
	})
	for _, v := range source {
		_ = collector.Visit(v)
	}
	collector.Wait()
	return result
}

/** 根据地址跟字段规则解析一条详情 **/
func (service *DefaultCollectService) ExtractDocumentMatching(url string, matching []*collect.Matching) (map[string]string, error) {
	domain := new(DefaultCollectService).queueUrlDomain(url)
	collector := service.collectInstance(5, 1, domain, true)
	result := make(map[string]string)
	var message error
	collector.OnHTML("html", func(e *colly.HTMLElement) {
		matching = append(matching, &collect.Matching{Field: "meta_title", Selector: "head > title", Filtration: 1, Form: 0})
		for _, v := range matching {
			if value := service.extractMatchingField(v, e); value != "" {
				result[v.Field] = value
			} else {
				message = errors.New("字段匹配不齐全；已被强制过滤！")
				break
			}
		}
	})
	_ = collector.Visit(url)
	collector.Wait()
	return result, message
}

/** 对接socket入口 **/
func (service *DefaultCollectService) InstanceBegin(instruct map[string]interface{}) {
	switch instruct["field"].(string) {
	case "push_status":
		service.pushStatus(instruct["id"].(int), instruct["status"].(int8), instruct["uid"].(int))
	case "status":
		service.collectStatus(instruct["id"].(int), instruct["status"].(int8), instruct["uid"].(int))
	}
}

/********************* 以下为私有方法 *****************************8*/

/**  更新发布结果 **/
func (service *DefaultCollectService) pushDetailAPIResult(id int, status int8, message error) {
	if _, err := orm.NewOrm().Update(&collect.Result{Id: id, Status: status, Logs: error.Error(message)}, "Status", "Logs"); err != nil {
		logs.Warn("更新发布结果失败；失败原因：%s", error.Error(err))
	}
}

/** 判断MAP是否存在某个值 **/
func (service *DefaultCollectService) inMapHref(need string, maps []map[string]string) (is bool) {
	for _, v := range maps {
		if need == v["href"] {
			is = true
			break
		}
	}
	return is
}

/** 检查链接是否合法 **/
func (service *DefaultCollectService) checkSourceRule(rule, url string) bool {
	return regexp.MustCompile(rule).MatchString(url)
}

/** 实例化采集器容器 todo 开启代理跟Redis储存 **/
func (service *DefaultCollectService) collectInstance(interval, depth int, domain string, async bool) *colly.Collector {
	collector := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{Output: logs.GetLogger().Writer()}),
		colly.Async(async),
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
		colly.MaxDepth(depth),
	)
	_ = collector.Limit(&colly.LimitRule{DomainGlob: domain, Parallelism: interval})
	collector.WithTransport(&http.Transport{DisableKeepAlives: true})
	collector.OnRequest(func(r *colly.Request) { r.Headers.Set("User-Agent", table.RandomString()) })
	return collector
}

/** 获得连接URL根域名[不带http] 用于限制抓取域 **/
func (service *DefaultCollectService) queueUrlDomain(s string) string {
	var domain string
	if o := strings.Split(s, "//"); len(o) >= 2 {
		if o := strings.Split(o[1], "/"); len(o) >= 1 {
			domain = o[0]
		}
	}
	return domain
}

/** 去除连续两个以上的空格 **/
func (service *DefaultCollectService) deleteExtraSpace(s string) string {
	s1 := strings.Replace(s, "	", " ", -1)
	regString := "\\s{2,}"
	reg, _ := regexp.Compile(regString)
	s2 := make([]byte, len(s1))
	copy(s2, s1)
	spcIndex := reg.FindStringIndex(string(s2))
	for len(spcIndex) > 0 {
		s2 = append(s2[:spcIndex[0]+1], s2[spcIndex[1]:]...)
		spcIndex = reg.FindStringIndex(string(s2))
	}
	return strings.Trim(string(s2), " ")
}

/** 剔除字符规定字符 **/
func (service *DefaultCollectService) eliminateTrim(str string, symbol []string) string {
	for _, v := range symbol {
		str = strings.Replace(str, v, "", -1)
	}
	return str
}

/** 过滤SCRIPT **/
func (service *DefaultCollectService) deleteScriptSpace(src string) string {
	re, _ := regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	return re.ReplaceAllString(src, "")
}

/** 过滤A标签 **/
func (service *DefaultCollectService) deleteHrefSpace(src string) string {
	re, _ := regexp.Compile("\\<a[\\S\\s]+?\\</a\\>")
	return re.ReplaceAllString(src, "")
}

/** 过滤Image标签 **/
func (service *DefaultCollectService) deleteImageSpace(src string) string {
	re, _ := regexp.Compile("\\<img [\\S\\s]+?/\\>")
	return re.ReplaceAllString(src, "")
}

/** 过滤换行标签 **/
func (service *DefaultCollectService) deleteBrSpace(src string) string {
	re, _ := regexp.Compile("\\<br[\\S\\s]+?\\>")
	return re.ReplaceAllString(src, "")
}

/** 根据链接获取一条结果数据 **/
func (service *DefaultCollectService) oneResultLink(url string) collect.Result {
	var item collect.Result
	_ = orm.NewOrm().QueryTable(new(collect.Result)).Filter("source", url).One(&item)
	return item
}

/** 根据字段规则过滤一个基本值; **/
func (service *DefaultCollectService) extractMatchingField(m *collect.Matching, e *colly.HTMLElement) string {
	var stringHtml string
	if m.Form == 1 || m.Form == 2 {
		stringHtml, _ = e.DOM.Find(m.Selector).Attr(m.AttrName)
		return stringHtml
	}
	if m.Filtration == 1 {
		stringHtml = service.eliminateTrim(e.ChildText(m.Selector), []string{" ", "\n"})
	} else {
		stringHtml, _ = e.DOM.Find(m.Selector).Html()
		stringHtml = service.eliminateTrim(stringHtml, []string{"\n"})
		stringHtml = service.deleteExtraSpace(stringHtml)
		stringHtml = strings.Trim(stringHtml, " ")
		stringHtml = service.deleteScriptSpace(stringHtml)
		stringHtml = service.deleteHrefSpace(stringHtml)
		stringHtml = service.deleteBrSpace(stringHtml)
		if m.Image == 1 { //  剔除Image
			stringHtml = service.deleteImageSpace(stringHtml)
		}
	}
	if m.Eliminate != "" { // 剔除规定字符
		eliminate := strings.Split(m.Eliminate, "|")
		stringHtml = service.eliminateTrim(stringHtml, eliminate)
	}
	return stringHtml
}

/** 创建一条结果数据 */
func (service *DefaultCollectService) createOneResult(title, source string, collectId int, result map[string]string) (err error) {
	str, _ := json.Marshal(result)
	verify := DefaultBaseVerify{}
	item := &collect.Result{Title: title, Source: source, Result: string(str), Collect: collectId, Status: 0}
	if message := verify.Begin().Struct(item); message != nil {
		return errors.New(verify.Translate(message.(validator.ValidationErrors)))
	}
	if one := service.oneResultLink(item.Source); one.Id > 0 {
		item.Id = one.Id
		_, err = orm.NewOrm().Update(item)
	} else {
		_, err = orm.NewOrm().Insert(item)
	}
	return err
}

/** 发送采集爬虫通知 **/
func (service *DefaultCollectService) createLogsInform(status, receiver, objectId int, name, message, url string) {
	var htmlInfo string
	nowTime := beego.Date(time.Now(), "m月d日 H:i")
	if status == 1 {
		htmlInfo = fmt.Sprintf(`名为<a href="javascript:void(0);">%s</a>的采集器；
					在%s爬取到一条数据；该数据的标题为:<a href="javascript:void(0);">%s</a>。
					您可以<a target="_blank" href="%s">查看原文</a>;`,
			name, nowTime, message, url,
		)
	} else {
		htmlInfo = fmt.Sprintf(
			`名为<a href="javascript:void(0);">%s</a>的采集器；
					在%s爬取数据失败；失败原因为:<a href="javascript:void(0);">%s</a>。
					您可以<a target="_blank" href="%s">查看原文</a>纠正过滤器或者直接忽略！`,
			name, nowTime, message, url,
		)
	}
	inform := DefaultInformService{}
	inform.SendSocketInform([]int{receiver}, objectId, 0, 1, htmlInfo)
}

/** 发送采集爬虫状态通知 **/
func (service *DefaultCollectService) createLogsInformStatus(status string, receiver, objectId int, name string) {
	var htmlInfo string
	nowTime := beego.Date(time.Now(), "m月d日 H:i")
	htmlInfo = fmt.Sprintf(`
		名为<a href="javascript:void(0);">%s</a>的采集器;在%s的时候%s；
		您可以<a target="_blank" href="%s">查看爬取结果</a>;`, name, nowTime, status, beego.URLFor("Collector.Index", ":id", objectId),
	)
	inform := DefaultInformService{}
	inform.SendSocketInform([]int{receiver}, objectId, 0, 1, htmlInfo)
}

/** 发送发布器状态通知 **/
func (service *DefaultCollectService) createLogsInformPushStatus(status string, receiver, objectId int, name string) {
	var htmlInfo string
	nowTime := beego.Date(time.Now(), "m月d日 H:i")
	htmlInfo = fmt.Sprintf(`
		名为<a href="javascript:void(0);">%s</a>的发布器;在%s的时候%s；
		您可以<a target="_blank" href="%s">查看发布结果</a>;`, name, nowTime, status, beego.URLFor("Collector.Index", ":id", objectId),
	)
	inform := DefaultInformService{}
	inform.SendSocketInform([]int{receiver}, objectId, 0, 1, htmlInfo)
}

/** todo 伪原创 **/
func (service *DefaultCollectService) fieldHandle(result map[string]string, detail collect.DetailCollect) map[string]string {
	for _, v := range detail.MatchingJson {
		if detail.Translate > 0 && v.Form == 0 || detail.Translate > 0 && v.Form == 2 { // 语言转换
			result[v.Field] = service.translate(result[v.Field], detail.Translate)
		}
		if detail.Download == 1 && v.Form == 1 { // 资源下载单个
			completion := service.completionSrc(detail.Source, result[v.Field])
			result[v.Field] = service.download(completion)
		}
		if detail.Download == 1 && v.Filtration == 0 && v.Form == 0 && v.AttrName != "" { // 下载html里面的src
			result[v.Field] = service.downloadStringResource(detail.Source, result[v.Field], v.AttrName)
		}
	}
	return result
}

/** 检索下载Html里面的src标签 **/
func (service *DefaultCollectService) downloadStringResource(source, html, name string) string {
	if dom, err := goquery.NewDocumentFromReader(strings.NewReader(html)); err == nil {
		dom.Find("img[" + name + "]").Each(func(i int, selection *goquery.Selection) {
			imageSrc, _ := selection.Attr(name)
			completion := service.completionSrc(source, imageSrc)
			selection.SetAttr("src", service.download(completion))
		})
		html, _ = dom.Find("body").Html()
	}
	return html
}

/** 组装资源完整路径 todo 此处应该在优化 **/
func (service *DefaultCollectService) completionSrc(source, urls string) string {
	index := strings.Index(urls, "//")
	if index == 0 {
		urls = "https:" + urls
	}
	if index == -1 {
		sourceMap := table.ParseAttrConfigArray(source)
		domain := service.queueUrlDomain(sourceMap[0])
		prefix := strings.Split(sourceMap[0], "//")
		urls = prefix[0] + "//" + domain + urls
	}
	return urls
}

/** 语言转换 **/
func (service *DefaultCollectService) translate(s string, t int8) string {
	if t == 1 {
		return ccst.S2T(s)
	} else {
		return ccst.T2S(s)
	}
}

/** 资源下载单个 **/
func (service *DefaultCollectService) download(url string) string {
	adjunct := DefaultAdjunctService{}
	status, _ := beego.AppConfig.Int("collect_collect_cloud_open")
	if status == 1 {
		n := adjunct.DownloadFileCloud(url, map[string]string{
			"bucket":   beego.AppConfig.String("collect_collect_cloud_bucket"),
			"name":     beego.AppConfig.String("collect_collect_cloud_name"),
			"password": beego.AppConfig.String("collect_collect_cloud_password"),
		})
		domain := beego.AppConfig.String("collect_collect_cloud_static")
		if n != url {
			url = domain + n
		}
	} else {
		n := adjunct.DownloadFileLocal(url)
		if n != url {
			url = table.DomainStatic("local") + n
		}
	}
	return url
}

/** 处理一条url进入库  **/
func (service *DefaultCollectService) handleSourceRuleBody(url string, uid int, detail collect.DetailCollect) {
	matching := detail.MatchingJson
	matching = append(matching, &collect.Matching{
		Field: "meta_title", Selector: "head > title", Filtration: 1, Form: 0,
	})
	collector := service.collectInstance(5, 1, service.queueUrlDomain(url), true)
	collector.OnHTML("html", func(e *colly.HTMLElement) {
		result := make(map[string]string)
		var message error
		for _, v := range matching {
			if value := service.extractMatchingField(v, e); value != "" {
				result[v.Field] = value
			} else {
				message = errors.New(v.Field + "字段匹配不齐全；已被强制过滤！")
				break
			}
		}
		if message == nil {
			resultField := service.fieldHandle(result, detail)
			message = service.createOneResult(result["meta_title"], url, detail.Id, resultField)
			if message != nil {
				service.createLogsInform(0, uid, detail.Id, detail.Name, error.Error(message), url)
			} else {
				// todo 因为采集数据过多，停用采集成功通知
				//service.createLogsInform(1, uid, detail.Id, detail.Name, result["meta_title"], url)
			}
		} else {
			// todo 采集过程字段空值被过滤的太多，停用采集通知
			logs.Error("在"+url+"中没有采集到全部数据：具体原因为：", error.Error(message))
			//service.createLogsInform(0, uid, detail.Id, detail.Name, error.Error(message), url)
		}
	})
	_ = collector.Visit(url)
	collector.Wait()
}

/** 获取采集器状态 todo 是否需要改为redis **/
func (service *DefaultCollectService) acquireCollectStatus(id int, field string) int8 {
	var item collect.Collect
	_ = orm.NewOrm().QueryTable(new(collect.Collect)).Filter("id", id).One(&item, field)
	return item.Status
}

/** 获取发布器状态 todo 是否需要改为redis **/
func (service *DefaultCollectService) acquirePushStatus(id int, field string) int8 {
	var item collect.Collect
	_ = orm.NewOrm().QueryTable(new(collect.Collect)).Filter("id", id).One(&item, field)
	return item.PushStatus
}

/** 采集分支 **/
func (service *DefaultCollectService) collectStatus(id int, status int8, uid int) {
	if status == 1 {
		service.collectStart(id, uid)
	} else {
		service.collectStop(id, uid, "停止爬取")
	}
}

/** 启动采集器 todo 停止无效，需要在研究 **/
func (service *DefaultCollectService) collectStart(id, uid int) {
	detail := service.FindOneDetail(id)
	source := table.ParseAttrConfigArray(detail.Source)
	collector := service.collectInstance(int(detail.Interval), 2, service.queueUrlDomain(source[0]), true)
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if service.checkSourceRule(detail.SourceRule, e.Attr("href")) {
			_ = e.Request.Visit(service.completionSrc(source[0], e.Attr("href")))
		}
	})
	collector.OnRequest(func(r *colly.Request) {
		if service.acquireCollectStatus(id, "Status") == 1 {
			if service.checkSourceRule(detail.SourceRule, r.URL.String()) {
				service.handleSourceRuleBody(r.URL.String(), uid, detail)
			}
		} else {
			r.Abort()
			logs.Warn("取消一条采集器链接；链接为：[%s]", r.URL.String())
		}
	})
	collector.OnError(func(r *colly.Response, err error) {
		if service.checkSourceRule(detail.SourceRule, r.Request.URL.String()) {
			service.createLogsInform(0, uid, detail.Id, detail.Name, error.Error(err), r.Request.URL.String())
		}
	})
	if _, err := orm.NewOrm().Update(&collect.Collect{Id: id, Status: 1}, "Status"); err != nil {
		logs.Warn("启动[%s]采集器失败！失败原因:%s", detail.Name, error.Error(err))
	} else {
		for _, v := range source {
			_ = collector.Visit(v)
		}
		service.createLogsInformStatus("开始爬取", uid, detail.Id, detail.Name)
		collector.Wait()
	}
	defer func() { service.collectStop(id, uid, "爬取任务完成") }()
}

/** 停止采集器 **/
func (service *DefaultCollectService) collectStop(id, uid int, s string) {
	detail := service.FindOneDetail(id)
	if _, err := orm.NewOrm().Update(&collect.Collect{Id: id, Status: 0}, "Status"); err != nil {
		logs.Warn("停止采集器失败！失败原因:%s", error.Error(err))
	} else {
		service.createLogsInformStatus(s, uid, detail.Id, detail.Name)
	}
}

/** 发布器分支 **/
func (service *DefaultCollectService) pushStatus(id int, status int8, uid int) {
	if status == 1 {
		service.pushStart(id, uid)
	} else {
		service.pushStop(id, uid, "停止发布数据")
	}
}

/** 获取一条可以发布的数据;从ID升序发布 **/
func (service *DefaultCollectService) pushDetail(c int) collect.Result {
	var item collect.Result
	_ = orm.NewOrm().QueryTable(new(collect.Result)).Filter("collect", c).Filter("status", 0).OrderBy("id").One(&item)
	return item
}

/** 启动发布器 **/
func (service *DefaultCollectService) pushStart(id, uid int) {
	detail := service.FindOneDetail(id)
	if _, err := orm.NewOrm().Update(&collect.Collect{Id: id, PushStatus: 1}, "PushStatus"); err != nil {
		logs.Warn("启动[%s]发布器失败！失败原因:%s", detail.Name, error.Error(err))
	} else {
		service.createLogsInformPushStatus("开始发布数据", uid, detail.Id, detail.Name)
		for {
			if service.acquirePushStatus(id, "PushStatus") == 1 {
				service.PushDetailAPI(service.pushDetail(id))
				time.Sleep(time.Duration(detail.PushTime*60*60) * time.Second)
			} else {
				logs.Warn("[%s]停止了发布任务器！已成功退出！", detail.Name)
				break
			}
		}
	}
	defer func() { service.pushStop(id, uid, "发布任务完成") }()
}

/** 停止发布器 **/
func (service *DefaultCollectService) pushStop(id, uid int, s string) {
	detail := service.FindOneDetail(id)
	if _, err := orm.NewOrm().Update(&collect.Collect{Id: id, PushStatus: 0}, "PushStatus"); err != nil {
		logs.Warn("停止发布器失败！失败原因:%s", error.Error(err))
	} else {
		service.createLogsInformPushStatus(s, uid, detail.Id, detail.Name)
	}
}

/************************************************表格渲染机制 ************************************************************/

/** 获取需要渲染的Column **/
func (service *DefaultCollectService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "爬虫标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "爬虫名称", "name": "name", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "最大并发", "name": "interval", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "发布间隔(H)", "name": "push_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "采集状态", "name": "status", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "发布状态", "name": "push_status", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "更新时间", "name": "update_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "爬虫操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultCollectService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "switch", "switch", "date"},
		"fieldName": {"id", "name", "interval", "push_time", "status", "push_status", "update_time"},
		"action":    {"", "", "", "", beego.URLFor("Collect.Status"), beego.URLFor("Collect.Status"), ""},
	}
	return result
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultCollectService) DataTableButtons() []*table.TableButtons {
	var array []*table.TableButtons
	array = append(array, &table.TableButtons{
		Text:      "创建爬虫",
		ClassName: "btn btn-sm btn-alt-success mt-1 jump_urls",
		Attribute: map[string]string{"data-action": beego.URLFor("Collect.Create")},
	})
	array = append(array, &table.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Collect.Delete")},
	})
	return array
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultCollectService) TableButtonsType() []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "修改爬虫",
			ClassName: "btn btn-sm btn-alt-primary jump_urls",
			Attribute: map[string]string{"data-action": beego.URLFor("Collect.Edit", ":id", "__ID__")},
		},
		{
			Text:      "采集结果",
			ClassName: "btn btn-sm btn-alt-primary jump_urls",
			Attribute: map[string]string{"data-action": beego.URLFor("Collector.Index", ":id", "__ID__")},
		},
		{
			Text:      "删除爬虫",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Collect.Delete"),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}

/** 处理分页 **/
func (service *DefaultCollectService) PageListItems(length, draw, page int, search string) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(collect.Collect))
	recordsTotal, _ := qs.Count()
	if search != "" {
		qs = qs.Filter("name__icontains", search)
	}
	_, _ = qs.Limit(length, length*(page-1)).ValuesList(&lists, "id", "name", "interval", "push_time", "status", "push_status", "update_time")
	data := map[string]interface{}{
		"draw":            draw,         // 请求次数
		"recordsFiltered": recordsTotal, // 从多少条里面筛选
		"recordsTotal":    recordsTotal, // 总条数
		"data":            lists,        // 筛选结果
	}
	return data
}
