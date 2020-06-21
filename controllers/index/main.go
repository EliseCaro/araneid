package index

import "github.com/beatrice950201/araneid/controllers"

type Main struct {
	controllers.Base
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
