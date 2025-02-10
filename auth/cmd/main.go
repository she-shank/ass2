package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/she-shank/ass2/server"
)

func main() {
	authServer, err := server.NewRestApi()
	if err != nil {
		slog.Error(fmt.Sprint(err))
	}

	err = authServer.Start()
	if err != nil {
		slog.Error(fmt.Sprint(err))
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	err = authServer.Close()
	if err != nil {
		//TODO: add logger panic/fatal and exit
		slog.Error(fmt.Sprint(err))
	}
}
