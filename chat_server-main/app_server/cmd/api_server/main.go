package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"app_server/http/docs"
	"app_server/http/file"
	"app_server/pkg/cbind"
	"app_server/pkg/cfg"
	"app_server/pkg/db"
	"app_server/pkg/jwt"
	"app_server/pkg/openaic"
	"app_server/pkg/ossc"
	"app_server/service/auth"
	"app_server/service/chat"
	"app_server/service/config"
	"app_server/service/ctx"
	"app_server/service/message"
	"app_server/service/profile"
	"app_server/service/translate"
	"app_server/service/user"

	"app_server/proto/chat/chatconnect"
	"app_server/proto/config/configconnect"
	"app_server/proto/message/messageconnect"
	"app_server/proto/profile/profileconnect"
	"app_server/proto/translate/translateconnect"
	"app_server/proto/user/userconnect"

	connect "connectrpc.com/connect"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

func main() {
	cfg.Init(*cfgFile)
	lo.Must0(db.Init(cfg.Viper().GetString("db.dsn"), cfg.Viper().GetBool("db.debug")))
	lo.Must0(ossc.Init(ossc.Cfg{
		PublicEndpoint:  cfg.Viper().GetString("aliyun.oss.public_endpoint"),
		Endpoint:        cfg.Viper().GetString("aliyun.oss.endpoint"),
		AccessKeyId:     cfg.Viper().GetString("aliyun.oss.access_key_id"),
		AccessKeySecret: cfg.Viper().GetString("aliyun.oss.access_key_secret"),
		UserFileBucket:  cfg.Viper().GetString("aliyun.oss.user_file_bucket"),
	}))
	openaic.Init(cfg.UnmarshalKey[openaic.Config]("ai.volces"))
	jwt.Init([]byte(cfg.Viper().GetString("jwt.secret")))
	log.Fatal(route().Run(lo.Ternary(*port != "", fmt.Sprintf(":%s", *port), cfg.Viper().GetString("server.address"))))
}

func route() *gin.Engine {
	root := gin.Default()

	// 添加CORS中间件
	root.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, connect-protocol-version, connect-timeout-ms")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	fileGroup := root.Group("/file")
	fileGroup.POST("/wx_upload", file.WxFileUpload)
	fileGroup.POST("/upload", file.WxFileUpload)

	docs.Register(root)

	binder := cbind.NewBinder(root.Group("/"))

	// 绑定业务服务
	binder.Bind(userconnect.NewUserServiceHandler(&user.UserService{}))
	binder.Bind(configconnect.NewConfigServiceHandler(&config.ConfigService{}))

	// 注册新的统一消息服务
	binder.Bind(messageconnect.NewChatMessageServiceHandler(&message.ChatMessageService{},
		connect.WithInterceptors(
			connect.UnaryInterceptorFunc(auth.AuthInterceptor),
			connect.UnaryInterceptorFunc(ctx.CtxInterceptor),
		),
	))
	binder.Bind(chatconnect.NewChatServiceHandler(&chat.ChatService{},
		connect.WithInterceptors(
			connect.UnaryInterceptorFunc(auth.AuthInterceptor),
			connect.UnaryInterceptorFunc(ctx.CtxInterceptor),
		),
	))
	binder.Bind(translateconnect.NewTranslateServiceHandler(&translate.TranslateService{},
		connect.WithInterceptors(
			connect.UnaryInterceptorFunc(auth.AuthInterceptor),
			connect.UnaryInterceptorFunc(ctx.CtxInterceptor),
		),
	))
	binder.Bind(profileconnect.NewProfileServiceHandler(&profile.ProfileService{},
		connect.WithInterceptors(
			connect.UnaryInterceptorFunc(auth.AuthInterceptor),
			connect.UnaryInterceptorFunc(ctx.CtxInterceptor),
		),
	))

	root.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	return root
}

var (
	cfgFile = flag.String("c", "config.yaml", "config file")
	port    = flag.String("p", "", "port")
)

func init() {
	flag.Parse()
}
