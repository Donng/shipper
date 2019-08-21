package main

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/Donng/shipper/consignment-service/proto/consignment"
	vesselProto "github.com/Donng/shipper/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
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
	repo         IRepository
	vesselClient vesselProto.VesselServiceClient
}

// 创建货物的接口实现
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, resp *pb.Response) error {
	// 调用货船服务的客户端实例，使用货运重量和此次货运的集装箱的数量作为容量值
	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	if err != nil {
		return err
	}
	fmt.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)

	req.VesselId = vesselResponse.Vessel.Id

	consignment, err := s.repo.Create(req)
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

	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

	// 注册 handler
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})

	// 运行服务器
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
