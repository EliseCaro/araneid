package service

import (
	"errors"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego/logs"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/gocolly/colly"
	"strings"
	"unicode"
)

type DefaultAdapterService struct {
	filtration []string // 过滤规则
	extract    []string // 提取规则
	length     int      // 提取长度
}

/** 提取站长工具文件 **/
func (service *DefaultAdapterService) ZhanzhangExtract(id, length int, filtration, extract string) (result []*spider.Class, err error) {
	path := new(DefaultAdjunctService).FindId(id).Path
	service.createInitialized(length, filtration, extract)
	if f, message := excelize.OpenFile("." + path); message == nil {
		for _, sheet := range f.GetSheetList() {
			rows, _ := f.GetRows(sheet)
			for _, row := range rows {
				if len(row) == 9 && row[0] != "" && row[7] != "" && row[8] != "" {
					service.ExtractUrlsContext(row[8], func(t, k, d string) {
						if k == "" {
							k = row[0]
						}
						if d == "" {
							d = row[7]
						}
						if item, message := service.ruleHandle(row[0], k, d); message == nil {
							result = append(result, item)
						}
					})
				} else {
					err = errors.New("解析文件格式失败！该文件不是站长工具的格式～")
				}
			}
		}
	} else {
		err = errors.New("解析文件失败！未监测到合法的导入文件～")
	}
	return result, err
}

/** 初始化规则 **/
func (service *DefaultAdapterService) createInitialized(length int, filtration, extract string) {
	service.extract = strings.Split(extract, "|")
	service.filtration = strings.Split(filtration, "|")
	service.length = length
}

/** 使用爬虫提取标题关键词跟描述 **/
func (service *DefaultAdapterService) ExtractUrlsContext(urls string, f func(t, k, d string)) {
	var (
		title    string
		keyword  string
		describe string
	)
	domain := new(DefaultCollectService).queueUrlDomain(urls)
	collector := new(DefaultCollectService).collectInstance(5, 1, domain, true)
	collector.OnHTML("html", func(e *colly.HTMLElement) {
		title = e.DOM.Find("head > title").Text()
		keyword, _ = e.DOM.Find("head > meta[name=keywords]").Attr("content")
		describe, _ = e.DOM.Find("head > meta[name=description]").Attr("content")
		f(title, keyword, describe)
	})
	collector.OnError(func(r *colly.Response, err error) {
		logs.Error("采集栏目数据失败；失败原因：" + err.Error())
		f("", "", "")
	})
	_ = collector.Visit(urls)
	collector.Wait()
}

/** 统计字数 **/
func (service *DefaultAdapterService) chineseCount(str1 string) (count int) {
	for _, char := range str1 {
		if unicode.Is(unicode.Han, char) {
			count++
		}
	}
	return
}

/** 规则处理 **/
func (service *DefaultAdapterService) ruleHandle(title, keyword, description string) (*spider.Class, error) {
	if service.length > 0 {
		title = service.ruleHandleTitle(title)
	}
	var message error
	var result spider.Class
	maps := map[string]string{"title": title, "keyword": keyword, "description": description}
	for k, v := range maps {
		v = service.ruleHandleFiltration(v) // 过滤
		if service.ruleHandleExtract(v) {
			maps[k] = v
		} else {
			message = errors.New("不匹配规则～")
		}
	}
	if message == nil {
		result.Keywords = maps["keyword"]
		result.Description = maps["description"]
		result.Title = maps["title"]
	}
	return &result, message
}

/** 标题检测**/
func (service *DefaultAdapterService) ruleHandleTitle(title string) string {
	var str string
	if service.length == service.chineseCount(title) {
		str = title
	}
	return str
}

/** 过滤字符 **/
func (service *DefaultAdapterService) ruleHandleFiltration(str string) string {
	for _, v := range service.filtration {
		str = strings.Replace(str, v, "", 1)
	}
	return strings.Replace(str, " ", "", -1)
}

/** 是否包含 **/
func (service *DefaultAdapterService) ruleHandleExtract(str string) bool {
	var status bool
	for _, v := range service.extract {
		if strings.Index(str, v) > 0 {
			status = true
			break
		}
	}
	return status
}
