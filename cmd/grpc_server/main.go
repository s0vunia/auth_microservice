package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/s0vunia/auth_microservices_course_boilerplate/pkg/auth_v1"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedAuthV1Server
}

// Get user by id
func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())
	log.Printf("ctx: %+v", ctx)
	return &desc.GetResponse{
		User: &desc.User{
			Id: req.Id,
			Info: &desc.UserInfo{
				Name:  gofakeit.BeerName(),
				Email: gofakeit.Email(),
				Role:  1,
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

// Create user
func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("User email: %s", req.GetInfo().GetEmail())
	log.Printf("ctx: %+v", ctx)
	return &desc.CreateResponse{
		Id: int64(gofakeit.Number(0, 100)),
	}, nil
}

// Update user credentials
func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("User email: %s", req.GetInfo().GetEmail())
	log.Printf("ctx: %+v", ctx)
	return nil, nil
}

// Delete user
func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("User id: %d", req.GetId())
	log.Printf("ctx: %+v", ctx)
	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
