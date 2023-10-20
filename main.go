package main

import (
	context2 "context"
	"fmt"
	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/filter/cors"
	"github.com/douguohai/onlyoffice-golang/base"
	"github.com/douguohai/onlyoffice-golang/models"
	_ "github.com/douguohai/onlyoffice-golang/routers"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ctx2 context2.Context

var auth bool
var localCache cache.Cache

func init() {
	auth, _ = web.AppConfig.Bool("auth")
	if auth {
		fmt.Println("[important] 检测到已经开启认证机制，请在数据库中填写相关认证信息")
	}
	var err error
	localCache, err = cache.NewCache("memory", `{"interval":30}`)
	if err != nil {
		panic("内存换存常见失败:" + err.Error())
	}

}

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
	web.InsertFilter("/office/**", web.BeforeExec, func(ctx *context.Context) {
		authCheck(ctx)
	})

	// 添加静态资源拦截器
	web.InsertFilter("/download/**", web.BeforeExec, func(ctx *context.Context) {
		authCheck(ctx)
	})

	// 添加静态资源拦截器
	web.InsertFilter("/static", web.BeforeStatic, func(ctx *context.Context) {
		if ctx.Request.RequestURI == "/static" || ctx.Request.RequestURI == "/static/" {
			ctx.Abort(http.StatusUnauthorized, "权限不足")
		}
	})

	web.SetStaticPath("/static", "static")
	web.Run()
}

// authCheck 接口权限核验，防止未授权访问
func authCheck(ctx *context.Context) {
	if auth {
		var tokens = make([]string, 0)
		var appIds = make([]string, 0)
		if ctx.Request.Method == "GET" && strings.HasPrefix(ctx.Request.RequestURI, "/office/") {
			token := ctx.Input.Query("Token")
			if "" != token {
				tokens = append(tokens, token)
			}
			appId := ctx.Input.Query("Appid")
			if "" != appId {
				appIds = append(appIds, appId)
			}
		} else {
			tokens = ctx.Request.Header["Token"]
			appIds = ctx.Request.Header["Appid"]
		}
		if len(tokens) == 0 || len(appIds) == 0 {
			_ = ctx.JSONResp(base.Result{
				Code:       401,
				ErrMessage: "权限不足",
			})
			return
		}
		token := tokens[0]
		appId := appIds[0]

		cacheKey := fmt.Sprintf("%s-%s", appId, token)

		checkResult := base.CheckResult{}

		val, err := localCache.Get(ctx2, cacheKey)

		if err != nil {
			app, err := models.GetAppById(appId)
			if err != nil {
				fmt.Println(err.Error())
				_ = ctx.JSONResp(base.Result{
					Code:       401,
					ErrMessage: "权限不足",
				})
				return
			}

			checkUrl, err := url.ParseRequestURI(app.CheckUrl + token)
			if err != nil {
				_ = ctx.JSONResp(base.Result{
					Code:       401,
					ErrMessage: "权限不足-非法核查api",
				})
				return
			}
			// Create a Resty Client
			client := resty.New()
			client.SetTimeout(5 * time.Second).SetRetryCount(3)

			_, err = client.R().
				EnableTrace().SetResult(&checkResult).
				Get(checkUrl.String())

			if err != nil {
				_ = ctx.JSONResp(base.Result{
					Code:       401,
					ErrMessage: "权限不足-核查api无法访问，请联系运维人员",
				})
				return
			}
			localCache.Put(ctx2, cacheKey, checkResult, 30*time.Second)
		} else {
			checkResult = val.(base.CheckResult)
		}

		fmt.Println(fmt.Sprintf("核验返回: %v", checkResult))

		if checkResult.Code != 200 {
			_ = ctx.JSONResp(base.Result{
				Code:       401,
				ErrMessage: "权限不足，非法访问",
			})
			return
		}
	}
}
