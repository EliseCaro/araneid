package service

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_func "github.com/beatrice950201/araneid/extend/func"
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

/** 重制一个域名数据 **/
func (service *DefaultDomainService) InitializedDomain(model, arachnid int, prefix, domain string) *spider.Domain {
	arachnidInfo, _ := new(DefaultArachnidService).Find(arachnid)
	matchModel, _ := new(DefaultMatchService).Find(arachnidInfo.Matching)
	item := &spider.Domain{
		Cate:   service.modelCateRandom(model),
		Name:   service.domainNameRandom(model, prefix),
		Title:  service.domainTagsRandom(model, arachnid, prefix, matchModel.IndexTitle),
		Links:  service.modelIndexRandom(arachnidInfo.Link, model, arachnid, arachnidInfo.Domain),
		Status: int8(1), Domain: fmt.Sprintf(`%s.%s`, prefix, domain),
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
	result = append(result, &map[string]string{"title": indexes.Title, "urls": indexes.Urls})
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
	var main = fmt.Sprintf(`%s.%s`, prefix, str)
	if _ = orm.NewOrm().QueryTable(new(spider.Domain)).Filter("domain", main).One(&domain); domain.Id > 0 {
		maps["title"] = domain.Title
	} else {
		maps["title"] = service.replaceSiteName(model, prefix, name)
	}
	maps["urls"] = fmt.Sprintf("http://%s/", main)
	return &maps
}

/** 获取一条域名配置不存在则创建 **/
func (service *DefaultDomainService) AcquireDomain(model, arachnid int, prefix, domain string) spider.Domain {
	var maps spider.Domain
	var main = fmt.Sprintf(`%s.%s`, prefix, domain)
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
