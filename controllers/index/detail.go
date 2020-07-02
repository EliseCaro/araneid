package index

type Detail struct{ Main }

// @router /detail-:id([0-9]+).html [get]
func (c *Detail) Index() {
	if c.spiderExtend {
		c.assignDetail()
	}
}

/** 提交详情数据 **/
func (c *Detail) assignDetail() {
	id := c.GetMustInt(":id", "非法请求～")
	detail := c.detailService.AcquireDetail(c.DomainCache.Id, c.DomainCache.Arachnid, id)
	c.Data["cateInfo"] = c.categoryService.AcquireCategory(c.DomainCache.Id, c.DomainCache.Arachnid, detail.Cid)
	c.Data["detail"] = detail
}
