package index

import (
	"encoding/json"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/beatrice950201/araneid/extend/service"
)

type Index struct{ Main }

// @router / [get]
func (c *Index) Index() {
	if c.spiderExtend {
		c.spiderIndex()
	}
}

/**  蜘蛛池首页 **/
func (c *Index) spiderIndex() {
	var cate []*spider.Class
	var links []*spider.Indexes
	_ = json.Unmarshal([]byte(c.DomainCache.Cate), &cate)
	_ = json.Unmarshal([]byte(c.DomainCache.Links), &links)
	c.Data["cate"] = cate
	c.Data["idsCate"] = new(service.DefaultDomainService).CateForIds(cate)
	c.Data["links"] = links
	c.Data["model"] = c.Model
	c.Data["arachnid"], _ = new(service.DefaultArachnidService).Find(c.DomainCache.Arachnid)
	c.Data["info"] = map[string]string{"name": c.DomainCache.Name, "title": c.DomainCache.Title, "keywords": c.DomainCache.Keywords, "description": c.DomainCache.Description}
}
