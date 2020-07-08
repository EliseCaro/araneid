package admin

import (
	"github.com/astaxie/beego"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
)

/** 蜘蛛池统计面板 **/
type Statistics struct{ Main }

// @router /statistics/index [get,post]
func (c *Statistics) Index() {
	w := c.journalService.CachedHandleGetWeek()
	beego.Warn(w)
	beego.Info("__________________1222____________________________")
	d := c.journalService.CachedHandleGetDya()
	beego.Warn(d)

	_ = _func.SetCache("test", 1233)
	test := _func.GetCache("test")
	beego.Error(test)
	beego.Info("__________________1222____________________________")

	_ = _func.SetRedisCache("class", c.classService.One(1), 86400)
	var r spider.Class
	class := _func.GetRedisCache("class")
	beego.Error(r)
	beego.Error(class)
}
