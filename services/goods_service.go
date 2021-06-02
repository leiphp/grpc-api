package services

import (
	"context"
)

type GoodsService struct {

}

func (this *GoodsService) GetGoodsStock(ctx context.Context, in *GoodsRequest) (*GoodsResponse, error)  {
	return &GoodsResponse{GoodsStock:25},nil
}


