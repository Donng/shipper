build:
	protoc -I. --go_out=plugins=grpc:$(GOPATH)/src/shipper/consignment-service \
	  proto/consignment/consignment.proto