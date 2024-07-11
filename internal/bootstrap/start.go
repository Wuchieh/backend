package bootstrap

import (
	"log"
	"os"
	"os/signal"
	"string_backend_0001/internal/logger"
	"string_backend_0001/internal/web"
)

var (
	stop = make(chan os.Signal, 1)
)

func pause() {
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}

func Start() {
	defer Close()

	if err := Init(); err != nil {
		log.Fatal(err)
	}

	go func() {
		err := web.Run()
		if err != nil {
			logger.Error("web server error: %+v", err)
			return
		}
	}()

	pause()
}
