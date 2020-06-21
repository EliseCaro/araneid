package users

import (
	"time"
)

type Users struct {
	Id          int       `orm:"pk;auto;column(id);type(int);description(主键)" json:"id" form:"id"`
	Username    string    `orm:"unique;column(username);size(32);type(char);default();description(账号)" json:"username" form:"username" validate:"required" label:"账号"`
	Nickname    string    `orm:"column(nickname);size(50);type(char);default();description(昵称)" json:"nickname" form:"nickname" validate:"required" label:"昵称"`
	Password    string    `orm:"column(password);size(96);type(char);default();description(密码)" json:"password" form:"password" validate:"required" label:"密码"`
	Email       string    `orm:"column(email);size(64);type(char);default();description(邮箱)" json:"email" form:"email" validate:"required,email" label:"邮箱"`
	Mobile      string    `orm:"column(mobile);size(11);type(char);default();description(手机)"json:"mobile" form:"mobile" validate:"required,gte=0,lte=11" label:"手机号码"`
	Avatar      int       `orm:"column(avatar);type(int);default(0);description(头像)" json:"avatar" form:"avatar" validate:"required,min=1" label:"头像"`
	Role        int       `orm:"index;column(role);type(int);default(0);description(角色)"  json:"role" form:"role" validate:"required,min=1" label:"角色"`
	CreateTime  time.Time `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime  time.Time `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
	LastLoginIp int64     `orm:"column(last_login_ip);default(0);description(最后IP)" json:"last_login_ip"`
	Sort        int       `orm:"column(sort);type(int);default(0);description(排序)" json:"sort" form:"sort"`
	Status      int8      `orm:"index;column(status);default(0);description(启用状态)" json:"status" form:"status"`
}

/** 设置引擎 **/
func (m *Users) TableEngine() string {
	return "INNODB"
}

/** 设置表名 **/
func (m *Users) TableName() string {
	return "admin_users"
}
