package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/micro/go-micro"
	pb "shipper/consignment-service/proto/consignment"
)

// 定义需要实现的服务接口
type IRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// 模拟数据库存储的功能，之后将用真正的实现代替
type Repository struct {
	// 读写锁，读锁下多个goroutine可读，写锁下其他goroutine不可读写
	mu           sync.RWMutex
	consignments []*pb.Consignment
}

// 创建新货物的本地实现
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()
	repo.consignments = append(repo.consignments, consignment)
	repo.mu.Unlock()
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// service 要实现在proto中定义的所有方法
type service struct {
	repo IRepository
}

// 创建货物的接口实现
func (s *service) CreateConsignment(ctx context.Context, consignment *pb.Consignment, resp *pb.Response) error {
	// 保存 consignment
	consignment, err := s.repo.Create(consignment)
	if err != nil {
		return err
	}

	resp.Created = true
	resp.Consignment = consignment
	return nil
}

// 获取所有货物的接口实现
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, resp *pb.Response) error {
	consignments := s.repo.GetAll()

	resp.Consignments = consignments
	return nil
}

func main() {
	repo := &Repository{}

	// 创建新的服务，其中包含一些可选的配置
	srv := micro.NewService(
		// name 方法必须与protobuf定义的package name相匹配
		micro.Name("go.micro.srv.consignment"),
	)

	// init 解析命令行参数
	srv.Init()

	// 注册 handler
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo})

	// 运行服务器
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
