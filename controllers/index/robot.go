package index

import (
	"github.com/astaxie/beego"
	"github.com/beatrice950201/araneid/extend/model/spider"
)

/** 伪原创输出接口 **/
type Robot struct {
	Main
	key    string
	secret string
	ones   *spider.Disguise
}

/** 构造函数  **/
func (c *Robot) NextPrepare() {
	key := c.GetMustString("appid", "非法请求！appid错误！")
	secret := c.GetMustString("secret", "非法请求！secret错误！")
	c.ones = c.disguiseService.KeyOne(key, secret)
	beego.Error(c.ones)
}
