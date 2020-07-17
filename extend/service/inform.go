package service

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/inform"
	"github.com/go-playground/validator"
	"time"
)

type DefaultInformService struct{}

/** 根据ID查找一条详情数据 **/
func (service *DefaultInformService) FindOneContext(id int) inform.Context {
	var item inform.Context
	_ = orm.NewOrm().QueryTable(new(inform.Context)).Filter("id", id).One(&item)
	return item
}

/** 获取一条数据 **/
func (service *DefaultInformService) FindOne(id int) inform.Inform {
	var item inform.Inform
	_ = orm.NewOrm().QueryTable(new(inform.Inform)).Filter("id", id).One(&item)
	return item
}

/** 创建一条站内信详情 **/
func (service *DefaultInformService) CreateContext(context string) (inform.Context, error) {
	verifyBase := DefaultBaseVerify{}
	informItem := inform.Context{Context: context}
	if message := verifyBase.Begin().Struct(informItem); message != nil {
		return informItem, errors.New(verifyBase.Translate(message.(validator.ValidationErrors)))
	}
	if id, message := orm.NewOrm().Insert(&informItem); message != nil {
		return informItem, message
	} else {
		return service.FindOneContext(int(id)), message
	}
}

/**  创建接受者 **/
func (service *DefaultInformService) CreateInform(receiver []int, context, objectId, expediter int, form int8) ([]*inform.Inform, error) {
	var items []*inform.Inform
	var message error
	for _, v := range receiver {
		items = append(items, &inform.Inform{
			Form:      form,
			Receiver:  v,
			ObjectId:  objectId,
			ContextId: context,
			Expediter: expediter,
		})
	}
	if _, message = orm.NewOrm().InsertMulti(100, items); message == nil {
		_, _ = orm.NewOrm().QueryTable(new(inform.Inform)).Filter("context_id", context).All(&items)
	}
	return items, message
}

/**  同时创建表数据 **/
func (service *DefaultInformService) CreateInformAndContext(receiver []int, objectId, expediter int, form int8, context string) ([]*inform.Inform, error) {
	var (
		items   []*inform.Inform
		con     inform.Context
		message error
	)
	_ = orm.NewOrm().Begin()
	if con, message = service.CreateContext(context); message != nil {
		return items, message
	}
	if items, message = service.CreateInform(receiver, con.Id, objectId, expediter, form); message == nil {
		_ = orm.NewOrm().Commit()
	} else {
		_ = orm.NewOrm().Rollback()
	}
	return items, message
}

/** 获取一条消息结构的数据 **/
func (service *DefaultInformService) ItemSocketInformMassage(id int) (item inform.Message) {
	prefix := beego.AppConfig.String("db_prefix")
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select(prefix+"admin_inform.id", prefix+"admin_inform.receiver", prefix+"admin_inform.statue", prefix+"admin_inform_context.context", prefix+"admin_inform_context.create_time").
		From(prefix + "admin_inform").
		InnerJoin(prefix + "admin_inform_context").
		On(prefix + "admin_inform.context_id = " + prefix + "admin_inform_context.id").
		Where(prefix + "admin_inform.id = ?")
	_ = orm.NewOrm().Raw(qb.String(), id).QueryRow(&item)
	item.Context = beego.Substr(beego.HTML2str(item.Context), 0, 37)
	item.StringTime = beego.Date(item.CreateTime, "Y-m-d H:i:s")
	return item
}

/** 发送通知到socket上 **/
func (service *DefaultInformService) SendSocketInform(receiver []int, objectId, expediter int, form int8, context string) {
	items, message := service.CreateInformAndContext(receiver, objectId, expediter, form, context)
	if message == nil {
		for _, v := range items {
			item := service.ItemSocketInformMassage(v.Id)
			SocketInstanceGet().InformHandle(item)
		}
	} else {
		logs.Warn("创建通知表数据失败！失败原因:%s", error.Error(message))
	}
}

/** 发送临时通知通知到socket上 **/
func (service *DefaultInformService) SendSocketInformTemp(receiver int, context string) {
	SocketInstanceGet().InformHandle(inform.Message{
		Id: 0, Statue: 0, Receiver: receiver, Context: context,
		StringTime: beego.Date(time.Now(), "Y-m-d H:i:s"),
	})
}

/** 标记已读 **/
func (service *DefaultInformService) StatusArray(array []int, status int8) (e error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, e = orm.NewOrm().Update(&inform.Inform{Id: v, Statue: status}, "Statue"); e != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if e == nil {
		_ = orm.NewOrm().Commit()
	}
	return e
}

/** 批量删除 **/
func (service *DefaultInformService) DeleteArray(array []int) (e error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		find := service.FindOne(v)
		if _, e = orm.NewOrm().Delete(&inform.Inform{Id: find.Id}); e != nil {
			_ = orm.NewOrm().Rollback()
			break
		} else {
			count, _ := orm.NewOrm().QueryTable(new(inform.Inform)).Filter("context_id", find.ContextId).Count()
			if count == 0 {
				_, _ = orm.NewOrm().Delete(&inform.Context{Id: find.ContextId})
			}
		}
	}
	if e == nil {
		_ = orm.NewOrm().Commit()
	}
	return e
}

/** 获取头部消息中心 **/
func (service *DefaultInformService) HeaderItems() map[string]interface{} {
	var items []*inform.Message
	prefix := beego.AppConfig.String("db_prefix")
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select(prefix+"admin_inform.id", prefix+"admin_inform.statue", prefix+"admin_inform_context.context", prefix+"admin_inform_context.create_time").
		From(prefix + "admin_inform").
		InnerJoin(prefix + "admin_inform_context").
		On(prefix + "admin_inform.context_id = " + prefix + "admin_inform_context.id").Where(prefix + "admin_inform.statue = ?").OrderBy(prefix + "admin_inform.id DESC").Limit(3)
	_, _ = orm.NewOrm().Raw(qb.String(), 0).QueryRows(&items)
	count, _ := orm.NewOrm().QueryTable(new(inform.Inform)).Filter("statue", 0).Count()
	for _, v := range items {
		v.Context = beego.Substr(beego.HTML2str(v.Context), 0, 37)
	}
	return map[string]interface{}{"items": items, "count": count}
}

/************************************************表格渲染机制 ************************************************************/

/** 获取需要渲染的Column **/
func (service *DefaultInformService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "消息标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "是否已读", "name": "statue", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "消息类型", "name": "form", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "类型对象", "name": "object_id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "投递用户", "name": "expediter", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "消息操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultInformService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string"},
		"fieldName": {"id", "statue", "form", "object_id", "expediter"},
		"action":    {"", "", "", "", ""},
	}
	return result
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultInformService) DataTableButtons() []*table.TableButtons {
	var array []*table.TableButtons
	array = append(array, &table.TableButtons{
		Text:      "标记已读",
		ClassName: "btn btn-sm btn-alt-primary mt-1 ids_enables",
		Attribute: map[string]string{"data-action": beego.URLFor("Inform.Status"), "data-field": "status"},
	})
	array = append(array, &table.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Inform.Delete")},
	})
	array = append(array, &table.TableButtons{
		Text:      "全部已读",
		ClassName: "btn btn-sm btn-alt-success mt-1 jump_urls",
		Attribute: map[string]string{"data-action": beego.URLFor("Inform.Index", ":statue", 1)},
	})
	array = append(array, &table.TableButtons{
		Text:      "全部未读",
		ClassName: "btn btn-sm btn-alt-warning mt-1 jump_urls",
		Attribute: map[string]string{"data-action": beego.URLFor("Inform.Index", ":statue", 0)},
	})
	return array
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultInformService) TableButtonsType() []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "查看",
			ClassName: "btn btn-sm btn-alt-primary open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Inform.Check", ":id", "__ID__", ":popup", 1),
				"data-area": "380px,240px",
			},
		},
		{
			Text:      "删除",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Inform.Delete"),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}

/** 处理分页 **/
func (service *DefaultInformService) PageListItems(length, draw, page int, search string, statue int8) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(inform.Inform))
	recordsTotal, _ := qs.Count()
	if search != "" {
		qs = qs.Filter("form", search)
	}
	if statue > -1 {
		recordsTotal, _ = qs.Filter("statue", statue).Count()
		qs = qs.Filter("statue", statue)
	}
	_, _ = qs.Limit(length, length*(page-1)).OrderBy("statue", "-id").ValuesList(&lists, "id", "statue", "form", "object_id", "expediter")

	for _, v := range lists {
		v[1] = service.PageListItemStatus(v[1].(int64))
		v[3] = service.PageListItemObject(v[2].(int64), v[3].(int64))
		v[2] = service.PageListItemForm(v[2].(int64))
		v[4] = service.PageListItemUserName(v[4].(int64))
	}
	data := map[string]interface{}{
		"draw":            draw,         // 请求次数
		"recordsFiltered": recordsTotal, // 从多少条里面筛选
		"recordsTotal":    recordsTotal, // 总条数
		"data":            lists,        // 筛选结果
	}
	return data
}

/** 解析是否已读状态 **/
func (service *DefaultInformService) PageListItemStatus(statue int64) string {
	if statue == 1 {
		return `<span class="badge badge-primary badge-pill pt-1 pb-1">已读</span>`
	} else {
		return `<span class="badge badge-danger badge-pill pt-1 pb-1">未读</span>`
	}
}

/** 解析通知类型 **/
func (service *DefaultInformService) PageListItemForm(form int64) string {
	maps := []string{"未知类型", "爬虫日志", "查字词典", "栏目转换"}
	return maps[form]
}

/** 解析类型对象 **/
func (service *DefaultInformService) PageListItemObject(form, object int64) string {
	nameValue := "未知对象"
	switch form {
	case 1:
		service := DefaultCollectService{}
		if oneFind := service.One(int(object)); oneFind.Id > 0 {
			nameValue = oneFind.Name
		}
	case 2:
		nameValue = "查字词典"
	case 3:
		nameValue = "栏目转换"
	}
	return nameValue
}

/** 解析类型投递对象 **/
func (service *DefaultInformService) PageListItemUserName(expediter int64) string {
	username := "系统投递"
	if expediter > 0 {
		service := DefaultUsersService{}
		if oneFind, _ := service.Find(int(expediter)); oneFind.Id > 0 {
			username = oneFind.Nickname
		}
	}
	return username
}

/** 获取一条解析好的详细数据 **/
func (service *DefaultInformService) DetailInform(id int) map[string]interface{} {
	result := make(map[string]interface{})
	item := service.FindOne(id)
	context := service.FindOneContext(item.ContextId)
	result["id"] = item.Id
	result["form"] = service.PageListItemForm(int64(item.Form))
	result["statue"] = service.PageListItemStatus(int64(item.Statue))
	result["expediter"] = service.PageListItemUserName(int64(item.Expediter))
	result["object_id"] = service.PageListItemObject(int64(item.Form), int64(item.ObjectId))
	result["context"] = context.Context
	result["create_time"] = beego.Date(context.CreateTime, "Y年m月d日 H:i:s")
	return result
}
