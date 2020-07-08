package admin

import (
	"encoding/json"
	"github.com/astaxie/beego"
)

/** 蜘蛛池统计面板 **/
type Statistics struct{ Main }

// @router /statistics/index [get,post]
func (c *Statistics) Index() {
	c.journalService.CachedHandleSetDebug() //todo 功能做好需要删除掉
	week, _ := json.Marshal(c.journalService.CachedHandleGetWeek())
	c.Data["week"] = string(week)
	day := c.journalService.CachedHandleGetDya()
	beego.Warn(day)
}
