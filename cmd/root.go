package cmd

import (
	"fmt"
	"g09-social-todo-list/common"
	"g09-social-todo-list/demogrpc/demo"
	"g09-social-todo-list/memcache"
	"g09-social-todo-list/middleware"
	ginitem "g09-social-todo-list/module/item/transport/gin"
	"g09-social-todo-list/module/upload"
	userstorage "g09-social-todo-list/module/user/storage"
	ginuser "g09-social-todo-list/module/user/transport/gin"
	"g09-social-todo-list/module/userlikeitem/storage"
	ginuserlikeitem "g09-social-todo-list/module/userlikeitem/transport/gin"
	"g09-social-todo-list/module/userlikeitem/transport/rpc"
	"g09-social-todo-list/plugin/appredis"
	"g09-social-todo-list/plugin/nats"
	"g09-social-todo-list/plugin/rpccaller"
	"g09-social-todo-list/plugin/sdkgorm"
	"g09-social-todo-list/plugin/simple"
	"g09-social-todo-list/plugin/tokenprovider/jwt"
	"g09-social-todo-list/subscriber"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
	"os"
)

func newService() goservice.Service {
	service := goservice.New(
		goservice.WithName("social-todo-list"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(sdkgorm.NewGormDB("main.mysql", common.PluginDBMain)),
		goservice.WithInitRunnable(jwt.NewJWTProvider(common.PluginJWT)),
		//goservice.WithInitRunnable(pubsub.NewPubSub(common.PluginPubSub)),
		goservice.WithInitRunnable(nats.NewNATS(common.PluginPubSub)),
		goservice.WithInitRunnable(rpccaller.NewApiItemCaller(common.PluginItemAPI)),
		goservice.WithInitRunnable(appredis.NewRedisDB("redis", common.PluginRedis)),
		goservice.WithInitRunnable(simple.NewSimplePlugin("simple")),
	)

	return service
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start social TODO service",
	Run: func(cmd *cobra.Command, args []string) {
		service := newService()

		serviceLogger := service.Logger("service")

		if err := service.Init(); err != nil {
			serviceLogger.Fatalln(err)
		}

		// Set up gRPC
		address := "0.0.0.0:50051"
		lis, err := net.Listen("tcp", address)

		if err != nil {
			log.Fatalf("Error %v", err)
		}
		fmt.Printf("Server is listening on %v ...", address)

		s := grpc.NewServer()

		db := service.MustGet(common.PluginDBMain).(*gorm.DB)

		store := storage.NewSQLStore(db)
		demo.RegisterItemLikeServiceServer(s, rpc.NewRPCService(store))

		go func() {
			if err := s.Serve(lis); err != nil {
				log.Fatalln(err)
			}
		}()

		opts := grpc.WithTransportCredentials(insecure.NewCredentials())

		cc, err := grpc.Dial("localhost:50051", opts)
		if err != nil {
			log.Fatal(err)
		}

		client := demo.NewItemLikeServiceClient(cc)
		//////

		service.HTTPServer().AddHandler(func(engine *gin.Engine) {
			engine.Use(middleware.Recover())

			// Example for Simple Plugin
			type CanGetValue interface {
				GetValue() string
			}
			log.Println(service.MustGet("simple").(CanGetValue).GetValue())
			/////////

			db := service.MustGet(common.PluginDBMain).(*gorm.DB)

			authStore := userstorage.NewSQLStore(db)
			authCache := memcache.NewUserCaching(memcache.NewRedisCache(service), authStore)

			middlewareAuth := middleware.RequiredAuth(authCache, service)

			v1 := engine.Group("/v1")
			{
				v1.PUT("/upload", upload.Upload(service))

				v1.POST("/register", ginuser.Register(service))
				v1.POST("/login", ginuser.Login(service))
				v1.GET("/profile", middlewareAuth, ginuser.Profile())

				items := v1.Group("/items", middlewareAuth)
				{
					items.POST("", ginitem.CreateItem(service))
					items.GET("", ginitem.ListItem(service, client))
					items.GET("/:id", ginitem.GetItem(service))
					items.PATCH("/:id", ginitem.UpdateItem(service))
					items.DELETE("/:id", ginitem.DeleteItem(service))

					items.POST("/:id/like", ginuserlikeitem.LikeItem(service))
					items.DELETE("/:id/unlike", ginuserlikeitem.UnlikeItem(service))
					items.GET("/:id/liked-users", ginuserlikeitem.ListUserLiked(service))
				}

				rpc := v1.Group("rpc")
				{
					rpc.POST("/get_item_likes", ginuserlikeitem.GetItemLikes(service))
				}
			}

			engine.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})
		})

		je, err := jaeger.NewExporter(jaeger.Options{
			AgentEndpoint: "localhost:6831",
			Process:       jaeger.Process{ServiceName: "Todo-List-Service"},
		})

		if err != nil {
			log.Println(err)
		}

		trace.RegisterExporter(je)
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(0.001)})

		_ = subscriber.NewEngine(service).Start()

		if err := service.Start(); err != nil {
			serviceLogger.Fatalln(err)
		}
	},
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)
	rootCmd.AddCommand(cronDemoCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
