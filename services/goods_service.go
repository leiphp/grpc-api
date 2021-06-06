package services

import (
	"context"
	"fmt"
)

type GoodsService struct {

}

//最终实现GetGoodsStock
func (this *GoodsService) GetGoodsStock(ctx context.Context, request *GoodsRequest) (*GoodsResponse, error)  {
	fmt.Println("GoodsRequest:",request)
	fmt.Println("GoodsID:",request.GoodsId)
	var stock int32=0
	if request.GoodsArea == GoodsAreas_A {
		stock = 30
	}else if request.GoodsArea == GoodsAreas_B {
		stock = 31
	}else {
		stock = 32
	}
	return &GoodsResponse{GoodsStock: request.GoodsId + stock},nil
}

func (this *GoodsService) GetGoodsStocks(ctx context.Context, size *GoodsSize) (*GoodsResponseList, error) {
	Goodsres := []*GoodsResponse{
		&GoodsResponse{GoodsStock: 10},
		&GoodsResponse{GoodsStock: 12},
		&GoodsResponse{GoodsStock: 15},
		&GoodsResponse{GoodsStock: 28},
	}
	return &GoodsResponseList{
		Goodsres: Goodsres,
	}, nil
}


