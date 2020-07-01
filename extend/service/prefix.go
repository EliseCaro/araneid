package service

import (
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/extend/model/spider"
)

type DefaultPrefixService struct{}

/** 获取一条 **/
func (service *DefaultPrefixService) One(id int) spider.Prefix {
	var item spider.Prefix
	_ = orm.NewOrm().QueryTable(new(spider.Prefix)).Filter("id", id).One(&item)
	return item
}

/** 获取模型所有域名前缀 **/
func (service *DefaultPrefixService) Select(model int) []*spider.Prefix {
	var maps []*spider.Prefix
	_, _ = orm.NewOrm().QueryTable(new(spider.Prefix)).Filter("model", model).OrderBy("-id").All(&maps)
	return maps
}

/** 根据模型跟前缀或一条数 **/
func (service *DefaultPrefixService) OnePrefix(model int, prefix string) spider.Prefix {
	var item spider.Prefix
	_ = orm.NewOrm().QueryTable(new(spider.Prefix)).Filter("tags", prefix).Filter("model", model).One(&item)
	return item
}
