package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/attachment"
	"mime"
	"path"
	"strconv"
	"unsafe"
)

type Attachment struct {
	Main
	imageName string
	filesName string
}

/** 构造函数  **/
func (c *Attachment) NextPrepare() {
	c.imageName = "image"
	c.filesName = "files"
}

// @router /attachment/index [get,post]
func (c *Attachment) Index() {
	if c.IsAjax() {
		length, _ := c.GetInt("length", _func.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.adjunctService.PageListItems(length, draw, c.PageNumber(), search)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.adjunctService.TableColumnsType(), c.adjunctService.TableButtonsType())
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.adjunctService.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.adjunctService.DataTableButtons() {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, 10)
}

// @router /attachment/space [post]
func (c *Attachment) Space() {
	var maps []orm.Params
	var prefix = beego.AppConfig.String("db_prefix")
	_, _ = orm.NewOrm().Raw("SELECT SUM(size) AS size_count FROM " + prefix + "admin_attachment").Values(&maps)
	size, _ := strconv.ParseInt(maps[0]["size_count"].(string), 10, 64)
	size = size / 1024
	unit := "KB"
	if (size / 1024) > 0 {
		size = size / 1024
		unit = "MB"
		if (size / 1024) > 0 {
			size = size / 1024
			unit = "GB"
		}
	}
	c.Succeed(&controllers.ResultJson{
		Data:    size,
		Message: "目前在所有驱动共使用约" + strconv.FormatInt(size, 10) + unit + "空间",
	})
}

// @router /attachment/edit [get,post]
func (c *Attachment) Edit() {
	id := c.GetMustInt(":id", "非法访问；已被拒绝！")
	c.Data["info"] = c.adjunctService.FindId(id)
}

// @router /attachment/delete [get,post]
func (c *Attachment) Delete() {
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.adjunctService.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除附件成功！马上返回中。。。",
			Url:     beego.URLFor("Attachment.Index"),
		})
	}
}

/**************************************以下为外部方法**********************/

// @router /attachment/upload_image [post]
func (c *Attachment) UploadImage() {
	uploadDriver := beego.AppConfig.String("upload_driver")
	uploadMimeType := beego.AppConfig.String("upload_image_mime")
	if uploadDriver == "local" {
		c.Succeed(&controllers.ResultJson{Data: c.localUpload(c.imageName, uploadMimeType)})
	} else {
		c.Succeed(&controllers.ResultJson{Data: c.cloudUpload(c.imageName, uploadMimeType)})
	}
}

// todo 云端上传驱动未完成
func (c *Attachment) cloudUpload(fileKey, uploadMimeType string) attachment.Attachment {
	return attachment.Attachment{}
}

/** 本地驱动上传 **/
func (c *Attachment) localUpload(fileKey, uploadMimeType string) attachment.Attachment {
	file, header, err := c.GetFile(fileKey)
	if err != nil {
		c.Fail(&controllers.ResultJson{Message: "未接收到任何文件对象"})
	}
	uploadMaxSize, _ := beego.AppConfig.Int64("upload_max_size")
	if !(uploadMaxSize > (header.Size / 1024)) {
		c.Fail(&controllers.ResultJson{Message: "文件大小超过限度值～"})
	}
	extend := path.Ext(header.Filename) // 检测文件类型是否合法
	if isOk, err := c.adjunctService.IsAllowExt(extend, uploadMimeType); isOk == false {
		c.Fail(&controllers.ResultJson{Message: error.Error(err)})
	}
	fileSha1 := c.adjunctService.Md5(file) // 检测文件快传
	if one := c.adjunctService.FindMd5(fileSha1); one.Id > 0 {
		return one
	}
	pathString, err := c.adjunctService.DateFolder(beego.AppConfig.String("upload_folder") + "/image/")
	if err != nil {
		c.Fail(&controllers.ResultJson{Message: "创建文件夹失败～"})
	}
	filePath := pathString + "/" + fileSha1 + extend
	if c.SaveToFile(c.imageName, filePath) != nil {
		c.Fail(&controllers.ResultJson{Message: "创建文件失败～"})
	}
	id, err := c.adjunctService.CreateFileLogs(&attachment.Attachment{
		Sha1: fileSha1,
		Size: header.Size,
		Name: header.Filename,
		Mime: mime.TypeByExtension(extend),
		Ext:  extend,
		Path: filePath,
	})
	if err != nil {
		c.Fail(&controllers.ResultJson{Message: "创建文件记录失败～"})
	}
	return c.adjunctService.FindId(*(*int)(unsafe.Pointer(&id)))
}
