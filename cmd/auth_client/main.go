package main

import (
	userAPI "auth/pkg/user_v1"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const (
	address = "localhost:50051"
	userID  = 12
)

func main() {
	// создаём новое gRPC-соединение с сервером по указанному адресу
	// соединение устанавливается без шифрования (insecure)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to ", address, "; ", err)
	}
	defer conn.Close()
	// создаётся новый клиент для вызова gRPC-методов, определённых в `UserAPI`.
	// этот клиент получает соединение `conn` как параметр, что позволяет ему отправлять запросы на сервер через установленное соединение
	c := userAPI.NewUserAPIClient(conn)
	// управление временем ожидания запроса к gRPC-серверу. Если ответ не будет получен за указанное время, запрос будет отменён
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// вызов gRPC-метода `Get` на сервере, который описан в `UserAPI`
	r, err := c.Get(ctx, &userAPI.GetRequest{Id: userID})
	if err != nil {
		log.Fatal("Failed to get response")
	}
	log.Println(r)
}
