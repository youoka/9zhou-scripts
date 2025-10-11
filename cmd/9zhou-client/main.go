package main

import (
	"9zhou-scripts/client"
	"9zhou-scripts/internal"
	"9zhou-scripts/pkg/config"
	"9zhou-scripts/pkg/utils"
	"fmt"
)

func main() {
	// 从key.yaml或用户输入读取账户信息
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("无法加载配置:", err)
		return
	}
	if !utils.CheckIsValid(cfg.AuthKey, cfg.ShopAccount, cfg.ReclaimAccount) {
		fmt.Println("授权码无效")
		return
	}
	fmt.Println("授权码有效")
	shop := client.NewShopAccount(cfg.ShopAccount, cfg.ShopPassword)
	reclaim := client.NewReclaimAccount(cfg.ReclaimAccount, cfg.ReclaimPassword)
	service := internal.NewService(shop, reclaim)
	err = service.Login()
	if err != nil {
		fmt.Println("登录失败:", err.Error())
		return
	}
	service.Hx()
	err = service.Transfer()
	if err != nil {
		fmt.Println("转账失败:", err.Error())
		return
	}
	err = service.Pay()
	if err != nil {
		fmt.Println("购买失败:", err.Error())
		return
	}
}
