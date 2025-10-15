package main

import (
	"9zhou-scripts/pkg/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConfigRequest struct {
	Type    int    `json:"type"`
	Num1000 int    `json:"num1000"`
	Num500  int    `json:"num500"`
	Num200  int    `json:"num200"`
	Num100  int    `json:"num100"`
	Date    string `json:"date"`
}

type ShopAccountRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type HxAccountRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Key      string `json:"key"`
}

func main() {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.Use(gin.Recovery())
	// 提供静态文件服务
	r.Static("/static", "./static")
	r.LoadHTMLFiles("./static/config.html")

	// 首页重定向到配置页面
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/config.html")
	})

	// 配置管理路由
	configGroup := r.Group("/config")
	{
		configGroup.GET("", getConfig)
		configGroup.POST("", updateConfig)
	}

	// 订单统计路由（仅查询功能）
	orderStatsGroup := r.Group("/order-stats")
	{
		orderStatsGroup.GET("/:date", getOrderStatistics)
	}

	// 商店账号管理路由
	shopAccountGroup := r.Group("/shop-account")
	{
		shopAccountGroup.GET("", getShopAccounts)
		shopAccountGroup.POST("", createShopAccount)
	}

	// 核销账号管理路由（只有一个账号）
	hxAccountGroup := r.Group("/hx-account")
	{
		hxAccountGroup.GET("", getHxAccount)
		hxAccountGroup.POST("", updateHxAccount)
	}

	r.Run(":8080")
}

// getConfig 获取配置信息
func getConfig(c *gin.Context) {
	config, err := database.Db.GetConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, config)
}

// updateConfig 更新配置信息
func updateConfig(c *gin.Context) {
	var req ConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := &database.Config{
		Type:    req.Type,
		Num1000: req.Num1000,
		Num500:  req.Num500,
		Num200:  req.Num200,
		Num100:  req.Num100,
		Date:    req.Date,
	}

	// 先尝试获取现有配置
	existingConfig, err := database.Db.GetConfig()
	if err != nil || existingConfig == nil {
		// 如果没有现有配置，则创建新配置
		if err := database.Db.CreateConfig(config); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 如果存在配置，则更新
		if err := database.Db.UpdateConfig(config); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, config)
}

// getOrderStatistics 获取指定日期的订单统计信息
func getOrderStatistics(c *gin.Context) {
	date := c.Param("date")
	stats, err := database.Db.GetOrderStatistics(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// getShopAccounts 获取所有商店账号
func getShopAccounts(c *gin.Context) {
	accounts, err := database.Db.GetAllShopAccount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

// createShopAccount 创建商店账号
func createShopAccount(c *gin.Context) {
	var req ShopAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account := &database.ShopAccount{
		Account:  req.Account,
		Password: req.Password,
	}

	if err := database.Db.CreateShopAccount(account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

// getHxAccount 获取核销账号
func getHxAccount(c *gin.Context) {
	account, err := database.Db.GetHxAccount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, account)
}

// updateHxAccount 更新核销账号
func updateHxAccount(c *gin.Context) {
	var req HxAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account := &database.HxAccount{
		Account:  req.Account,
		Password: req.Password,
		Key:      req.Key,
	}

	// 先尝试获取现有账号
	existingAccount, err := database.Db.GetHxAccount()
	if err != nil || existingAccount == nil {
		// 如果没有现有账号，则创建新账号
		if err := database.Db.CreateHxAccount(account); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 如果存在账号，则更新
		if err := database.Db.UpdateHxAccount(account); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, account)
}
