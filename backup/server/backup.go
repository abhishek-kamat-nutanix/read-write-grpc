package main

import (
	"io"
	"log"
	"os"

	pb "github.com/abhishek-kamat-nutanix/read-write-grpc/backup/proto"
)
func (s *Server)BackupBlock(stream pb.BackupService_BackupBlockServer) error{
	log.Println("BackupBlock function was invoked")

	res := "Writer has completed writing"

	devicePath := "/dev/loop1"
	device, err := os.OpenFile(devicePath, os.O_WRONLY, 0644)

	if err != nil {
        log.Printf("unable to open block device %s: %v", devicePath, err)
    }

	for {
		req, err := stream.Recv()
		
		if err == io.EOF {
			return stream.SendAndClose(&pb.DataResponse{
				Result: res,
			})
		}
		
		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n",err)
		}

		n, err := device.Write(req.Data)
		if n > 0 {
			// log.Printf("wrote %v bytes on disk",n)
		}
		if err!= nil {
			log.Printf("error writing to disk : %v",err)
		}

	}
}
