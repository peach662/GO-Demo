package main

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/database"
	"awesomeProject/internal/todo"
	"awesomeProject/internal/user"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"awesomeProject/internal/auth"
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
	jwtService := auth.NewJWT(cfg.JWTSecret)

	repo := todo.NewMySQLRepository(db)
	service := todo.NewService(repo)
	userRepo := user.NewMySQLRepository(db)
	userService := user.NewService(userRepo)

	r := newRouter(service, userService, jwtService)
	if err := r.Run(":9090"); err != nil {
		panic(err)
	}
}

func newRouter(service *todo.Service, userService *user.Service, jwtService *auth.JWT) *gin.Engine {
	r := gin.Default()
	// ... (路由定义)
	r.POST("/users/login", func(c *gin.Context) {
		var req user.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"data":    nil,
			})
			return
		}
		login, err := userService.Login(c.Request.Context(), req.Username, req.Password)
		if errors.Is(err, user.ErrInvalidCredentials) {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "用户名或密码错误",
				"data":    nil,
			})
			return
		}
		if errors.Is(err, user.ErrInvalidUsername) {
			c.JSON(400, gin.H{
				"code":    400,
				"message": "用户名不能为空",
				"data":    nil,
			})
			return
		}
		if errors.Is(err, user.ErrInvalidPassword) {
			c.JSON(400, gin.H{
				"code":    400,
				"message": "密码不能为空",
				"data":    nil,
			})
			return
		}
		if err != nil {
			c.JSON(500, gin.H{
				"code":    500,
				"message": "服务内部错误",
				"data":    nil,
			})
			return
		}
		token, err := jwtService.GenerateToken(login.ID)
		if err != nil {
			c.JSON(500, gin.H{
				"code":    500,
				"message": "生成token失败",
				"data":    nil,
			})
			return
		}
		c.JSON(200, gin.H{
			"code":    0,
			"message": "ok",
			"data":    gin.H{
				"token": token,
				"user":  login,
			},
		})
	})
	r.POST("/users/register", func(c *gin.Context) {
		var req user.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"data":    nil,
			})
			return
		}
		registered, err := userService.Register(c.Request.Context(), req.Username, req.Password)
		if errors.Is(err, user.ErrUsernameTaken) {
			c.JSON(400, gin.H{
				"code":    400,
				"message": "用户名已存在",
				"data":    nil,
			})
			return
		}
		if errors.Is(err, user.ErrInvalidPassword) {
			c.JSON(400, gin.H{
				"code":    400,
				"message": "密码不能为空",
				"data":    nil,
			})
			return
		}
		if errors.Is(err, user.ErrInvalidUsername) {
			c.JSON(400, gin.H{
				"code":    400,
				"message": "用户名不能为空",
				"data":    nil,
			})
			return
		}
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
			"data":    registered,
		})
	})
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
		updatedTodo, found, err := service.UpdateStatus(c.Request.Context(), id, *req.Status)
		if err != nil {
			if errors.Is(err, todo.ErrInvalidTransition) {
				c.JSON(400, gin.H{
					"code":    400,
					"message": "非法状态流转",
					"data":    nil,
				})
				return
			}
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
