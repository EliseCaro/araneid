package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/extend/model/automatic"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"io/ioutil"
	"net/http"
	"strconv"
)

/** 自动推送url链接层 **/
type DefaultAutomaticService struct{}

/** 根据文档ID推送一条链接 **/
func (service *DefaultAutomaticService) AutomaticDocument(object int) {
	detail := new(DefaultArticleService).One(object)
	var domain spider.Domain
	if detail.Class > 0 {
		if cate := new(DefaultCategoryService).AcquireCategoryOldWhere(detail.Class); cate.Domain > 0 {
			domain, _ = new(DefaultDomainService).Find(cate.Domain)
		}
	}
	if domain.Domain != "" && domain.Submit != "" {
		urls := fmt.Sprintf(`http://%s/index/detail-%d.html`, domain.Domain, detail.Id)
		service.baiduAutomatic(domain.Domain, domain.Submit, urls)
	}
}

/** 百度提交 **/
func (service *DefaultAutomaticService) baiduAutomatic(domain, token, urls string) {
	action := fmt.Sprintf(beego.AppConfig.String("baidu_automatic_urls"), domain, token)
	result := make(map[string]interface{})
	client := &http.Client{}
	r, _ := http.NewRequest("POST", action, bytes.NewBuffer([]byte(urls)))
	resp, err := client.Do(r)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err == nil {
		if body, message := ioutil.ReadAll(resp.Body); message == nil {
			message = json.Unmarshal(body, &result)
			if result["success"] != nil && result["success"].(int64) == 1 {
				service.CreateLogs(domain, urls, strconv.Itoa(result["remain"].(int)), 1)
			} else {
				service.CreateLogs(domain, urls, result["message"].(string), 0)
			}
		} else {
			service.CreateLogs(domain, urls, message.Error(), 0)
		}
		defer resp.Body.Close()
	} else {
		service.CreateLogs(domain, urls, err.Error(), 0)
	}
}

/** 根据老id查询新ID数据 **/
func (service *DefaultAutomaticService) AcquireUrlsOne(urls string) automatic.Automatic {
	var maps automatic.Automatic
	_ = orm.NewOrm().QueryTable(new(automatic.Automatic)).Filter("urls", urls).One(&maps)
	return maps
}

/** 创建推送记录 **/
func (service *DefaultAutomaticService) CreateLogs(domain, urls, remark string, status int8) {
	item := automatic.Automatic{Domain: domain, Urls: urls, Remark: remark, Status: status}
	if index := service.AcquireUrlsOne(urls); index.Id > 0 {
		_, _ = orm.NewOrm().Update(&item)
	} else {
		_, _ = orm.NewOrm().Insert(&item)
	}
}
