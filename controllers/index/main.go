package index

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/beatrice950201/araneid/controllers"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/beatrice950201/araneid/extend/service"
)

type Main struct {
	controllers.Base
	Arachnid        spider.Arachnid
	Model           spider.Models
	DomainCache     spider.Domain
	DomainPrefix    string
	DomainMain      string
	spiderExtend    bool
	articleService  service.DefaultArticleService
	classService    service.DefaultClassService
	disguiseService service.DefaultDisguiseService
	modelsService   service.DefaultModelsService
	arachnidService service.DefaultArachnidService
	prefixService   service.DefaultPrefixService
	templateService service.DefaultTemplateService
	domainService   service.DefaultDomainService
	categoryService service.DefaultCategoryService
	detailService   service.DefaultDetailService
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
	c.DomainCheck(c.mainCheckDomain)
	if c.spiderExtend {
		c.assignVolt()
	}
	if app, ok := c.AppController.(NextPreparer); ok {
		app.NextPrepare()
	}
}

/** 提交蜘蛛池默认数据 **/
func (c *Main) assignVolt() {
	var (
		cate  []*spider.Class
		links []*spider.Indexes
	)
	json.Unmarshal([]byte(c.DomainCache.Cate), &cate)
	json.Unmarshal([]byte(c.DomainCache.Links), &links)
	c.Data["cate"] = cate
	c.Data["cids"] = new(service.DefaultDomainService).CateForIds(cate)
	c.Data["links"] = links
	c.Data["model"] = c.Model
	c.Data["domain"] = map[string]string{"name": c.DomainCache.Name, "title": c.DomainCache.Title, "keywords": c.DomainCache.Keywords, "description": c.DomainCache.Description}
	c.Data["arachnid"], _ = new(service.DefaultArachnidService).Find(c.DomainCache.Arachnid)
}

/** 区分蜘蛛池跟主站 **/
func (c *Main) mainCheckDomain(prefix, main string) (bool, string, string) {
	adminDomain := beego.AppConfig.String("system_admin_domain")
	c.DomainPrefix = prefix
	c.DomainMain = main
	if adminDomain != fmt.Sprintf("%s.%s", prefix, main) {
		c.spiderExtend = true
		return c.indexCheck()
	} else {
		c.spiderExtend = false
		return true, "success!!", "success!!"
	}
}

/** 检测前台域名 蜘蛛池绑定 **/
func (c *Main) indexCheck() (bool, string, string) {
	message := &checkInitialize{}
	arachnid := c.arachnidService.FindDomain(c.DomainMain)
	if arachnid.Id >= 0 && arachnid.Status == 1 {
		c.Arachnid = arachnid
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
		return c.indexCheckTemplate()
	} else {
		message = &checkInitialize{
			Title:   "初始化错误",
			Message: fmt.Sprintf(`在[%s]中未找到%s域名下已经挂载的%s前缀`, c.Model.Name, c.DomainMain, c.DomainPrefix),
		}
	}
	return message.Status, message.Title, message.Message
}

/** 检测模板分配 **/
func (c *Main) indexCheckTemplate() (bool, string, string) {
	var message checkInitialize
	if result := c.templateService.Items(c.Model.Template); len(result) > 0 {
		return c.indexSetTemplate()
	} else {
		message = checkInitialize{
			Title:   "初始化错误",
			Message: fmt.Sprintf(`在[%s]中没有找到可以渲染模板`, c.Model.Name),
		}
	}
	return message.Status, message.Title, message.Message
}

/** 分配模板分类数据 **/
func (c *Main) indexSetTemplate() (bool, string, string) {
	var message checkInitialize
	if result := c.domainService.AcquireDomain(c.Model.Id, c.Arachnid.Id, c.DomainPrefix, c.DomainMain); result.Id > 0 {
		if result.Status == 0 {
			message = checkInitialize{Title: "初始化错误", Message: fmt.Sprintf(`域名 %s 在蜘蛛池中已被停用；请联系管理员检查权限～`, result.Domain)}
		} else {
			c.SetTemplate(fmt.Sprintf("%s/%s", beego.AppConfig.String("spider_views_path"), result.Template))
			c.DomainCache = result
			message = checkInitialize{Status: true}
		}
	} else {
		message = checkInitialize{Title: "初始化错误", Message: fmt.Sprintf(`初始化%s模板实例失败；请检查数据合法性`, c.Model.Name)}
	}
	return message.Status, message.Title, message.Message
}
