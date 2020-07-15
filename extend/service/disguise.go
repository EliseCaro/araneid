package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/go-playground/validator"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tencent "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	nlp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/nlp/v20190408"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type DefaultDisguiseService struct{}

/** 获取一条数据 **/
func (service *DefaultDisguiseService) Find(id int) (spider.Disguise, error) {
	item := spider.Disguise{
		Id: id,
	}
	return item, orm.NewOrm().Read(&item)
}

/** 获取输出模型 **/
func (service *DefaultDisguiseService) KeyOne(k, s string) spider.Disguise {
	var items spider.Disguise
	_ = orm.NewOrm().QueryTable(new(spider.Disguise)).Filter("api_key", k).Filter("api_secret", s).One(&items)
	return items
}

/** 获取所有分组 **/
func (service *DefaultDisguiseService) Groups() []*spider.Disguise {
	var items []*spider.Disguise
	_, _ = orm.NewOrm().QueryTable(new(spider.Disguise)).All(&items)
	return items
}

/** 批量删除结果 **/
func (service *DefaultDisguiseService) DeleteArray(array []int) (message error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, message = orm.NewOrm().Delete(&spider.Disguise{Id: v}); message != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if message == nil {
		_ = orm.NewOrm().Commit()
	}
	return message
}

/** 更新状态 **/
func (service *DefaultDisguiseService) UpdateStatus(id int, status int8, field string) (e error) {
	if field == "keyword" {
		_, e = orm.NewOrm().Update(&spider.Disguise{Id: id, Keyword: status}, "Keyword")
	}
	if field == "description" {
		_, e = orm.NewOrm().Update(&spider.Disguise{Id: id, Description: status}, "Description")
	}
	if field == "context" {
		_, e = orm.NewOrm().Update(&spider.Disguise{Id: id, Context: status}, "Context")
	}
	return e
}

/** 计数器++ **/
func (service *DefaultDisguiseService) Inc(id, old int) (errorMsg error) {
	if old != 0 {
		errorMsg = service.Dec(old)
	}
	_, errorMsg = orm.NewOrm().QueryTable(new(spider.Disguise)).Filter("id", id).Update(orm.Params{
		"usage": orm.ColValue(orm.ColAdd, 1),
	})
	return errorMsg
}

/** 计数器-- **/
func (service *DefaultDisguiseService) Dec(id int) error {
	_, errorMessage := orm.NewOrm().QueryTable(new(spider.Disguise)).Filter("id", id).Update(orm.Params{
		"usage": orm.ColValue(orm.ColMinus, 1),
	})
	return errorMessage
}

/***************************自然语言算法处理中心 *****************************************8/
/** 获取百度翻译结果 **/
func (service *DefaultDisguiseService) baiduTranslation(s string) ([]interface{}, string, string) {
	salt := strconv.FormatInt(time.Now().Unix(), 10)
	appId := beego.AppConfig.String("baidu_translation_id")
	secret := beego.AppConfig.String("baidu_translation_sign")
	domain := beego.AppConfig.String("baidu_translation_url")
	result, _ := service.requestPostForm(domain, url.Values{
		"q": {s}, "from": {"auto"}, "to": {"auto"}, "appid": {appId},
		"sign": {service.md5Value(appId + s + salt + secret)}, "salt": {salt},
	})
	beego.Error(result)
	if result["error_code"] != nil && result["error_code"].(string) == "54003" || result["error_code"] != nil && result["error_code"].(string) == "54005" {
		return service.baiduTranslation(s)
	}
	if result["error_code"] != nil {
		return []interface{}{}, "", ""
	} else {
		return result["trans_result"].([]interface{}), result["from"].(string), result["to"].(string)
	}
}

/** 作为等价交换手段 **/
func (service *DefaultDisguiseService) modifierContext(s string) string {
	translation, _, to := service.baiduTranslation(s)
	var result string
	for _, v := range translation {
		result = v.(map[string]interface{})["dst"].(string)
	}
	if result != "" && to != "zh" {
		return service.modifierContext(result)
	}
	return result
}

/** 公用请求函数 **/
func (service *DefaultDisguiseService) requestPostForm(domain string, v url.Values) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	var message error
	if resp, message := http.PostForm(domain, v); message == nil {
		if body, message := ioutil.ReadAll(resp.Body); message == nil {
			message = json.Unmarshal(body, &result)
		}
		defer resp.Body.Close()
	}
	return result, message
}

/** 字符串生成md5**/
func (service *DefaultDisguiseService) md5Value(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

/** 自然语言处理入口 **/
func (service *DefaultDisguiseService) DisguiseHandleManage(disguise int, module *spider.HandleModule) (*spider.HandleModule, error) {
	var verify = DefaultBaseVerify{}
	if message := verify.Begin().Struct(module); message == nil {
		return service.handleManageBegin(disguise, module)
	} else {
		return module, errors.New(verify.Translate(message.(validator.ValidationErrors)))
	}
}

/** 将map转为json字符串 **/
func (service *DefaultDisguiseService) jsonFormString(maps map[string]string) string {
	j, _ := json.Marshal(maps)
	return string(j)
}

/** 将文本过滤成纯文本 **/
func (service *DefaultDisguiseService) contextFiltration(s string) string {
	return strings.Replace(beego.HTML2str(s), "\n", "", -1)
}

/** 根据机器ID获取一个腾讯实例 **/
func (service *DefaultDisguiseService) nlpInstance(disguise int) (*nlp.Client, error) {
	var (
		client  *nlp.Client
		message error
		object  spider.Disguise
	)
	if object, message = service.Find(disguise); message == nil {
		credential := common.NewCredential(object.ApiKey, object.ApiSecret)
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "nlp.tencentcloudapi.com"
		client, _ = nlp.NewClient(credential, "ap-guangzhou", cpf)
	}
	return client, message
}

/** 初始化关键字环境 **/
func (service *DefaultDisguiseService) handleManageBeginKeyword(disguise int, module *spider.HandleModule) ([]*nlp.Keyword, error) {
	if client, message := service.nlpInstance(disguise); message == nil {
		var response *nlp.KeywordsExtractionResponse
		request := nlp.NewKeywordsExtractionRequest()
		_ = request.FromJsonString(service.jsonFormString(map[string]string{
			"Text": service.contextFiltration(module.Context),
		}))
		response, message = client.KeywordsExtraction(request)
		if _, ok := message.(*tencent.TencentCloudSDKError); ok == false && message == nil {
			return response.Response.Keywords, message
		} else {
			return nil, errors.New(fmt.Sprintf("An API error has returned: %s", message))
		}
	} else {
		return nil, message
	}
}

/** 获取真是关键字**/
func (service *DefaultDisguiseService) robotKeywordManage(keyword []*nlp.Keyword, disguise int) string {
	var result string
	var object, _ = service.Find(disguise)
	for i, v := range keyword {
		if object.Length >= int8(i) {
			result += *v.Word + ","
		}
	}
	return strings.Trim(result, ",")
}

/*** 从训练模型提取不到关键词；需要写入; **/
func (service *DefaultDisguiseService) setRobotKeywordManage(keyword *string, disguise int) string {
	maps, message := service.resemblanceTags(keyword, disguise)
	if message == nil || len(maps) > 0 {
		maps = new(DefaultRobotService).Insert(*keyword, maps)
	}
	var result []string
	for _, value := range maps {
		result = append(result, *value)
	}
	return strings.Join(result, ",")
}

/*** 获得同义词; **/
func (service *DefaultDisguiseService) resemblanceTags(s *string, disguise int) ([]*string, error) {
	if client, message := service.nlpInstance(disguise); message == nil {
		var response *nlp.SimilarWordsResponse
		request := nlp.NewSimilarWordsRequest()
		_ = request.FromJsonString(service.jsonFormString(map[string]string{"Text": *s}))
		response, message = client.SimilarWords(request)
		if _, ok := message.(*tencent.TencentCloudSDKError); ok == false && message == nil {
			return response.Response.SimilarWords, message
		} else {
			return nil, errors.New(fmt.Sprintf("An API error has returned: %s", message))
		}
	} else {
		return nil, message
	}
}

/** 开始处理 **/
func (service *DefaultDisguiseService) handleManageBegin(disguise int, module *spider.HandleModule) (*spider.HandleModule, error) {
	keyword, message := service.handleManageBeginKeyword(disguise, module)
	config, _ := service.Find(disguise)
	module.Title = service.robotTitleManage(module.Title, disguise)
	if config.Keyword == 1 {
		module.Keywords = service.robotKeywordManage(keyword, disguise)
	}
	if config.Description == 1 {
		module.Description = service.robotDescriptionManage(module, disguise)
	}
	if config.Context == 1 {
		module.Context = service.robotContextManage(module.Context, disguise)
	}
	return module, message
}

/** 提取文档描述 todo 同义词替换未做 **/
func (service *DefaultDisguiseService) robotDescriptionManage(module *spider.HandleModule, disguise int) string {
	if client, message := service.nlpInstance(disguise); message == nil {
		var response *nlp.AutoSummarizationResponse
		request := nlp.NewAutoSummarizationRequest()
		text := fmt.Sprintf(`%s。%s`,
			service.contextFiltration(module.Description),
			service.contextFiltration(module.Context))
		_ = request.FromJsonString(service.jsonFormString(map[string]string{"Text": beego.Substr(text, 0, 2000)}))
		response, message = client.AutoSummarization(request)
		if _, ok := message.(*tencent.TencentCloudSDKError); ok == false && message == nil {
			module.Description = *response.Response.Summary
		}
	}
	object, _ := service.Find(disguise)
	if object.Modifier == 1 { //  等价交换手段
		module.Description = service.modifierContext(module.Description)
	}
	return module.Description
}

/** 提取内容；重建结构 todo 同义词替换未做 ***/
func (service *DefaultDisguiseService) robotContextManage(s string, d int) string {
	object, _ := service.Find(d)
	s = service.contextFiltration(s)
	if object.Modifier == 1 {
		s = service.modifierContext(s)
	}
	justify := strings.Split(s, "。")
	news := "<p class='justify'>"
	for k, value := range justify {
		if k%2 == 0 {
			news += value + "。</p>"
		} else {
			news += "<p>" + value
		}
	}
	return news
}

/** 文档标题处理 todo 同义词替换未做 ***/
func (service *DefaultDisguiseService) robotTitleManage(s string, d int) string {
	object, _ := service.Find(d)
	s = service.contextFiltration(s)
	if object.Modifier == 1 {
		s = service.modifierContext(s)
	}
	return s
}

/************************************************表格渲染机制 ************************************************************/

/** 获取需要渲染的Column **/
func (service *DefaultDisguiseService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "机器名称", "name": "name", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "挂载次数", "name": "usage", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "自然单词", "name": "keyword", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "自然描述", "name": "description", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "自然内容", "name": "context", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "接口KEY", "name": "api_key", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "接口SECRET", "name": "api_secret", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultDisguiseService) DataTableButtons() []*table.TableButtons {
	var array []*table.TableButtons
	array = append(array, &table.TableButtons{
		Text:      "创建机器",
		ClassName: "btn btn-sm btn-alt-success mt-1 open_iframe",
		Attribute: map[string]string{
			"href":      beego.URLFor("Disguise.Create", ":popup", 1),
			"data-area": "600px,400px",
		},
	})
	array = append(array, &table.TableButtons{
		Text:      "删除机器",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Disguise.Delete")},
	})
	return array
}

/** 处理分页 **/
func (service *DefaultDisguiseService) PageListItems(length, draw, page int, search string) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(spider.Disguise))
	recordsTotal, _ := qs.Count()
	if search != "" {
		qs = qs.Filter("name__icontains", search)
	}
	_, _ = qs.Limit(length, length*(page-1)).OrderBy("-id").ValuesList(&lists, "id", "name", "usage", "keyword", "description", "context", "api_key", "api_secret")
	for _, v := range lists {
		v[6] = service.Substr2HtmlSpan(v[6].(string), 0, 30)
		v[7] = service.Substr2HtmlSpan(v[7].(string), 0, 30)
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
func (service *DefaultDisguiseService) Substr2HtmlSpan(s string, start, end int) string {
	html := fmt.Sprintf(`<span class="badge badge-primary js-tooltip" data-placement="top" data-toggle="tooltip" data-original-title="%s">%s</span>`, s, beego.Substr(s, start, end))
	return html
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultDisguiseService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "switch", "switch", "switch", "string", "string"},
		"fieldName": {"id", "name", "usage", "keyword", "description", "context", "api_key", "api_secret"},
		"action":    {"", "", "", beego.URLFor("Disguise.Status"), beego.URLFor("Disguise.Status"), beego.URLFor("Disguise.Status"), "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultDisguiseService) TableButtonsType() []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "编辑",
			ClassName: "btn btn-sm btn-alt-primary open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Disguise.Edit", ":id", "__ID__", ":popup", 1),
				"data-area": "600px,400px",
			},
		},
		{
			Text:      "删除",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Disguise.Delete"),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}
