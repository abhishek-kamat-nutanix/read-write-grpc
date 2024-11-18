package main

import (
	"context"

	"log"
	"os"
	"os/exec"

	pb "github.com/abhishek-kamat-nutanix/read-write-grpc/backup/proto"
)

func (s *Server)SendJSONData(ctx context.Context,in *pb.JSONDataRequest) ( *pb.DataResponse, error) {
	jsonData := in.Jsondata

	if len(jsonData) == 0 {
        log.Println("No JSON data provided")
    }

	tFile, err := os.CreateTemp("","configs-*.json")
	if err!=nil{
		log.Printf("error creating temp file: %v\n",err)
	}
	defer os.Remove(tFile.Name())

	byteData := []byte(jsonData)

	_, err = tFile.Write(byteData)
	if err != nil {
		log.Printf("failed to write to temp file: %v", err)
	}

	err = tFile.Close()
	if err != nil {
		log.Printf("failed to close temp file: %v", err)
	}

	//contents, _ := os.ReadFile(tFile.Name())
    //log.Printf("File contents for kubectl apply: %s\n", string(contents))

	cmd := exec.Command("kubectl","apply","-f",tFile.Name())
	output, err := cmd.CombinedOutput()
	if err != nil{
		log.Printf("kubectl apply failed with error %v\n %s",err,string(output))
	} else {log.Printf("kubectl applied %s",string(output))}
	
	return &pb.DataResponse{Result: string(output)},nil
}