package index

import (
	"github.com/astaxie/beego/orm"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/index"
	"github.com/beatrice950201/araneid/extend/model/spider"
)

type Lists struct{ Main }

// @router /column-:id([0-9]+)-:page([0-9]+).html [get]
func (c *Lists) Index() {
	if c.spiderExtend {
		c.assignCate()
	} else {
		c.assignItem()
	}
}

/** 提交主站列表 **/
func (c *Lists) assignItem() {
	pageSize := 2
	pageNum := c.getPage()
	qm := orm.NewOrm().QueryTable(new(index.News)).Filter("status", 1)
	var list []*index.News
	_, _ = qm.OrderBy("-id").Limit(pageSize, pageNum*pageSize-pageSize).All(&list)
	count, _ := qm.Count()
	c.Data["pageList"] = _func.PageUtil(int(count), pageNum, pageSize, list)
}

/** 提交分类数据 **/
func (c *Lists) assignCate() {
	id := c.GetMustInt(":id", "非法请求～")
	pageSize := 10
	pageNum := c.getPage()
	qm := orm.NewOrm().QueryTable(new(spider.Article)).Filter("class", id)
	var list []*spider.Article
	_, _ = qm.OrderBy("-id").Limit(pageSize, pageNum*pageSize-pageSize).All(&list)
	count, _ := qm.Count()
	c.Data["pageList"] = _func.PageUtil(int(count), pageNum, pageSize, list)
	c.Data["cateInfo"] = c.categoryService.AcquireCategory(c.DomainCache.Id, c.DomainCache.Arachnid, id)
}

/** 解析页码 **/
func (c *Lists) getPage() int {
	page, _ := c.GetInt(":page", 1)
	if page == 0 {
		page += 1
	}
	return page
}
