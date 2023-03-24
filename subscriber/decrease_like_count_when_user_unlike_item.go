package subscriber

import (
	"context"
	"g09-social-todo-list/common"
	"g09-social-todo-list/module/item/storage"
	"g09-social-todo-list/pubsub"
	goservice "github.com/200Lab-Education/go-sdk"
	"gorm.io/gorm"
)

func DecreaseLikeCountAfterUserUnlikeItem(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
		Title: "Decrease like count after user unlikes item",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

			//data := message.Data().(HasItemId)
			data := message.Data().(map[string]interface{})

			itemId := data["item_id"].(float64)

			if err := storage.NewSQLStore(db).DecreaseLikeCount(ctx, int(itemId)); err != nil {
				return err
			}

			_ = message.Ack()
			return nil
		},
	}
}
