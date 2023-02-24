package model

import (
	"errors"
	"g09-social-todo-list/common"
	"strings"
)

var (
	ErrTitleCannotBeEmpty = errors.New("title cannot be empty")
	ErrItemIsDeleted      = errors.New("item is deleted")
)

const (
	EntityName = "Item"
)

type TodoItem struct {
	common.SQLModel
	UserId      int                `json:"-" gorm:"column:user_id;"`
	Title       string             `json:"title" gorm:"column:title;"`
	Description string             `json:"description" gorm:"column:description;"`
	Status      string             `json:"status" gorm:"column:status;"`
	Image       *common.Image      `json:"image" gorm:"column:image;"`
	LikedCount  int                `json:"liked_count" gorm:"-"`
	Owner       *common.SimpleUser `json:"owner" gorm:"foreignKey:UserId;"`
}

func (TodoItem) TableName() string { return "todo_items" }

func (i *TodoItem) Mask() {
	i.SQLModel.Mask(common.DbTypeItem)

	if v := i.Owner; v != nil {
		v.Mask()
	}
}

type TodoItemCreation struct {
	Id          int           `json:"id" gorm:"column:id;"`
	UserId      int           `json:"-" gorm:"column:user_id;"`
	Title       string        `json:"title" gorm:"column:title;"`
	Description string        `json:"description" gorm:"column:description;"`
	Image       *common.Image `json:"image" gorm:"column:image;"`
}

func (i *TodoItemCreation) Validate() error {
	i.Title = strings.TrimSpace(i.Title)

	if i.Title == "" {
		return ErrTitleCannotBeEmpty
	}

	return nil
}

func (TodoItemCreation) TableName() string { return TodoItem{}.TableName() }

type TodoItemUpdate struct {
	Title       *string `json:"title" gorm:"column:title;"`
	Description *string `json:"description" gorm:"column:description;"`
	Status      *string `json:"status" gorm:"column:status;"`
}

func (TodoItemUpdate) TableName() string { return TodoItem{}.TableName() }
