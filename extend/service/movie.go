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

/** 批量删除结果 **/
func (service *DefaultMovieService) ArrayDelete(array []int) (message error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if service.extArticle(v).Id == 0 {
			if _, message = orm.NewOrm().Delete(&movie.Movie{Id: v}); message != nil {
				_ = orm.NewOrm().Rollback()
				break
			}
		} else {
			message = errors.New("被挂载或者存在文档数据，不允许删除使用中的栏目！")
			break
		}
	}
	if message == nil {
		_ = orm.NewOrm().Commit()
	}
	return message
}

/*** 获取一条详情 **/
func (service *DefaultMovieService) DetailOne(id int) map[string]string {
	result := make(map[string]string)
	detail := service.One(id)
	result["logs"] = detail.Logs
	result["title"] = detail.Title
	result["status"] = strconv.Itoa(int(detail.Status))
	result["update_time"] = beego.Date(detail.UpdateTime, "Y年m月d H:i:s")
	maps := map[string]string{
		"name": detail.Name, "actor": detail.Actor,
		"director": detail.Director, "year": detail.Year,
		"district": detail.District, "genre": detail.Genre,
		"cover": detail.Cover, "context": detail.Context,
		"actor_ctx": detail.ActorCtx, "title": detail.Title,
		"keywords": detail.Keywords, "description": detail.Description,
		"short": detail.Short, "source": detail.Source,
	}
	stringJson, _ := json.Marshal(maps)
	result["result"] = string(stringJson)
	return result
}

/** 检测分类是否存在文章 **/
func (service *DefaultMovieService) extArticle(id int) movie.EpisodeMovie {
	var item movie.EpisodeMovie
	_ = orm.NewOrm().QueryTable(new(movie.EpisodeMovie)).Filter("pid", id).One(&item)
	return item
}

/** 根据name获取一条数据 **/
func (service *DefaultMovieService) ConfigByName(name string) movie.ConfigMovie {
	var item movie.ConfigMovie
	_ = orm.NewOrm().QueryTable(new(movie.ConfigMovie)).Filter("name", name).One(&item)
	return item
}

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

/** 根据链接获取一条详情结果数据 **/
func (service *DefaultMovieService) oneResultLinkDetail(url string) movie.EpisodeMovie {
	var item movie.EpisodeMovie
	_ = orm.NewOrm().QueryTable(new(movie.EpisodeMovie)).Filter("source", url).One(&item)
	return item
}

/** 获取一条 **/
func (service *DefaultMovieService) One(id int) movie.Movie {
	var item movie.Movie
	_ = orm.NewOrm().QueryTable(new(movie.Movie)).Filter("id", id).One(&item)
	return item
}

/** 编码转换 **/
func (service *DefaultMovieService) coverGBKToUTF8(src string) string {
	return mahonia.NewDecoder("gbk").ConvertString(src)
}

/** 解析通知错误 **/
func (service *DefaultMovieService) analysisMessage(title string, err error) (int, string) {
	if err != nil {
		return 0, err.Error()
	} else {
		return 1, title
	}
}

/******************************************************电影作品提取 ****************************************/
/******************************************************电影作品提取 ****************************************/

/** 解析基本值 **/
func (service *DefaultMovieService) collectOnResultMovie(e *colly.HTMLElement) *movie.Movie {
	return &movie.Movie{
		Name:        service.collectOnResultMovieName(e),
		Short:       service.collectOnResultMovieShort(e),
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
		ClassType:   0,
	}
}

/** 名称提取 **/
func (service *DefaultMovieService) collectOnResultMovieName(e *colly.HTMLElement) string {
	str := e.ChildText(".Main3 > .m_T3 > h1")
	str = service.coverGBKToUTF8(str)
	str = new(DefaultCollectService).eliminateTrim(str, []string{"剧情介绍"})
	maps := strings.Split(str, "（")
	if strings.Index(str, "(") >= 0 {
		maps = strings.Split(str, "(")
	}
	return maps[0]
}

/** 提取短标题 **/
func (service *DefaultMovieService) collectOnResultMovieShort(e *colly.HTMLElement) string {
	var short string
	str := e.ChildText(".Main3 > .m_T3 > h1")
	str = service.coverGBKToUTF8(str)
	str = new(DefaultCollectService).eliminateTrim(str, []string{"剧情介绍"})
	maps := strings.Split(str, "（")
	if strings.Index(str, "(") >= 0 {
		maps = strings.Split(str, "(")
	}
	if len(maps) > 1 {
		short = strings.Trim(strings.Trim(maps[1], ")"), "）")
	}
	return short
}

/** 主演提取 **/
func (service *DefaultMovieService) collectOnResultMovieActor(e *colly.HTMLElement) string {
	str := e.ChildText(".g-box7 > .box > .txt > p:nth-child(1)")
	str = service.coverGBKToUTF8(str)
	str = new(DefaultCollectService).eliminateTrim(str, []string{"演：", "主", "/"})
	str = new(DefaultCollectService).deleteExtraSpace(str)
	var maps []string
	var result string
	if strings.Index(str, ",") >= 0 {
		maps = strings.Split(str, ",")
	} else {
		maps = strings.Split(str, " ")
	}
	if len(maps) > 0 {
		res, _ := json.Marshal(maps)
		result = string(res)
	}
	return result
}

/** 导演提取 **/
func (service *DefaultMovieService) collectOnResultMovieDirector(e *colly.HTMLElement) string {
	str := e.ChildText(".g-box7 > .box > .txt > p:nth-child(2) span:nth-child(2)")
	str = service.coverGBKToUTF8(str)
	str = new(DefaultCollectService).eliminateTrim(str, []string{"导", "演："})
	return new(DefaultCollectService).deleteExtraSpace(str)
}

/** 年份提取 **/
func (service *DefaultMovieService) collectOnResultMovieYear(e *colly.HTMLElement) string {
	str := e.ChildText(".g-box7 > .box > .txt > p:nth-child(4) span:nth-child(1)")
	str = service.coverGBKToUTF8(str)
	str = new(DefaultCollectService).eliminateTrim(str, []string{"份：", "年"})
	return new(DefaultCollectService).deleteExtraSpace(str)
}

/** 地区提取 **/
func (service *DefaultMovieService) collectOnResultMovieDistrict(e *colly.HTMLElement) string {
	str := e.ChildText(".g-box7 > .box > .txt > p:nth-child(4) span:nth-child(2)")
	str = service.coverGBKToUTF8(str)
	str = new(DefaultCollectService).eliminateTrim(str, []string{" 区：", "/", "地 "})
	maps := strings.Split(new(DefaultCollectService).deleteExtraSpace(str), " ")
	result, _ := json.Marshal(maps)
	return string(result)
}

/** 提取类型 **/
func (service *DefaultMovieService) collectOnResultMovieGenre(e *colly.HTMLElement) string {
	str := e.ChildText(".g-box7 > .box > .txt > p:nth-child(6)")
	str = service.coverGBKToUTF8(str)
	str = new(DefaultCollectService).eliminateTrim(str, []string{"型：", "类"})
	str = new(DefaultCollectService).deleteExtraSpace(str)
	var maps []string
	var result string
	if strings.Index(str, ",") >= 0 {
		maps = strings.Split(str, ",")
	} else {
		maps = strings.Split(str, " ")
	}
	if strings.Index(str, "/") >= 0 {
		maps = strings.Split(str, "/")
	}
	if len(maps) > 0 {
		res, _ := json.Marshal(maps)
		result = string(res)
	}
	return result
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
	str = new(DefaultCollectService).eliminateTrim(str, []string{"　", "\n"})
	return str
}

/** 演员详情 **/
func (service *DefaultMovieService) collectOnResultMovieActorCtx(e *colly.HTMLElement) string {
	var content string
	var urls string
	e.ForEach(".Main3 > .m_T3 > .Sub > ul > li", func(_ int, el *colly.HTMLElement) {
		if service.coverGBKToUTF8(el.ChildText("a[href]")) == "演员表" {
			urls = el.ChildAttr("a[href]", "href")
		}
	})
	if urls != "" {
		collector := new(DefaultCollectService).collectInstance(0, 1, "www.juqingba.cn", true)
		collector.OnHTML("html", func(e *colly.HTMLElement) {
			str := e.ChildText(".m_Left3 > .m_jq")
			str = service.coverGBKToUTF8(str)
			re, _ := regexp.Compile("\\<table[\\S\\s]+?\\</table\\>")
			str = re.ReplaceAllString(str, "")
			content = new(DefaultCollectService).eliminateTrim(str, []string{"　", "\n"})
		})
		_ = collector.Visit(urls)
		collector.Wait()
	}
	return content
}

/** 提取标题 **/
func (service *DefaultMovieService) collectOnResultMovieTitle(e *colly.HTMLElement) string {
	str := e.ChildText("head > title")
	str = service.coverGBKToUTF8(str)
	str = new(DefaultCollectService).eliminateTrim(str, []string{"_剧情吧"})
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

/** 入库一条数据 **/
func (service *DefaultMovieService) collectOnResultMovieInsert(res *movie.Movie, callback func(int64, error)) {
	if message := new(DefaultBaseVerify).Begin().Struct(res); message != nil {
		callback(0, errors.New(new(DefaultBaseVerify).Translate(message.(validator.ValidationErrors))))
	}
	if one := service.oneResultLink(res.Source); one.Id <= 0 {
		callback(orm.NewOrm().Insert(res))
	} else {
		res.Id = one.Id
		res.Status = 0
		_, err := orm.NewOrm().Update(res)
		callback(int64(one.Id), err)
	}
}

/************************* 剧集提取 *****************************/

/** 提取剧情链接加入ID提取 **/
func (service *DefaultMovieService) collectOnResultMovieExtractEpisode(e *colly.HTMLElement, id int64) {
	e.ForEach(".m_Left3 > .m_Box15 a[href]", func(_ int, el *colly.HTMLElement) {
		hrefUrls := el.Attr("href")
		if regexp.MustCompile(`https://www.juqingba.cn/zjuqing/(.*?).html$`).MatchString(hrefUrls) {
			hrefUrls = hrefUrls + "?sindhu=" + strconv.FormatInt(id, 10)
			_ = e.Request.Visit(hrefUrls)
		}
	})
}

/** 提取详情短标题 **/
func (service *DefaultMovieService) collectOnResultMovieDetailShort(e *colly.HTMLElement) string {
	urls := e.Request.URL.String()
	if strings.Index(urls, "_") >= 0 && len(strings.Split(urls, "_")) > 1 {
		str := strings.Split(urls, "_")[1]
		return strings.Split(str, ".")[0]
	} else {
		return "1"
	}
}

/** 提取详情内容 **/
func (service *DefaultMovieService) collectOnResultMovieContextShort(e *colly.HTMLElement) string {
	str, _ := e.DOM.Find(".m_jq > .ndfj").Html()
	str = service.coverGBKToUTF8(str)
	str = beego.HTML2str(str)
	str = new(DefaultCollectService).eliminateTrim(str, []string{"　", "\n"})
	return str
}

/** 提取一条剧情 **/
func (service *DefaultMovieService) collectOnResultMovieDetail(e *colly.HTMLElement, id int) *movie.EpisodeMovie {
	cate := service.One(id)
	return &movie.EpisodeMovie{
		Pid:         cate.Id,
		Name:        cate.Name,
		Short:       service.collectOnResultMovieDetailShort(e),
		Title:       service.collectOnResultMovieTitle(e),
		Keywords:    service.collectOnResultMovieKeywords(e),
		Description: service.collectOnResultMovieDescription(e),
		Context:     service.collectOnResultMovieContextShort(e),
		Source:      e.Request.URL.String(),
	}
}

/** 入库剧集详情 **/
func (service *DefaultMovieService) collectOnResultMovieDetailInsert(res *movie.EpisodeMovie, callback func(int64, error)) {
	if message := new(DefaultBaseVerify).Begin().Struct(res); message != nil {
		callback(0, errors.New(new(DefaultBaseVerify).Translate(message.(validator.ValidationErrors))))
	}
	if one := service.oneResultLinkDetail(res.Source); one.Id <= 0 {
		callback(orm.NewOrm().Insert(res))
	} else {
		res.Id = one.Id
		res.Status = 0
		_, err := orm.NewOrm().Update(res)
		callback(int64(one.Id), err)
	}
}

/******************************************************作品提取END ****************************************/
/******************************************************作品提取END ****************************************/

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
	collector.OnHTML("html", func(e *colly.HTMLElement) {
		if regexp.MustCompile("https://www.juqingba.cn/dianshiju/([0-9]+).html$").MatchString(e.Request.URL.String()) {
			result := service.collectOnResultMovie(e)
			go service.collectOnResultMovieInsert(result, func(id int64, err error) {
				status, message := service.analysisMessage(result.Title, err)
				new(DefaultCollectService).createLogsInform(status, uid, "电视剧情", message, e.Request.URL.String())
				if err == nil && id > 0 {
					service.collectOnResultMovieExtractEpisode(e, id)
				}
			})
		}
	})
	collector.OnHTML("html", func(e *colly.HTMLElement) { // 剧情剧集
		if regexp.MustCompile(`https://www.juqingba.cn/zjuqing/(.*?).html\?sindhu=([0-9]+)`).MatchString(e.Request.URL.String()) {
			cid := e.Request.URL.Query()["sindhu"]
			reagent, _ := strconv.Atoi(cid[0])
			result := service.collectOnResultMovieDetail(e, reagent)
			go service.collectOnResultMovieDetailInsert(result, func(id int64, err error) {
				status, message := service.analysisMessage(result.Title, err)
				new(DefaultCollectService).createLogsInform(status, uid, "剧集详情", message, e.Request.URL.String())
			})
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
		Text:      "爬虫配置",
		ClassName: "btn btn-sm btn-alt-warning mt-1 open_iframe",
		Attribute: map[string]string{
			"href":      beego.URLFor("Movie.Config", ":popup", 1),
			"data-area": "600px,300px",
		},
	})
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
	_, _ = qs.Limit(length, length*(page-1)).OrderBy("-id").ValuesList(&lists, "id", "name", "year", "class_type", "source", "status", "update_time")
	for _, v := range lists {
		v[1] = service.substr2HtmlHref(v[1].(string), v[4].(string), 0, 20)
		v[3] = service.substr2HtmlClass(v[3].(int64))
		v[4] = service.substr2HtmlHref(v[4].(string), v[4].(string), 0, 40)
		v[5] = service.Int2HtmlStatus(v[5], v[0], beego.URLFor("Movie.Push"))
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

/**  转为pop提示 **/
func (service *DefaultMovieService) substr2HtmlHref(s, urls string, start, end int) string {
	html := fmt.Sprintf(`<a href="%s" target="_blank" class="badge badge-primary js-tooltip" data-placement="top" data-toggle="tooltip" data-original-title="%s">%s</a>`, urls, s, beego.Substr(s, start, end))
	return html
}

func (service *DefaultMovieService) substr2HtmlClass(c int64) string {
	// 内容类型[0:电视 1:电影 2:综艺]
	maps := []string{"电视", "电影", "综艺"}
	return maps[c]
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultMovieService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string", "string", "date"},
		"fieldName": {"id", "name", "year", "class_type", "source", "status", "update_time"},
		"action":    {"", "", "", "", "", "", ""},
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
