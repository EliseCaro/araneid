package service

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/go-playground/validator"
)

type DefaultDetailService struct{}

/** 获取一条详情数据； 不存在则创建 **/
func (service *DefaultDetailService) AcquireDetail(did, aid, oid int) spider.Detail {
	var maps spider.Detail
	if _ = orm.NewOrm().QueryTable(new(spider.Detail)).Filter("domain", did).Filter("arachnid", aid).Filter("oid", oid).One(&maps); maps.Id <= 0 {
		maps = *service.InitializedDetail(did, aid, oid)
	}
	_ = new(DefaultArticleService).Inc(oid, "view")
	return maps
}

/** 初始化一个详情 **/
func (service *DefaultDetailService) InitializedDetail(did, aid, oid int) *spider.Detail {
	detail, _ := new(DefaultArachnidService).Find(aid)
	match, _ := new(DefaultMatchService).Find(detail.Matching)
	articleOne := new(DefaultArticleService).One(oid)
	domain, _ := new(DefaultDomainService).Find(did)
	object := new(DefaultCategoryService)
	item := &spider.Detail{
		Oid:         oid,
		Domain:      did,
		Cid:         articleOne.Class,
		Model:       detail.Models,
		Arachnid:    aid,
		Title:       object.TagsRandom(aid, domain.Name, articleOne.Title, match.DetailTitle),
		Keywords:    object.TagsRandom(aid, domain.Name, articleOne.Keywords, match.DetailKeyword),
		Description: object.TagsRandom(aid, domain.Name, articleOne.Description, match.DetailDescription),
	}
	item.Context = object.TagsRandomContext(detail.Id, articleOne.Context, item.Keywords)
	if message := new(DefaultBaseVerify).Begin().Struct(item); message != nil {
		logs.Error(`初始化文章[ %s ]失败；失败原因：%s`, articleOne.Title, new(DefaultBaseVerify).Translate(message.(validator.ValidationErrors)))
	} else {
		if _, err := orm.NewOrm().Insert(item); err == nil {
			*item = service.AcquireDetail(did, aid, oid)
			_ = new(DefaultArticleService).Inc(oid, "usage")
		} else {
			logs.Error(`初始化分类[ %s ]失败；失败原因：%s`, articleOne.Title, err.Error())
		}
	}
	return item
}

/** 批量删除 **/
func (service *DefaultDetailService) DeleteArray(array []int) (message error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, message = orm.NewOrm().Delete(&spider.Detail{Id: v}); message != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if message == nil {
		_ = orm.NewOrm().Commit()
	}
	return message
}

/** 清空 **/
func (service *DefaultDetailService) EmptyDelete(domain int) {
	var item []*spider.Detail
	qs := orm.NewOrm().QueryTable(new(spider.Detail))
	if domain > 0 {
		qs = qs.Filter("domain", domain)
	}
	_, _ = qs.All(&item)
	for _, v := range item {
		_, _ = orm.NewOrm().Delete(&spider.Detail{Id: v.Id})
	}
}

/** 获取一条数据 **/
func (service *DefaultDetailService) Find(id int) (spider.Detail, error) {
	item := spider.Detail{Id: id}
	return item, orm.NewOrm().Read(&item)
}

/************************************************表格渲染机制 ************************************************************/

/** 获取需要渲染的Column **/
func (service *DefaultDetailService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "域名", "name": "domain", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "标题", "name": "title", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "项目", "name": "arachnid", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "关键字", "name": "keywords", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "更新时间", "name": "update_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultDetailService) DataTableButtons(id, arachnid int) []*table.TableButtons {
	var array []*table.TableButtons
	if id > 0 {
		array = append(array, &table.TableButtons{
			Text:      "返回域名",
			ClassName: "btn btn-sm btn-alt-success mt-1 jump_urls",
			Attribute: map[string]string{"data-action": beego.URLFor("Domain.Index", ":arachnid", arachnid)},
		})
	}
	array = append(array, &table.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Detail.Delete", ":parent", id)},
	})
	array = append(array, &table.TableButtons{
		Text:      "清空缓存",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Detail.Empty", ":parent", id)},
	})
	return array
}

/** 处理分页 **/
func (service *DefaultDetailService) PageListItems(length, draw, page int, search string, id int) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(spider.Detail))
	if search != "" {
		qs = qs.Filter("title__icontains", search)
	}
	if id > 0 {
		qs = qs.Filter("domain", id)
	}
	recordsTotal, _ := qs.Count()
	_, _ = qs.Limit(length, length*(page-1)).OrderBy("-id").ValuesList(&lists, "id", "domain", "title", "arachnid", "keywords", "update_time")
	for _, v := range lists {
		domain, _ := new(DefaultDomainService).Find(int(v[1].(int64)))
		arachnid, _ := new(DefaultArachnidService).Find(int(v[3].(int64)))
		v[2] = service.substr2HtmlHref(fmt.Sprintf("http://%s/index/detail-%d.html", domain.Domain, v[0].(int64)), v[2].(string), 0, 20)
		v[4] = service.substr2HtmlHref(fmt.Sprintf("http://%s/index/detail-%d.html", domain.Domain, v[0].(int64)), v[4].(string), 0, 20)
		v[1] = domain.Domain
		v[3] = arachnid.Name
	}
	data := map[string]interface{}{
		"draw":            draw,         // 请求次数
		"recordsFiltered": recordsTotal, // 从多少条里面筛选
		"recordsTotal":    recordsTotal, // 总条数
		"data":            lists,        // 筛选结果
	}
	return data
}

/**  转为pop提示 **/
func (service *DefaultDetailService) substr2HtmlHref(u, s string, start, end int) string {
	html := fmt.Sprintf(`<a href="%s" target="_blank" class="badge badge-primary js-tooltip" data-placement="top" data-toggle="tooltip" data-original-title="%s">%s...</a>`, u, s, beego.Substr(s, start, end))
	return html
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultDetailService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string", "date"},
		"fieldName": {"id", "domain", "title", "arachnid", "keywords", "update_time"},
		"action":    {"", "", "", "", "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultDetailService) TableButtonsType(id int) []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "编辑",
			ClassName: "btn btn-sm btn-alt-warning jump_urls",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Detail.Edit", ":id", "__ID__", ":parent", id),
			},
		},
		{
			Text:      "删除",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Detail.Delete", ":parent", id),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}
