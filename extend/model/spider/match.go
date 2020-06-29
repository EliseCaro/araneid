package spider

/*** KTDB调用模板 **/
// #关键词# [ 替换一个随机关键词 ]
// #原文# [ 替换成原文 ]

// 内容通过轮询关键词模糊查询链接模型或者索引池里面是否存在相似标题的链接；找到则替换原文词语为a标签；实现相互传导

type Match struct {
	Id                int    `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Name              string `orm:"column(name);size(30);type(char);default();description(模板名称)" json:"name" form:"name" validate:"required" label:"模板名称"`
	IndexTitle        string `orm:"column(index_title);size(150);type(char);default();description(首页标题)" json:"index_title" form:"index_title" validate:"required" label:"首页标题"`
	IndexKeyword      string `orm:"column(index_keyword);size(200);type(char);default();description(首页关键字)" json:"index_keyword" form:"index_keyword" validate:"required" label:"首页关键字"`
	IndexDescription  string `orm:"column(index_description);size(255);type(char);default();description(首页描述)" json:"index_description" form:"index_description" validate:"required" label:"首页描述"`
	CateTitle         string `orm:"column(cate_title);size(150);type(char);default();description(列表页标题)" json:"cate_title" form:"cate_title" validate:"required" label:"列表页标题"`
	CateKeyword       string `orm:"column(cate_keyword);size(200);type(char);default();description(列表页关键字)" json:"cate_keyword" form:"cate_keyword" validate:"required" label:"列表页关键字"`
	CateDescription   string `orm:"column(cate_description);size(255);type(char);default();description(列表页描述)" json:"cate_description" form:"cate_description" validate:"required" label:"列表页描述"`
	DetailTitle       string `orm:"column(detail_title);size(150);type(char);default();description(详情页标题)" json:"detail_title" form:"detail_title" validate:"required" label:"详情页标题"`
	DetailKeyword     string `orm:"column(detail_keyword);size(200);type(char);default();description(详情页关键字)" json:"detail_keyword" form:"detail_keyword" validate:"required" label:"详情页关键字"`
	DetailDescription string `orm:"column(detail_description);size(255);type(char);default();description(详情页描述)" json:"detail_description" form:"detail_description" validate:"required" label:"详情页描述"`
}

/** 设置引擎 **/
func (m *Match) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Match) TableName() string {
	return "spider_match"
}
