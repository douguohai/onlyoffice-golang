package main

import (
	context2 "context"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/filter/cors"
	"github.com/douguohai/onlyoffice-golang/base"
	_ "github.com/douguohai/onlyoffice-golang/routers"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
)

var ctx2 = context2.Background()

var rdb *redis.Client

var redisAuth bool

func init() {
	redisAuth, _ = web.AppConfig.Bool("redisAuth")
	if redisAuth {
		redisHost, _ := web.AppConfig.String("redisHost")
		redisPassword, _ := web.AppConfig.String("redisPassword")
		redisDB, _ := web.AppConfig.Int("redisDB")
		rdb = redis.NewClient(&redis.Options{
			Addr:     redisHost,
			Password: redisPassword, // no password set
			DB:       redisDB,       // use default DB
		})
		// 尝试连接 Redis
		_, err := rdb.Ping(ctx2).Result()
		if err != nil {
			// 连接失败，记录日志
			fmt.Println("redisAuth 为true ，开启redis 认证机制，redis 连接失败 ，请核查 redis 相关配置 ", err)
		}
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
	if redisAuth {
		var tokens = make([]string, 0)
		if ctx.Request.Method == "GET" && strings.HasPrefix(ctx.Request.RequestURI, "/office/") {
			token := ctx.Input.Query("Token")
			if "" != token {
				tokens = append(tokens, token)
			}
		} else {
			tokens = ctx.Request.Header["Token"]
		}
		if len(tokens) == 0 {
			_ = ctx.JSONResp(base.Result{
				Code:       401,
				ErrMessage: "权限不足",
			})
		}
		_, err := rdb.Get(ctx2, tokens[0]).Result()
		if err == redis.Nil {
			_ = ctx.JSONResp(base.Result{
				Code:       401,
				ErrMessage: "权限不足",
			})
		} else if err != nil {
			fmt.Println(fmt.Sprintf("系统内部错误 redis 相关 %s", err))
			_ = ctx.JSONResp(base.Result{
				Code:       500,
				ErrMessage: fmt.Sprintf("系统内部错误 redis 相关 %s", err),
			})
		}
	}
}
