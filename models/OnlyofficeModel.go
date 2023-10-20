package models

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
	"time"
)

// OnlyAttachment 附件
type OnlyAttachment struct {
	Id         int64
	FileName   string    `orm:"description(文件名称)"`
	OriginName string    `orm:"description(文件原始名称)"`
	Created    time.Time `orm:"auto_now_add;type(datetime);description(创建时间)"`
	Updated    time.Time `orm:"auto_now;type(datetime);description(修改时间)"`
}

// OnlyHistory 历史版本
type OnlyHistory struct {
	Id            int64
	AttachId      int64     `orm:"description(关联附件id)"`
	ServerVersion string    `orm:"description(文档服务器服务版本)"`
	FileUrl       string    `orm:"description(文件下载连接)"`
	ChangesUrl    string    `orm:"description(文件变化下载连接)"`
	HistoryKey    string    `orm:"sie(19);description(文档服务器key)"`
	Expires       time.Time `orm:"type(datetime);description(过期时间)"`
	Created       time.Time `orm:"type(datetime);description(创建时间)"`
}

// App 应用信息
type App struct {
	Id       int64
	AppId    string    `orm:"unique;description(应用名称)"`
	CheckUrl string    `orm:"description(核查url)"`
	Status   int       `orm:"auto_now_add;default:0;description(状态删除 0有效 1无效)"`
	UpdateAt time.Time `orm:"auto_now_add;type(datetime); description(修改时间)"`
	CreateAt time.Time `orm:"auto_now_add;type(datetime);description(创建时间)"`
}

func (flow *App) TableName() string {
	return "app"
}

func init() {

	dbHost, _ := web.AppConfig.String("dbHost")
	dbName, _ := web.AppConfig.String("dbName")
	dbPort, _ := web.AppConfig.String("dbPort")
	dbUser, _ := web.AppConfig.String("dbUser")
	dbPassword, _ := web.AppConfig.String("dbPassword")

	orm.RegisterDataBase("default", "postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", dbUser, dbPassword, dbName, dbHost, dbPort))
	// set default database
	// register model
	orm.RegisterModel(new(OnlyAttachment), new(OnlyHistory), new(App))
	// create table
	orm.RunSyncdb("default", false, true)
}

// AddOnlyAttachment 添加附件到成果id下
// 如果附件名称已经存在，则不再追加写入数据库
// 应该用ReadOrCreate尝试从数据库读取，不存在的话就创建一个
func AddOnlyAttachment(originName, fileName string) (id int64, err error) {
	o := orm.NewOrm()
	attachment := &OnlyAttachment{
		OriginName: originName,
		FileName:   fileName,
		Created:    time.Now(),
		Updated:    time.Now(),
	}
	return o.Insert(attachment)
}

// GetOnlyAttachmentById 根据附件id查询附件
func GetOnlyAttachmentById(Id int64) (attach OnlyAttachment, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable("OnlyAttachment")
	err = qs.Filter("id", Id).One(&attach)
	return attach, err
}

// UpdateOnlyAttachmentTime 修改附件的日期和changesurl修改记录地址
func UpdateOnlyAttachmentTime(cid int64) (err error) {
	o := orm.NewOrm()
	attachment := &OnlyAttachment{Id: cid}
	if o.Read(attachment) == nil {
		attachment.Updated = time.Now()
		_, err = o.Update(attachment, "Updated")
		if err != nil {
			return err
		}
	}
	return err
}

// AddOnlyHistory 添加历史版本
func AddOnlyHistory(docId int64, serverVersion string, key, fileUrl, changesUrl string, expires, created time.Time) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(&OnlyHistory{
		AttachId:      docId,
		ServerVersion: serverVersion,
		HistoryKey:    key,
		FileUrl:       fileUrl,
		ChangesUrl:    changesUrl,
		Expires:       expires.Local(),
		Created:       created,
	})
	return id, err
}

// GetLastHistoryByAttachId 根据附件id获取最后一次历史版本
func GetLastHistoryByAttachId(attachId int64) (history OnlyHistory, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable("OnlyHistory")
	err = qs.Filter("attach_id", attachId).OrderBy("-id").One(&history)
	return history, err
}

// UpdateFileUrlAndUpdateTime 更新文件地址和更新时间
func UpdateFileUrlAndUpdateTime(attachId int64, fileName string) (int64, error) {
	o := orm.NewOrm()
	tx, err := o.Begin()
	if err != nil {
		logs.Error("start the transaction failed")
		return 0, err
	}

	_, err = o.Delete(&OnlyHistory{AttachId: attachId}, "AttachId")
	if err != nil {
		logs.Error("清除历史数据异常")
		return 0, err
	}
	id, err := o.Update(&OnlyAttachment{Id: attachId, FileName: fileName, Updated: time.Now()}, "FileName", "Updated")
	if err != nil {
		err = tx.Rollback()
		logs.Error("更新附件异常")
		return 0, err
	}
	err = tx.Commit()
	return id, err
}

// GetAppById 根据appid查询相关服务
func GetAppById(appId string) (App, error) {
	app := App{}
	o := orm.NewOrm()
	err := o.QueryTable("app").Filter("app_id", appId).Filter("status", 0).One(&app)
	if err != nil {
		return app, errors.New("非法appId")
	} else {
		return app, nil
	}
}
