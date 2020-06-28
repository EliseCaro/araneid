package service

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"strings"
)

type DefaultRobotService struct{}

/** 根据原词提取近义词数据 **/
func (service *DefaultRobotService) OneString(name *string) spider.Robot {
	var item spider.Robot
	var resemblance []string
	cond := orm.NewCondition().And("name", name).Or("resemblance__icontains", fmt.Sprintf(`"%s"`, *name))
	_ = orm.NewOrm().QueryTable(new(spider.Robot)).SetCond(cond).One(&item)
	_ = json.Unmarshal([]byte(item.Resemblance), &resemblance)
	item.Resemblance = strings.Join(resemblance, ",")
	return item
}

/** 添加一条训练结果；返回json字符串 **/
func (service *DefaultRobotService) Insert(s string, maps []*string) []*string {
	jsonByte, _ := json.Marshal(maps)
	_, _ = orm.NewOrm().Insert(&spider.Robot{
		Name:        s,
		Resemblance: string(jsonByte),
	})
	return maps
}
