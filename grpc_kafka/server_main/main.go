package main

import (
	"github.com/GRPC/proto"
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type server struct{}

func main() {
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	proto.RegisterAddServiceServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}

}

func (s *server) AddtoKafka(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	a, b := request.GetUsername(), request.GetName()
	data := a + "&" + b
	fmt.Print(data)
	
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}

	// Produce messages to topic (asynchronously)
	topic := "data-topic"
	for _, word := range []string{string(data)} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: 0},
			Value:          []byte(word),
		}, nil)
	}

	return &proto.Response{Status: data}, nil
}
