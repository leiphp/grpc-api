package services

import (
	"context"
	"fmt"
)

type OrdersService struct {

}

//创建订单
func (this *OrdersService) CreateOrder(ctx context.Context, request *OrdersRequest) (*OrdersResponse, error)  {
	fmt.Println("request:",request)
	return &OrdersResponse{Status: "ok", Message:"成功"},nil
}


