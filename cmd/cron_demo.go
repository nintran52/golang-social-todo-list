package cmd

import (
	"g09-social-todo-list/common"
	"g09-social-todo-list/plugin/sdkgorm"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"log"
)

// Ví dụ cho cronjob update lại số like count trên toàn bộ table todo_items

// UPDATE todo_items ti INNER JOIN (
// SELECT item_id, COUNT(item_id) as `count` FROM `user_like_items`
// GROUP BY item_id
// ) c ON c.item_id = ti.id SET ti.liked_count = c.count

var cronDemoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run demo cron job",
	Run: func(cmd *cobra.Command, args []string) {

		service := goservice.New(
			goservice.WithName("social-todo-list"),
			goservice.WithVersion("1.0.0"),
			goservice.WithInitRunnable(sdkgorm.NewGormDB("main.mysql", common.PluginDBMain)),
		)

		if err := service.Init(); err != nil {
			log.Fatalln(err)
		}

		db := service.MustGet(common.PluginDBMain).(*gorm.DB)

		log.Println("I am demo cron with DB connection:", db)

	},
}
