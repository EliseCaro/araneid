package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/axgle/mahonia"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/movie"
	"github.com/go-playground/validator"
	"github.com/gocolly/colly"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type DefaultMovieService struct{}

/** 获取爬虫配置解析成map **/
func (service *DefaultMovieService) ConfigMaps() map[string]interface{} {
	var item []*movie.ConfigMovie
	maps := make(map[string]interface{})
	_, _ = orm.NewOrm().QueryTable(new(movie.ConfigMovie)).All(&item)
	for _, v := range item {
		maps[v.Name] = v.Value
	}
	if maps["status"] == nil {
		maps = service.ConfigMaps()
	}
	return maps
}

/** 根据链接获取一条结果数据 **/
func (service *DefaultMovieService) oneResultLink(url string) movie.Movie {
	var item movie.Movie
	_ = orm.NewOrm().QueryTable(new(movie.Movie)).Filter("source", url).One(&item)
	return item
}

/** 作品名称提取 **/
func (service *DefaultMovieService) collectOnResultMovieName(e *colly.HTMLElement) string {
	str := e.ChildText(".Main3 > .m_T3 > h1")
	str = service.coverGBKToUTF8(str)
	str = new(DefaultCollectService).eliminateTrim(str, []string{"剧情介绍"})
	return str
}

/** 主演提取返回json字符串 **/
func (service *DefaultMovieService) collectOnResultMovieActor(e *colly.HTMLElement) string {
	str := e.ChildText(".g-box7 > .box > .txt > p:nth-child(1)")
	str = service.coverGBKToUTF8(str)
	str = new(DefaultCollectService).eliminateTrim(str, []string{"主                演："})
	maps := strings.Split(str, " ")
	result, _ := json.Marshal(maps)
	return string(result)
}

/** 作品导演提取 **/
func (service *DefaultMovieService) collectOnResultMovieDirector(e *colly.HTMLElement) string {
	str := e.ChildText(".g-box7 > .box > .txt > p:nth-child(3) span:nth-child(1)")
	str = service.coverGBKToUTF8(str)
	str = new(DefaultCollectService).eliminateTrim(str, []string{"导                演："})
	return str
}

/** 作品年份提取 **/
func (service *DefaultMovieService) collectOnResultMovieYear(e *colly.HTMLElement) string {
	str := e.ChildText(".g-box7 > .box > .txt > p:nth-child(4) span:nth-child(1)")
	str = service.coverGBKToUTF8(str)
	str = new(DefaultCollectService).eliminateTrim(str, []string{"年                份："})
	return str
}

/** 作品地区提取 **/
func (service *DefaultMovieService) collectOnResultMovieDistrict(e *colly.HTMLElement) string {
	str := e.ChildText(".g-box7 > .box > .txt > p:nth-child(4) span:nth-child(2)")
	str = service.coverGBKToUTF8(str)
	str = new(DefaultCollectService).eliminateTrim(str, []string{"地                区："})
	maps := strings.Split(str, " ")
	result, _ := json.Marshal(maps)
	return string(result)
}

/** 提取类型 **/
func (service *DefaultMovieService) collectOnResultMovieGenre(e *colly.HTMLElement) string {
	str := e.ChildText(".g-box7 > .box > .txt > p:nth-child(7)")
	str = service.coverGBKToUTF8(str)
	str = new(DefaultCollectService).eliminateTrim(str, []string{"类                型："})
	maps := strings.Split(str, " ")
	result, _ := json.Marshal(maps)
	return string(result)
}

/** 提取封面 **/
func (service *DefaultMovieService) collectOnResultMovieCover(e *colly.HTMLElement) string {
	str, _ := e.DOM.Find(".g-box7 > .pic > img").Attr("src")
	if str != "" {
		completion := new(DefaultCollectService).completionSrc(e.Request.URL.String(), str)
		if strings.Index(completion, "{") == -1 {
			str = new(DefaultCollectService).download(completion)
		} else {
			str = ""
			logs.Error("提取封面图失败！URL地址为：%s", e.Request.URL.String())
		}
	}
	return str
}

/** 提取剧情 **/
func (service *DefaultMovieService) collectOnResultMovieContext(e *colly.HTMLElement) string {
	str := e.ChildText(".m_Left3 > .m_jq")
	str = service.coverGBKToUTF8(str)
	return str
}

/** 提取演员详情 **/
func (service *DefaultMovieService) collectOnResultMovieActorCtx(e *colly.HTMLElement) string {
	return ""
}

/** 提取标题 **/
func (service *DefaultMovieService) collectOnResultMovieTitle(e *colly.HTMLElement) string {
	str := e.ChildText("head > title")
	str = new(DefaultCollectService).eliminateTrim(str, []string{"_剧情吧"})
	str = service.coverGBKToUTF8(str)
	return str
}

/** 提取关键字 **/
func (service *DefaultMovieService) collectOnResultMovieKeywords(e *colly.HTMLElement) string {
	str, _ := e.DOM.Find("head > meta[name=keywords]").Attr("content")
	str = service.coverGBKToUTF8(str)
	return str
}

/** 提取描述 **/
func (service *DefaultMovieService) collectOnResultMovieDescription(e *colly.HTMLElement) string {
	str, _ := e.DOM.Find("head > meta[name=description]").Attr("content")
	str = service.coverGBKToUTF8(str)
	return str
}

/** 编码转换 **/
func (service *DefaultMovieService) coverGBKToUTF8(src string) string {
	return mahonia.NewDecoder("gbk").ConvertString(src)
}

/** 解析一部作品基本值 **/
func (service *DefaultMovieService) collectOnResultMovie(e *colly.HTMLElement) *movie.Movie {
	return &movie.Movie{
		Name:        service.collectOnResultMovieName(e),
		Actor:       service.collectOnResultMovieActor(e),
		Director:    service.collectOnResultMovieDirector(e),
		Year:        service.collectOnResultMovieYear(e),
		District:    service.collectOnResultMovieDistrict(e),
		Genre:       service.collectOnResultMovieGenre(e),
		Cover:       service.collectOnResultMovieCover(e),
		Context:     service.collectOnResultMovieContext(e),
		ActorCtx:    service.collectOnResultMovieActorCtx(e),
		Title:       service.collectOnResultMovieTitle(e),
		Keywords:    service.collectOnResultMovieKeywords(e),
		Description: service.collectOnResultMovieDescription(e),
		Source:      e.Request.URL.String(),
	}
}

/** 入库一条数据 **/
func (service *DefaultMovieService) collectOnResultMovieInsert(res *movie.Movie, callback func(int64, error)) {
	if message := new(DefaultBaseVerify).Begin().Struct(res); message != nil {
		callback(0, errors.New(new(DefaultBaseVerify).Translate(message.(validator.ValidationErrors))))
	}
	if one := service.oneResultLink(res.Source); one.Id <= 0 {
		callback(orm.NewOrm().Insert(res))
	} else {
		callback(int64(one.Id), nil)
	}
}

/** 启动爬虫 **/
func (service *DefaultMovieService) Start(uid int) {
	var (
		config      = service.ConfigMaps()
		interval, _ = strconv.Atoi(config["interval"].(string))
	)
	collector := new(DefaultCollectService).collectInstance(interval, 3, "www.juqingba.cn", true)
	collector.OnRequest(func(r *colly.Request) {
		if status, _ := strconv.Atoi(service.ConfigMaps()["status"].(string)); status == 0 {
			r.Abort()
		}
	})
	collector.OnHTML("html", func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) { _ = e.Request.Visit(el.Attr("href")) })
	})
	collector.OnHTML("html", func(e *colly.HTMLElement) { // 监听电影
		if regexp.MustCompile("https://www.juqingba.cn/dianyingjuqing/([0-9]+).html$").MatchString(e.Request.URL.String()) {
			result := service.collectOnResultMovie(e)
			result.ClassType = 1
			go service.collectOnResultMovieInsert(result, func(id int64, err error) {
				if id > 0 {
					new(DefaultCollectService).createLogsInform(1, uid, "剧情", result.Title, e.Request.URL.String())
				} else {
					new(DefaultCollectService).createLogsInform(0, uid, "剧情", error.Error(err), e.Request.URL.String())
				}
			})
		}
	})
	collector.OnHTML("html", func(e *colly.HTMLElement) { // 监听电视
		if regexp.MustCompile("https://www.juqingba.cn/zjuqing/([0-9]+).html$").MatchString(e.Request.URL.String()) {

		}
	})
	collector.OnError(func(r *colly.Response, err error) {
		new(DefaultCollectService).createLogsInform(0, uid, "剧情采集", error.Error(err), r.Request.URL.String())
	})
	if _, err := orm.NewOrm().Update(&movie.ConfigMovie{Id: 1, Value: "1"}, "Value"); err == nil {
		_ = collector.Visit("https://www.juqingba.cn/")
		service.createLogsInformStatus("开始爬取剧情采集", uid)
		collector.Wait()
	} else {
		logs.Warn("启动剧情采集器失败！失败原因:%s", error.Error(err))
	}
	defer func() { service.Stop(uid, "剧情采集爬取任务完成") }()
}

/** 停止爬虫 **/
func (service *DefaultMovieService) Stop(uid int, message string) {
	if _, err := orm.NewOrm().Update(&movie.ConfigMovie{Id: 1, Value: "0"}, "Value"); err != nil {
		logs.Warn("停止剧情采集器失败！失败原因:%s", error.Error(err))
	} else {
		service.createLogsInformStatus(message, uid)
	}
}

/** 发送采集爬虫状态通知 **/
func (service *DefaultMovieService) createLogsInformStatus(status string, receiver int) {
	var htmlInfo string
	nowTime := beego.Date(time.Now(), "m月d日 H:i")
	htmlInfo = fmt.Sprintf(`
		名为<a href="javascript:void(0);">剧情采集</a>的采集器;在%s的时候%s；
		您可以<a target="_blank" href="%s">查看爬取结果</a>;`,
		nowTime, status, beego.URLFor("Movie.Index"),
	)
	go new(DefaultInformService).SendSocketInform([]int{receiver}, 0, 0, 4, htmlInfo)
}

/** 爬虫操纵命令 **/
func (service *DefaultMovieService) CateInstanceBegin(instruct map[string]interface{}) {
	if instruct["field"].(string) == "status" {
		if instruct["status"].(int) == 1 {
			service.Start(instruct["uid"].(int))
		} else {
			service.Stop(instruct["uid"].(int), "停止爬取剧情采集")
		}
	}
	//if instruct["field"].(string) == "send_status" {
	//	if instruct["status"].(int) == 1 {
	//		service.StartPush(instruct["uid"].(int))
	//	} else {
	//		service.StopPush(instruct["uid"].(int), "停止发布剧情采集")
	//	}
	//}
}

/*********************************以下为表格渲染*******************************************8/

/** 获取需要渲染的Column **/
func (service *DefaultMovieService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "作品标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "作品名称", "name": "name", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "作品年份", "name": "year", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "内容类型", "name": "class_type", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "简短标题", "name": "short", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "对方地址", "name": "source", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "发布状态", "name": "status", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "更新时间", "name": "update_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "爬虫操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultMovieService) DataTableButtons() []*table.TableButtons {
	var config = service.ConfigMaps()
	var array []*table.TableButtons
	array = append(array, &table.TableButtons{
		Text:      "发布选中",
		ClassName: "btn btn-sm btn-alt-primary mt-1 ids_enables",
		Attribute: map[string]string{
			"data-action": beego.URLFor("Movie.Push"),
			"data-field":  "status",
		},
	})
	array = append(array, &table.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Movie.Delete")},
	})
	for k, v := range []map[string]string{
		{"title": "关闭爬虫", "statusCls": "btn-alt-danger"},
		{"title": "启动爬虫", "statusCls": "btn-alt-success"},
	} {
		if intStatus, err := strconv.Atoi(config["status"].(string)); err == nil && intStatus != k {
			array = append(array, &table.TableButtons{
				Text:      v["title"],
				ClassName: "btn btn-sm " + v["statusCls"] + " mt-1 handle_collect",
				Attribute: map[string]string{
					"data-action": beego.URLFor("Movie.Status"),
					"data-status": strconv.Itoa(k),
					"data-field":  "status",
				},
			})
		}
	}
	return array
}

/*** 获取分页数据 **/
func (service *DefaultMovieService) PageListItems(length, draw, page int, search string) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(movie.Movie))
	if search != "" {
		qs = qs.Filter("title__icontains", search)
	}
	recordsTotal, _ := qs.Count()
	_, _ = qs.Limit(length, length*(page-1)).ValuesList(&lists, "id", "name", "year", "class_type", "short", "source", "status", "update_time")
	for _, v := range lists {
		v[4] = service.Int2HtmlStatus(v[4], v[0], beego.URLFor("Movie.Push"))
	}
	data := map[string]interface{}{
		"draw":            draw,         // 请求次数
		"recordsFiltered": recordsTotal, // 从多少条里面筛选
		"recordsTotal":    recordsTotal, // 总条数
		"data":            lists,        // 筛选结果
	}
	return data
}

/** 将状态码转为html **/
func (service *DefaultMovieService) Int2HtmlStatus(val interface{}, id interface{}, url string) string {
	t := table.BuilderTable{}
	return t.AnalysisTypeSwitch(map[string]interface{}{"status": val, "id": id}, "status", url, map[int64]string{-1: "已失败", 0: "待发布", 1: "已发布"})
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultMovieService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string", "string", "string", "date"},
		"fieldName": {"id", "name", "year", "class_type", "short", "source", "status", "update_time"},
		"action":    {"", "", "", "", "", "", "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultMovieService) TableButtonsType() []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "查看",
			ClassName: "btn btn-sm btn-alt-primary open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Movie.Detail", ":id", "__ID__", ":popup", 1),
				"data-area": "620px,400px",
			},
		},
		{
			Text:      "删除",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Movie.Delete"),
				"data-ids":    "__ID__",
			},
		},
		{
			Text:      "剧情",
			ClassName: "btn btn-sm btn-alt-success jump_urls",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Movie.Index", ":id", "__ID__"),
			},
		},
	}
	return buttons
}
