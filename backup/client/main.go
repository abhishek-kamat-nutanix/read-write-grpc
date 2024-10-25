package main

import (
	"log"
	"os"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/abhishek-kamat-nutanix/read-write-grpc/backup/proto"
	   
)

func main() {

	addr := os.Getenv("GRPC_SERVER_ADDR")
    if addr == "" {
        log.Fatalf("GRPC_SERVER_ADDR environment variable is not set")
    }

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err!=nil {
		log.Fatalf("Failed to connect %v\n", err)
	}

	defer conn.Close()

	c := pb.NewBackupServiceClient(conn)

	doBackup(c)
}