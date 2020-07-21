package index

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	"github.com/beatrice950201/araneid/extend/model/index"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/go-playground/validator"
)

type Index struct{ Main }

// @router / [get]
func (c *Index) Index() {
	if !c.spiderExtend {
		c.Data["advantage"] = c.advantage()   // 产品优势
		c.Data["news"] = c.newsList()         //  新闻资讯
		c.Data["team"] = c.teamList()         //  团队信息
		c.Data["Template"] = c.templateList() //  模板信息
		c.Data["message"] = c.messageList()   //  反馈信息
	}
}

// @router /index/message [post]
func (c *Index) Message() {
	if c.IsAjax() {
		item := index.Message{}
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Insert(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "提交成功；请耐心等待", Url: beego.URLFor("Index.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "提交失败，请稍后再试！"})
		}
	}
}

/** 获取产品优势列表**/
func (c *Index) advantage() []*index.Advantage {
	var item []*index.Advantage
	_, _ = orm.NewOrm().QueryTable(new(index.Advantage)).Filter("status", 1).All(&item)
	return item
}

/** 获取新闻动态 **/
func (c *Index) newsList() []*index.News {
	var item []*index.News
	_, _ = orm.NewOrm().QueryTable(new(index.News)).Filter("status", 1).OrderBy("-id").Limit(3, 0).All(&item)
	return item
}

/** 获取团队信息 **/
func (c *Index) teamList() []*index.Team {
	var item []*index.Team
	_, _ = orm.NewOrm().QueryTable(new(index.Team)).Filter("status", 1).All(&item)
	return item
}

/** 获取团队信息 **/
func (c *Index) messageList() []*index.Message {
	var item []*index.Message
	_, _ = orm.NewOrm().QueryTable(new(index.Message)).Filter("status", 1).OrderBy("-id").Limit(6, 0).All(&item)
	return item
}

/** 获取模板信息 **/
func (c *Index) templateList() []*spider.TemplateConfig {
	var group []*spider.TemplateClass
	var items []*spider.TemplateConfig
	class := c.templateService.Groups()
	if j, _ := json.Marshal(class); len(j) > 0 {
		_ = json.Unmarshal(j, &group)
	}
	for _, v := range group {
		items = append(items, c.templateService.Items(v.Id)...)
	}
	return items
}
