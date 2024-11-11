package main

import (
	"context"
	"log"
	pb "github.com/abhishek-kamat-nutanix/read-write-grpc/backup/proto"


)
func (s *Server)SendName(ctx context.Context,in *pb.NameRequest) (*pb.NameResponse, error) {
	volumeName = in.Name
	namespace = in.Namespace
	log.Printf("received volume name : %v \n",volumeName)
	log.Printf("received namespace name : %v \n",namespace)
	return &pb.NameResponse{Message: "volume and namespace name set"}, nil
}