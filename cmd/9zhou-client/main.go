package main

import (
	"9zhou-scripts/client"
	"9zhou-scripts/pkg/config"
	"9zhou-scripts/pkg/utils"
	"fmt"
	"math"
	"strconv"
	"time"
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
	shopCredentials := []string{cfg.ShopAccount, cfg.ShopPassword}
	reclaimCredentials := []string{cfg.ReclaimAccount, cfg.ReclaimPassword}

	account := client.NewShopAccount(shopCredentials[0], shopCredentials[1])
	reclaim := client.NewReclaimAccount(reclaimCredentials[0], reclaimCredentials[1])
	reclaim.Login()
	account.Login()
	reclaim.Info()
	for {
		order, err := account.GetHXOrder()
		if err != nil {
			return
		}
		reclaim.Hx(order)
		if len(order) < 100 {
			break
		}
	}
	reclaim.Transfer(shopCredentials[0])
	name, err := account.Info()
	if err != nil {
		return
	}
	go func() {
		for {
			time.Sleep(time.Second * 3600)
			account.Login()
			reclaim.Login()
		}
	}()
	fmt.Println(name)
	f, err := strconv.ParseFloat(account.UserInfo.Data.Wallet.Balance, 64)
	if err != nil {
		fmt.Println("余额转换错误:", err)
		return
	}
	if f < 1000 {
		fmt.Println("余额不足")
		return
	}
	fmt.Println("余额:", f)
	fmt.Println("开始下单")
	count := math.Floor(f / 1000)
	for i := 0; i < int(count); i++ {
		time.Sleep(time.Second * 2)
		orderId, err := account.CreateOrder(client.P1000)
		if err != nil {
			return
		}
		if orderId != "" {
			msg, err := account.PayOrder(orderId)
			fmt.Println(msg)
			if err != nil {
				return
			}
		}
	}
}
