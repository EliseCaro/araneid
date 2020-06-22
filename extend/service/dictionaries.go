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
	"github.com/go-playground/validator"
	"github.com/gocolly/colly"
	"regexp"
	"strconv"
	"time"
)

type DefaultDictionariesService struct{}

/** 根据name获取一条数据 **/
func (service *DefaultDictionariesService) ConfigByName(name string) dictionaries.DictConfig {
	var item dictionaries.DictConfig
	_ = orm.NewOrm().QueryTable(new(dictionaries.DictConfig)).Filter("name", name).One(&item)
	return item
}

/** 根据name获取一条数据 **/
func (service *DefaultDictionariesService) OneCate(id int) dictionaries.DictCate {
	var item dictionaries.DictCate
	_ = orm.NewOrm().QueryTable(new(dictionaries.DictCate)).Filter("id", id).One(&item)
	return item
}

/** 获取爬虫配置解析成map **/
func (service *DefaultDictionariesService) ConfigMaps() map[string]interface{} {
	var item []*dictionaries.DictConfig
	maps := make(map[string]interface{})
	_, _ = orm.NewOrm().QueryTable(new(dictionaries.DictConfig)).All(&item)
	for _, v := range item {
		maps[v.Name] = v.Value
	}
	return maps
}

/**************以下为分类爬虫 *****************/

/** 爬虫操纵命令 **/
func (service *DefaultDictionariesService) CateInstanceBegin(instruct map[string]interface{}) {
	if instruct["status"].(int) == 1 {
		service.StartCate(instruct["uid"].(int))
	} else {
		service.StopCate(instruct["uid"].(int), "停止爬取查字词典分类")
	}
}

/** 启动爬虫 **/
func (service *DefaultDictionariesService) StartCate(uid int) {
	object := DefaultCollectService{}
	config := service.ConfigMaps()
	interval, _ := strconv.Atoi(config["interval"].(string))
	collector := object.collectInstance(interval, 2, "www.chazidian.com")
	collector.OnResponse(func(r *colly.Response) {
		if status, _ := strconv.Atoi(service.ConfigMaps()["status"].(string)); status == 0 {
			r.Request.Abort()
		}
	})
	collector.OnHTML(".pinyin_sou", func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
			_ = e.Request.Visit(el.Attr("href"))
		})
	})
	collector.OnHTML("html", func(e *colly.HTMLElement) {
		if regexp.MustCompile("https://www.chazidian.com/ci_([a-z]+)_([0-9]+)/").MatchString(e.Request.URL.String()) {
			keywords, _ := e.DOM.Find("head > meta[name=keywords]").Attr("content")
			description, _ := e.DOM.Find("head > meta[name=description]").Attr("content")
			message := service.createOneResultCate(&dictionaries.DictCate{
				Title:       object.eliminateTrim(e.DOM.Find("head > title").Text(), []string{" - 查字典"}),
				Source:      e.Request.URL.String(),
				Name:        e.ChildText(".main_left > div:last-child > .box_head > .box_title > h2 > b"),
				Initial:     e.ChildText(".main_left > div:nth-child(3) > .box_head > .box_title > h2 > b"),
				Status:      0,
				Keywords:    keywords,
				Description: description,
			})
			if message != nil {
				service.createLogsInformStatusCate("查字词典分类采集到一条非法数据；<a href='"+e.Request.URL.String()+"'>查看原文</a>;非法原因："+error.Error(message)+";", uid)
			}
		}
	})
	collector.OnError(func(r *colly.Response, err error) {
		service.createLogsInformStatusCate("查字词典分类在采集过程中出现一次失败的采集；已经跳过错误继续采集；", uid)
	})
	if _, err := orm.NewOrm().Update(&dictionaries.DictConfig{Id: 5, Value: "1"}, "Value"); err != nil {
		logs.Warn("启动查字词典采集器失败！失败原因:%s", error.Error(err))
	} else {
		_ = collector.Visit("https://www.chazidian.com/ci_pinyin/")
		service.createLogsInformStatusCate("开始爬取查字词典分类", uid)
		collector.Wait()
	}
	defer func() { service.StopCate(uid, "查字词典分类爬取任务完成") }()
}

/** 停止爬虫 **/
func (service *DefaultDictionariesService) StopCate(uid int, message string) {
	if _, err := orm.NewOrm().Update(&dictionaries.DictConfig{Id: 5, Value: "0"}, "Value"); err != nil {
		logs.Warn("停止查字词典采集器失败！失败原因:%s", error.Error(err))
	} else {
		service.createLogsInformStatusCate(message, uid)
	}
}

/** 批量删除分类采集结果 **/
func (service *DefaultDictionariesService) DeleteArrayCate(array []int) (message error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, message = orm.NewOrm().Delete(&dictionaries.DictCate{Id: v}); message != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if message == nil {
		_ = orm.NewOrm().Commit()
	}
	return message
}

/** 查看分类结果 **/
func (service *DefaultDictionariesService) DetailCateOne(id int) map[string]string {
	result := make(map[string]string)
	detail := service.OneCate(id)
	result["logs"] = detail.Logs
	result["title"] = detail.Title
	result["status"] = strconv.Itoa(int(detail.Status))
	result["update_time"] = beego.Date(detail.UpdateTime, "Y年m月d H:i:s")
	maps := map[string]string{
		"name": detail.Name, "title": detail.Title,
		"keywords": detail.Keywords, "description": detail.Description,
		"source": detail.Source, "initial": detail.Initial,
	}
	stringJson, _ := json.Marshal(maps)
	result["result"] = string(stringJson)
	return result
}

/*****************以下为私有方法 *************************/
/** 发送采集爬虫状态通知 **/
func (service *DefaultDictionariesService) createLogsInformStatusCate(status string, receiver int) {
	var htmlInfo string
	nowTime := beego.Date(time.Now(), "m月d日 H:i")
	htmlInfo = fmt.Sprintf(`
		名为<a href="javascript:void(0);">查字词典</a>的采集器;在%s的时候%s；
		您可以<a target="_blank" href="%s">查看爬取结果</a>;`,
		nowTime, status, beego.URLFor("Dictionaries.Cate"),
	)
	inform := DefaultInformService{}
	inform.SendSocketInform([]int{receiver}, 0, 0, 2, htmlInfo)
}

/** 创建一条分类结果数据 */
func (service *DefaultDictionariesService) createOneResultCate(item *dictionaries.DictCate) (err error) {
	verify := DefaultBaseVerify{}
	if message := verify.Begin().Struct(item); message != nil {
		return errors.New(verify.Translate(message.(validator.ValidationErrors)))
	}
	if one := service.oneResultLinkCate(item.Source); one.Id > 0 {
		item.Id = one.Id
		_, err = orm.NewOrm().Update(item)
	} else {
		_, err = orm.NewOrm().Insert(item)
	}
	return err
}

/** 根据链接获取一条结果数据 **/
func (service *DefaultDictionariesService) oneResultLinkCate(url string) dictionaries.DictCate {
	var item dictionaries.DictCate
	_ = orm.NewOrm().QueryTable(new(dictionaries.DictCate)).Filter("source", url).One(&item)
	return item
}

/** 将状态码转为html **/
func (service *DefaultDictionariesService) int2HtmlStatus(val interface{}, id interface{}, url string) string {
	t := table.BuilderTable{}
	return t.AnalysisTypeSwitch(map[string]interface{}{"status": val, "id": id}, "status", url, map[int64]string{-1: "已失败", 0: "待发布", 1: "已发布"})
}

/************************************************分类表格渲染机制 ************************************************************/

/** 获取需要渲染的Column **/
func (service *DefaultDictionariesService) DataTableCateColumns() []map[string]interface{} {
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

/** 返回表单结构字段如何解析 **/
func (service *DefaultDictionariesService) TableColumnsTypeCate() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string", "date"},
		"fieldName": {"id", "name", "initial", "source", "status", "update_time"},
		"action":    {"", "", "", "", "", ""},
	}
	return result
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultDictionariesService) DataTableCateButtons() []*table.TableButtons {
	var array []*table.TableButtons
	config := service.ConfigMaps()
	array = append(array, &table.TableButtons{
		Text:      "爬虫配置",
		ClassName: "btn btn-sm btn-alt-warning mt-1 open_iframe",
		Attribute: map[string]string{
			"href":      beego.URLFor("Dictionaries.Config", ":popup", 1),
			"data-area": "600px,400px",
		},
	})
	if config["status"].(string) == "1" {
		array = append(array, &table.TableButtons{
			Text:      "关闭爬虫",
			ClassName: "btn btn-sm btn-alt-danger mt-1 handle_collect",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Dictionaries.StatusCate"),
				"data-status": "0",
			},
		})
	} else {
		array = append(array, &table.TableButtons{
			Text:      "启动爬虫",
			ClassName: "btn btn-sm btn-alt-success mt-1 handle_collect",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Dictionaries.StatusCate"),
				"data-status": "1",
			},
		})
	}
	array = append(array, &table.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Dictionaries.DeleteCate")},
	})
	array = append(array, &table.TableButtons{
		Text:      "发布选中",
		ClassName: "btn btn-sm btn-alt-primary mt-1 ids_enables",
		Attribute: map[string]string{
			"data-action": beego.URLFor("Dictionaries.Status"),
			"data-field":  "status",
		},
	})
	return array
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultDictionariesService) TableButtonsTypeCate() []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "查看结果",
			ClassName: "btn btn-sm btn-alt-primary open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Dictionaries.DetailCate", ":id", "__ID__", ":popup", 1),
				"data-area": "620px,400px",
			},
		},
		{
			Text:      "删除结果",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Dictionaries.DeleteCate"),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}

/** 处理分页 **/
func (service *DefaultDictionariesService) PageListItemsCate(length, draw, page int, search string) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(dictionaries.DictCate))
	if search != "" {
		qs = qs.Filter("initial__icontains", search)
	}
	recordsTotal, _ := qs.Count()
	_, _ = qs.Limit(length, length*(page-1)).ValuesList(&lists, "id", "name", "initial", "source", "status", "update_time")
	for _, v := range lists {
		v[4] = service.int2HtmlStatus(v[4], v[0], beego.URLFor("Dictionaries.PushCate"))
	}
	data := map[string]interface{}{
		"draw":            draw,         // 请求次数
		"recordsFiltered": recordsTotal, // 从多少条里面筛选
		"recordsTotal":    recordsTotal, // 总条数
		"data":            lists,        // 筛选结果
	}
	return data
}
