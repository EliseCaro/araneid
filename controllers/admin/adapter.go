package admin

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/beatrice950201/araneid/controllers"
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
	filtration := c.GetString("filtration", "")
	extract := c.GetString("extract", "")
	if res, err := c.adapterService.ZhanzhangExtract(file, length, filtration, extract); err == nil {
		beego.Warn(res) // 将分类数据写入文件
	} else {
		c.Fail(&controllers.ResultJson{Message: fmt.Sprintf(`处理错误：%s`, err.Error())})
	}
}
