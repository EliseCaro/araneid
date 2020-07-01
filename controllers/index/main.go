package index

import (
	"fmt"
	"github.com/beatrice950201/araneid/controllers"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/beatrice950201/araneid/extend/service"
)

type Main struct {
	controllers.Base
	Arachnid        spider.Arachnid
	Model           spider.Models
	DomainPrefix    string
	DomainMain      string
	articleService  service.DefaultArticleService
	classService    service.DefaultClassService
	disguiseService service.DefaultDisguiseService
	modelsService   service.DefaultModelsService
	arachnidService service.DefaultArachnidService
	prefixService   service.DefaultPrefixService
}

/** 准备下一级构造函数 **/
type NextPreparer interface {
	NextPrepare()
}

/** 初始化message **/
type checkInitialize struct {
	Status  bool
	Title   string
	Message string
}

/** 实现上一级构造 **/
func (c *Main) NestPrepare() {
	c.DomainCheck(c.indexCheck)
	if app, ok := c.AppController.(NextPreparer); ok {
		app.NextPrepare()
	}
}

/** 检测前台域名 蜘蛛池绑定 **/
func (c *Main) indexCheck(prefix, main string) (bool, string, string) {
	message := &checkInitialize{}
	arachnid := c.arachnidService.FindDomain(main)
	if arachnid.Id >= 0 && arachnid.Status == 1 {
		c.Arachnid = arachnid
		c.DomainPrefix = prefix
		c.DomainMain = main
		return c.indexCheckModel()
	} else {
		message = &checkInitialize{
			Title:   "初始化错误",
			Message: "该蜘蛛池未启用或者该域名未注册到蜘蛛池",
		}
	}
	return message.Status, message.Title, message.Message
}

/** 检测前台域名模型 todo 需要修改被蜘蛛池挂载的模型不能删除 **/
func (c *Main) indexCheckModel() (bool, string, string) {
	message := &checkInitialize{}
	if module := c.modelsService.One(c.Arachnid.Models); module.Id > 0 && module.Template > 0 {
		c.Model = module
		return c.indexCheckPrefix()
	} else {
		message = &checkInitialize{
			Title:   "初始化错误",
			Message: fmt.Sprintf(`在%s项目未挂载模型或者模型被删除或者模型未挂载模板分组`, c.Arachnid.Name),
		}
	}
	return message.Status, message.Title, message.Message
}

/** 检测域名前缀 **/
func (c *Main) indexCheckPrefix() (bool, string, string) {
	message := &checkInitialize{}
	if module := c.prefixService.OnePrefix(c.Model.Id, c.DomainPrefix); module.Id > 0 {
		message = &checkInitialize{Status: true}
	} else {
		message = &checkInitialize{
			Title:   "初始化错误",
			Message: fmt.Sprintf(`在[%s]中未找到%s域名下已经挂载的%s前缀`, c.Model.Name, c.DomainMain, c.DomainPrefix),
		}
	}
	return message.Status, message.Title, message.Message
}
