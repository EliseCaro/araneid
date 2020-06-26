package attachment

import "time"

type Attachment struct {
	Id         int       `orm:"pk;auto;column(id);type(int);description(主键,自增)" json:"id" form:"id"`
	Name       string    `orm:"column(name);size(50);type(char);default();description(附件名称)" json:"name" form:"name"`
	Path       string    `orm:"column(path);size(150);type(char);default();description(附件路径)" json:"path" form:"path"`
	Thumb      string    `orm:"column(thumb);size(150);type(char);default();description(附件缩略)" json:"thumb" form:"thumb"`
	Mime       string    `orm:"column(mime);size(200);type(char);default();description(附件类型)" json:"mime" form:"mime"`
	Ext        string    `orm:"column(ext);size(4);type(char);default();description(附件后缀名)" json:"ext" form:"ext"`
	Size       int64     `orm:"column(size);type(int);default(0);description(附件大小)" json:"size" form:"size"`
	Usage      int       `orm:"column(usage);type(int);default(0);description(使用数量)" json:"usage" form:"usage"`
	Sha1       string    `orm:"column(sha1);size(40);type(char);default();description(散列值)" json:"sha1" form:"sha1"`
	Driver     string    `orm:"column(driver);size(16);type(char);default();description(上传驱动)" json:"driver" form:"driver"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
}

/** 设置引擎 **/
func (m *Attachment) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Attachment) TableName() string {
	return "admin_attachment"
}
