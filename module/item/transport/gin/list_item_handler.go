package ginitem

import (
	"g09-social-todo-list/common"
	"g09-social-todo-list/demogrpc/demo"
	"g09-social-todo-list/module/item/biz"
	"g09-social-todo-list/module/item/model"
	"g09-social-todo-list/module/item/repository"
	"g09-social-todo-list/module/item/storage"
	"g09-social-todo-list/module/item/storage/rpc"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func ListItem(serviceCtx goservice.ServiceContext, client demo.ItemLikeServiceClient) func(*gin.Context) {
	return func(c *gin.Context) {
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		//apiItemCaller := serviceCtx.MustGet(common.PluginItemAPI).(interface {
		//	GetServiceURL() string
		//})

		var queryString struct {
			common.Paging
			model.Filter
		}

		if err := c.ShouldBind(&queryString); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		queryString.Paging.Process()

		requester := c.MustGet(common.CurrentUser).(common.Requester)

		store := storage.NewSQLStore(db)
		//likeStore := restapi.New(apiItemCaller.GetServiceURL(), serviceCtx.Logger("restapi.itemlikes"))

		likeStore := rpc.NewClient(client)
		repo := repository.NewListItemRepo(store, likeStore, requester)
		business := biz.NewListItemBiz(repo, requester)

		result, err := business.ListItem(c.Request.Context(), &queryString.Filter, &queryString.Paging)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		for i := range result {
			result[i].Mask()
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, queryString.Paging, queryString.Filter))
	}
}
