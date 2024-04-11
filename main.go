package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/Tiger-Coders/tigerlily-payment/api/rest/router"
	"github.com/Tiger-Coders/tigerlily-payment/api/rpc"
	"github.com/Tiger-Coders/tigerlily-payment/internal/db"
	"github.com/Tiger-Coders/tigerlily-payment/internal/pkg/env"
	"github.com/Tiger-Coders/tigerlily-payment/internal/pkg/logger"
	"github.com/Tiger-Coders/tigerlily-payment/internal/service/checkout"
	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

func main() {
	logs := logger.NewLogger()
	logs.InfoLogger.Println("Starting up server ...")

	// Set ENV vars
	env.SetEnv()

	// Spin up the main server instance
	lis, err := net.Listen("tcp", ":8001")
	if err != nil {
		logs.ErrorLogger.Println("Something went wrong in the server startup")
		log.Fatalf("Error connecting tcp port 8000")
	}
	logs.InfoLogger.Println("Successfull server init")

	// Start a new multiplexer passing in the main server
	m := cmux.New(lis)

	// Listen for HTTP requests first
	// If request headers don't specify HTTP, next mux would handle the request
	httpListener := m.Match(cmux.HTTP1Fast())
	grpclistener := m.Match(cmux.Any())

	// Run GO routine to run both servers at diff processes at the same time
	go serveGRPC(grpclistener)
	go serveHTTP(httpListener)

	fmt.Printf("Payment Service Running@%v\n", lis.Addr())

	if err := m.Serve(); !strings.Contains(err.Error(), "use of closed network connection") {
		log.Fatalf("MUX ERR : %+v", err)
	}

}

// GRPC Server initialisation
func serveGRPC(l net.Listener) {
	grpcServer := grpc.NewServer()

	// Register GRPC stubs (pass the GRPCServer and the initialisation of the service layer)
	rpc.RegisterCheckoutServiceServer(grpcServer, checkout.NewCheckoutService(db.NewDB()))

	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("error running GRPC server %+v", err)
	}
}

// HTTP Server initialisation (using gin)
func serveHTTP(l net.Listener) {
	h := gin.Default()
	router.Router(h)
	s := &http.Server{
		Handler: h,
	}
	if err := s.Serve(l); err != cmux.ErrListenerClosed {
		log.Fatalf("error serving HTTP : %+v", err)
	}
}