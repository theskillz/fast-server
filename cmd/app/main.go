package main

import (
	"context"
	"log"
	"sync"

	"github.com/theskillz/fast-server/internal/config"
	"github.com/theskillz/fast-server/internal/services/stats"
	"github.com/theskillz/fast-server/internal/web"
)

func main() {
	ctx := context.Background()
	cfg := config.NewConfig()
	statsSrv := stats.NewStats(cfg)
	webSrv := web.NewWebServer(cfg, statsSrv)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err := webSrv.Run(ctx); err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		if err := statsSrv.Run(ctx); err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()
	wg.Wait()
}
