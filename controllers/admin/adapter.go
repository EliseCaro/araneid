package admin

import (
	"fmt"
	"github.com/beatrice950201/araneid/controllers"
	"github.com/beatrice950201/araneid/extend/service"
)

/** 蜘蛛池栏目转换工具 **/
type Adapter struct {
	Main
}

// @router /adapter/index [get,post]
func (c *Adapter) Index() {}

/** 站长工具处理 **/
// @router /adapter/zhanzhang [post]
func (c *Adapter) Zhanzhang() {
	file := c.GetMustInt("files", "未检测到导入的正确文件文件")
	length, _ := c.GetInt("length", 0)
	form, _ := c.GetInt("form", 0)
	filtration := c.GetString("filtration", "")
	extract := c.GetString("extract", "")
	if form == 1 {
		go c.adapterService.SocketContextRuleHandle(file, length, form, c.UserInfo.Id, filtration, extract)
		c.Succeed(&controllers.ResultJson{Message: "已转为后台处理；结果将以站内信的形式通知！请注意查收！"})
	} else {
		if res, err := c.adapterService.ZhanzhangExtract(file, length, form, filtration, extract); err == nil {
			if path, err := c.adapterService.CreateXLSXFile(res); err == nil {
				name := new(service.DefaultAdjunctService).FindId(file).Name
				c.adapterService.CreateLogsInformStatus(name, "已经全部转换完成；请及时下载；系统将不会保留太长时间", "/"+path, "立即下载", c.UserInfo.Id)
				c.Succeed(&controllers.ResultJson{Message: "请及时下载您的文件；刷新页面后将无法还原！！", Data: path})
			} else {
				c.Fail(&controllers.ResultJson{Message: fmt.Sprintf(`处理错误：%s`, err.Error())})
			}
		} else {
			c.Fail(&controllers.ResultJson{Message: fmt.Sprintf(`处理错误：%s`, err.Error())})
		}
	}
}
