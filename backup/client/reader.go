package main

import (
    "os"
    "io"
    "log"

    pb "github.com/abhishek-kamat-nutanix/read-write-grpc/backup/proto"
)

func HandleBlock(devicePath string) []*pb.DataRequest {

    res := make([]*pb.DataRequest ,0)

    device, err := os.Open(devicePath)

    if err != nil {
        log.Printf("unable to open block device %s: %v", devicePath, err)
    }

    defer device.Close()

    // Buffer to hold data read from the block device
    buffer := make([]byte, 4096) // Adjust buffer size as needed

    for {
        n, err := device.Read(buffer)
        if err != nil {
            if err == io.EOF {
                break // End of file reached
            }
            log.Printf("error reading from block device: %v", err)
        }
        if n > 0 {
            data := make([]byte, n)
            copy(data, buffer[:n])

            res = append(res, &pb.DataRequest{Data: data})
        }
    }
    return res
}