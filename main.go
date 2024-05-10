package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.homebank.kz/homebank-oauth/halykid-events/configuration"
	"git.homebank.kz/homebank-oauth/halykid-events/events"
	"git.homebank.kz/homebank-oauth/halykid-events/queue"
	"git.homebank.kz/homebank-oauth/halykid-events/server"
	"git.homebank.kz/homebank-oauth/halykid-events/service"
)

func main() {

	appLoger, err := service.NewAppLogger(&configuration.Config.Logger)
	if err != nil {
		panic(err)
	}

	if err := events.Init(appLoger); err != nil {
		panic(err)
	}

	if err := queue.Init(appLoger, configuration.Config.Queue); err != nil {
		panic(err)
	}

	srv, err := server.NewServer(appLoger, configuration.Config.WEBServer)
	if err != nil {
		panic(err)
	}

	startServerErrorCH := srv.Start()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err = <-startServerErrorCH:
		{
			panic(err)
		}
	case q := <-quit:
		{
			fmt.Printf("receive signal %s, stopping server...\n", q.String())
			appLoger.ServerInfo("main", fmt.Sprintf("receive signal %s, stopping server...\n", q.String()))
			if err = srv.Stop(); err != nil {
				fmt.Printf("stop server error: %s\n", err.Error())
				appLoger.ServerError("main", err.Error(), "stop server error")
			}
		}
	}
}
