package index

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/extend/model/index"
)

type Detail struct{ Main }

// @router /detail-:id([0-9]+).html [get]
func (c *Detail) Index() {
	if c.spiderExtend {
		c.assignDetail()
	} else {
		c.assignDetailIndex()
		c.Data["rand"] = c.indexArticleLimit(3)
	}
}

/** 提交详情数据 **/
func (c *Detail) assignDetailIndex() {
	id := c.GetMustInt(":id", "非法请求～")
	item := index.News{Id: id}
	_ = orm.NewOrm().Read(&item)
	c.Data["info"] = item
}

/** 获取随机文章 **/
func (c *Detail) indexArticleLimit(limit int) []index.News {
	var prefix = beego.AppConfig.String("db_prefix")
	var maps []index.News
	sql := fmt.Sprintf(`SELECT id,title,description,cover FROM %sindex_news WHERE status=1 ORDER BY RAND() LIMIT %d`, prefix, limit)
	_, _ = orm.NewOrm().Raw(sql).QueryRows(&maps)
	return maps
}

/** 提交详情数据 **/
func (c *Detail) assignDetail() {
	id := c.GetMustInt(":id", "非法请求～")
	detail := c.detailService.AcquireDetail(c.DomainCache.Id, c.DomainCache.Arachnid, id)
	c.Data["cateInfo"] = c.categoryService.AcquireCategory(c.DomainCache.Id, c.DomainCache.Arachnid, detail.Cid)
	c.Data["detail"] = detail
}
