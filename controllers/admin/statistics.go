package admin

/** 蜘蛛池统计面板 **/
type Statistics struct{ Main }

// @router /statistics/index [get,post]
func (c *Statistics) Index() {
	c.journalService.CachedHandleSetDebug() //todo 功能做好需要删除掉
	//蜘蛛统计
	items := c.journalService.CachedHandleAnalysisWeek()
	dayRes := c.journalService.CachedHandleGetDya()
	class := c.journalService.CachedHandleAnalysisClass(dayRes)
	c.Data["class"] = class
	c.Data["week"] = items
	// 推送统计
	c.Data["autoWeek"] = c.automaticService.AnalysisWeek()
}
