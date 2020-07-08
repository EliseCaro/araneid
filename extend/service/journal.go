package service

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	_func "github.com/beatrice950201/araneid/extend/func"
	"net/http"
	"strings"
)

type DefaultJournalService struct{}

/** 获取一条数 **/
func (service *DefaultJournalService) HandleInstantiation(ctx *context.Context) {
	beego.Warn(ctx.Input.Domain())
	beego.Warn(service.GetURL(ctx.Request))
	beego.Warn(ctx.Input.IP())
	beego.Warn(ctx.Input.UserAgent())
	service.UserAgent()
}

/** 获取完整url **/
func (service *DefaultJournalService) GetURL(r *http.Request) (Url string) {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	return strings.Join([]string{scheme, r.Host, r.RequestURI}, "")
}

/** 获取蜘蛛标识 **/
func (service *DefaultJournalService) UserAgent() {
	class := _func.ParseAttrConfigArray(beego.AppConfig.String("system_spider_class"))
	beego.Warn(class)
}
