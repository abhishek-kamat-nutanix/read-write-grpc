package main

import (
	"google.golang.org/grpc"
	"log"
	"net"

	pb "github.com/abhishek-kamat-nutanix/read-write-grpc/backup/proto"
)
var addr string = "0.0.0.0:50051"

type Server struct {
	pb.BackupServiceServer
}

func main() {
	lis, err := net.Listen("tcp", addr)

	if err!=nil {
		log.Fatalf("Failed to listen on: %v\n",err)
	}

	log.Printf("Listening on %s\n", addr)

	s:= grpc.NewServer()
	
	pb.RegisterBackupServiceServer(s, &Server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen on: %v\n",err)
	}

}
