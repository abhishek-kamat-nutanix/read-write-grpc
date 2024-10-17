package main

import (
	"log"
	
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/abhishek-kamat-nutanix/read-write-grpc/backup/proto"
	   
)

var addr string = "10.46.60.85:50051"

func main() {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err!=nil {
		log.Fatalf("Failed to connect %v\n", err)
	}

	defer conn.Close()

	c := pb.NewBackupServiceClient(conn)

	doBackup(c)
}