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
	counter int32
	log     *slog.Logger
}

func New(cndHost string, log *slog.Logger) *Router {
	return &Router{
		cdnHost: cndHost,
		counter: 0,
		log:     log,
	}
}

func (r *Router) GetTargetURL(originURL string) (redirect string, err error) {
	requestCount := atomic.AddInt32(&r.counter, 1)

	if requestCount%10 == 0 {
		r.counter = 0
		r.log.Info("Redirect to origin", slog.Any("counter", r.counter))
		return originURL, nil
	}

	u, err := url.Parse(originURL)
	if err != nil {
		r.log.Error("Failed to parse URL, redirect to originURL", slog.String("url", originURL))
		return originURL, nil
	}

	cacheServer := strings.Split(u.Host, ".")[0]
	redirect = fmt.Sprintf("%s/%s%s", r.cdnHost, cacheServer, u.EscapedPath())
	r.log.Info("Redirect to CDN", slog.Any("counter", r.counter))
	return redirect, nil
}
