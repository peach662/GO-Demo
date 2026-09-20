package main

import (
	"awesomeProject/internal/database"
	"awesomeProject/internal/todo"
	"github.com/gin-gonic/gin"
	"strconv"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>

	db, err := database.OpenMySQL()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	service := todo.NewService([]todo.Todo{
		{ID: 1, Title: "Learn Go", Done: false},
		{ID: 2, Title: "Build a web app", Done: false},
	})

	r := newRouter(service)

	r.Run(":9090")
}

func newRouter(service *todo.Service) *gin.Engine {
	r := gin.Default()
	// ... (路由定义)
	r.GET("/todos", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    0,
			"message": "ok",
			"data":    service.List()})
	})
	r.GET("/todos/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{
				"code":    400,
				"message": "id 必须是数字",
				"data":    nil,
			})
			return
		}
		item, found := service.GetByID(id)
		if !found {
			c.JSON(404, gin.H{
				"code":    404,
				"message": "todo 不存在",
				"data":    nil,
			})
			return
		}
		c.JSON(200, gin.H{
			"code":    0,
			"message": "ok",
			"data":    item,
		})
	})
	r.POST("/todos", func(c *gin.Context) {
		var req todo.CreateTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"data":    nil,
			})
			return
		}

		newTodo := service.Create(req.Title)
		c.JSON(200, gin.H{
			"code":    0,
			"message": "ok",
			"data":    newTodo,
		})
	})
	r.PATCH("/todos/:id", func(c *gin.Context) {

		var req todo.UpdateTodoStatusRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"data":    nil,
			})
			return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{
				"code":    400,
				"message": "id 必须是数字",
				"data":    nil,
			})
			return
		}
		updatedTodo, found := service.UpdateStatus(id, *req.Done)
		if !found {
			c.JSON(404, gin.H{
				"code":    404,
				"message": "todo 不存在",
				"data":    nil,
			})
			return
		}
		c.JSON(200, gin.H{
			"code":    0,
			"message": "ok",
			"data":    updatedTodo,
		})
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    0,
			"message": "ok",
			"data":    "pong"})
	})
	return r
}
