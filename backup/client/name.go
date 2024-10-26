package main

import (
	"context"
	"log"

	pb "github.com/abhishek-kamat-nutanix/read-write-grpc/backup/proto"
)

func doSendName(c pb.BackupServiceClient) {
	log.Println("doSendName was invoked")

	req := &pb.NameRequest{Name: volumeName}

	res, err := c.SendName(context.Background(), req)

	if err!= nil {
		log.Fatalf("error while calling SendName : %v\n",err)
	}
	
	log.Printf("response from server: %v\n",res.Message)

}