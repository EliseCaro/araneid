package index

import (
	"github.com/beatrice950201/araneid/controllers"
	"github.com/beatrice950201/araneid/extend/model/spider"
)

/** 伪原创输出接口 **/
type Robot struct {
	Main
	key    string
	secret string
	ones   spider.Disguise
}

/** 构造函数  **/
func (c *Robot) NextPrepare() {
	key := c.GetMustString("appid", "非法请求！appid错误！")
	secret := c.GetMustString("secret", "非法请求！secret错误！")
	c.ones = c.disguiseService.KeyOne(key, secret)
}

// @router /robot/robot [post]
func (c *Robot) Robot() {
	title := c.GetMustString("title", "原始标题字段不能为空！")
	context := c.GetMustString("context", "原始内容字段不能为空！")
	Keyword := c.GetMustString("Keyword", "原始关键字不能为空！")
	description := c.GetMustString("description", "原始描述不能为空！")
	res, err := c.disguiseService.DisguiseHandleManage(c.ones.Id, &spider.HandleModule{
		Title: title, Context: context, Keywords: Keyword, Description: description,
	})
	if err != nil {
		c.Fail(&controllers.ResultJson{Message: "伪原创处理失败！失败原因：" + err.Error(), Data: c.ones})
	} else {
		c.Succeed(&controllers.ResultJson{Message: "ok", Data: res})
	}
}
