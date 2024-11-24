package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"net/http"
	_ "net/http/pprof"

	"github.com/GarryStalker/loadBalancer/internal/config"
	"github.com/GarryStalker/loadBalancer/internal/logger"
	"github.com/GarryStalker/loadBalancer/internal/server"
	lbv1 "github.com/GarryStalker/loadBalancer_protos/gen/go"
	"google.golang.org/grpc"
)

func main() {
	runtime.SetBlockProfileRate(1)
	cfg := config.MustLoad()

	log := logger.InitLogger(cfg.Env)
	log.Info("Starting LoadBalancer service", slog.String("Port", cfg.Port))

	grpcServer := grpc.NewServer()
	srv := server.New(cfg, log)
	lbv1.RegisterLoadBalancerServer(grpcServer, srv)

	go func() {
		http.ListenAndServe(":8080", nil)
	}()

	go func() {
		l, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
		if err != nil {
			log.Error("Failed to listen", slog.String("address", cfg.Port), slog.Any("error", err))
			return
		}

		log.Info("Server is running", slog.String("address", cfg.Port))
		if err := grpcServer.Serve(l); err != nil {
			log.Error("Failed to serve", slog.Any("error", err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	signalName := <-stop
	log.Info("stopping application", slog.String("signal", signalName.String()))
	grpcServer.Stop()
	log.Info("application stoped")

}
