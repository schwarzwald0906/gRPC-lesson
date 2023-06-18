package main

import (
	"context"
	"fmt"
	"grpc-lesson/pb"
	"io"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFileServiceServer
}

func (*server) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	fmt.Println("ListFiles was invoked")

	dir := "/Users/kimsukhong/git/grpc-lesson/storage"
	paths, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var filenames []string
	for _, path := range paths {
		if !path.IsDir() {
			filenames = append(filenames, path.Name())
		}
	}

	res := &pb.ListFilesResponse{
		Filenames: filenames,
	}
	return res, nil
}

func (*server) UploadAndNotifyProgress(stream pb.FileService_UploadAndNotifyProgressServer) error {
	fmt.Println("UploadAndNotifyProgress was invoked")

	size := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		data := req.GetData()
		log.Printf("received data: %v", data)
		size += len(data)

		res := &pb.UploadAndNotifyProgressResponse{
			Msg: fmt.Sprintf("received %v bytes", size),
		}
		err = stream.Send(res)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &server{})

	fmt.Println("server listening on")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
