/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"log"
	"net"

	pb "grpc_test"
	//google.golang.org/grpc/examples/helloworld/helloworld

	"database/sql"

	"google.golang.org/grpc"

	_ "github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
)

const (
	port   = ":50051"
	DB_DSN = "postgres://postgres:gnsecret@gononeterp.cud7jbsftjfi.ap-southeast-1.rds.amazonaws.com:5432/choukash_erp_1_0_0"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedUserserviceServer
}
type User struct {
	tenant_id       string
	tenant_email    string
	tenant_username string
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {

	db, err := sql.Open("postgres", DB_DSN)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	var myUser User
	userSql := "SELECT tenant_id, tenant_username, tenant_email FROM imraan_auth.tenants"

	err = db.QueryRow(userSql).Scan(&myUser.tenant_id, &myUser.tenant_username, &myUser.tenant_email)
	if err != nil {
		log.Fatal("Failed to execute query: ", err)
	}

	log.Printf("Received: %v", in.GetName())
	log.Printf("Hi %s, welcome back!\n", myUser.tenant_username)
	log.Printf("Hi your email is: %s \n", myUser.tenant_email)
	return &pb.HelloReply{Message: "Hello Mr. " + in.GetName()}, nil
}

func main() {
	//db connecion

	//grpc part

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserserviceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
