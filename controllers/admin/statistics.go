package admin

import "github.com/astaxie/beego"

/** 蜘蛛池统计面板 **/
type Statistics struct{ Main }

// @router /statistics/index [get,post]
func (c *Statistics) Index() {
	w := c.journalService.CachedHandleGetWeek()
	beego.Warn(w)
	beego.Info("__________________1222____________________________")
	d := c.journalService.CachedHandleGetDya()
	beego.Warn(d)
}
