package service

import (
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"strings"
)

type DefaultRobotService struct{}

/** 根据原词提取近义词数据 **/
func (service *DefaultRobotService) One(name *string) spider.Robot {
	var item spider.Robot
	_ = orm.NewOrm().QueryTable(new(spider.Robot)).Filter("name", name).One(&item)
	return item
}

/** 添加一条训练结果 **/
func (service *DefaultRobotService) Insert(s string, maps []string) string {
	maps = append(maps, s)
	_, _ = orm.NewOrm().Insert(&spider.Robot{
		Name:        s,
		Resemblance: strings.Join(maps, ","),
	})
	return strings.Join(maps, ",")
}
