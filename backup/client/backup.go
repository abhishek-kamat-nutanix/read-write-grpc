package main

import (
	"context"
	"log"

	pb "github.com/abhishek-kamat-nutanix/read-write-grpc/backup/proto"
)

func doBackup(c pb.BackupServiceClient) {
	log.Println("doBackup was invoked")

	reqs := HandleBlock("/dev/loop0")

	stream, err := c.BackupBlock(context.Background())

	if err != nil {
		log.Fatalf("error while calling BackupBlock : %v", err)
	}

	for _, req := range reqs{
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()

	if err!=nil {
		log.Fatalf("Error while receiving response from BackupBlock: %v\n", err)
	}

	log.Printf("BackupBlock: %v\n",res.Result)
}