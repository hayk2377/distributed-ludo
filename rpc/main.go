package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	pb "github.com/hayk2377/distributed-ludo/rpc/LoadBalancer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
var (
	servers = []string{"192.168.1.1", "192.168.1.2", "192.168.1.3", "192.168.1.4", "192.168.1.5"}
	queue   = make([]string, len(servers))
	mutex   sync.Mutex
	games map[string]string
	gamesMutex sync.Mutex
	
)
type Response struct {
	Status string `json:"status"`
	Message string `json:"Message,omitempty"`
}
type LoadBalancer struct{
	pb.UnimplementedLoadBalancerServer
}
func (l *LoadBalancer) NewServer(ctx context.Context, req *pb.ServerRequest) (*pb.ServerResponse, error) {
    fmt.Println("Received newServer request from IP:", req.Ip)
	if req.Password!="whatever"{
		return nil, status.Error(codes.InvalidArgument,"you are blacklisted!")
	}
	mutex.Lock()
	queue=append(queue, req.Ip)
	mutex.Unlock()
    return &pb.ServerResponse{ServerId: "1234"}, nil
}


func (l *LoadBalancer) HeartBeat(ctx context.Context, req *pb.Heartreq) (*pb.ServerResponse, error) {
	value:=req.Status
    fmt.Println("Received heartBeat request from IP:", value)
	mutex.Lock()
	if len(queue)==1{
		mutex.Unlock()
		return &pb.ServerResponse{ServerId: "ok"}, nil
	}
	for i := range queue {
		if queue[i] == value {
			// If the value is found, remove it from the current position
			queue = append(queue[:i], queue[i+1:]...)
			// Insert the value at the 0 index
			queue = append([]string{value}, queue...)
		}
	}
	mutex.Unlock()
    return &pb.ServerResponse{ServerId: "ok"}, nil
}


func (l *LoadBalancer) Notify(ctx context.Context, req *pb.NotifyRequest) (*pb.ServerResponse, error) {
	gameid:=req.GameId
	serverip:=req.ServerIp
    fmt.Printf("Received notify request for GameID: %s, ServerIP: %s\n", gameid, serverip)
	gamesMutex.Lock()
	ip,thereIs:=games[gameid]
	if !thereIs ||ip!=serverip {
		return nil,status.Error(codes.Unauthenticated,"loadbalancer didnt assign that")
	}
    return &pb.ServerResponse{ServerId: "ok"}, nil
}

func startrpc(){
lis, err := net.Listen("tcp",":1234")
if err!=nil{
	log.Fatal("error listnenig on port 1234")
}
serverReg :=grpc.NewServer()
service:=&LoadBalancer{}
pb.RegisterLoadBalancerServer(serverReg,service)
fmt.Println("rpc server has started on port 1234")
err=serverReg.Serve(lis)
if err!=nil{
	log.Fatal("error when serving")
}

}
func gameHandler(w http.ResponseWriter, r *http.Request) {

	// Check if the request has a game ID
	gameID := r.FormValue("gameid")
	if gameID == "" {
		errorResponse(w, http.StatusBadRequest, "Missing game ID")
		return
	}
	mutex.Lock()
	IP := queue[0]
	queue = append(queue[1:], IP)
	mutex.Unlock()
	// Check if the game ID is in the games map
	gamesMutex.Lock()
	serverIP, exists := games[gameID]

	if exists {
		IP=serverIP
	}else{
	games[gameID] = IP}
	gamesMutex.Unlock()
	go func() {
		// Wait for thirty seconds
		time.Sleep(30 * time.Second)

		// remove the record.
		removeGameId(gameID)
	}()
	response := Response{
		Status: "success",
		Message: IP,
	}
	jsonResponse(w, response)
}
func removeGameId(game string){
	gamesMutex.Lock()
	delete(games,game)
	gamesMutex.Unlock()
}
func newConnect(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()

	if len(queue) == 0 {
		// No available servers in the queue
		http.Error(w, "No available servers", http.StatusServiceUnavailable)
		return
	}

	// Get the server IP from the front of the queue
	serverIP := queue[0]
	queue = append(queue[1:], serverIP)

	// Unlock the mutex before setting up the reverse proxy
	mutex.Unlock()

	// Create a new URL with the server IP
	serverURL, _ := url.Parse("http://" + serverIP)

	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(serverURL)

	// Serve the request using the reverse proxy
	proxy.ServeHTTP(w, r)
}

func jsonResponse(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func errorResponse(w http.ResponseWriter, status int, Message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := Response{
		Status: "error",
		Message: Message,
	}
	json.NewEncoder(w).Encode(response)
}
func startHTTP() {
	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		go newConnect(w, r)
	})

	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		go gameHandler(w, r)
	})

	// Start the load balancer server in a goroutine
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println("Error starting load balancer server:", err)
		}
	}()

	fmt.Println("Load balancer server started on :8080")

	// Block the main goroutine to keep the server running
	select {}
}

func main(){
	fmt.Println("starting server")
	go startrpc()
	startHTTP()
}