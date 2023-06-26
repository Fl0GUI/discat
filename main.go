package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/Fl0GUI/discat/discord"
)

var reset = flag.Bool("reset", false, "reset the discord bot commands")
var register = flag.Bool("register", true, "register the discord bot commands")

func main() {
	flag.Parse()
	ctx := context.Background()
	dctx, cancel := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	c, err := discord.New(dctx)
	defer c.Close()
	if err != nil {
		log.Fatalf("Couldn't create a client: %s\n", err)
	}

	if *reset {
		if err = c.ResetCommands(); err != nil {
			log.Fatalf("Couldn't delete commands: %s\n", err)
		}
		return
	} else if *register {
		if err = c.RegisterCommands(); err != nil {
			log.Fatalf("Couldn't register commands: %s\n", err)
		}
	}

	for _ = range dctx.Done() {
	}
}
