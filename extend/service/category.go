package service

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"github.com/go-playground/validator"
	"regexp"
	"strings"
)

type DefaultCategoryService struct{}

/** 获取一条分类不存在则创建 **/
func (service *DefaultCategoryService) AcquireCategory(did, aid, cid int) spider.Category {
	var maps spider.Category
	if _ = orm.NewOrm().QueryTable(new(spider.Category)).Filter("domain", did).Filter("arachnid", aid).Filter("cid", cid).One(&maps); maps.Id <= 0 {
		maps = *service.InitializedCategory(did, aid, cid)
	}
	return maps
}

/** 初始化一个分类 **/
func (service *DefaultCategoryService) InitializedCategory(did, aid, cid int) *spider.Category {
	detail, _ := new(DefaultArachnidService).Find(aid)
	match, _ := new(DefaultMatchService).Find(detail.Matching)
	cateOne := new(DefaultClassService).One(cid)
	domain, _ := new(DefaultDomainService).Find(did)
	item := &spider.Category{
		Domain:      did,
		Cid:         cid,
		Model:       detail.Models,
		Title:       service.TagsRandom(aid, domain.Name, cateOne.Title, match.CateTitle),
		Arachnid:    aid,
		Keywords:    service.TagsRandom(aid, domain.Name, cateOne.Keywords, match.CateKeyword),
		Description: service.TagsRandom(aid, domain.Name, cateOne.Description, match.CateDescription),
	}
	if message := new(DefaultBaseVerify).Begin().Struct(item); message != nil {
		logs.Error(`初始化分类[ %s ]失败；失败原因：%s`, cateOne.Title, new(DefaultBaseVerify).Translate(message.(validator.ValidationErrors)))
	} else {
		if _, err := orm.NewOrm().Insert(item); err == nil {
			*item = service.AcquireCategory(did, aid, cid)
			_ = service.InitializeDomainCate(did) // 重建缓存分区
		} else {
			logs.Error(`初始化分类[ %s ]失败；失败原因：%s`, cateOne.Title, err.Error())
		}
	}
	return item
}

/** 处理内容标签 内容通过轮询关键词模糊查询链接模型或者索引池里面是否存在相似标题的链接；找到则替换原文词语为a标签；实现相互传导 **/
func (service *DefaultCategoryService) TagsRandomContext(id int, result, keyword string) string {
	detail, _ := new(DefaultArachnidService).Find(id)
	keywords := strings.Split(keyword, ",")
	var indexes []map[string]string
	for _, v := range keywords { // 处理索引先
		if item := service.indexesMatchingLike(id, v); item.Id > 1 && detail.Indexes > len(indexes)+1 {
			indexes = append(indexes, map[string]string{"title": item.Title, "urls": item.Urls, "keyword": v})
		}
	}
	var articles []map[string]string
	for _, v := range keywords { // 处理轮链先
		if item := service.articlesMatchingLike(id, v); item.Id > 1 && detail.Zoology > len(articles)+1 {
			articles = append(articles, map[string]string{"title": item.Title, "urls": item.Urls, "keyword": v})
		}
	}
	articles = append(articles, indexes...)
	for _, v := range articles {
		n := fmt.Sprintf(`<a href="%s" title="%s">%s</a>`, v["urls"], v["title"], v["keyword"])
		result = strings.Replace(result, v["keyword"], n, 1)
	}
	return result
}

/** 根据关键词提取索引 **/
func (service *DefaultCategoryService) indexesMatchingLike(arachnid int, keyword string) spider.Indexes {
	var maps spider.Indexes
	_ = orm.NewOrm().QueryTable(new(spider.Indexes)).Filter("arachnid", arachnid).Filter("title__icontains", keyword).OrderBy("-sort", "usage").One(&maps)
	return maps
}

/** 根据关键词获取轮链索引**/
func (service *DefaultCategoryService) articlesMatchingLike(arachnid int, keyword string) *spider.Indexes {
	var maps spider.Detail
	cond := orm.NewCondition().
		AndCond(orm.NewCondition().And("arachnid", arachnid).And("title__icontains", keyword)).
		OrCond(orm.NewCondition().And("arachnid", arachnid).And("context__icontains", keyword)).
		OrCond(orm.NewCondition().And("arachnid", arachnid).And("description__icontains", keyword)).
		OrCond(orm.NewCondition().And("arachnid", arachnid).And("keywords__icontains", keyword))
	_ = orm.NewOrm().QueryTable(new(spider.Detail)).SetCond(cond).OrderBy("-id").One(&maps)
	if maps.Id > 0 {
		d, _ := new(DefaultDomainService).Find(maps.Domain)
		return &spider.Indexes{
			Title: maps.Title,
			Urls:  fmt.Sprintf("http://%s/index/detail-%d.html", d.Domain, maps.Oid),
			Id:    maps.Id,
		}
	} else {
		return &spider.Indexes{}
	}
}

/** 处理标签 **/
func (service *DefaultCategoryService) TagsRandom(aid int, title, result, match string) string {
	result = service.replaceOldText(result, match)
	result = service.replaceRandomKeyword(aid, result)
	result = service.replaceSiteName(title, result)
	return result
}

/** 替换标签 #关键词#**/
func (service *DefaultCategoryService) replaceRandomKeyword(aid int, str string) string {
	var dbPrefix = beego.AppConfig.String("db_prefix")
	var items []*spider.Keyword
	sql := fmt.Sprintf(`SELECT title FROM %sspider_keyword WHERE arachnid=? ORDER BY RAND() LIMIT %d`, dbPrefix, strings.Count(str, "#关键词#"))
	_, _ = orm.NewOrm().Raw(sql, aid).QueryRows(&items)
	for _, v := range items {
		str = strings.Replace(str, "#关键词#", v.Title, 1)
	}
	return str
}

/** 替换标签 #站点名# **/
func (service *DefaultCategoryService) replaceSiteName(result, match string) string {
	re, _ := regexp.Compile(`#站点名#`)
	return re.ReplaceAllString(match, result) // 在match中将#站点名#替换为result
}

/** 替换标签 #原文# **/
func (service *DefaultCategoryService) replaceOldText(result, match string) string {
	re, _ := regexp.Compile(`#原文#`)
	return re.ReplaceAllString(match, result) //   在match中将#站点名#替换为result
}

/** 获取一条数据 **/
func (service *DefaultCategoryService) Find(id int) (spider.Category, error) {
	item := spider.Category{Id: id}
	return item, orm.NewOrm().Read(&item)
}

/** 重建域名分类json **/
func (service *DefaultCategoryService) InitializeDomainCate(d int) error {
	var maps []*spider.Category
	var result []*spider.Class
	_, _ = orm.NewOrm().QueryTable(new(spider.Category)).Filter("domain", d).All(&maps)
	for _, v := range maps {
		class := new(DefaultClassService).One(v.Cid)
		class.Title = v.Title
		class.Keywords = v.Keywords
		class.Description = v.Description
		result = append(result, &class)
	}
	stringJson, _ := json.Marshal(result)
	_, err := orm.NewOrm().Update(&spider.Domain{Id: d, Cate: string(stringJson)}, "Cate")
	return err
}

/** 批量删除 **/
func (service *DefaultCategoryService) DeleteArray(array []int) (message error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, message = orm.NewOrm().Delete(&spider.Category{Id: v}); message != nil {
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
func (service *DefaultCategoryService) EmptyDelete(domain int) {
	var item []*spider.Category
	qs := orm.NewOrm().QueryTable(new(spider.Category))
	if domain > 0 {
		qs = qs.Filter("domain", domain)
	}
	_, _ = qs.All(&item)
	for _, v := range item {
		_, _ = orm.NewOrm().Delete(&spider.Category{Id: v.Id})
	}
}

/************************************************表格渲染机制 ************************************************************/
/** 获取需要渲染的Column **/
func (service *DefaultCategoryService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "标题", "name": "title", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "域名", "name": "domain", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "项目", "name": "arachnid", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "关键字", "name": "keywords", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "更新时间", "name": "update_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultCategoryService) DataTableButtons(id, arachnid int) []*table.TableButtons {
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
		Attribute: map[string]string{"data-action": beego.URLFor("Category.Delete", ":parent", id)},
	})
	array = append(array, &table.TableButtons{
		Text:      "清空缓存",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Category.Empty", ":parent", id)},
	})
	return array
}

/** 处理分页 **/
func (service *DefaultCategoryService) PageListItems(length, draw, page int, search string, id int) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(spider.Category))
	if search != "" {
		qs = qs.Filter("title__icontains", search)
	}
	if id > 0 {
		qs = qs.Filter("domain", id)
	}
	recordsTotal, _ := qs.Count()
	_, _ = qs.Limit(length, length*(page-1)).OrderBy("-id").ValuesList(&lists, "id", "title", "domain", "arachnid", "keywords", "update_time")
	for _, v := range lists {
		domain, _ := new(DefaultDomainService).Find(int(v[2].(int64)))
		arachnid, _ := new(DefaultArachnidService).Find(int(v[3].(int64)))
		v[1] = service.substr2HtmlHref(fmt.Sprintf("http://%s/index/column-%d-0.html", domain.Domain, v[0].(int64)), v[1].(string), 0, 20)
		v[4] = service.substr2HtmlHref(fmt.Sprintf("http://%s/index/column-%d-0.html", domain.Domain, v[0].(int64)), v[4].(string), 0, 20)
		v[2] = domain.Domain
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
func (service *DefaultCategoryService) substr2HtmlHref(u, s string, start, end int) string {
	html := fmt.Sprintf(`<a href="%s" target="_blank" class="badge badge-primary js-tooltip" data-placement="top" data-toggle="tooltip" data-original-title="%s">%s...</a>`, u, s, beego.Substr(s, start, end))
	return html
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultCategoryService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string", "date"},
		"fieldName": {"id", "title", "domain", "arachnid", "keywords", "update_time"},
		"action":    {"", "", "", "", "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultCategoryService) TableButtonsType(id int) []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "编辑",
			ClassName: "btn btn-sm btn-alt-warning open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Category.Edit", ":id", "__ID__", ":popup", 1, ":parent", id),
				"data-area": "580px,380px",
			},
		},
		{
			Text:      "删除",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Category.Delete", ":parent", id),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}
