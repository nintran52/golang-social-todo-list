package rpc

import (
	"context"
	"g09-social-todo-list/demogrpc/demo"
)

type rpcClient struct {
	client demo.ItemLikeServiceClient
}

func NewClient(client demo.ItemLikeServiceClient) *rpcClient {

	return &rpcClient{client: client}
}

func (c *rpcClient) GetItemLikes(ctx context.Context, ids []int) (map[int]int, error) {
	reqIds := make([]int32, len(ids))

	for i := range ids {
		reqIds[i] = int32(ids[i])
	}

	resp, err := c.client.GetItemLikes(ctx, &demo.GetItemLikesReq{Ids: reqIds})

	if err != nil {
		return nil, err
	}

	rs := make(map[int]int)

	for k, v := range resp.Result {
		rs[int(k)] = int(v)
	}

	return rs, nil
}
