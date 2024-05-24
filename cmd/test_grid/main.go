package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"

	"test_grid/internal"
	"test_grid/internal/battery"
	"test_grid/internal/grid"
	"test_grid/internal/http_wrapper"
)

const (
	defaultHTTPTimeout = 2 * time.Second
	nationalGridDomain = "https://api.carbonintensity.org.uk"
)

func main() {
	mainCtx, mainCtxCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer mainCtxCancel()

	batteryController := battery.NewController([]battery.Charger{battery.NewStandard("battery-a"), battery.NewStandard("battery-b")})
	httpClient := http_wrapper.NewClient(defaultHTTPTimeout)
	gridClient := grid.NewClient(nationalGridDomain, httpClient)
	processor := internal.NewProcessor(batteryController, gridClient, mainCtx, internal.DefaultActionMap)
	// Since there wasn't requirement :"don't use external libraries", it allows me to avoid ticker's part implementation.
	// briefly: NewTicker just once. calculate duration after each time based on time.Now + ticker.Reset after
	c := cron.New(cron.WithLocation(time.UTC))
	_, err := c.AddFunc("*/30 * * * *", processor.Run)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Start()

	select {
	case <-mainCtx.Done():
		fmt.Println("shutdown started")
		jobsContext := c.Stop()
		fmt.Println("wait all jobs")
		<-jobsContext.Done()
		fmt.Println("shutdown finished")
	}

}
