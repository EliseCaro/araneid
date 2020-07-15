package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_func "github.com/beatrice950201/araneid/extend/func"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/go-playground/validator"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

type DefaultDomainService struct{}

/** 将分类数据转为字符串id组 **/
func (service *DefaultDomainService) CateForIds(cate []*spider.Class) string {
	var id []string
	for _, i := range cate {
		id = append(id, strconv.Itoa(i.Id))
	}
	return fmt.Sprintf(`(%s)`, strings.Join(id, ","))
}

/** 获取一条数据 **/
func (service *DefaultDomainService) Find(id int) (spider.Domain, error) {
	item := spider.Domain{Id: id}
	return item, orm.NewOrm().Read(&item)
}

/** 根据URL获取一条数据 **/
func (service *DefaultDomainService) OneDomain(domain string) spider.Domain {
	var maps spider.Domain
	_ = orm.NewOrm().QueryTable(new(spider.Domain)).Filter("domain", domain).One(&maps)
	return maps
}

/** 获取全部站群数量 **/
func (service *DefaultDomainService) AliveAllNum() int64 {
	index, _ := orm.NewOrm().QueryTable(new(spider.Domain)).Count()
	return index
}

/** 返回带前缀或不带前缀的域名 **/
func (service *DefaultDomainService) getDomain(d, p string) string {
	var domain string
	if p != "" {
		domain = fmt.Sprintf(`%s.%s`, p, d)
	} else {
		domain = d
	}
	return domain
}

/** 重制一个域名数据 **/
func (service *DefaultDomainService) InitializedDomain(model, arachnid int, prefix, domain string) *spider.Domain {
	arachnidInfo, _ := new(DefaultArachnidService).Find(arachnid)
	matchModel, _ := new(DefaultMatchService).Find(arachnidInfo.Matching)
	item := &spider.Domain{
		Cate:   service.modelCateRandom(model),
		Name:   service.domainNameRandom(model, prefix),
		Title:  service.domainTagsRandom(model, arachnid, prefix, matchModel.IndexTitle),
		Links:  service.modelIndexRandom(arachnidInfo.Link, model, arachnid, arachnidInfo.Domain),
		Status: int8(1), Domain: service.getDomain(domain, prefix),
		Template:    service.modelTemplateRandom(model),
		Keywords:    service.domainTagsRandom(model, arachnid, prefix, matchModel.IndexKeyword),
		Arachnid:    arachnid,
		Description: service.domainTagsRandom(model, arachnid, prefix, matchModel.IndexDescription),
	}
	if message := new(DefaultBaseVerify).Begin().Struct(item); message != nil {
		logs.Error(`创建站点失败；失败原因：%s`, new(DefaultBaseVerify).Translate(message.(validator.ValidationErrors)))
	} else {
		if _, err := orm.NewOrm().Insert(item); err == nil {
			*item = service.AcquireDomain(model, arachnid, prefix, domain)
		} else {
			logs.Error(`创建站点失败；失败原因：%s`, err.Error())
		}
	}
	return item
}

/** 底部友情链接 **/
func (service *DefaultDomainService) modelIndexRandom(number, model, arachnid int, textarea string) string {
	maps := _func.ParseAttrConfigArray(textarea)
	count := int(math.Floor(float64(number / len(maps))))
	var result []*map[string]string
	for _, domain := range maps {
		result = append(result, service.modelLinksRange(model, count, domain)...)
	}
	indexes := new(DefaultIndexesService).UsageOneIndexes(arachnid)
	if indexes.Title != "" {
		result = append(result, &map[string]string{"title": indexes.Title, "urls": indexes.Urls})
	}
	bytes, _ := json.Marshal(result)
	return string(bytes)
}

/** 获取一组友情链接 **/
func (service *DefaultDomainService) modelLinksRange(model, count int, domain string) []*map[string]string {
	var prefix = beego.AppConfig.String("db_prefix")
	var items []*spider.Prefix
	var result []*map[string]string
	sql := fmt.Sprintf(`SELECT title,tags FROM %sspider_prefix WHERE model=? ORDER BY RAND() LIMIT %d`, prefix, count)
	_, _ = orm.NewOrm().Raw(sql, model).QueryRows(&items)
	for _, item := range items {
		result = append(result, service.webSiteLinks(model, item.Tags, domain, item.Title))
	}
	return result
}

/** 解析一个友情链接 **/
func (service *DefaultDomainService) webSiteLinks(model int, prefix, str, name string) *map[string]string {
	var maps = make(map[string]string)
	var domain spider.Domain
	beego.Warn(str, prefix)
	var main = service.getDomain(str, prefix)
	if _ = orm.NewOrm().QueryTable(new(spider.Domain)).Filter("domain", main).One(&domain); domain.Id > 0 {
		maps["title"] = domain.Name
	} else {
		maps["title"] = service.replaceSiteName(model, prefix, name)
	}
	maps["urls"] = fmt.Sprintf("http://%s/", main)
	return &maps
}

/** 获取一条域名配置不存在则创建 **/
func (service *DefaultDomainService) AcquireDomain(model, arachnid int, prefix, domain string) spider.Domain {
	var maps spider.Domain
	var main = service.getDomain(domain, prefix)
	if _ = orm.NewOrm().QueryTable(new(spider.Domain)).Filter("domain", main).One(&maps); maps.Id <= 0 {
		maps = *service.InitializedDomain(model, arachnid, prefix, domain)
	}
	return maps
}

/** 分配模板 **/
func (service *DefaultDomainService) modelTemplateRandom(model int) string {
	modelDetail := new(DefaultModelsService).One(model)
	result := new(DefaultTemplateService).Items(modelDetail.Template)
	item := result[rand.Intn(len(result))]
	return item.Name
}

/** 处理标签 **/
func (service *DefaultDomainService) domainTagsRandom(model, arachnid int, prefix, str string) string {
	resultTitle := service.replaceSiteName(model, prefix, str)
	return service.replaceRandomKeyword(arachnid, resultTitle)
}

/** 替换标签 #关键词# todo 此sql效率极低**/
func (service *DefaultDomainService) replaceRandomKeyword(arachnid int, str string) string {
	var dbPrefix = beego.AppConfig.String("db_prefix")
	var items []*spider.Keyword
	sql := fmt.Sprintf(`SELECT title FROM %sspider_keyword WHERE arachnid=? ORDER BY RAND() LIMIT %d`, dbPrefix, strings.Count(str, "#关键词#"))
	_, _ = orm.NewOrm().Raw(sql, arachnid).QueryRows(&items)
	for _, v := range items {
		str = strings.Replace(str, "#关键词#", v.Title, 1)
	}
	return str
}

/** 替换标签 #站点名# **/
func (service *DefaultDomainService) replaceSiteName(model int, prefix, str string) string {
	siteName := service.domainNameRandom(model, prefix)
	re, _ := regexp.Compile(`#站点名#`)
	return re.ReplaceAllString(str, siteName)
}

/** 获取分类返回json字符串;按照挂载次数从低到高获取;todo 分类加数否需要后台控制 **/
func (service *DefaultDomainService) modelCateRandom(model int) string {
	var maps []*spider.Class
	_, _ = orm.NewOrm().QueryTable(new(spider.Class)).Filter("model", model).OrderBy("usage").Limit(10).All(&maps)
	for _, v := range maps {
		_ = new(DefaultClassService).Inc(v.Id)
	}
	result, _ := json.Marshal(maps)
	return string(result)
}

/** 获取站点名称；todo 考虑从后台控制是否匹配随机关键字 **/
func (service *DefaultDomainService) domainNameRandom(model int, prefix string) string {
	var maps spider.Prefix
	_ = orm.NewOrm().QueryTable(new(spider.Prefix)).Filter("model", model).Filter("tags", prefix).One(&maps)
	return maps.Title
}

/** 获取是否存在分类池 **/
func (service *DefaultDomainService) extArticleCategory(domain int) spider.Category {
	var maps spider.Category
	_ = orm.NewOrm().QueryTable(new(spider.Category)).Filter("domain", domain).One(&maps)
	return maps
}

/** 获取是否存在文章 **/
func (service *DefaultDomainService) extArticleDetail(domain int) spider.Detail {
	var maps spider.Detail
	_ = orm.NewOrm().QueryTable(new(spider.Detail)).Filter("domain", domain).One(&maps)
	return maps
}

/************后台控制器*******************8/

/** 更新状态 **/
func (service *DefaultDomainService) StatusArray(array []int, status int8) (e error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, e = orm.NewOrm().Update(&spider.Domain{Id: v, Status: status}, "Status"); e != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if e == nil {
		_ = orm.NewOrm().Commit()
	}
	return e
}

/** 批量删除 **/
func (service *DefaultDomainService) DeleteArray(array []int) (message error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if service.extArticleCategory(v).Id == 0 && service.extArticleDetail(v).Id == 0 {
			if _, message = orm.NewOrm().Delete(&spider.Domain{Id: v}); message != nil {
				_ = orm.NewOrm().Rollback()
				break
			}
		} else {
			message = errors.New("该项目下还有分类缓存池或文章缓存池！")
			break
		}
	}
	if message == nil {
		_ = orm.NewOrm().Commit()
	}
	return message
}

/** 清空 **/
func (service *DefaultDomainService) EmptyDelete(arachnid int) {
	var item []*spider.Domain
	qs := orm.NewOrm().QueryTable(new(spider.Domain))
	if arachnid > 0 {
		qs = qs.Filter("arachnid", arachnid)
	}
	_, _ = qs.All(&item)
	for _, v := range item {
		new(DefaultCategoryService).EmptyDelete(v.Id)
		new(DefaultDetailService).EmptyDelete(v.Id)
		_, _ = orm.NewOrm().Delete(&spider.Domain{Id: v.Id})
	}
}

/************************************************表格渲染机制 ************************************************************/

/** 获取需要渲染的Column **/
func (service *DefaultDomainService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "名称", "name": "name", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "域名", "name": "domain", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "项目", "name": "arachnid", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "模板", "name": "template", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "状态", "name": "status", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "更新时间", "name": "update_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultDomainService) DataTableButtons(id int) []*table.TableButtons {
	var array []*table.TableButtons
	if id > 0 {
		array = append(array, &table.TableButtons{
			Text:      "返回上级",
			ClassName: "btn btn-sm btn-alt-success mt-1 jump_urls",
			Attribute: map[string]string{"data-action": beego.URLFor("Arachnid.Index")},
		})
	}
	array = append(array, &table.TableButtons{
		Text:      "重建缓存",
		ClassName: "btn btn-sm btn-alt-danger mt-1 js-tooltip ids_deletes",
		Attribute: map[string]string{
			"data-action":         beego.URLFor("Domain.Empty", ":parent", id),
			"data-toggle":         "tooltip",
			"data-original-title": "重建缓存后所有站点数据将会重建；所有链接模板将出现变动；请务必考虑清楚在重建；",
		},
	})
	array = append(array, &table.TableButtons{
		Text:      "启用选中",
		ClassName: "btn btn-sm btn-alt-primary mt-1 ids_enables",
		Attribute: map[string]string{"data-action": beego.URLFor("Domain.Status", ":parent", id), "data-field": "status"},
	})
	array = append(array, &table.TableButtons{
		Text:      "禁用选中",
		ClassName: "btn btn-sm btn-alt-warning mt-1 ids_disables",
		Attribute: map[string]string{"data-action": beego.URLFor("Domain.Status", ":parent", id), "data-field": "status"},
	})
	array = append(array, &table.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Domain.Delete", ":parent", id)},
	})
	return array
}

/** 处理分页 **/
func (service *DefaultDomainService) PageListItems(length, draw, page int, search string, id int) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(spider.Domain))
	if search != "" {
		qs = qs.Filter("name__icontains", search)
	}
	if id > 0 {
		qs = qs.Filter("arachnid", id)
	}
	recordsTotal, _ := qs.Count()
	_, _ = qs.Limit(length, length*(page-1)).OrderBy("-id").ValuesList(&lists, "id", "name", "domain", "arachnid", "template", "status", "update_time")
	for _, v := range lists {
		one, _ := new(DefaultArachnidService).Find(int(v[3].(int64)))
		v[3] = one.Name
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
func (service *DefaultDomainService) TableColumnsType(id int) map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string", "switch", "date"},
		"fieldName": {"id", "name", "domain", "arachnid", "template", "status", "update_time"},
		"action":    {"", "", "", "", "", beego.URLFor("Domain.Status", ":parent", id), ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultDomainService) TableButtonsType(id int) []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "友情链接",
			ClassName: "btn btn-sm btn-alt-primary open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Domain.Links", ":id", "__ID__", ":popup", 1, ":parent", id),
				"data-area": "600px,400px",
			},
		},
		{
			Text:      "挂载分类",
			ClassName: "btn btn-sm btn-alt-success jump_urls",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Category.Index", ":parent", "__ID__"),
			},
		},
		{
			Text:      "挂载文章",
			ClassName: "btn btn-sm btn-alt-info jump_urls",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Detail.Index", ":parent", "__ID__"),
			},
		},
		{
			Text:      "编辑",
			ClassName: "btn btn-sm btn-alt-warning open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Domain.Edit", ":id", "__ID__", ":popup", 1, ":parent", id),
				"data-area": "600px,375px",
			},
		},
		{
			Text:      "删除",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Domain.Delete", ":parent", id),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}
