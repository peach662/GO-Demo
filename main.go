package main

import (
	"awesomeProject/internal/config"
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

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	db, err := database.OpenMySQL(cfg.MySQLDSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repo := todo.NewMySQLRepository(db)
	service := todo.NewService(repo)
	r := newRouter(service)

	if err := r.Run(":9090"); err != nil {
		panic(err)
	}
}

func newRouter(service *todo.Service) *gin.Engine {
	r := gin.Default()
	// ... (路由定义)
	r.GET("/todos", func(c *gin.Context) {
		items, err := service.List(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{
				"code":    500,
				"message": "服务内部错误",
				"data":    nil,
			})
			return
		}

		c.JSON(200, gin.H{
			"code":    0,
			"message": "ok",
			"data":    items,
		})
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
		item, found, err := service.GetByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(500, gin.H{
				"code":    500,
				"message": "服务内部错误",
				"data":    nil,
			})
			return
		}
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

		newTodo, err := service.Create(c.Request.Context(), req.Title)
		if err != nil {
			c.JSON(500, gin.H{
				"code":    500,
				"message": "服务内部错误",
				"data":    nil,
			})
			return
		}
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
		updatedTodo, found, err := service.UpdateStatus(c.Request.Context(), id, *req.Done)
		if err != nil {
			c.JSON(500, gin.H{
				"code":    500,
				"message": "服务内部错误",
				"data":    nil,
			})
			return
		}
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
