package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"os"

	pb "github.com/abhishek-kamat-nutanix/read-write-grpc/backup/proto"
)
var addr string = "0.0.0.0:50051"
var volumeName string
var namespace string
var snapClass string
var storageClassName string

type Server struct {
	pb.BackupServiceServer
}

func main() {

	snapClass = os.Getenv("SNAPCLASS")
    if addr == "" {
        log.Fatalf("SNAPCLASS environment variable is not set")
    }

	storageClassName = os.Getenv("STORAGECLASS")
    if addr == "" {
        log.Fatalf("STORAGECLASS environment variable is not set")
    }
	
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
