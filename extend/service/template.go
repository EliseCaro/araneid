package service

import (
	"encoding/base64"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	//"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/extend/model/spider"
	"os"
	"path/filepath"
)

type DefaultTemplateService struct{}

/** 获取一条分组数据 **/
func (service *DefaultTemplateService) One(id int) spider.Template {
	var item spider.Template
	_ = orm.NewOrm().QueryTable(new(spider.Template)).Filter("id", id).One(&item)
	return item
}

/** 批量删除结果 **/
func (service *DefaultTemplateService) DeleteArray(array []int) (message error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if list := service.Items(v); len(list) <= 0 {
			if _, message = orm.NewOrm().Delete(&spider.Template{Id: v}); message != nil {
				_ = orm.NewOrm().Rollback()
				break
			}
		} else {
			message = errors.New("该分组还有模板存在；不能删除！")
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if message == nil {
		_ = orm.NewOrm().Commit()
	}
	return message
}

/** 获取所有分组 **/
func (service *DefaultTemplateService) Groups() []*spider.Template {
	var items []*spider.Template
	_, _ = orm.NewOrm().QueryTable(new(spider.Template)).All(&items)
	return items
}

/** 根据Id获取该分组下的所有模板 **/
func (service *DefaultTemplateService) Items(id int) []*spider.TemplateConfig {
	var result []*spider.TemplateConfig
	fs, _ := service.PathList(beego.AppConfig.String("spider_template"))
	for _, v := range fs {
		ini := v + "/" + "config.conf"
		if b, e := service.PathExists(ini); e == nil && b == true {
			if item := service.ReadConfig(ini); item.Class > 0 && item.Class == id {
				result = append(result, item)
			}
		}
	}
	return result
}

/** 读取一个模板配置 **/
func (service *DefaultTemplateService) ReadConfig(path string) *spider.TemplateConfig {
	result := &spider.TemplateConfig{}
	if con, err := config.NewConfig("ini", path); err == nil {
		result.Title = con.String("title")
		result.Remark = con.String("remark")
		result.Cover = con.String("cover")
		result.Name = con.String("name")
		result.Class, _ = con.Int("class")
	}
	if result.Cover != "" {
		result.Cover = service.ImageToBase64(result.Cover)
	}
	return result
}

/** 将图片转为base64显示 **/
func (service *DefaultTemplateService) ImageToBase64(path string) string {
	ff, _ := os.Open(path)
	defer ff.Close()
	buffer := make([]byte, 500000)
	n, _ := ff.Read(buffer)
	return base64.StdEncoding.EncodeToString(buffer[:n])
}

/** 文件夹提取 **/
func (service *DefaultTemplateService) PathList(dir string) ([]string, error) {
	var list []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			list = append(list, path)
			return nil
		}
		return nil
	})
	return list, err
}

/** 检测文件是否存在 **/
func (service *DefaultTemplateService) PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
