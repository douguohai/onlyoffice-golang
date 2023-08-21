// Package routers @APIVersion 1.0.0
// @Title mobile API
// @Description mobile has every tool to get any job done, so codename for the new mobile APIs.
// @Contact astaxie@gmail.com
package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/douguohai/onlyoffice-golang/controllers"
)

func init() {
	//回调函数
	web.Router("/url-to-callback", &controllers.OnlyController{}, "*:DocCallback")
	//添加一个文档
	web.Router("/office/attachment", &controllers.OnlyController{}, "put:AddOnlyAttachment")
	//访问一个文档
	web.Router("/office/:id:string", &controllers.OnlyController{}, "get:OnlyOffice")
	//获取下载文件连接
	web.Router("/download/office/:id", &controllers.OnlyController{}, "get:DownloadDoc")
	//强行覆盖文件内容
	web.Router("/office/overwrite/:id", &controllers.OnlyController{}, "post:OverwriteDoc")
	// 设置路由
	web.Router("/ws/**", &controllers.WebsocketController{}, "get:WsHandler")
}
