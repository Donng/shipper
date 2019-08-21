package main

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/Donng/shipper/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
)

type IRepository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

type Repository struct {
	vessels []*pb.Vessel
}

// 根据特定条件依次对比轮船，如果要求的容量和最大重量都低于轮船，则返回此轮船。
func (repo *Repository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	for _, vessel := range repo.vessels {
		if spec.Capacity <= vessel.Capacity && spec.MaxWeight <= vessel.Capacity {
			return vessel, nil
		}
	}

	return nil, errors.New("No vessel found by that spec")
}

// gRPC 服务处理者
type service struct {
	repo IRepository
}

func (s *service) FindAvailable(ctx context.Context, req *pb.Specification, resp *pb.Response) error {
	vessel, err := s.repo.FindAvailable(req)
	if err == nil {
		resp.Vessel = vessel
	}

	return nil
}

func main() {
	vessels := []*pb.Vessel{
		&pb.Vessel{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}
	repo := &Repository{vessels}

	srv := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
	)

	srv.Init()

	pb.RegisterVesselServiceHandler(srv.Server(), &service{repo})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
