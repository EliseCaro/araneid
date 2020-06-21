package admin

import _func "github.com/beatrice950201/araneid/extend/func"

/** 查词典采集器 **/
type Dictionaries struct {
	Main
}

// 采集分类
// @router /dictionaries/cate [get,post]
func (c *Dictionaries) Cate() {
	if c.IsAjax() {

	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.dictionariesService.DataTableCateColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.dictionariesService.DataTableCateButtons() {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, _func.WebPageSize())
}
