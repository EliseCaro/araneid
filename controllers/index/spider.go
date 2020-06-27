package index

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	"github.com/beatrice950201/araneid/extend/model/spider"
)

type Spider struct {
	Main
	Password string
	DbPrefix string
}

/** 构造函数  **/
func (c *Spider) NextPrepare() {
	c.Password = "$2a$04$fA.06wUF6PbGZFY/lrQWeuOEqdvIkfsbPp5jmXjSIMzrkllhNbPkK"
	c.DbPrefix = beego.AppConfig.String("db_prefix")
}

// @router /spider/api [get,post]
func (c *Spider) Api() {
	password := c.GetMustString("password", "验证器密码是空的！")
	if password == c.Password {
		if e := c.ResultInsert(c.ResultMaps()); e != nil {
			c.Fail(&controllers.ResultJson{Message: "入库发生错误；错误原因：" + e.Error()})
		} else {
			c.Succeed(&controllers.ResultJson{Message: "文章已经发布到对应平台！"})
		}
	} else {
		c.Fail(&controllers.ResultJson{Message: "验证器密码错误！"})
	}
}

/*** 接受内容包解析 **/
func (c *Spider) ResultMaps() spider.Article {
	result := c.GetMustString("result", "采集资源包是空的！")
	var mapsResult spider.Article
	e := json.Unmarshal([]byte(result), &mapsResult)
	if e != nil {
		c.Fail(&controllers.ResultJson{Message: "解析内容包失败！"})
	}
	mapsResult.Object = c.GetMustInt("object", "文档ID是空的！")
	mapsResult.Model = c.GetMustInt(":module", "蜘蛛池模型ID是空的！")
	mapsResult.Class = c.randomClass(mapsResult.Model).Id
	return mapsResult
}

/** 更新或者插入 **/
func (c *Spider) ResultInsert(res spider.Article) (err error) {
	ext := c.articleService.OneByObject(res.Object)
	if ext.Id > 0 {
		res.Id = ext.Id
		_, err = orm.NewOrm().Update(&res)
	} else {
		_, err = orm.NewOrm().Insert(&res)
	}
	return err
}

/** 获取一条随机ID**/
func (c *Spider) randomClass(module int) spider.Class {
	var item spider.Class
	sql := `SELECT * FROM ` + c.DbPrefix + `spider_class AS t1 JOIN (SELECT ROUND(RAND()*(SELECT MAX(id) FROM ` + c.DbPrefix + `spider_class)) AS id) AS t2 WHERE t1.id>=t2.id AND t1.model=? ORDER BY t1.id LIMIT 1`
	_ = orm.NewOrm().Raw(sql, module).QueryRow(&item)
	return item
}
