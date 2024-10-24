package main

import (
	"log"
	"flag"
	
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/abhishek-kamat-nutanix/read-write-grpc/backup/proto"
	   
)


func main() {
	address := flag.String("addr","10.46.60.86:50051","writer server ip")
	flag.Parse()
	conn, err := grpc.NewClient(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err!=nil {
		log.Fatalf("Failed to connect %v\n", err)
	}

	defer conn.Close()

	c := pb.NewBackupServiceClient(conn)

	doBackup(c)
}