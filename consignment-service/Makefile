build:
	protoc -I. --go_out=plugins=micro:. \
	  proto/consignment/consignment.proto
	GOARCH=amd64
	docker build --network host -t consignment-service .
run:
	docker run -p 50051:50051 \
	-e MICRO_SERVER_ADDRESS=:50051 \
	-e MICRO_REGISTRY=mdns consignment-service
