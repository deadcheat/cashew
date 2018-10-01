package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/deadcheat/cashew/timer"

	"github.com/deadcheat/cashew/executor/expiration"
	"github.com/deadcheat/cashew/foundation"
	"github.com/kawasin73/htask/cron"
)

func main() {
	var err error
	var configFile string

	flag.StringVar(&configFile, "c", "config.yml", "specify config file path")
	flag.Parse()

	if err = foundation.PrepareApp(configFile); err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup
	workers := 1

	// get current time
	now := timer.Local.Now()
	startHour := now.Hour()
	startMin := now.Minute()
	// TODO make these logic to be an interface
	if startMin > 30 {
		startHour = startHour + 1
		startMin = 0
	} else {
		startMin = 30
	}
	// create start time of current hour
	startTime := time.Date(now.Year(), now.Month(), now.Day(), startHour, startMin, 0, 0, now.Location())
	log.Printf("cashew-isolator will start at %s \n", startTime.String())
	c := cron.NewCron(&wg, cron.Option{
		Workers: workers,
	})
	defer c.Close()
	e := expiration.New()

	// every n minute from start time
	// for example, if cmd was run at 19:20 and run every 30 minute,
	// 19:30 is a first run and every 30 minutes from then
	var cancel func()
	cancel, err = c.Every(foundation.App().ExpirationCheckInterval).Minute().From(startTime).Run(e.Execute)
	if err != nil {
		log.Fatal(err)
	}

	// wait signal to be sent
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)
	<-stop

	log.Println("received interrupt, process will be stopped soon")
	cancel()
}
