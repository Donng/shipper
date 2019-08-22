package main

import (
	"context"
	"fmt"

	pb "github.com/Donng/shipper/consignment-service/proto/consignment"
	vesselProto "github.com/Donng/shipper/vessel-service/proto/vessel"
)

type handler struct {
	repository
	vesselClient vesselProto.VesselServiceClient
}

func (s *handler) CreateConsignment(ctx context.Context, req *pb.Consignment, resp *pb.Response) error {
	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	if err != nil {
		return err
	}
	fmt.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)

	req.VesselId = vesselResponse.Vessel.Id

	if err = s.repository.Create(req); err != nil {
		return err
	}

	resp.Created = true
	resp.Consignment = req
	return nil
}

func (s *handler) GetConsignments(ctx context.Context, req *pb.GetRequest, resp *pb.Response) error {
	consignments, err := s.repository.GetAll()
	if err != nil {
		return err
	}

	resp.Consignments = consignments
	return nil
}
