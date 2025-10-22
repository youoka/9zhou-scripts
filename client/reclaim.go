package client

import (
	"9zhou-scripts/pkg/http_client"
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	ReclaimApi      = "https://reclaim-api.9zhou.shop"
	ReclaimLoginURL = "%s/auth/login"
	HxURL           = "%s/reclaim"
	TransferURL     = "%s/transfer"
)

type ReclaimAccount struct {
	Account     string
	PassWord    string
	Token       string
	ReclaimInfo *ReclaimAccountInfo
}

func NewReclaimAccount(account string, passWord string) *ReclaimAccount {
	return &ReclaimAccount{Account: account, PassWord: passWord}
}
func (r *ReclaimAccount) Login() error {
	req := LoginReq{
		Account:  r.Account,
		PassWord: r.PassWord,
	}
	response := LoginResponse{}
	err := http_client.Post(fmt.Sprintf(ReclaimLoginURL, ReclaimApi), "", req, &response)
	if err != nil {
		return err
	}
	if response.Code != 0 && response.Code != 200 {
		return fmt.Errorf("核销登录失败: %s", response.Msg)
	}
	r.Token = response.Data.Token
	return err
}

type ReclaimAccountInfo struct {
	XTraceId string `json:"x_trace_id"`
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	Data     struct {
		Id        string `json:"id"`
		Account   string `json:"account"`
		NickName  string `json:"nickName"`
		Email     string `json:"email"`
		Password  string `json:"Password"`
		ShopUsers []struct {
			Account string `json:"account"`
		} `json:"shopUsers"`
		TotpSecret  string    `json:"TotpSecret"`
		Status      int       `json:"status"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
		LastLoginAt time.Time `json:"lastLoginAt"`
		Wallet      struct {
			Id      string `json:"id"`
			Account string `json:"account"`
			Balance string `json:"balance"`
		} `json:"wallet"`
		Distributor struct {
			Id             string  `json:"id"`
			ReclaimAccount string  `json:"reclaimAccount"`
			CommissionRate float64 `json:"commissionRate"`
			Uplink         struct {
				Id             string `json:"id"`
				ReclaimAccount string `json:"reclaimAccount"`
			} `json:"uplink"`
		} `json:"distributor"`
	} `json:"data"`
}

func (r *ReclaimAccount) Info() error {
	response := ReclaimAccountInfo{}
	err := http_client.Get(fmt.Sprintf("%s/user/info", ReclaimApi), r.Token, nil, &response)
	if err != nil {
		return err
	}
	if response.Code != 0 && response.Code != 200 {
		return fmt.Errorf("获取用户信息失败: %s", response.Msg)
	}
	r.ReclaimInfo = &response
	return nil
}

type HxReq struct {
	Cards []string `json:"cards"`
}

func (r *ReclaimAccount) Hx(orders []string) error {
	req := HxReq{
		Cards: orders,
	}
	return http_client.Post(fmt.Sprintf(HxURL, ReclaimApi), r.Token, req, nil)
}

type TransferToShopUserRequest struct {
	ShopUserAccount string  `json:"shop_user_account" binding:"required"`               // 商城用户账号
	Amount          float64 `json:"amount" binding:"required" binding:"required,min=1"` // 转账金额
	Password        string  `json:"password" binding:"required"`
	//TotpCode        string  `json:"totp_code" binding:"required"`
}

func (r *ReclaimAccount) Transfer(account string) error {
	f, err := strconv.ParseFloat(r.ReclaimInfo.Data.Wallet.Balance, 64)
	if err != nil {
		return errors.New("余额转换错误")
	}
	if f < 100 {
		return errors.New("核销余额不足，跳过转账")
	}
	req := TransferToShopUserRequest{
		ShopUserAccount: account,
		Amount:          f,
		Password:        r.PassWord,
	}
	var response BaseResponse
	err = http_client.Post(fmt.Sprintf(TransferURL, ReclaimApi), r.Token, req, &response)
	if err != nil {
		return err
	}
	if response.Code != 0 && response.Code != 200 {
		return fmt.Errorf("转账失败: %s", response.Msg)
	}
	fmt.Println("转账成功", f)
	return nil
}
