package main

import (
	userAPI "auth/pkg/user_v1"
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
)

const grpcPort = 50051

type server struct {
	userAPI.UnimplementedUserAPIServer
}

func (s *server) Create(ctx context.Context, req *userAPI.CreateRequest) (*userAPI.CreateResponse, error) {
	log.Println(req.Name)
	return &userAPI.CreateResponse{
		Id: gofakeit.Int64(),
	}, nil
}

func (s *server) Get(ctx context.Context, req *userAPI.GetRequest) (*userAPI.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())
	now := timestamppb.Now()
	return &userAPI.GetResponse{
		Id:        req.GetId(),
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      userAPI.Role_admin,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
func (s *server) Update(ctx context.Context, req *userAPI.UpdateRequest) (*emptypb.Empty, error) {
	log.Println(req.Name)
	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *userAPI.DeleteRequest) (*emptypb.Empty, error) {
	log.Println(req.Id)
	return &emptypb.Empty{}, nil
}

func main() {
	// создаём слушатель (listener) для TCP соединений на порту, который указан в переменной `grpcPort`
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatal("Server fell down while listening...\n", err)
	}
	// создаём новый экземпляр gRPC сервера.
	// это объект, который будет обрабатывать входящие запросы от клиентов
	s := grpc.NewServer()
	// регистрируем рефлексию, что позволяет клиентам получить информацию о доступных методах и сервисах на сервере
	reflection.Register(s)
	// регистрируем конкретную реализацию API
	// вторым аргументом передаём структуру, которая должна содержать реализацию логики, необходимой для каждого метода API
	userAPI.RegisterUserAPIServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	// метод Serve запускает сервер и начинает обрабатывать входящие запросы
	if err = s.Serve(lis); err != nil {
		log.Fatal("Server fell down while serving...\n", err)
	}

}
