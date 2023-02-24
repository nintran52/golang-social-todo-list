package storage

import (
	"context"
	"g09-social-todo-list/common"
	"g09-social-todo-list/module/item/model"
)

func (s *sqlStore) ListItem(
	ctx context.Context,
	filter *model.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]model.TodoItem, error) {
	var result []model.TodoItem

	db := s.db.Where("status <> ?", "Deleted")

	// Get items of requester only
	//requester := ctx.Value(common.CurrentUser).(common.Requester)
	//db = db.Where("user_id = ?", requester.GetUserId())

	if f := filter; f != nil {
		if v := f.Status; v != "" {
			db = db.Where("status = ?", v)
		}
	}

	if err := db.Select("id").Table(model.TodoItem{}.TableName()).Count(&paging.Total).Error; err != nil {
		return nil, err
	}

	for i := range moreKeys {
		db = db.Preload(moreKeys[i])
	}

	if v := paging.FakeCursor; v != "" {
		uid, err := common.FromBase58(v)

		if err != nil {
			return nil, common.ErrDB(err)
		}

		db = db.Where("id < ?", uid.GetLocalID())
	} else {
		db = db.Offset((paging.Page - 1) * paging.Limit)
	}

	if err := db.Select("*").
		Order("id desc").
		Limit(paging.Limit).
		Find(&result).Error; err != nil {

		return nil, common.ErrDB(err)
	}

	if len(result) > 0 {
		result[len(result)-1].Mask()
		paging.NextCursor = result[len(result)-1].FakeId.String()
	}

	return result, nil
}
