package main

import (
	"context"
	"errors"
	pb "github.com/Donng/shipper/vessel-service/proto/vessel"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(*pb.Vessel) error
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

type VesselRepository struct {
	collection *mongo.Collection
}

func (repo *VesselRepository) Create(vessel *pb.Vessel) error {
	_, err := repo.collection.InsertOne(context.Background(), vessel)

	return err
}

// 根据特定条件依次对比轮船，如果要求的容量和最大重量都低于轮船，则返回此轮船。
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	cur, err := repo.collection.Find(context.Background(), nil, nil)
	if err != nil {
		return nil, err
	}

	//todo 待优化：使用MongoDB更高阶的查询方式
	var vessels []*pb.Vessel
	for cur.Next(context.Background()) {
		var vessel *pb.Vessel
		err := cur.Decode(&vessel)
		if err != nil {
			return nil, err
		}
		vessels = append(vessels, vessel)
	}

	for _, vessel := range vessels {
		if spec.MaxWeight <= vessel.MaxWeight && spec.Capacity <= vessel.Capacity {
			return vessel, nil
		}
	}

	return nil, errors.New("No available vessel.")
}

