package main

import (
	"context"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "shipper/consignment-service/proto/consignment"
)

const (
	port = ":50051"
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
func (s *service) CreateConsignment(ctx context.Context, consignment *pb.Consignment) (*pb.Response, error) {
	// 保存 consignment
	consignment, err := s.repo.Create(consignment)
	if err != nil {
		return nil, err
	}

	return &pb.Response{Created: true, Consignment: consignment}, nil
}

// 获取所有货物的接口实现
func (s *service) GetConsignments(ctx context.Context,req *pb.GetRequest) (*pb.Response, error) {
	consignments := s.repo.GetAll()

	return &pb.Response{Consignments:consignments}, nil
}

func main() {
	repo := &Repository{}

	// 设置 gRPC 服务
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// 注册服务到 gRPC
	pb.RegisterShippingServiceServer(s, &service{repo})

	// 注册反射服务到 gRPC 服务器
	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
