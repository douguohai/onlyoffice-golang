package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/douguohai/onlyoffice-golang/base"
	"github.com/douguohai/onlyoffice-golang/models"
	"github.com/douguohai/onlyoffice-golang/utils"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"
)

type OnlyController struct {
	web.Controller
}

// CommonCallback 通用回调状态
type CommonCallback struct {
	Key    string `json:"key"`
	Status int    `json:"status"`
}

// EditHasPrepareSave 2- 文档已准备好保存，
type EditHasPrepareSave struct {
	Key        string `json:"key"`
	Status     int    `json:"status"`
	Url        string `json:"url"`
	ChangesUrl string `json:"changesurl"`
	History    struct {
		ServerVersion string `json:"serverVersion"`
		Changes       []struct {
			Created string `json:"created"`
			User    struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			} `json:"user"`
		} `json:"changes"`
	} `json:"history"`
	Users   []string `json:"users"`
	Actions []struct {
		Type   int    `json:"type"`
		Userid string `json:"userid"`
	} `json:"actions"`
	LastSave    time.Time `json:"lastsave"`
	NotModified bool      `json:"notmodified"`
	Filetype    string    `json:"filetype"`
}

// EditHasSaved 6 文档正在编辑，但当前文档状态已保存，
type EditHasSaved struct {
	Key        string `json:"key"`
	Status     int    `json:"status"`
	Url        string `json:"url"`
	ChangesUrl string `json:"changesurl"`
	History    struct {
		ServerVersion string `json:"serverVersion"`
		Changes       []struct {
			Created string `json:"created"`
			User    struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			} `json:"user"`
		} `json:"changes"`
	} `json:"history"`
	Users   []string `json:"users"`
	Actions []struct {
		Type   int    `json:"type"`
		Userid string `json:"userid"`
	} `json:"actions"`
	LastSave      time.Time `json:"lastsave"`
	ForceSaveType int       `json:"forcesavetype"`
	Filetype      string    `json:"filetype"`
}

// savePath 保存路径
const savePath = "./static/"

var serverUrl = ""
var wsServer = ""
var documentServer = ""

func init() {
	//取得客户端用户名
	err := os.MkdirAll(savePath, 0777) //..代表本当前exe文件目录的上级，.表示当前目录，没有.表示盘的根目录
	if err != nil {
		logs.Error("[文件存储目录初始化异常 %s]", err)
	}
	val, err := config.String("serverUrl")
	if err != nil {
		logs.Error("[获取文件预览服务外部访问地址 %s]", val)
	} else {
		serverUrl = val
		logs.Info("[获取文件预览服务外部访问地址 %s]", val)
	}
	val, err = config.String("wsServer")
	if err != nil {
		logs.Error("[ws 访问地址  %s]")
	} else {
		wsServer = val
		logs.Info("[ws 访问地址  %s]", val)
	}
	val, err = config.String("documentServer")
	if err != nil {
		logs.Error("[文档服务器访问地址 %s]")
	} else {
		documentServer = val
		logs.Info("[文档服务器访问地址 %s]", val)
	}
}

// OnlyOffice 协作页面的显示
// 补充权限判断
// 补充token
func (c *OnlyController) OnlyOffice() {
	//pid转成64为
	docId, err := c.GetInt64(":id")
	if err != nil {
		logs.Error(err)
		c.Data["json"] = base.BuildResult(-1, "文件id参数存在问题")
		c.ServeJSON()
		return
	}
	//根据附件id取得附件的详细信息
	attachment, err := models.GetOnlyAttachmentById(docId)
	if err != nil {
		logs.Error(err)
		c.Data["json"] = base.BuildResult(-1, "根据文档id未查询到当前文档")
		c.ServeJSON()
		return
	}

	c.Data["Username"] = "tianwen"

	c.Data["Mode"] = "edit"
	c.Data["Edit"] = true
	c.Data["Review"] = false

	c.Data["Doc"] = attachment
	c.Data["Key"] = strconv.FormatInt(attachment.Updated.UnixNano(), 10)
	c.Data["Doc.FileName"] = attachment.FileName
	c.Data["serverUrl"] = serverUrl
	c.Data["wsServer"] = wsServer
	c.Data["documentServer"] = documentServer

	if path.Ext(attachment.FileName) == ".docx" || path.Ext(attachment.FileName) == ".DOCX" {
		c.Data["fileType"] = "docx"
		c.Data["documentType"] = "text"
	} else if path.Ext(attachment.FileName) == ".XLSX" || path.Ext(attachment.FileName) == ".xlsx" {
		c.Data["fileType"] = "xlsx"
		c.Data["documentType"] = "spreadsheet"
	} else if path.Ext(attachment.FileName) == ".pptx" || path.Ext(attachment.FileName) == ".PPTX" {
		c.Data["fileType"] = "pptx"
		c.Data["documentType"] = "presentation"
	} else if path.Ext(attachment.FileName) == ".doc" || path.Ext(attachment.FileName) == ".DOC" {
		c.Data["fileType"] = "doc"
		c.Data["documentType"] = "text"
	} else if path.Ext(attachment.FileName) == ".txt" || path.Ext(attachment.FileName) == ".TXT" {
		c.Data["fileType"] = "txt"
		c.Data["documentType"] = "text"
	} else if path.Ext(attachment.FileName) == ".XLS" || path.Ext(attachment.FileName) == ".xls" {
		c.Data["fileType"] = "xls"
		c.Data["documentType"] = "spreadsheet"
	} else if path.Ext(attachment.FileName) == ".csv" || path.Ext(attachment.FileName) == ".CSV" {
		c.Data["fileType"] = "csv"
		c.Data["documentType"] = "spreadsheet"
	} else if path.Ext(attachment.FileName) == ".ppt" || path.Ext(attachment.FileName) == ".PPT" {
		c.Data["fileType"] = "ppt"
		c.Data["documentType"] = "presentation"
	} else if path.Ext(attachment.FileName) == ".pdf" || path.Ext(attachment.FileName) == ".PDF" {
		c.Data["fileType"] = "pdf"
		c.Data["documentType"] = "text"
		c.Data["Mode"] = "view"
	}

	u := c.Ctx.Input.UserAgent()
	matched, err := regexp.MatchString("AppleWebKit.*Mobile.*", u)
	if err != nil {
		logs.Error(err)
	}
	if matched == true {
		c.Data["Type"] = "mobile"
	} else {
		c.Data["Type"] = "desktop"
	}
	c.TplName = "onlyoffice/onlyoffice.tpl"
}

// getExpireTime 从文件url中提取过期时间
func getExpireTime(url string) time.Time {
	//获取连接到下载的预期过期事件
	paramsMap, err := utils.GetParams(url)
	if err != nil {
		logs.Error(err)
	}
	var expiresTime int64
	expires, ok := paramsMap["expires"]
	if ok {
		expiresTime, err = strconv.ParseInt(expires, 10, 64)
		if err != nil {
			expiresTime = time.Now().Unix()
		}
	}
	//时间戳转日期
	dataTimeStr := time.Unix(expiresTime, 0)
	return dataTimeStr
}

// AddOnlyAttachment 批量添加一对一模式
// 要避免同名覆盖的严重bug！！！！
func (c *OnlyController) AddOnlyAttachment() {
	_, h, err := c.GetFile("file")
	if err != nil {
		c.Data["json"] = base.BuildResult(-1, "获取上传文件失败")
		c.ServeJSON()
		return
	}
	if h != nil {
		fileName := strconv.FormatInt(time.Now().Unix(), 10) + path.Ext(h.Filename)
		if c.SaveToFile("file", savePath+fileName) != nil {
			c.Data["json"] = base.BuildResult(-1, "转储上传文件失败")
			c.ServeJSON()
			return
		}
		//保存附件
		id, err := models.AddOnlyAttachment(h.Filename, fileName)
		if err != nil {
			logs.Error(err)
			c.Data["json"] = base.BuildResult(-1, "更新文档信息失败")
		} else {
			c.Data["json"] = &base.AddOnlyAttachmentResult{
				AddOnlyAttachmentVo: base.AddOnlyAttachmentVo{
					FileUrl: fmt.Sprintf("%s/office/%d", serverUrl, id),
					FileId:  id,
				},
				Result: base.Result{
					Code:       0,
					ErrMessage: "操作成功",
				},
			}
		}
		c.ServeJSON()
	}
}

// DownloadDoc 协作页面下载的文档，采用绝对路径型式
func (c *OnlyController) DownloadDoc() {
	//pid转成64为
	docId, err := c.GetInt64(":id")
	if err != nil {
		c.Data["json"] = base.BuildResult(-1, "文件id参数存在问题")
		c.ServeJSON()
		return
	}
	//根据附件id取得附件的详细信息
	attachment, err := models.GetOnlyAttachmentById(docId)
	if err != nil {
		c.Data["json"] = base.BuildResult(-1, "根据id查询不到该文件")
		c.ServeJSON()
		return
	}
	var downloadDocVo = base.DownloadDocVo{
		Url:     fmt.Sprintf("%s/static/%s", serverUrl, attachment.FileName),
		Expires: attachment.Updated,
	}
	//根据docId获取附件信息
	history, err := models.GetLastHistoryByAttachId(attachment.Id)
	if err == nil {
		//判定下载连接是否过期，未过期返回服务端下载连接，已经过期返回本地备份文档
		if history.Expires.After(time.Now()) {
			downloadDocVo = base.DownloadDocVo{
				Url:     history.FileUrl,
				Expires: history.Expires,
			}
		}
	} else {
		if err != orm.ErrNoRows {
			c.Data["json"] = base.BuildResult(-1, "获取文档信息失败")
			c.ServeJSON()
			return
		}
	}
	c.Data["json"] = &base.DownloadDocResult{
		DownloadDocVo: downloadDocVo,
		Result: base.Result{
			Code:       0,
			ErrMessage: "操作成功",
		},
	}
	c.ServeJSON()
	return
}

// OverwriteDoc 协作页面覆盖的文档，采用绝对路径型式
func (c *OnlyController) OverwriteDoc() {
	//pid转成64为
	docId, err := c.GetInt64(":id")
	if err != nil {
		logs.Error(err)
		c.Data["json"] = base.BuildResult(-1, "文件id参数存在问题")
		c.ServeJSON()
		return
	}
	_, h, err := c.GetFile("file")
	if err != nil {
		c.Data["json"] = base.BuildResult(-1, "获取上传文件失败")
		c.ServeJSON()
		return
	}
	if h != nil {
		fileName := strconv.FormatInt(time.Now().Unix(), 10) + path.Ext(h.Filename)
		if c.SaveToFile("file", savePath+fileName) != nil {
			c.Data["json"] = base.BuildResult(-1, "转储上传文件失败")
			c.ServeJSON()
			return
		}
		//根据id修改附件名称
		//保存附件
		_, err := models.UpdateFileUrlAndUpdateTime(docId, fileName)
		if err != nil {
			logs.Error(err)
			c.Data["json"] = base.BuildResult(-1, "更新文档信息失败")
		} else {
			c.Data["json"] = base.BuildResult(0, "更新文档信息成功")
		}
		c.ServeJSON()
	}
	msg := base.Message{
		Type: 0,
		Data: docId,
	}
	str, err := json.Marshal(msg)
	if err == nil {
		hub.broadcast <- str
	}

}

// DocCallback 协作页面的保存和回调
// 关闭浏览器标签后获取最新文档保存到文件夹
func (c *OnlyController) DocCallback() {
	//pid转成64为
	docId, err := c.GetInt64("id")
	if err != nil {
		logs.Error(err)
	}
	//根据附件id取得附件的详细信息
	attachment, err := models.GetOnlyAttachmentById(docId)
	if err != nil {
		logs.Error("[通用-回调函数 解析回调异常] 查询不到该文档 %v", docId)
		c.Data["json"] = map[string]interface{}{"error": 0}
		c.ServeJSON()
		return
	}
	var callback CommonCallback
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &callback)
	if err != nil {
		logs.Error("[通用-回调函数 解析回调异常]： %v", err)
	} else {
		logs.Info("[通用-回调函数 解析回调] :%v", callback)
	}
	if callback.Status == 1 || callback.Status == 4 {
		c.Data["json"] = map[string]interface{}{"error": 0}
		c.ServeJSON()
		return
	} else if callback.Status == 6 {
		var editHasSaved EditHasSaved
		_ = json.Unmarshal(c.Ctx.Input.RequestBody, &editHasSaved)
		//下载文件到本地
		err := utils.DownloadFile(editHasSaved.Url, savePath+attachment.FileName)
		if err != nil {
			logs.Error(err)
		} else {
			dataTimeStr := getExpireTime(editHasSaved.Url)
			//更新附件的时间
			err = models.UpdateOnlyAttachmentTime(docId)
			if err != nil {
				logs.Error(err)
			}
			//写入历史版本
			_, err = models.AddOnlyHistory(attachment.Id, editHasSaved.History.ServerVersion, callback.Key, editHasSaved.Url, editHasSaved.ChangesUrl, dataTimeStr, editHasSaved.LastSave)
			if err != nil {
				logs.Error(err)
			}
		}
		c.Data["json"] = map[string]interface{}{"error": 0}
		c.ServeJSON()
	} else if callback.Status == 2 {
		var editHasPrepareSave EditHasPrepareSave
		err = json.Unmarshal(c.Ctx.Input.RequestBody, &editHasPrepareSave)
		if err != nil {
			logs.Error("[状态2-回调函数 解析回调异常]： %v", err)
		}
		//下载文件到本地
		err := utils.DownloadFile(editHasPrepareSave.Url, savePath+attachment.FileName)
		if err != nil {
			logs.Error(err)
		} else {
			dataTimeStr := getExpireTime(editHasPrepareSave.Url)
			//更新附件的时间
			err = models.UpdateOnlyAttachmentTime(docId)
			if err != nil {
				logs.Error(err)
			}
			//写入历史版本
			_, err = models.AddOnlyHistory(attachment.Id, editHasPrepareSave.History.ServerVersion, callback.Key, editHasPrepareSave.Url, editHasPrepareSave.ChangesUrl, dataTimeStr, editHasPrepareSave.LastSave)
			if err != nil {
				logs.Error(err)
			}
		}
		c.Data["json"] = map[string]interface{}{"error": 0}
		c.ServeJSON()
		return
	} else if callback.Status == 3 {
		err = models.UpdateOnlyAttachmentTime(docId)
		if err != nil {
			logs.Error(err)
		}
		c.Data["json"] = map[string]interface{}{"error": 0}
		c.ServeJSON()
	} else {
		c.Data["json"] = map[string]interface{}{"error": 0}
		c.ServeJSON()
	}
}
