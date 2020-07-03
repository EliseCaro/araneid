package service

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type DefaultAdapterService struct {
	filtration []string // 过滤规则
	extract    []string // 提取规则
	length     int      // 提取长度
	form       int      // 提取类型
}

/** 提取站长工具文件 **/
func (service *DefaultAdapterService) ZhanzhangExtract(id, length, form int, filtration, extract string) (result []*spider.Class, err error) {
	path := new(DefaultAdjunctService).FindId(id).Path
	service.createInitialized(length, form, filtration, extract)
	if f, message := excelize.OpenFile("." + path); message == nil {
		for _, sheet := range f.GetSheetList() {
			rows, _ := f.GetRows(sheet)
			for _, row := range rows {
				if len(row) == 9 && row[0] != "" && row[7] != "" && row[8] != "" {
					if item, message := service.extractUrlsContextRuleHandle(row[0], row[0], row[7], row[8]); message == nil {
						result = append(result, item)
					}
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

/** 发送处理通知 **/
func (service *DefaultAdapterService) createLogsInformStatus(name, status, urls, message string, receiver int) {
	var htmlInfo string
	nowTime := beego.Date(time.Now(), "m月d日 H:i")
	htmlInfo = fmt.Sprintf(`
		名为<a href="javascript:void(0);">%s</a>的栏目转换文件;在%s的时候%s；
		您可以<a href="%s">%s</a>;`,
		name, nowTime, status, urls, message,
	)
	new(DefaultInformService).SendSocketInform([]int{receiver}, 0, 0, 3, htmlInfo)
}

/**  携程处理方式 **/
func (service *DefaultAdapterService) SocketContextRuleHandle(id, length, form, uid int, filtration, extract string) {
	name := new(DefaultAdjunctService).FindId(id).Name
	if items, err := service.ZhanzhangExtract(id, length, form, filtration, extract); err == nil {
		if path, err := service.CreateXLSXFile(items); err == nil {
			service.createLogsInformStatus(name, "已经全部转换完成；请及时下载；系统将不会保留太长时间；", "/"+path, "点击下载", uid)
		} else {
			service.createLogsInformStatus(name, "创建文件发生错误；错误原因为："+err.Error(), beego.URLFor("Adapter.Index"), "检查系统目录权限,重新导入", uid)
		}
	} else {
		service.createLogsInformStatus(name, "初始化发生错误；错误原因为："+err.Error(), beego.URLFor("Adapter.Index"), "检查源文件,重新导入", uid)
	}
}

/** 匹配规则并远程提取 **/
func (service *DefaultAdapterService) extractUrlsContextRuleHandle(t, k, d, u string) (*spider.Class, error) {
	var message error
	var result *spider.Class
	if result, message = service.ruleHandle(t, k, d); message == nil {
		if service.form == 1 {
			service.ExtractUrlsContext(u, k, d, func(title, keyword, description string) {
				result, message = service.ruleHandle(t, service.emptyCheck(keyword, d), service.emptyCheck(description, d))
			})
		}
	}
	return result, message
}

/** 检测交换值 **/
func (service *DefaultAdapterService) emptyCheck(n, o string) string {
	if n == "" {
		return o
	} else {
		return n
	}
}

/** 创建一个文件写入栏目数据 **/
func (service *DefaultAdapterService) CreateXLSXFile(result []*spider.Class) (string, error) {
	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")
	for k, v := range result {
		indexes := k + 1
		_ = f.SetCellValue("Sheet1", "A"+strconv.Itoa(indexes), v.Title)
		_ = f.SetCellValue("Sheet1", "B"+strconv.Itoa(indexes), v.Keywords)
		_ = f.SetCellValue("Sheet1", "C"+strconv.Itoa(indexes), v.Description)
	}
	f.SetActiveSheet(index)
	path, _ := new(DefaultAdjunctService).DateFolder(beego.AppConfig.String("upload_folder") + "/sandbox/")
	name := fmt.Sprintf("%s/%d.xlsx", path, time.Now().Unix())
	if err := f.SaveAs(name); err != nil {
		return "", errors.New("创建导出文件失败；失败原因：" + err.Error())
	} else {
		return name, nil
	}
}

/** 初始化规则 **/
func (service *DefaultAdapterService) createInitialized(length, form int, filtration, extract string) {
	service.extract = strings.Split(extract, "|")
	service.filtration = strings.Split(filtration, "|")
	service.length = length
	service.form = form
}

/** 使用爬虫提取标题关键词跟描述 **/
func (service *DefaultAdapterService) ExtractUrlsContext(urls string, k string, d string, f func(t string, k string, d string)) {
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
