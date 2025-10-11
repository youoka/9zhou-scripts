package main

import (
	"9zhou-scripts/pkg/utils"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// KeyGenerateRequest 定义生成密钥的请求参数
type KeyGenerateRequest struct {
	ShopAccount    string `json:"shop_account" binding:"required"`
	ReclaimAccount string `json:"reclaim_account" binding:"required"`
	ExpiryDate     string `json:"expiry_date" binding:"required"` // 格式: YYYY-MM-DD
}

// KeyGenerateResponse 定义生成密钥的响应
type KeyGenerateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Key     string `json:"key,omitempty"`
}

func main() {
	// 设置为发布模式
	gin.SetMode(gin.ReleaseMode)

	// 创建 Gin 路由器
	router := gin.Default()

	// 静态文件服务
	router.Static("/static", "./static")
	router.LoadHTMLFiles("./static/index.html")

	// 定义路由
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router.POST("/api/key/generate", generateKeyHandler)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// 启动服务器
	router.Run(":8080")
}

// generateKeyHandler 处理生成密钥的请求
func generateKeyHandler(c *gin.Context) {
	var req KeyGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, KeyGenerateResponse{
			Success: false,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证日期格式
	_, err := time.Parse("2006-01-02", req.ExpiryDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, KeyGenerateResponse{
			Success: false,
			Message: "日期格式错误，应为 YYYY-MM-DD",
		})
		return
	}

	// 构造授权码明文
	// 格式: 商城号|核销号(到期时间)
	authCode := req.ShopAccount + "|" + req.ReclaimAccount + "|(" + req.ExpiryDate + ")"

	// 加密授权码
	encryptedKey, err := utils.HashAuthorization(authCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, KeyGenerateResponse{
			Success: false,
			Message: "密钥生成失败: " + err.Error(),
		})
		return
	}

	// 保存记录到 log.txt
	logEntry := fmt.Sprintf("[%s] 商城账号: %s, 回收账号: %s, 到期时间: %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		req.ShopAccount,
		req.ReclaimAccount,
		req.ExpiryDate)

	// 以追加模式打开文件
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// 记录日志失败不应该影响主要功能
		fmt.Printf("无法打开日志文件: %v\n", err)
	} else {
		// 写入日志
		_, err = file.WriteString(logEntry)
		if err != nil {
			fmt.Printf("写入日志文件失败: %v\n", err)
		}
		file.Close()
	}

	// 返回成功响应
	c.JSON(http.StatusOK, KeyGenerateResponse{
		Success: true,
		Message: "密钥生成成功",
		Key:     encryptedKey,
	})
}
