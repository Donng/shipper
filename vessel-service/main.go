package main

import (
	"fmt"
	"log"
	"os"

	pb "github.com/Donng/shipper/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
)

const (
	defaultHost = "localhost:27017"
)

func createDummyData(repo Repository) {
	vessels := []*pb.Vessel{
		{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}

	for _, v := range vessels {
		repo.Create(v)
	}
}

func main() {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = defaultHost
	}

	client, err := CreateClient(host)
	if err != nil {
		log.Fatalf("Create client error: %v", err)
	}

	repo := &VesselRepository{client.Database(dbName).Collection(vesselCollection)}

	// 添加模拟数据
	createDummyData(repo)

	srv := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
	)

	srv.Init()

	pb.RegisterVesselServiceHandler(srv.Server(), &service{repo})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
