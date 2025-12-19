package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func main() {
	// 创建一个默认的 Gin 引擎
	// gin.Default() 会返回一个包含了 Logger 和 Recovery 中间件的引擎
	router := gin.Default()

	//基础的“Hello“路由
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello,Gin",
		})
	})

	//获取路径参数
	//:id是一个占位符，表示这个部分是动态的
	router.GET("/users/:id", func(c *gin.Context) {
		userID := c.Param("id")

		c.JSON(http.StatusOK, gin.H{
			"message": "你正在查询的ID是",
			"user_id": userID,
		})
	})

	//获取查询参数
	//访问/search？type=article&q=golang
	router.GET("/search", func(c *gin.Context) {
		//c.Query()获取URL？后面的参数
		queryType := c.Query("type")

		//c.DefaultQuery()和c.Query()类似，但如果参数不存在，会返回一个默认值
		queryWord := c.DefaultQuery("q", "default_keyword")
		c.JSON(http.StatusOK, gin.H{
			"message":     "你正在搜索",
			"search_type": queryType,
			"keyword":     queryWord,
		})
	})

	//处理POST请求，绑定JSON数据
	router.POST("/login", func(c *gin.Context) {
		var req LoginRequest

		//c.ShouldBindJSON()会尝试把请求体中的JSON解析并填充到req结构体中
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Json格式错误或缺少必要字段",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"error":    false,
			"message":  "login success",
			"username": req.Username,
		})
	})

	// 启动 HTTP 服务
	// 默认在 8080 端口启动服务
	// router.Run() 会一直阻塞，监听和处理请求
	router.Run(":8080")
	// 你也可以指定其他端口，例如 router.Run(":9000")
}
