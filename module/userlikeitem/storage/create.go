package storage

import (
	"context"
	"g09-social-todo-list/common"
	"g09-social-todo-list/module/userlikeitem/model"
)

func (s *sqlStore) Create(ctx context.Context, data *model.Like) error {
	if err := s.db.Create(data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
