package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	cli "github.com/jawher/mow.cli"

	"github.com/raoulh/binky-nfc/app"
)

var (
	macAddress *string
)

func exit(err error, exit int) {
	fmt.Fprintln(os.Stderr, err)
	cli.Exit(exit)
}

func handleSignals() {
	sigint := make(chan os.Signal, 1)

	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-sigint

	log.Println("Shuting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		exit(err, 1)
	}
}

func main() {

	cApp := cli.App("binky-nfc", "Binky radio software")
	cApp.Spec = "[-m]"

	macAddress = cApp.String(cli.StringOpt{
		Name:   "m mac",
		Desc:   "Set MAC address",
		EnvVar: "MAC_ADDRESS",
	})

	cApp.Action = func() {
		if err := app.Init(*macAddress); err != nil {
			exit(err, 1)
		}

		app.Run()

		// This will block until a signal is received
		handleSignals()
	}

	if err := cApp.Run(os.Args); err != nil {
		exit(err, 1)
	}
}
