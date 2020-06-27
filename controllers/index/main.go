package index

import (
	"github.com/beatrice950201/araneid/controllers"
	"github.com/beatrice950201/araneid/extend/service"
)

type Main struct {
	controllers.Base
	articleService service.DefaultArticleService
}

/** 准备下一级构造函数 **/
type NextPreparer interface {
	NextPrepare()
}

/** 实现上一级构造 **/
func (c *Main) NestPrepare() {
	if app, ok := c.AppController.(NextPreparer); ok {
		app.NextPrepare()
	}
}
