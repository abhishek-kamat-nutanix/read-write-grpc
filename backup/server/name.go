package main

import (
	"context"
	"log"
	pb "github.com/abhishek-kamat-nutanix/read-write-grpc/backup/proto"


)
func (s *Server)SendName(ctx context.Context,in *pb.NameRequest) (*pb.NameResponse, error) {
	volumeName = in.Name
	log.Printf("received volume name : %v \n",volumeName)
	return &pb.NameResponse{Message: "volume name set"}, nil
}