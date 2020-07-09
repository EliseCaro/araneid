package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/automatic"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

/** 自动推送url链接层 **/
type DefaultAutomaticService struct{}

/** 根据文档ID推送一条链接 **/
func (service *DefaultAutomaticService) AutomaticDocument(object int) {
	detail := new(DefaultArticleService).One(object)
	var domain spider.Domain
	if detail.Class > 0 {
		if cate := new(DefaultCategoryService).AcquireCategoryOldWhere(detail.Class); cate.Domain > 0 {
			domain, _ = new(DefaultDomainService).Find(cate.Domain)
		}
	}
	if domain.Domain != "" && domain.Submit != "" {
		urls := fmt.Sprintf(`http://%s/index/detail-%d.html`, domain.Domain, detail.Id)
		service.baiduAutomatic(domain.Domain, domain.Submit, urls)
	}
}

/** 统计日总蜘蛛数量 **/
func (service *DefaultAutomaticService) SumAll() int64 {
	sum, _ := orm.NewOrm().QueryTable(new(automatic.Automatic)).Filter("status", 1).Count()
	return sum
}

/** 统计日总蜘蛛数量 **/
func (service *DefaultAutomaticService) SumDay() int64 {
	date := time.Now().Format("2006-01-02")
	qs := orm.NewOrm().QueryTable(new(automatic.Automatic)).Filter("status", 1)
	qs = qs.Filter("update_time__gte", fmt.Sprintf(`%s 00:00:00`, date))
	sum, _ := qs.Filter("update_time__lte", fmt.Sprintf(`%s 23:59:59`, date)).Count()
	return sum
}

/** 解析一个星期的推送记录 **/
func (service *DefaultAutomaticService) AnalysisWeek() string {
	var result []int64
	var labels []string
	var current = time.Now().Unix()
	for i := 0; i <= 6; i++ {
		date := time.Unix(current-(int64(i)*86400), 0).Format("2006-01-02")
		qs := orm.NewOrm().QueryTable(new(automatic.Automatic)).Filter("status", 1)
		qs = qs.Filter("update_time__gte", fmt.Sprintf(`%s 00:00:00`, date))
		num, _ := qs.Filter("update_time__lte", fmt.Sprintf(`%s 23:59:59`, date)).Count()
		result = append(result, num)
		labels = append(labels, time.Unix(current-(int64(i)*86400), 0).Format("01-02"))
	}
	ret := map[string]interface{}{
		"date":  labels,
		"items": []map[string]interface{}{{"label": "自动推送量", "data": result, "fill": true, "backgroundColor": "rgba(75, 106, 199,0.8)"}},
	}
	byteStr, _ := json.Marshal(ret)
	return string(byteStr)
}

/** 分类自动提交 **/
func (service *DefaultAutomaticService) AutomaticClass(domain int) {
	if detail, _ := new(DefaultDomainService).Find(domain); detail.Domain != "" && detail.Submit != "" {
		var cate []*spider.Class
		_ = json.Unmarshal([]byte(detail.Cate), &cate)
		for _, v := range cate {
			urls := fmt.Sprintf(`http://%s/index/column-%d-0.html`, detail.Domain, v.Id)
			service.baiduAutomatic(detail.Domain, detail.Submit, urls)
		}
	}
}

/** 百度提交 **/
func (service *DefaultAutomaticService) baiduAutomatic(domain, token, urls string) {
	action := fmt.Sprintf(beego.AppConfig.String("baidu_automatic_urls"), domain, token)
	result := make(map[string]interface{})
	client := &http.Client{}
	r, _ := http.NewRequest("POST", action, bytes.NewBuffer([]byte(urls)))
	resp, err := client.Do(r)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err == nil {
		if body, message := ioutil.ReadAll(resp.Body); message == nil {
			message = json.Unmarshal(body, &result)
			if result["success"] != nil && result["success"].(int64) == 1 {
				service.CreateLogs(domain, token, urls, strconv.Itoa(result["remain"].(int)), 1)
			} else {
				service.CreateLogs(domain, token, urls, result["message"].(string), 0)
			}
		} else {
			service.CreateLogs(domain, token, urls, message.Error(), 0)
		}
		defer resp.Body.Close()
	} else {
		service.CreateLogs(domain, token, urls, err.Error(), 0)
	}
}

/** 根据老id查询新ID数据 **/
func (service *DefaultAutomaticService) AcquireUrlsOne(urls string) automatic.Automatic {
	var maps automatic.Automatic
	_ = orm.NewOrm().QueryTable(new(automatic.Automatic)).Filter("urls", urls).One(&maps)
	return maps
}

/** 创建推送记录 **/
func (service *DefaultAutomaticService) CreateLogs(domain, token, urls, remark string, status int8) {
	item := automatic.Automatic{Domain: domain, Token: token, Urls: urls, Remark: remark, Status: status}
	if index := service.AcquireUrlsOne(urls); index.Id > 0 {
		item.Id = index.Id
		_, _ = orm.NewOrm().Update(&item)
	} else {
		_, _ = orm.NewOrm().Insert(&item)
	}
}

/** 获取一条数 **/
func (service *DefaultAutomaticService) One(id int) automatic.Automatic {
	var maps automatic.Automatic
	_ = orm.NewOrm().QueryTable(new(automatic.Automatic)).Filter("id", id).One(&maps)
	return maps
}

/** 重新推送 **/
func (service *DefaultAutomaticService) StatusArray(array []int) {
	for _, v := range array {
		one := service.One(v)
		service.baiduAutomatic(one.Domain, one.Token, one.Urls)
	}
}

/** 批量删除 **/
func (service *DefaultAutomaticService) DeleteArray(array []int) (message error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, message = orm.NewOrm().Delete(&automatic.Automatic{Id: v}); message != nil {
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
func (service *DefaultAutomaticService) EmptyDelete() {
	var item []*automatic.Automatic
	_, _ = orm.NewOrm().QueryTable(new(automatic.Automatic)).All(&item)
	for _, v := range item {
		_, _ = orm.NewOrm().Delete(&automatic.Automatic{Id: v.Id})
	}
}

/****************** 以下为表格渲染  ***********************/

/** 获取需要渲染的Column **/
func (service *DefaultAutomaticService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "推送标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "推送URL", "name": "urls", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "推送域名", "name": "domain", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "推送日志", "name": "remark", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "推送状态", "name": "status", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "更新时间", "name": "update_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "数据操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultAutomaticService) DataTableButtons() []*_func.TableButtons {
	var array []*_func.TableButtons
	array = append(array, &_func.TableButtons{
		Text:      "重新推送",
		ClassName: "btn btn-sm btn-alt-primary mt-1 ids_enables",
		Attribute: map[string]string{
			"data-action": beego.URLFor("Automatic.Status"),
			"data-field":  "status",
		},
	})
	array = append(array, &_func.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-warning mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Automatic.Delete")},
	})
	array = append(array, &_func.TableButtons{
		Text:      "清空记录",
		ClassName: "btn btn-sm btn-alt-danger mt-1 js-tooltip ids_deletes",
		Attribute: map[string]string{
			"data-action":         beego.URLFor("Automatic.Empty"),
			"data-toggle":         "tooltip",
			"data-original-title": "清空后将无法恢复，请谨慎操作！",
		},
	})
	return array
}

/** 处理分页 **/
func (service *DefaultAutomaticService) PageListItems(length, draw, page int, search string) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(automatic.Automatic))
	if search != "" {
		qs = qs.Filter("domain", search)
	}
	recordsTotal, _ := qs.Count()
	_, _ = qs.Limit(length, length*(page-1)).OrderBy("-id").ValuesList(&lists, "id", "urls", "domain", "remark", "status", "update_time")
	for _, v := range lists {
		v[1] = service.substr2HtmlHref(v[1].(string), v[1].(string), 0, 25)
		v[4] = service.Int2HtmlStatus(v[4], v[0], beego.URLFor("Automatic.Status"))
	}
	data := map[string]interface{}{
		"draw":            draw,         // 请求次数
		"recordsFiltered": recordsTotal, // 从多少条里面筛选
		"recordsTotal":    recordsTotal, // 总条数
		"data":            lists,        // 筛选结果
	}
	return data
}

/**  转为pop提示 **/
func (service *DefaultAutomaticService) substr2HtmlHref(u, s string, start, end int) string {
	html := fmt.Sprintf(`<a href="%s" target="_blank" class="badge badge-primary js-tooltip" data-placement="top" data-toggle="tooltip" data-original-title="%s">%s...</a>`, u, s, beego.Substr(s, start, end))
	return html
}

/** 将状态码转为html **/
func (service *DefaultAutomaticService) Int2HtmlStatus(val interface{}, id interface{}, url string) string {
	t := _func.BuilderTable{}
	return t.AnalysisTypeSwitch(map[string]interface{}{"status": val, "id": id}, "status", url, map[int64]string{0: "推送失败", 1: "推送成功"})
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultAutomaticService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string", "date"},
		"fieldName": {"id", "urls", "domain", "remark", "status", "update_time"},
		"action":    {"", "", "", "", "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultAutomaticService) TableButtonsType() []*_func.TableButtons {
	buttons := []*_func.TableButtons{
		{
			Text:      "删除推送",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Automatic.Delete"),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}
