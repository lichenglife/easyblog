package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	apiv1 "github.com/lichenglife/easyblog/pkg/api/apiserver/v1"
)

// user service grpc 请求客户端

var (
	// 定义命令行参数
	addr = flag.String("addr", "localhost:9090", "grpc server address")
)

func main() {
	// 解析命令行参数
	flag.Parse()
	// 创建grpc连接
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("cannot connect to grpc server %v", err)
	}
	// 关闭连接，释放资源
	defer conn.Close()
	// 创建grpc客户端，调用grpc服务
	client := apiv1.NewEasyblogClient(conn)

	// 请求grpc服务，进行健康校验
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second)
	defer cancel()

	// 创建用户请求
	userCreateResponse, err := client.CreateUser(ctx, &apiv1.CreateUserRequest{
		Username: "wangwu11",
		Password: "123456Aa",
		Nickname: "wangwu11",
		Email:    "7098905856@qq.com",
		Phone:    "12564564561",
	})
	if err != nil {
		log.Printf("用户创建失败: %v", err)
	}
	fmt.Println("======创建用户成功=======")
	fmt.Println(userCreateResponse)
	// 登录请求
	userLoginResponse, err := client.Login(ctx, &apiv1.LoginRequest{Username: "wangwu11", Password: "123456Aa"})
	if err != nil {
		log.Fatalf("用户登录失败: %v", err)
		return
	}
	fmt.Printf("======用户%s成功======= \n", userCreateResponse.UserID)
	// 创建 metadata，用于传递 Token
	md := metadata.Pairs("Authorization", "Bearer "+userLoginResponse.Token)
	// 将 metadata 附加到上下文中
	ctx = metadata.NewOutgoingContext(ctx, md)

	// 查询用户
	user, err := client.GetUser(ctx, &apiv1.GetUserRequest{UserID: userCreateResponse.UserID})
	if err != nil {
		log.Fatalf("cannot call grpc server %v", err)
		return
	}
	// 打印响应结果
	jsonData, _ := json.Marshal(user)
	fmt.Println(string(jsonData))

	// 查询所有用户
	userList, err := client.ListUser(ctx, &apiv1.ListUserRequest{})
	if err != nil {
		log.Fatalf("request ListUser error  %v", err)
		return
	}
	for _, user := range userList.Users {
		fmt.Println(user)
	}
	// 更新用户
	response, err := client.UpdateUser(ctx, &apiv1.UpdateUserRequest{UserID: user.User.UserID, Email: "88888888@qq.com"})
	if err != nil {
		log.Fatalf("request UpdateUser error  %v", err)
		return
	}
	fmt.Println(response)

	// 查询更新后的用户信息
	user, err = client.GetUser(ctx, &apiv1.GetUserRequest{UserID: userCreateResponse.UserID})
	if err != nil {
		log.Fatalf("cannot call grpc server %v", err)
		return
	}
	// 打印响应结果
	jsonData, _ = json.Marshal(user)
	fmt.Println(string(jsonData))
	// 删除用户

	// 登录请求
	userLoginResponse, err = client.Login(ctx, &apiv1.LoginRequest{Username: "root", Password: "123456Aa"})
	if err != nil {
		log.Fatalf("用户登录失败: %v", err)
		return
	}
	fmt.Printf("======用户%s登录成功======= \n", "root")
	// 创建 metadata，用于传递 Token
	md = metadata.Pairs("Authorization", "Bearer "+userLoginResponse.Token)
	// 将 metadata 附加到上下文中
	ctx = metadata.NewOutgoingContext(ctx, md)

	deleteResponse, err := client.DeleteUser(ctx, &apiv1.DeleteUserRequest{UserID: userCreateResponse.UserID})
	if err != nil {
		log.Fatalf("request DeleteUser error  %v", err)
		return
	}
	fmt.Printf("======删除用户%s成功======= \n", userCreateResponse.UserID)
	fmt.Println(deleteResponse)

}
