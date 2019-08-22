package main

import (
	"context"
	"fmt"
	pb "github.com/Donng/shipper/vessel-service/proto/vessel"
)

type service struct {
	repo Repository
}

func (s *service) Create(ctx context.Context, req *pb.Vessel, res *pb.Response) error {
	if err := s.repo.Create(req); err != nil {
		return err
	}

	res.Vessel = req
	res.Created = true
	return nil
}

func (s *service) FindAvailable(ctx context.Context, req *pb.Specification, resp *pb.Response) error {
	vessel, err := s.repo.FindAvailable(req)
	if err != nil {
		fmt.Printf("err: %s, weight: %d, capacity: %d", err, req.MaxWeight, req.Capacity)
		return err
	}
	resp.Vessel = vessel

	return nil
}