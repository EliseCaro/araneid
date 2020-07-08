package admin

/** 蜘蛛池统计面板 **/
type Statistics struct{ Main }

// @router /statistics/index [get,post]
func (c *Statistics) Index() {
	c.journalService.CachedHandleSetDebug() //todo 功能做好需要删除掉
	items := c.journalService.CachedHandleAnalysisWeek()
	class := c.journalService.CachedHandleAnalysisClass(c.journalService.CachedHandleGetDya())
	c.Data["class"] = class
	c.Data["week"] = items
}
