package services

import (
	"context"
	"fmt"
)

type GoodsService struct {

}

//最终实现GetGoodsStock
func (this *GoodsService) GetGoodsStock(ctx context.Context, in *GoodsRequest) (*GoodsResponse, error)  {
	fmt.Println("GoodsRequest:",in)
	fmt.Println("GoodsID:",in.GoodsId)
	return &GoodsResponse{GoodsStock: in.GoodsId*10},nil
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


