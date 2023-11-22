package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/grantjforrester/go-ticket/app"

	"github.com/spf13/viper"
)

func main() {
	config := viper.New()
	config.AutomaticEnv()
	app := app.NewApp(config)
	app.Start()

    sig := make(chan os.Signal, 1)
    done := make(chan bool, 1)

    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        sig := <-sig
        log.Printf("Received signal: %s", sig)
		app.Stop()
        done <- true
    }()

    log.Println("Started")
    <-done
    log.Println("Exiting...")
}
