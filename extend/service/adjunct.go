package service

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_func "github.com/beatrice950201/araneid/extend/func"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/attachment"
	"github.com/go-playground/validator"
	"github.com/upyun/go-sdk/upyun"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type DefaultAdjunctService struct{}

/** 获取一个文件详情[根据ID] **/
func (service *DefaultAdjunctService) FindId(id int) (item attachment.Attachment) {
	item.Id = id
	_ = orm.NewOrm().Read(&item)
	item.Path = _func.DomainStatic(item.Driver) + item.Path
	return item
}

/** 云盘全部文件数量 **/
func (service *DefaultAdjunctService) aliveNum() int64 {
	index, _ := orm.NewOrm().QueryTable(new(attachment.Attachment)).Count()
	return index
}

/**  下载远程附件下载到第三方 **/
func (service *DefaultAdjunctService) DownloadFileCloud(url string, cloud map[string]string) string {
	up := upyun.NewUpYun(&upyun.UpYunConfig{Bucket: cloud["bucket"], Operator: cloud["name"], Password: cloud["password"]})
	local := service.DownloadFileLocal(url)
	if local != url {
		localName := path.Base(local)
		if err := up.Put(&upyun.PutObjectConfig{Path: "/" + localName, LocalPath: local}); err != nil {
			logs.Warn("远程资源拉取上传到云端失败！失败原因:%s", error.Error(err))
		} else {
			url = localName
			_ = os.Remove(local)
		}
	}
	return url
}

/**  下载远程附件下载到本地**/
func (service *DefaultAdjunctService) DownloadFileLocal(url string) string {
	if resp, err := http.Get(url); err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		filename := fmt.Sprintf(`%x`, md5.Sum(body))
		pathString, err := service.DateFolder(beego.AppConfig.String("upload_folder") + "/collect/")
		if err == nil {
			name := pathString + "/" + filename + path.Ext(path.Base(url))
			out, _ := os.Create(name)
			if _, err := io.Copy(out, bytes.NewReader(body)); err == nil {
				url = name
			} else {
				logs.Warn("远程资源拉取在io.Copy()处失败！失败原因:%s", error.Error(err))
			}
		} else {
			logs.Warn("远程资源拉取创建文件夹失败！失败原因:%s", error.Error(err))
		}
	} else {
		logs.Warn("检测远程资源拉取失败！失败原因:%s", error.Error(err))
	}
	return url
}

/** 图片格式检测 **/
func (service *DefaultAdjunctService) IsAllowExt(ext, extString string) (bool, error) {
	extArray := strings.Split(extString, "|")
	AllowExtMap := make(map[string]bool)
	for _, v := range extArray {
		AllowExtMap[v] = true
	}
	if _, ok := AllowExtMap[ext]; !ok {
		return false, errors.New("不允许上传该类型的文件！")
	} else {
		return true, nil
	}
}

/** 获取文件md5 **/
func (service *DefaultAdjunctService) Md5(file multipart.File) string {
	md5Object := md5.New()
	_, _ = io.Copy(md5Object, file)
	defer file.Close()
	return hex.EncodeToString(md5Object.Sum(nil))
}

/** 根据IO Reader 获取MD5 **/
func (service *DefaultAdjunctService) ReaderMd5(r io.Reader) string {
	body, _ := ioutil.ReadAll(r)
	return fmt.Sprintf(`%x`, md5.Sum(body))
}

/** 获取一个文件详情[根据md5] **/
func (service *DefaultAdjunctService) FindMd5(md5 string) attachment.Attachment {
	var item attachment.Attachment
	_ = orm.NewOrm().QueryTable(new(attachment.Attachment)).Filter("sha1", md5).One(&item)
	return item
}

/** 创建一个文件记录 **/
func (service *DefaultAdjunctService) CreateFileLogs(item *attachment.Attachment) (int64, error) {
	item.Driver = beego.AppConfig.String("upload_driver")
	v := DefaultBaseVerify{}
	if message := v.Begin().Struct(item); message != nil {
		return 0, errors.New(v.Translate(message.(validator.ValidationErrors)))
	}
	return orm.NewOrm().Insert(item)
}

/** 获取完整路径 **/
func (service *DefaultAdjunctService) DateFolder(path string) (string, error) {
	folderName := time.Now().Format("20060102")
	folderPath := filepath.Join(path, folderName)
	err := service.createFolder(folderPath)
	return folderPath, err
}

/** 计数器++ **/
func (service *DefaultAdjunctService) Inc(new, old int) (errorMsg error) {
	if old != 0 {
		errorMsg = service.Dec(old)
	}
	one := service.FindId(new)
	one.Usage += 1
	_, errorMsg = orm.NewOrm().Update(&one, "Usage")
	return errorMsg
}

/** 计数器-- **/
func (service *DefaultAdjunctService) Dec(new int) error {
	one := service.FindId(new)
	if one.Usage > 0 {
		one.Usage -= 1
	}
	_, errorMessage := orm.NewOrm().Update(&one, "Usage")
	return errorMessage
}

/** 批量删除 **/
func (service *DefaultAdjunctService) DeleteArray(array []int) (e error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		one := service.FindId(v)
		if one.Usage > 0 {
			e = errors.New("不允许删除使用中的附件！")
			_ = orm.NewOrm().Rollback()
			break
		} else {
			if e = service.deleteFile(one.Id); e != nil {
				_ = orm.NewOrm().Rollback()
				e = errors.New("删除原始附件失败！")
				break
			} else {
				if _, e = orm.NewOrm().Delete(&attachment.Attachment{Id: v}); e != nil {
					_ = orm.NewOrm().Rollback()
					break
				}
			}
		}
	}
	if e == nil {
		_ = orm.NewOrm().Commit()
	}
	return e
}

/************************************************表格渲染机制 ************************************************************/

/** 获取需要渲染的Column **/
func (service *DefaultAdjunctService) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "文件标识", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "文件名称", "name": "name", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "文件后缀", "name": "ext", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "文件大小(KB)", "name": "size", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "挂载次数", "name": "usage", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "驱动类型", "name": "driver", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "更新时间", "name": "update_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 返回表单结构字段如何解析 **/
func (service *DefaultAdjunctService) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "fileSize", "string", "fileDriver", "date"},
		"fieldName": {"id", "name", "ext", "size", "usage", "driver", "update_time"},
		"action":    {"", "", "", "", "", "", ""},
	}
	return result
}

/** 获取需要渲染的按钮组 **/
func (service *DefaultAdjunctService) DataTableButtons() []*table.TableButtons {
	var array []*table.TableButtons
	array = append(array, &table.TableButtons{
		Text:      "占用统计",
		ClassName: "btn btn-sm btn-alt-primary mt-1 space_count",
		Attribute: map[string]string{"data-action": beego.URLFor("Attachment.Space")},
	})
	array = append(array, &table.TableButtons{
		Text:      "删除选中",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Attachment.Delete")},
	})
	return array
}

/** 返回右侧按钮数据结构 **/
func (service *DefaultAdjunctService) TableButtonsType() []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "预览详情",
			ClassName: "btn btn-sm btn-alt-primary open_iframe",
			Attribute: map[string]string{
				"data-area": "580px,390px",
				"href":      beego.URLFor("Attachment.Edit", ":id", "__ID__", ":popup", 1),
			},
		},
	}
	return buttons
}

/** 处理分页 **/
func (service *DefaultAdjunctService) PageListItems(length, draw, page int, search string) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(attachment.Attachment))
	recordsTotal, _ := qs.Count()
	if search != "" {
		qs = qs.Filter("name__icontains", search)
	}
	_, _ = qs.Limit(length, length*(page-1)).ValuesList(&lists, "id", "name", "ext", "size", "usage", "driver", "update_time")
	data := map[string]interface{}{
		"draw":            draw,         // 请求次数
		"recordsFiltered": recordsTotal, // 从多少条里面筛选
		"recordsTotal":    recordsTotal, // 总条数
		"data":            lists,        // 筛选结果
	}
	return data
}

/******************以下方法为私有方法***********************/

/**  创建检测文件夹  **/
func (service *DefaultAdjunctService) isExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return os.IsExist(err)
	}
	return true
}

/**  创建检测生成文件夹  **/
func (service *DefaultAdjunctService) createFolder(path string) error {
	if !service.isExist(path) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

/** 删除文件 **/
func (service *DefaultAdjunctService) deleteFile(id int) error {
	var file attachment.Attachment
	file.Id = id
	_ = orm.NewOrm().Read(&file)
	if file.Driver == "local" {
		return os.Remove(file.Path)
	} else {
		//todo 第三方删除文件流程
		return nil
	}
}
