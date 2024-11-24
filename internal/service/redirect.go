package service

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"sync/atomic"
)

type Router struct {
	cdnHost string
	counter atomic.Int32
	log     *slog.Logger
}

func New(cndHost string, log *slog.Logger) *Router {
	return &Router{
		cdnHost: cndHost,
		counter: atomic.Int32{},
		log:     log,
	}
}

func (r *Router) GetTargetURL(originURL string) (redirect string, err error) {
	requestsCount := r.counter.Add(1)

	if requestsCount%10 == 0 {
		r.counter.Store(0)
		r.log.Info("Redirect to origin", slog.String("target_url", originURL), slog.Any("counter", r.counter.Load()))
		return originURL, nil
	}

	targetURL, err := r.getCDNServer(originURL)
	if err != nil {
		r.log.Error("failed get CDN server", slog.String("err", err.Error()))
		return "", fmt.Errorf("failed get CDN server: %w", err)
	}

	r.log.Info("Redirect to CDN", slog.String("target_url", targetURL), slog.Any("counter", requestsCount))
	return targetURL, nil
}

func (r *Router) getCDNServer(originURL string) (string, error) {
	u, err := url.Parse(originURL)
	if err != nil {
		r.log.Error("Failed to parse URL, redirect to originURL", slog.String("url", originURL))
		return originURL, nil
	}

	cacheServer := strings.Split(u.Host, ".")[0]
	redirect := fmt.Sprintf("%s/%s%s", r.cdnHost, cacheServer, u.EscapedPath())

	return redirect, nil
}
