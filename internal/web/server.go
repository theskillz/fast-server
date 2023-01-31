package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/theskillz/fast-server/internal/config"
	"github.com/theskillz/fast-server/internal/services/stats"
)

type WebServer interface {
	Run(ctx context.Context) error
}

type server struct {
	port  uint64
	stats stats.Stats
}

func NewWebServer(cfg *config.Config, statsSrv stats.Stats) WebServer {
	return &server{
		port:  cfg.Port,
		stats: statsSrv,
	}
}

func (s *server) Run(ctx context.Context) error {
	log.Printf("web server started on :%d", s.port)
	http.HandleFunc("/", s.handle)
	http.HandleFunc("/stats", s.handleStats)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

func (s *server) handle(w http.ResponseWriter, r *http.Request) {
	go func() {
		if err := s.stats.Add(r.RemoteAddr, r.UserAgent()); err != nil {
			log.Printf("%v", err)
			return
		}
	}()
}
func (s *server) handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	stats, err := s.stats.GetStatsForDay(time.Now())
	if err != nil {
		fmt.Fprintf(w, "some error: %v", err)
		w.WriteHeader(503)
		return
	}
	fmt.Fprint(w, "<pre>\ncount\tip_address\tuseragent\n")
	for i := range stats {
		fmt.Fprintf(w, "%d\t%s\t%s\n", stats[i].Count, stats[i].IPAddress, stats[i].Useragent)
	}
}
