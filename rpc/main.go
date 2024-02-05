package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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
	queue      = []string{}
	mutex      sync.Mutex
	games      = make(map[string]string)
	gamesMutex sync.Mutex
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"Message,omitempty"`
}
type LoadBalancer struct {
	pb.UnimplementedLoadBalancerServer
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (l *LoadBalancer) NewServer(ctx context.Context, req *pb.ServerRequest) (*pb.ServerResponse, error) {
	value := req.Ip
	fmt.Println("Received newServer request from IP:", value)
	if req.Password != "whatever" {
		return nil, status.Error(codes.InvalidArgument, "you are blacklisted!")
	}
	mutex.Lock()
	// exists := false
	// fmt.Println("Servers:", queue)
	// //Dont modify a server that is already in the queue
	// for _, ip := range queue {
	// 	fmt.Println("comparing", ip, req.Ip)
	// 	if ip == req.Ip {
	// 		exists = true
	// 		break
	// 	}
	// }

	// if !exists {
	// 	queue = append(queue, req.Ip)
	// }
	queue = append(queue, value)
	fmt.Println("Queue:", queue)
	mutex.Unlock()
	go removeIP(value)
	return &pb.ServerResponse{ServerId: "1234"}, nil
}

func (l *LoadBalancer) HeartBeat(ctx context.Context, req *pb.Heartreq) (*pb.ServerResponse, error) {
	value := req.Status
	fmt.Println("Received heartBeat request from IP:", value)
	mutex.Lock()
	// if len(queue) == 1 {
	// 	mutex.Unlock()
	// 	return &pb.ServerResponse{ServerId: "ok"}, nil
	// }
	// for i := range queue {
	// 	if queue[i] == value {
	// 		// If the value is found, remove it from the current position
	// 		queue = append(queue[:i], queue[i+1:]...)
	// 		// Insert the value at the 0 index
	// 		queue = append([]string{value}, queue...)
	// 	}
	// }
	queue = append(queue, value)
	mutex.Unlock()
	go removeIP(value)
	return &pb.ServerResponse{ServerId: "ok"}, nil
}
func removeIP(value string) {
	time.Sleep(time.Second)
	mutex.Lock()
	defer mutex.Unlock()
	if len(queue) == 1 {
		// mutex.Unlock()
		return
	}
	for i := range queue {
		if queue[i] == value {
			if len(queue) == i+1 {
				queue = queue[:i]
				return
			}
			queue = append(queue[:i], queue[i+1:]...)
			// mutex.Unlock()
			return
		}
	}
}
func (l *LoadBalancer) Notify(ctx context.Context, req *pb.NotifyRequest) (*pb.ServerResponse, error) {
	gameid := req.GameId
	serverip := req.ServerIp
	fmt.Printf("Received notify request for GameID: %s, ServerIP: %s\n", gameid, serverip)
	gamesMutex.Lock()
	defer gamesMutex.Unlock()
	ip, thereIs := games[gameid]
	if thereIs && ip != serverip {
		return nil, status.Error(codes.Unauthenticated, "loadbalancer didnt assign that")
	}
	games[gameid] = serverip
	return &pb.ServerResponse{ServerId: "ok"}, nil
}

func startrpc() {
	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("error listnenig on port 1234")
	}
	serverReg := grpc.NewServer()
	service := &LoadBalancer{}
	pb.RegisterLoadBalancerServer(serverReg, service)
	fmt.Println("rpc server has started on port 1234")
	err = serverReg.Serve(lis)
	if err != nil {
		log.Fatal("error when serving")
	}

}
func gameHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	fmt.Println("in game handler")

	// Check if the request has a game ID
	gameID := r.URL.Query().Get("gameId")
	if gameID == "" {
		errorResponse(w, http.StatusBadRequest, "Missing game ID")
		return
	}
	mutex.Lock()
	n := len(queue)

	if n == 0 {
		// No available servers in the queue
		http.Error(w, "No available servers", http.StatusServiceUnavailable)
		return
	}
	rand.Seed(time.Now().UnixNano()) //seed the rand number generator
	randomNumber := rand.Intn(n)     //get a random index
	IP := queue[randomNumber]
	//IP := queue[0]
	// if n != 1 {
	// 	queue = append(queue[1:], IP)
	// }
	mutex.Unlock()

	// Check if the game ID is in the games map
	gamesMutex.Lock()
	serverIP, exists := games[gameID]

	if exists {
		url := "https://" + serverIP + "/"
		_, err := http.Get(url)
		if err == nil {
			IP = serverIP
		}

	}
	games[gameID] = IP
	fmt.Println("Game ID:", gameID, "Server IP:", IP)
	gamesMutex.Unlock()

	response := Response{
		Status:  "success",
		Message: IP,
	}
	jsonResponse(w, response)
}

func newConnect(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	n := len(queue)
	if n == 0 {
		// No available servers in the queue
		http.Error(w, "No available servers", http.StatusServiceUnavailable)
		return
	}
	rand.Seed(time.Now().UnixNano()) //seed the rand number generator
	randomNumber := rand.Intn(n)     //get a random index
	IP := queue[randomNumber]
	// if n == 1 {
	// 	serverIP := queue[0]
	// 	mutex.Unlock()

	// 	serverURL, _ := url.Parse("http://" + serverIP + "/lobbies")
	// 	proxy := httputil.NewSingleHostReverseProxy(serverURL)

	// 	proxy.Director = func(req *http.Request) {
	// 		req.Header = r.Header
	// 		req.URL.Scheme = serverURL.Scheme
	// 		req.URL.Host = serverURL.Host
	// 	}

	// 	proxy.ServeHTTP(w, r)
	// 	return
	// }

	// Unlock the mutex before setting up the reverse proxy
	mutex.Unlock()

	// Create a new URL with the server IP
	fmt.Println("using server to server IP")
	serverURL, _ := url.Parse("http://" + IP + "/lobbies")

	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(serverURL)
	proxy.Director = func(req *http.Request) {
		req.Header = r.Header
		req.URL.Scheme = serverURL.Scheme
		req.URL.Host = serverURL.Host
	}

	// Serve the request using the reverse proxy
	proxy.ServeHTTP(w, r)
}

// Action can be from client or worker node
// client can login /signup
// worker node can verify if jwt token is valid
func Action(w http.ResponseWriter, r *http.Request, Act string) {
	authServerUrl, err := url.Parse("http://localhost:8000/" + Act) //this path will be changed
	fmt.Println(authServerUrl)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Internal server error")
		fmt.Println("Error parsing URL:", err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(authServerUrl)

	proxy.Director = func(req *http.Request) {
		req.Header = r.Header
		req.URL.Scheme = authServerUrl.Scheme
		req.URL.Host = authServerUrl.Host
	}

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
		Status:  "error",
		Message: Message,
	}
	json.NewEncoder(w).Encode(response)
}
func statusHandler(w http.ResponseWriter, r *http.Request) {
	responseString := "ok"
	res := Response{
		Status:  "ok",
		Message: responseString,
	}
	jsonResponse(w, res)
}
func startHTTP() {
	//new game creation
	http.HandleFunc("/lobbies", func(w http.ResponseWriter, r *http.Request) {
		// enableCors(&w)
		newConnect(w, r)
	})
	// trying to join the game and play
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {

		// enableCors(&w)
		gameHandler(w, r)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		Action(w, r, "login")
	})
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		Action(w, r, "signup")
	})
	//trying to verify the token
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		Action(w, r, "user")
	})
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		Action(w, r, "test")
	})
	http.HandleFunc("/status", statusHandler)

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

func main() {
	fmt.Println("starting server")
	go startrpc()
	startHTTP()
}
