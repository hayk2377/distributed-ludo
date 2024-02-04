package main

import (
	"context"
	"fmt"

	pb "github.com/hayk2377/distributed-ludo/rpc/loadbalancer"
	"google.golang.org/grpc"
)

var ip string = "127.0.0.1:124"

func main() {
	// Connect to the gRPC server
	fmt.Println("going to make the three rpc calls")
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	fmt.Println("successfull rpc connection")
	if err != nil {
		fmt.Println("Error connecting to gRPC server:", err)
		return
	}
	defer conn.Close()

	// Create a gRPC client
	client := pb.NewLoadBalancerClient(conn)

	// NewServer request
	newServerReq := &pb.ServerRequest{
		Ip:       ip,
		Password: "whatever",
	}
	fmt.Println("calling the 'NewServer' rpc. this is to register a new worker node.")
	newServerRes, err := client.NewServer(context.Background(), newServerReq)
	if err != nil {
		fmt.Println("Error calling NewServer:", err)
		return
	}
	fmt.Println("newServer response:", newServerRes.ServerId)

	// HeartBeat request
	heartBeatReq := &pb.Heartreq{
		Status: ip,
	}
	fmt.Println("calling the 'HeartBeat' rpc. this is to update the server health")
	heartBeatRes, err := client.HeartBeat(context.Background(), heartBeatReq)
	if err != nil {
		fmt.Println("Error calling HeartBeat:", err)
		return
	}
	fmt.Println("heartBeat response:", heartBeatRes.ServerId)

	// Notify request
	notifyReq := &pb.NotifyRequest{
		GameId:   "345",
		ServerIp: ip,
	}
	fmt.Println("calling the 'Notify' rpc. it is to check if loadbalancer assigned the game with 'gameId' to this workernode.")
	notifyRes, err := client.Notify(context.Background(), notifyReq)
	if err != nil {
		fmt.Println("Error calling Notify:", err)
		return
	}
	fmt.Println("notify response:", notifyRes.ServerId)
}
