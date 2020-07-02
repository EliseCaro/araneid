package service

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/go-playground/validator"
	"regexp"
	"strings"
)

type DefaultCategoryService struct{}

/** 获取一条分类不存在则创建 **/
func (service *DefaultCategoryService) AcquireCategory(did, aid, cid int) spider.Category {
	var maps spider.Category
	if _ = orm.NewOrm().QueryTable(new(spider.Category)).Filter("domain", did).Filter("arachnid", aid).Filter("cid", cid).One(&maps); maps.Id <= 0 {
		maps = *service.InitializedCategory(did, aid, cid)
	}
	return maps
}

/** 初始化一个分类 **/
func (service *DefaultCategoryService) InitializedCategory(did, aid, cid int) *spider.Category {
	detail, _ := new(DefaultArachnidService).Find(aid)
	match, _ := new(DefaultMatchService).Find(detail.Matching)
	cateOne := new(DefaultClassService).One(cid)
	domain, _ := new(DefaultDomainService).Find(did)
	item := &spider.Category{
		Domain:      did,
		Cid:         cid,
		Model:       detail.Models,
		Title:       service.TagsRandom(aid, domain.Name, cateOne.Title, match.CateTitle),
		Arachnid:    aid,
		Keywords:    service.TagsRandom(aid, domain.Name, cateOne.Keywords, match.CateKeyword),
		Description: service.TagsRandom(aid, domain.Name, cateOne.Description, match.CateDescription),
	}
	if message := new(DefaultBaseVerify).Begin().Struct(item); message != nil {
		logs.Error(`初始化分类[ %s ]失败；失败原因：%s`, cateOne.Title, new(DefaultBaseVerify).Translate(message.(validator.ValidationErrors)))
	} else {
		if _, err := orm.NewOrm().Insert(item); err == nil {
			*item = service.AcquireCategory(did, aid, cid)
		} else {
			logs.Error(`初始化分类[ %s ]失败；失败原因：%s`, cateOne.Title, err.Error())
		}
	}
	return item
}

/** 处理内容标签 内容通过轮询关键词模糊查询链接模型或者索引池里面是否存在相似标题的链接；找到则替换原文词语为a标签；实现相互传导 **/
func (service *DefaultCategoryService) TagsRandomContext(id int, result, keyword string) string {
	detail, _ := new(DefaultArachnidService).Find(id)
	keywords := strings.Split(keyword, ",")
	var indexes []map[string]string
	for _, v := range keywords { // 处理索引先
		if item := service.indexesMatchingLike(id, v); item.Id > 1 && detail.Indexes > len(indexes)+1 {
			indexes = append(indexes, map[string]string{"title": item.Title, "urls": item.Urls, "keyword": v})
		}
	}
	var articles []map[string]string
	for _, v := range keywords { // 处理轮链先
		if item := service.articlesMatchingLike(id, v); item.Id > 1 && detail.Zoology > len(articles)+1 {
			articles = append(articles, map[string]string{"title": item.Title, "urls": item.Urls, "keyword": v})
		}
	}
	articles = append(articles, indexes...)
	for _, v := range articles {
		n := fmt.Sprintf(`<a href="%s" title="%s">%s</a>`, v["urls"], v["title"], v["keyword"])
		result = strings.Replace(result, v["keyword"], n, 1)
	}
	return result
}

/** 根据关键词提取索引 **/
func (service *DefaultCategoryService) indexesMatchingLike(arachnid int, keyword string) spider.Indexes {
	var maps spider.Indexes
	_ = orm.NewOrm().QueryTable(new(spider.Indexes)).Filter("arachnid", arachnid).Filter("title__icontains", keyword).OrderBy("-sort", "usage").One(&maps)
	return maps
}

/** 根据关键词获取轮链索引**/
func (service *DefaultCategoryService) articlesMatchingLike(arachnid int, keyword string) *spider.Indexes {
	var maps spider.Detail
	cond := orm.NewCondition().
		AndCond(orm.NewCondition().And("arachnid", arachnid).And("title__icontains", keyword)).
		OrCond(orm.NewCondition().And("arachnid", arachnid).And("context__icontains", keyword)).
		OrCond(orm.NewCondition().And("arachnid", arachnid).And("description__icontains", keyword)).
		OrCond(orm.NewCondition().And("arachnid", arachnid).And("keywords__icontains", keyword))
	_ = orm.NewOrm().QueryTable(new(spider.Detail)).SetCond(cond).OrderBy("-id").One(&maps)
	if maps.Id > 0 {
		d, _ := new(DefaultDomainService).Find(maps.Domain)
		return &spider.Indexes{
			Title: maps.Title,
			Urls:  fmt.Sprintf("http://%s/index/detail-%d.html", d.Domain, maps.Oid),
			Id:    maps.Id,
		}
	} else {
		return &spider.Indexes{}
	}
}

/** 处理标签 **/
func (service *DefaultCategoryService) TagsRandom(aid int, title, result, match string) string {
	result = service.replaceOldText(result, match)
	result = service.replaceRandomKeyword(aid, result)
	result = service.replaceSiteName(title, result)
	return result
}

/** 替换标签 #关键词# todo 此sql效率极低**/
func (service *DefaultCategoryService) replaceRandomKeyword(aid int, str string) string {
	var dbPrefix = beego.AppConfig.String("db_prefix")
	var items []*spider.Keyword
	sql := fmt.Sprintf(`SELECT title FROM %sspider_keyword WHERE arachnid=? ORDER BY RAND() LIMIT %d`, dbPrefix, strings.Count(str, "#关键词#"))
	_, _ = orm.NewOrm().Raw(sql, aid).QueryRows(&items)
	for _, v := range items {
		str = strings.Replace(str, "#关键词#", v.Title, 1)
	}
	return str
}

/** 替换标签 #站点名# **/
func (service *DefaultCategoryService) replaceSiteName(result, match string) string {
	re, _ := regexp.Compile(`#站点名#`)
	return re.ReplaceAllString(match, result) // 在match中将#站点名#替换为result
}

/** 替换标签 #原文# **/
func (service *DefaultCategoryService) replaceOldText(result, match string) string {
	re, _ := regexp.Compile(`#原文#`)
	return re.ReplaceAllString(match, result) //   在match中将#站点名#替换为result
}
