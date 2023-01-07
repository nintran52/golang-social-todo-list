package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type TodoItem struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func main() {
	now := time.Now().UTC()

	item := TodoItem{
		Id:          1,
		Title:       "Task 1",
		Description: "Content 1",
		Status:      "Doing",
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}

	jsData, err := json.Marshal(item)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(jsData))

	jsString := "{\"id\":1,\"title\":\"Task 1\",\"description\":\"Content 1\",\"status\":\"Doing\",\"created_at\":\"2023-01-06T14:21:00.47733Z\",\"updated_at\":\"2023-01-06T14:21:00.47733Z\"}"

	var item2 TodoItem

	if err := json.Unmarshal([]byte(jsString), &item2); err != nil {
		log.Fatalln(err)
	}

	log.Println(item2)

	//////////////////

	r := gin.Default()

	v1 := r.Group("/v1")
	{
		items := v1.Group("/items")
		{
			items.POST("")
			items.GET("")
			items.GET("/:id")
			items.PATCH("/:id")
			items.DELETE("/:id")
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if err := r.Run(":3000"); err != nil {
		log.Fatalln(err)
	}
}
