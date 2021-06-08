package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"
)

type UserService struct {

}

//获取用户积分
func (this *UserService) GetUserScore(ctx context.Context, request *UserScoreRequest) (*UserScoreResponse, error)  {
	fmt.Println("request:",request)
	var score int32=101
	users := make([]*UserInfo, 0)
	for _, user := range request.Users {
		user.UserScore = score
		score++
		users = append(users, user)
	}
	return &UserScoreResponse{Users:users},nil
}

//服务端流模式返回积分
func (this *UserService) GetUserScoreByServerStream(request *UserScoreRequest, stream UserService_GetUserScoreByServerStreamServer) error {
	fmt.Println("server stream request:",request)
	var score int32=101
	users := make([]*UserInfo, 0)
	for index, user := range request.Users {
		user.UserScore = score
		score++
		users = append(users, user)
		if (index+1)%2 == 0 && index > 0 {
			err := stream.Send(&UserScoreResponse{Users: users})
			if err != nil {
				return err
			}
			users = (users)[0:0]
		}
		time.Sleep(time.Second * 1)
	}
	//如果还有值，再发送一次
	if len(users) > 0 {
		err := stream.Send(&UserScoreResponse{Users:users})
		if err != nil {
			return err
		}
	}
	return nil
}

//客户端流模式请求积分
func (this *UserService) GetUserScoreByClientStream(stream UserService_GetUserScoreByClientStreamServer) error {
	var score int32=101
	users := make([]*UserInfo, 0)
	for{
		req, err := stream.Recv()
		if err == io.EOF {//接收完了
			return stream.SendAndClose(&UserScoreResponse{Users:users})
		}
		if err != nil {
			return err
		}
		for _, user := range req.Users {
			user.UserScore = score
			score++
			users = append(users, user)
		}
	}
}

//双向流模式获取积分
func (this *UserService) GetUserScoreByTWStream(stream UserService_GetUserScoreByTWStreamServer) error {
	var score int32=101
	users := make([]*UserInfo, 0)
	for{
		req, err := stream.Recv()
		if err == io.EOF {//接收完了
			return nil
		}
		if err != nil {
			return err
		}
		for _, user := range req.Users {
			user.UserScore = score
			score++
			users = append(users, user)
		}
		err = stream.Send(&UserScoreResponse{Users:users})
		if err != nil {
			log.Println(err)
		}
		users = (users)[0:0]
	}
}


