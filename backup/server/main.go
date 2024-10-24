package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"flag"

	pb "github.com/abhishek-kamat-nutanix/read-write-grpc/backup/proto"
)
var addr string = "0.0.0.0:50051"
var kubeconfig string
var volumeName string

type Server struct {
	pb.BackupServiceServer
}

func main() {

	flag.StringVar(&kubeconfig, "kubeconfig", "/home/nutanix/nke-target.cfg", "location to your kubeconfig file")
	flag.StringVar(&volumeName, "pv", "migrated-pv", "Persistent Volume Claim name for backup")
	flag.Parse()
	
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
