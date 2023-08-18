package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
	"time"
)

// OnlyAttachment 附件
type OnlyAttachment struct {
	Id         int64
	FileName   string
	OriginName string
	Created    time.Time `orm:"auto_now_add;type(datetime)"`
	Updated    time.Time `orm:"auto_now;type(datetime)"`
}

// OnlyHistory 历史版本
type OnlyHistory struct {
	Id            int64
	AttachId      int64
	ServerVersion string
	FileUrl       string
	ChangesUrl    string    //`orm:"null"`
	HistoryKey    string    `orm:"sie(19)"`
	Expires       time.Time `orm:"type(datetime)"`
	Created       time.Time `orm:"type(datetime)"`
}

func init() {
	dbHost := beego.AppConfig.String("dbHost")
	dbName := beego.AppConfig.String("dbName")
	dbPort := beego.AppConfig.String("dbPort")
	dbUser := beego.AppConfig.String("dbUser")
	dbPassword := beego.AppConfig.String("dbPassword")

	orm.RegisterDataBase("default", "postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", dbUser, dbPassword, dbName, dbHost, dbPort))
	// set default database
	// register model
	orm.RegisterModel(new(OnlyAttachment), new(OnlyHistory))
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
	err := o.Begin()
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
		err = o.Rollback()
		logs.Error("更新附件异常")
		return 0, err
	}
	err = o.Commit()
	return id, err
}
