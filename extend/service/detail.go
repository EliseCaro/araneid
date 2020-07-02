package service

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/go-playground/validator"
)

type DefaultDetailService struct{}

/** 获取一条详情数据； 不存在则创建 **/
func (service *DefaultDetailService) AcquireDetail(did, aid, oid int) spider.Detail {
	var maps spider.Detail
	if _ = orm.NewOrm().QueryTable(new(spider.Detail)).Filter("domain", did).Filter("arachnid", aid).Filter("oid", oid).One(&maps); maps.Id <= 0 {
		maps = *service.InitializedDetail(did, aid, oid)
	}
	_ = new(DefaultArticleService).Inc(oid, "view")
	return maps
}

/** 初始化一个详情 **/
func (service *DefaultDetailService) InitializedDetail(did, aid, oid int) *spider.Detail {
	detail, _ := new(DefaultArachnidService).Find(aid)
	match, _ := new(DefaultMatchService).Find(detail.Matching)
	articleOne := new(DefaultArticleService).One(oid)
	domain, _ := new(DefaultDomainService).Find(did)
	object := new(DefaultCategoryService)
	item := &spider.Detail{
		Oid:         oid,
		Domain:      did,
		Cid:         articleOne.Class,
		Model:       detail.Models,
		Arachnid:    aid,
		Title:       object.TagsRandom(aid, domain.Name, articleOne.Title, match.DetailTitle),
		Keywords:    object.TagsRandom(aid, domain.Name, articleOne.Keywords, match.DetailKeyword),
		Description: object.TagsRandom(aid, domain.Name, articleOne.Description, match.DetailDescription),
	}
	item.Context = object.TagsRandomContext(detail.Id, articleOne.Context, item.Keywords)
	if message := new(DefaultBaseVerify).Begin().Struct(item); message != nil {
		logs.Error(`初始化文章[ %s ]失败；失败原因：%s`, articleOne.Title, new(DefaultBaseVerify).Translate(message.(validator.ValidationErrors)))
	} else {
		if _, err := orm.NewOrm().Insert(item); err == nil {
			*item = service.AcquireDetail(did, aid, oid)
			_ = new(DefaultArticleService).Inc(oid, "usage")
		} else {
			logs.Error(`初始化分类[ %s ]失败；失败原因：%s`, articleOne.Title, err.Error())
		}
	}
	return item
}
