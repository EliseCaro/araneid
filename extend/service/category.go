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
	if _ = orm.NewOrm().QueryTable(new(spider.Category)).Filter("arachnid", aid).Filter("cid", cid).One(&maps); maps.Id <= 0 {
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
