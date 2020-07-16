package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/dictionaries"
	ccst "github.com/go-cc/cc-jianfan"
	"github.com/go-playground/validator"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type DefaultDictionariesService struct{}

/** 获取爬虫配置解析成map **/
func (service *DefaultDictionariesService) ConfigMaps() map[string]interface{} {
	var item []*dictionaries.DictConfig
	maps := make(map[string]interface{})
	_, _ = orm.NewOrm().QueryTable(new(dictionaries.DictConfig)).All(&item)
	for _, v := range item {
		maps[v.Name] = v.Value
	}
	if maps["status"] == nil {
		maps = service.ConfigMaps()
	}
	return maps
}

/** 根据name获取一条数据 **/
func (service *DefaultDictionariesService) ConfigByName(name string) dictionaries.DictConfig {
	var item dictionaries.DictConfig
	_ = orm.NewOrm().QueryTable(new(dictionaries.DictConfig)).Filter("name", name).One(&item)
	return item
}

/** 获取一条分类数据 **/
func (service *DefaultDictionariesService) One(id int) dictionaries.Dictionaries {
	var item dictionaries.Dictionaries
	_ = orm.NewOrm().QueryTable(new(dictionaries.Dictionaries)).Filter("id", id).One(&item)
	return item
}

/*** 获取一条详情 **/
func (service *DefaultDictionariesService) DetailOne(id int) map[string]string {
	result := make(map[string]string)
	detail := service.One(id)
	result["logs"] = detail.Logs
	result["title"] = detail.Title
	result["status"] = strconv.Itoa(int(detail.Status))
	result["update_time"] = beego.Date(detail.UpdateTime, "Y年m月d H:i:s")
	maps := map[string]string{
		"name": detail.Name, "title": detail.Title,
		"keywords": detail.Keywords, "description": detail.Description,
		"source": detail.Source, "initial": detail.Initial,
		"context": detail.Context,
	}
	stringJson, _ := json.Marshal(maps)
	result["result"] = string(stringJson)
	return result
}

/** 批量删除结果 **/
func (service *DefaultDictionariesService) DeleteArray(array []int) (message error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, message = orm.NewOrm().Delete(&dictionaries.Dictionaries{Id: v}); message != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if message == nil {
		_ = orm.NewOrm().Commit()
	}
	return message
}

/** 爬虫操纵命令 **/
func (service *DefaultDictionariesService) CateInstanceBegin(instruct map[string]interface{}) {
	if instruct["field"].(string) == "status" {
		if instruct["status"].(int) == 1 {
			service.Start(instruct["uid"].(int))
		} else {
			service.Stop(instruct["uid"].(int), "停止爬取查字词典")
		}
	}
	if instruct["field"].(string) == "send_status" {
		if instruct["status"].(int) == 1 {
			service.StartPush(instruct["uid"].(int))
		} else {
			service.StopPush(instruct["uid"].(int), "停止发布查字词典")
		}
	}
}

/** 停止爬虫 **/
func (service *DefaultDictionariesService) Stop(uid int, message string) {
	if _, err := orm.NewOrm().Update(&dictionaries.DictConfig{Id: 5, Value: "0"}, "Value"); err != nil {
		logs.Warn("停止查字词典采集器失败！失败原因:%s", error.Error(err))
	} else {
		service.createLogsInformStatus(message, uid)
	}
}

/** 启动爬虫 **/
func (service *DefaultDictionariesService) Start(uid int) {
	var (
		object      = DefaultCollectService{}
		config      = service.ConfigMaps()
		interval, _ = strconv.Atoi(config["interval"].(string))
	)
	collector := object.collectInstance(interval, 3, "www.chazidian.com", true)
	collector.OnRequest(func(r *colly.Request) {
		if status, _ := strconv.Atoi(service.ConfigMaps()["status"].(string)); status == 0 {
			r.Abort()
		}
	})
	collector.OnHTML(".pinyin_sou", func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
			_ = e.Request.Visit(el.Attr("href"))
		})
	})
	collector.OnHTML("html", func(e *colly.HTMLElement) {
		if regexp.MustCompile("https://www.chazidian.com/ci_([a-z]+)_([0-9]+)/$").MatchString(e.Request.URL.String()) {
			if id, message := service.collectOnResultCate(e); message != nil {
				service.createLogsInformStatus("查字词典分类采集到一条非法数据；<a href='"+e.Request.URL.String()+"'>查看原文</a>;非法原因："+error.Error(message)+";", uid)
			} else {
				service.collectVisitDetailHref(id, e)
			}
		}
	})
	collector.OnHTML("html", func(e *colly.HTMLElement) {
		if regexp.MustCompile(`https://www.chazidian.com/([a-z]+)_([a-z]+)_(.*?)/\?dict=([0-9]+)$`).MatchString(e.Request.URL.String()) {
			if _, message := service.collectOnResultDetail(e); message != nil {
				service.createLogsInformStatus("查字词典详情采集到一条非法数据；<a href='"+e.Request.URL.String()+"'>查看原文</a>;非法原因："+error.Error(message)+";", uid)
			}
		}
	})
	collector.OnError(func(r *colly.Response, err error) {
		service.createLogsInformStatus("查字词典在采集过程中出现一次失败的采集；<a href='"+r.Request.URL.String()+"'>查看原文</a>;错误原因："+error.Error(err)+";已经跳过错误继续采集；", uid)
	})
	if _, err := orm.NewOrm().Update(&dictionaries.DictConfig{Id: 5, Value: "1"}, "Value"); err != nil {
		logs.Warn("启动查字词典采集器失败！失败原因:%s", error.Error(err))
	} else {
		_ = collector.Visit("https://www.chazidian.com/ci_pinyin/")
		service.createLogsInformStatus("开始爬取查字词典", uid)
		collector.Wait()
	}
	defer func() { service.Stop(uid, "查字词典爬取任务完成") }()
}

/** 启动数据发布器 **/
func (service *DefaultDictionariesService) StartPush(uid int) {
	if _, err := orm.NewOrm().Update(&dictionaries.DictConfig{Id: 8, Value: "1"}, "Value"); err != nil {
		logs.Warn("启动[%s]发布器失败！失败原因:%s", "查字词典", error.Error(err))
	} else {
		config := service.ConfigMaps()
		pushTime, _ := strconv.ParseFloat(config["push_time"].(string), 64)
		service.createLogsInformStatusPush("查字词典发布任务启动", uid)
		for {
			if status, _ := strconv.Atoi(service.ConfigMaps()["send_status"].(string)); status == 1 {
				service.PushDetailAPI(service.pushDetail())
				time.Sleep(time.Duration(pushTime*60*60) * time.Second)
			} else {
				logs.Warn("[%s]停止了发布任务器！已成功退出！", "查字词典")
				break
			}
		}
	}
	defer func() { service.StopPush(uid, "查字词典发布任务完成") }()
}

/** 停止分发布 **/
func (service *DefaultDictionariesService) StopPush(uid int, message string) {
	if _, err := orm.NewOrm().Update(&dictionaries.DictConfig{Id: 8, Value: "0"}, "Value"); err != nil {
		logs.Warn("停止查字词典发布器失败！失败原因:%s", error.Error(err))
	} else {
		service.createLogsInformStatusPush(message, uid)
	}
}

/** 发布一条分类数据;返回格式：{status:[false|true],message:[message]} **/
func (service *DefaultDictionariesService) PushDetailAPI(item dictionaries.Dictionaries) {
	config := service.ConfigMaps()
	detail := service.DetailOne(item.Id)
	if resp, err := http.PostForm(config["send_domain"].(string), url.Values{
		"password": {beego.AppConfig.String("collect_collect_password")},
		"title":    {item.Title}, "source": {item.Source}, "result": {detail["result"]},
		"pid": {strconv.Itoa(item.Pid)}, "id": {strconv.Itoa(item.Id)},
		"update": {strconv.FormatInt(item.UpdateTime.Unix(), 10)},
		"create": {strconv.FormatInt(item.CreateTime.Unix(), 10)},
	}); err == nil {
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			service.pushDetailAPIResult(item.Id, -1, err)
		} else {
			result := make(map[string]interface{})
			if err := json.Unmarshal(body, &result); err != nil && len(result) > 0 {
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

/*********************************以下为表格渲染*******************************************8/

/** 获取需要渲染的Column **/
func (service *DefaultDictionariesService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "名称", "name": "name", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "首字母", "name": "initial", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "源地址", "name": "source", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "发布状态", "name": "status", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "更新时间", "name": "update_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "爬虫操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/*** 获取分页数据 **/
func (service *DefaultDictionariesService) PageListItems(length, draw, page int, search string) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(dictionaries.Dictionaries)).Filter("pid", 0)
	if search != "" {
		qs = qs.Filter("initial__icontains", search)
	}
	recordsTotal, _ := qs.Count()
	_, _ = qs.Limit(length, length*(page-1)).ValuesList(&lists, "id", "name", "initial", "source", "status", "update_time")
	for _, v := range lists {
		v[4] = service.Int2HtmlStatus(v[4], v[0], beego.URLFor("Dictionaries.Push"))
	}
	data := map[string]interface{}{
		"draw":            draw,         // 请求次数
		"recordsFiltered": recordsTotal, // 从多少条里面筛选
		"recordsTotal":    recordsTotal, // 总条数
		"data":            lists,        // 筛选结果
	}
	return data
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultDictionariesService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string", "date"},
		"fieldName": {"id", "name", "initial", "source", "status", "update_time"},
		"action":    {"", "", "", "", "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultDictionariesService) TableButtonsType() []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "查看",
			ClassName: "btn btn-sm btn-alt-primary open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Dictionaries.Detail", ":id", "__ID__", ":popup", 1),
				"data-area": "620px,400px",
			},
		},
		{
			Text:      "删除",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Dictionaries.Delete"),
				"data-ids":    "__ID__",
			},
		},
		{
			Text:      "列表",
			ClassName: "btn btn-sm btn-alt-success jump_urls",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Lexicon.Index", ":id", "__ID__"),
			},
		},
	}
	return buttons
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultDictionariesService) DataTableButtons() []*table.TableButtons {
	var array []*table.TableButtons
	statusMaps := map[string]string{"0": "关闭爬虫", "1": "启动爬虫"}
	statusCls := map[string]string{"0": "btn-alt-danger", "1": "btn-alt-success"}
	sendMaps := map[string]string{"0": "停止发布", "1": "启动发布"}
	sendCls := map[string]string{"0": "btn-alt-danger", "1": "btn-alt-success"}
	config := service.ConfigMaps()
	array = append(array, &table.TableButtons{
		Text:      "爬虫配置",
		ClassName: "btn btn-sm btn-alt-warning mt-1 open_iframe",
		Attribute: map[string]string{
			"href":      beego.URLFor("Dictionaries.Config", ":popup", 1),
			"data-area": "600px,400px",
		},
	})
	array = append(array, &table.TableButtons{
		Text:      "发布选中",
		ClassName: "btn btn-sm btn-alt-primary mt-1 ids_enables",
		Attribute: map[string]string{
			"data-action": beego.URLFor("Dictionaries.Push"),
			"data-field":  "status",
		},
	})
	array = append(array, &table.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Dictionaries.Delete")},
	})
	for k, v := range statusMaps {
		if config["status"].(string) != k {
			array = append(array, &table.TableButtons{
				Text:      v,
				ClassName: "btn btn-sm " + statusCls[k] + " mt-1 handle_collect",
				Attribute: map[string]string{
					"data-action": beego.URLFor("Dictionaries.Status"),
					"data-status": k,
					"data-field":  "status",
				},
			})
		}
	}
	for k, v := range sendMaps {
		if config["send_status"].(string) != k {
			array = append(array, &table.TableButtons{
				Text:      v,
				ClassName: "btn btn-sm " + sendCls[k] + " mt-1 handle_collect",
				Attribute: map[string]string{
					"data-action": beego.URLFor("Dictionaries.Status"),
					"data-status": k,
					"data-field":  "send_status",
				},
			})

		}
	}
	return array
}

/*****************以下为私有方法 *************************/
/** 发送采集爬虫状态通知 **/
func (service *DefaultDictionariesService) createLogsInformStatus(status string, receiver int) {
	var htmlInfo string
	nowTime := beego.Date(time.Now(), "m月d日 H:i")
	htmlInfo = fmt.Sprintf(`
		名为<a href="javascript:void(0);">查字词典</a>的采集器;在%s的时候%s；
		您可以<a target="_blank" href="%s">查看爬取结果</a>;`,
		nowTime, status, beego.URLFor("Dictionaries.Cate"),
	)
	go new(DefaultInformService).SendSocketInform([]int{receiver}, 0, 0, 2, htmlInfo)
}

/** 发送发布分类状态通知 **/
func (service *DefaultDictionariesService) createLogsInformStatusPush(status string, receiver int) {
	var htmlInfo string
	nowTime := beego.Date(time.Now(), "m月d日 H:i")
	htmlInfo = fmt.Sprintf(`
		名为<a href="javascript:void(0);">查字词典</a>的发布器;在%s的时候%s；
		您可以<a target="_blank" href="%s">查看发布结果</a>;`,
		nowTime, status, beego.URLFor("Dictionaries.Cate"),
	)
	go new(DefaultInformService).SendSocketInform([]int{receiver}, 0, 0, 2, htmlInfo)
}

/** 创建一条分类结果数据 */
func (service *DefaultDictionariesService) createOneResult(item *dictionaries.Dictionaries) (id int64, err error) {
	verify := DefaultBaseVerify{}
	if message := verify.Begin().Struct(item); message != nil {
		return 0, errors.New(verify.Translate(message.(validator.ValidationErrors)))
	}
	if one := service.oneResultLink(item.Source); one.Id > 0 {
		item.Id = one.Id
		_, err = orm.NewOrm().Update(item)
		id = int64(one.Id)
	} else {
		id, err = orm.NewOrm().Insert(item)
	}
	return id, err
}

/** 根据链接获取一条结果数据 **/
func (service *DefaultDictionariesService) oneResultLink(url string) dictionaries.Dictionaries {
	var item dictionaries.Dictionaries
	_ = orm.NewOrm().QueryTable(new(dictionaries.Dictionaries)).Filter("source", url).One(&item)
	return item
}

/** 将状态码转为html **/
func (service *DefaultDictionariesService) Int2HtmlStatus(val interface{}, id interface{}, url string) string {
	t := table.BuilderTable{}
	return t.AnalysisTypeSwitch(map[string]interface{}{"status": val, "id": id}, "status", url, map[int64]string{-1: "已失败", 0: "待发布", 1: "已发布"})
}

/**  转为pop提示 **/
func (service *DefaultDictionariesService) Substr2HtmlSpan(s string, start, end int) string {
	html := fmt.Sprintf(`<span class="badge badge-primary js-tooltip" data-placement="top" data-toggle="tooltip" data-original-title="%s">%s...</span>`, s, beego.Substr(s, start, end))
	return html
}

/**  转为pop提示 **/
func (service *DefaultDictionariesService) Substr2HtmlHref(s string, start, end int) string {
	html := fmt.Sprintf(`<a href="%s" target="_blank" class="badge badge-primary js-tooltip" data-placement="top" data-toggle="tooltip" data-original-title="%s">%s...</a>`, s, s, beego.Substr(s, start, end))
	return html
}

/** 获取一条可以发布的数据;从ID升序发布 **/
func (service *DefaultDictionariesService) pushDetail() dictionaries.Dictionaries {
	var item dictionaries.Dictionaries
	_ = orm.NewOrm().QueryTable(new(dictionaries.Dictionaries)).Filter("status", 0).OrderBy("pid").One(&item)
	return item
}

/**  更新分类发布结果 **/
func (service *DefaultDictionariesService) pushDetailAPIResult(id int, status int8, message error) {
	if _, err := orm.NewOrm().Update(&dictionaries.Dictionaries{Id: id, Status: status, Logs: error.Error(message)}, "Status", "Logs"); err != nil {
		logs.Warn("更新发布结果失败；失败原因：%s", error.Error(err))
	}
}

/** 采集分类结果解析创建分类 **/
func (service *DefaultDictionariesService) collectOnResultCate(e *colly.HTMLElement) (int64, error) {
	object := DefaultCollectService{}
	result := map[string]string{
		"title":   object.eliminateTrim(e.DOM.Find("head > title").Text(), []string{" - 查字典"}),
		"name":    e.ChildText(".main_left > div:last-child > .box_head > .box_title > h2 > b"),
		"initial": e.ChildText(".main_left > div:nth-child(3) > .box_head > .box_title > h2 > b"),
		"context": "NONE",
	}
	result["keyword"], _ = e.DOM.Find("head > meta[name=keywords]").Attr("content")
	result["describe"], _ = e.DOM.Find("head > meta[name=description]").Attr("content")
	return service.createOneResultDict(result, e.Request.URL.String(), 0)
}

/** 创建一条数据库记录 **/
func (service *DefaultDictionariesService) createOneResultDict(res map[string]string, source string, pid int) (int64, error) {
	translate := service.Translate(res)
	item := &dictionaries.Dictionaries{
		Title:       translate["title"],
		Name:        translate["name"],
		Initial:     translate["initial"],
		Context:     translate["context"],
		Keywords:    translate["keyword"],
		Description: translate["describe"],
		Status:      0, Source: source, Pid: pid,
	}
	return service.createOneResult(item)
}

/** 语言转换 **/
func (service *DefaultDictionariesService) Translate(res map[string]string) map[string]string {
	config := service.ConfigMaps()
	translate, _ := strconv.Atoi(config["translate"].(string))
	for k, v := range res {
		if translate > 0 {
			if translate == 1 {
				res[k] = ccst.S2T(v)
			} else {
				res[k] = ccst.T2S(v)
			}
		}
		res[k] = service.deleteScriptSpace(service.deleteExtraSpace(service.eliminateTrim(res[k], []string{"\n"})))
	}
	return res
}

/** 剔除字符规定字符 **/
func (service *DefaultDictionariesService) eliminateTrim(str string, symbol []string) string {
	for _, v := range symbol {
		str = strings.Replace(str, v, "", -1)
	}
	return str
}

/** 过滤SCRIPT **/
func (service *DefaultDictionariesService) deleteScriptSpace(src string) string {
	re, _ := regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	return re.ReplaceAllString(src, "")
}

/** 去除连续两个以上的空格 **/
func (service *DefaultDictionariesService) deleteExtraSpace(s string) string {
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

/** 采集结果创建详情数据 **/
func (service *DefaultDictionariesService) collectOnResultDetail(e *colly.HTMLElement) (int64, error) {
	var (
		parent     = e.Request.URL.Query()["dict"]
		reagent, _ = strconv.Atoi(parent[0])
		object     = DefaultCollectService{}
	)
	result := map[string]string{
		"title":   object.eliminateTrim(e.DOM.Find("head > title").Text(), []string{" - 查字典"}),
		"name":    e.ChildText(".bktitle > h1"),
		"initial": "NONE",
	}
	result["context"], _ = e.DOM.Find(".content > div:last-child").Html()
	result["keyword"], _ = e.DOM.Find("head > meta[name=keywords]").Attr("content")
	result["describe"], _ = e.DOM.Find("head > meta[name=description]").Attr("content")
	return service.createOneResultDict(result, e.Request.URL.String(), reagent)
}

/** 匹配创建所有详情链接 **/
func (service *DefaultDictionariesService) collectVisitDetailHref(id int64, e *colly.HTMLElement) {
	e.ForEach(".main_left > div:last-child >.box_content a[href]", func(_ int, el *colly.HTMLElement) {
		href := el.Attr("href")
		if o := strings.Index(href, "https://www.chazidian.com"); o < 0 {
			href = "https://www.chazidian.com" + href
		}
		if regexp.MustCompile(`https://www.chazidian.com/([a-z]+)_([a-z]+)_([a-zA-Z0-9]{32})/$`).MatchString(href) {
			href = href + "?dict=" + strconv.FormatInt(id, 10)
		}
		_ = e.Request.Visit(href)
	})
}
