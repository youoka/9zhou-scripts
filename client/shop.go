package client

import (
	"9zhou-scripts/pkg/http_client"
	"fmt"
	"strconv"
	"time"
)

const (
	ShopAddr       = "https://shop-api.9zhou.shop"
	loginURL       = "%s/auth/login"
	CreateOrderURL = "%s/order"
	PayOrderURL    = "%s/order/%s/pay"
	HxOrderURL     = "%s/order?page=%s&page_size=100&startTime=%s&endTime=%s&status=%s"
)
const (
	P1000 = "68dea238d5e367848e140696"
	P500  = "68dea21fd5e367848e140695"
	P200  = "68dea1fdd5e367848e140694"
	P100  = "68dea1aad5e367848e140693"
)

type BaseResponse struct {
	XTraceId string `json:"x-trace-id"`
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
}

type LoginData struct {
	Token      string `json:"Token"`
	ExpireTime int64  `json:"expire_time"`
}
type LoginResponse struct {
	BaseResponse
	Data LoginData `json:"data"`
}
type LoginReq struct {
	Account  string `json:"account"`
	PassWord string `json:"password"`
}
type ShopAccount struct {
	Account  string
	PassWord string
	token    string
	UserInfo *ShopUserInfo
}

func NewShopAccount(account string, passWord string) *ShopAccount {
	return &ShopAccount{Account: account, PassWord: passWord}
}
func (s *ShopAccount) Login() error {
	req := LoginReq{
		Account:  s.Account,
		PassWord: s.PassWord,
	}
	response := LoginResponse{}
	err := http_client.Post(fmt.Sprintf(loginURL, ShopAddr), "", req, &response)
	if err != nil {
		return err
	}
	if response.Code != 0 {
		return fmt.Errorf("商城登录失败: %s", response.Msg)
	}
	s.token = response.Data.Token
	return nil
}

type UserInfoResponse struct {
	BaseResponse
	Data *ShopUserInfo `json:"data"`
}
type ShopUserInfo struct {
	XTraceId string `json:"x-trace-id"`
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	Data     struct {
		Id          string    `json:"id"`
		Account     string    `json:"account"`
		NickName    string    `json:"nickName"`
		Email       string    `json:"email"`
		Phone       string    `json:"phone"`
		Tag         int       `json:"tag"`
		Status      int       `json:"status"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
		LastLoginAt time.Time `json:"lastLoginAt"`
		Wallet      struct {
			Id       string `json:"id"`
			Account  string `json:"account"`
			Balance  string `json:"balance"`
			USDTAddr string `json:"USDTAddr"`
		} `json:"wallet"`
	} `json:"data"`
}

func (s *ShopAccount) Info() error {
	response := ShopUserInfo{}
	err := http_client.Get(fmt.Sprintf("%s/user/info", ShopAddr), s.token, nil, &response)
	if err != nil {
		return err
	}
	if response.Code != 0 {
		return fmt.Errorf("获取用户信息失败: %s", response.Msg)
	}
	s.UserInfo = &response
	return err
}

type CreateOrderRequest struct {
	ProductItems []OrderProductItem `json:"product_items" binding:"required"`
}

// OrderProductItem 订单商品项
type OrderProductItem struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}
type CreateOrderResp struct {
	XTraceId string `json:"x-trace-id"`
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	Data     struct {
		Id       string `json:"id"`
		Account  string `json:"account"`
		Products []struct {
			Id           string    `json:"id"`
			Name         string    `json:"name"`
			CoverURL     string    `json:"coverURL"`
			Description  string    `json:"description"`
			Sort         int       `json:"sort"`
			Price        int       `json:"price"`
			Quantity     int       `json:"quantity"`
			PaymentTypes []int     `json:"paymentTypes"`
			CategoryId   string    `json:"categoryId"`
			CreatedAt    time.Time `json:"createdAt"`
			UpdatedAt    time.Time `json:"updatedAt"`
		} `json:"products"`
		TotalPrice  int       `json:"totalPrice"`
		PaymentTime time.Time `json:"paymentTime"`
		Status      int       `json:"status"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
	} `json:"data"`
}

func (s *ShopAccount) CreateOrder(id string) (string, error) {
	request := CreateOrderRequest{
		ProductItems: []OrderProductItem{
			{
				ProductID: id,
				Quantity:  1,
			},
		},
	}
	resp := CreateOrderResp{}
	err := http_client.Post(fmt.Sprintf(CreateOrderURL, ShopAddr), s.token, request, &resp)
	if err != nil {
		return "", err
	}
	if resp.Code != 0 {
		return "", fmt.Errorf("创建订单失败: %s", resp.Msg)
	}
	return resp.Data.Id, err
}

type PayOrderRequest struct {
	PaymentType int `json:"payment_type" binding:"required"` // 1:余额支付, 2:微信支付, 3:支付宝支付
}

func (s *ShopAccount) PayOrder(orderId string) (string, error) {
	request := PayOrderRequest{
		PaymentType: 1,
	}
	resp := BaseResponse{}
	err := http_client.Post(fmt.Sprintf(PayOrderURL, ShopAddr, orderId), s.token, request, &resp)
	if err != nil {
		return "", err
	}
	if resp.Code != 0 {
		return "", fmt.Errorf("支付订单失败: %s", resp.Msg)
	}
	return resp.Msg, err
}

type GetHXOrderResp struct {
	XTraceId string `json:"x-trace-id"`
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	Data     []struct {
		Id       string `json:"id"`
		Account  string `json:"account"`
		Products []struct {
			Id           string    `json:"id"`
			Name         string    `json:"name"`
			CoverURL     string    `json:"coverURL"`
			Description  string    `json:"description"`
			Sort         int       `json:"sort"`
			Price        int       `json:"price"`
			Quantity     int       `json:"quantity"`
			PaymentTypes []int     `json:"paymentTypes"`
			CategoryId   string    `json:"categoryId"`
			CreatedAt    time.Time `json:"createdAt"`
			UpdatedAt    time.Time `json:"updatedAt"`
		} `json:"products"`
		TotalPrice  int       `json:"totalPrice"`
		PaymentType int       `json:"paymentType"`
		PaymentTime time.Time `json:"paymentTime"`
		Status      int       `json:"status"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
	} `json:"data"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalCount int `json:"total_count"`
	TotalPage  int `json:"total_page"`
}

func (s *ShopAccount) GetShippedOrder(startTime, endTime string, size int) (*GetHXOrderResp, error) {
	resp := GetHXOrderResp{}
	err := http_client.Get(fmt.Sprintf(HxOrderURL, ShopAddr, strconv.Itoa(size), startTime, endTime, "3"), s.token, nil, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("获取订单失败: %s", resp.Msg)
	}
	return &resp, err
}

// 获取已支付的订单
func (s *ShopAccount) GetPaidOrder(startTime, endTime string, size int) (*GetHXOrderResp, error) {
	resp := GetHXOrderResp{}
	err := http_client.Get(fmt.Sprintf(HxOrderURL, ShopAddr, strconv.Itoa(size), startTime, endTime, "2"), s.token, nil, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("获取订单失败: %s", resp.Msg)
	}
	return &resp, err
}
func (s *ShopAccount) GetCancelledOrder(startTime, endTime string, size int) (*GetHXOrderResp, error) {
	resp := GetHXOrderResp{}
	err := http_client.Get(fmt.Sprintf(HxOrderURL, ShopAddr, strconv.Itoa(size), startTime, endTime, "5"), s.token, nil, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("获取订单失败: %s", resp.Msg)
	}
	return &resp, err
}
