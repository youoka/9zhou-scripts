package main

import (
	"9zhou-scripts/client"
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// 从文件读取账户信息
	file, err := os.Open("9zhou.txt")
	if err != nil {
		fmt.Println("无法打开配置文件:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) < 2 {
		fmt.Println("配置文件格式不正确，至少需要两行")
		return
	}

	// 第一行为shop账号和密码
	shopCredentials := strings.Split(lines[0], " ")
	if len(shopCredentials) != 2 {
		fmt.Println("商城账号配置格式不正确，应为'账号 密码'")
		return
	}

	// 第二行为reclaim账号和密码
	reclaimCredentials := strings.Split(lines[1], " ")
	if len(reclaimCredentials) != 2 {
		fmt.Println("核销账号配置格式不正确，应为'账号 密码'")
		return
	}

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
