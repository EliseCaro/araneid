package index

import (
	"github.com/astaxie/beego"
	"github.com/beatrice950201/araneid/controllers"
	"github.com/beatrice950201/araneid/extend/service"
)

type Main struct {
	controllers.Base
	articleService  service.DefaultArticleService
	classService    service.DefaultClassService
	disguiseService service.DefaultDisguiseService
	modelsService   service.DefaultModelsService
}

/** 准备下一级构造函数 **/
type NextPreparer interface {
	NextPrepare()
}

/** 实现上一级构造 **/
func (c *Main) NestPrepare() {
	c.DomainCheck(c.indexCheck)
	if app, ok := c.AppController.(NextPreparer); ok {
		app.NextPrepare()
	}
}

/** 检测前台域名 **/
func (c *Main) indexCheck(prefix, main string) bool {
	beego.Info(prefix, main)
	return false
}
