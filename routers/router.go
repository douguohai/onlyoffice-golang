package routers

import (
	"github.com/astaxie/beego"
	"github.com/douguohai/onlyoffice-golang/controllers"
)

func init() {
	//回调函数
	beego.Router("/url-to-callback", &controllers.OnlyController{}, "*:DocCallback")
	//添加一个文档
	beego.Router("/office/attachment", &controllers.OnlyController{}, "put:AddOnlyAttachment")
	//访问一个文档
	beego.Router("/office/:id:string", &controllers.OnlyController{}, "get:OnlyOffice")
	//获取下载文件连接
	beego.Router("/download/office/:id", &controllers.OnlyController{}, "get:DownloadDoc")
	//强行覆盖文件内容
	beego.Router("/office/overwrite/:id", &controllers.OnlyController{}, "post:OverwriteDoc")
	// 设置路由
	beego.Router("/ws/**", &controllers.WebsocketController{}, "get:WsHandler")
}
