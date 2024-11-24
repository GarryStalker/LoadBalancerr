package server

import (
	"context"
	"log/slog"
	"strings"

	"github.com/GarryStalker/loadBalancer/internal/config"
	"github.com/GarryStalker/loadBalancer/internal/service"
	lbv1 "github.com/GarryStalker/loadBalancer_protos/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoadBalancerServer struct {
	lbv1.UnimplementedLoadBalancerServer
	cfg    *config.Config
	router *service.Router
	log    *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) *LoadBalancerServer {
	return &LoadBalancerServer{
		cfg:    cfg,
		router: service.New(cfg.CDNHost, log),
		log:    log,
	}
}

func (s *LoadBalancerServer) Redirect(ctx context.Context, in *lbv1.Request) (*lbv1.Response, error) {
	videoURL := strings.TrimSpace(in.GetVideo())
	if videoURL == "" {
		return nil, status.Error(codes.InvalidArgument, "url is required")
	}

	targetURL, err := s.router.GetTargetURL(videoURL)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get target url")
	}
	return &lbv1.Response{Redirect: targetURL}, nil
}
