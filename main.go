package main

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/filter/cors"
	_ "github.com/douguohai/onlyoffice-golang/routers"
	"net/http"
)

func main() {
	//开启orm调试模式
	orm.Debug = true

	//创建附件目录
	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	if web.BConfig.RunMode == "dev" {
		web.BConfig.WebConfig.DirectoryIndex = true
		web.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	// 添加静态资源拦截器
	web.InsertFilter("/static", web.BeforeStatic, func(ctx *context.Context) {
		if ctx.Request.RequestURI == "/static" || ctx.Request.RequestURI == "/static/" {
			ctx.Abort(http.StatusUnauthorized, "权限不足")
		}
	})

	web.SetStaticPath("/static", "static")
	web.Run()
}
