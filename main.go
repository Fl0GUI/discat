package main

import (
	"context"
	"discat/discord"
	"flag"
	"log"
)

var reset = flag.Bool("reset", false, "reset the discord bot commands")
var register = flag.Bool("register", true, "register the discord bot commands")

func main() {
	flag.Parse()
	ctx := context.Background()
	dctx, cancel := context.WithCancel(ctx)
	defer cancel()
	c, err := discord.New(dctx)
	if err != nil {
		log.Fatal("Couldn't create a client: ", err)
	}

	if *reset {
		if err = c.ResetCommands(); err != nil {
      c.Close()
			log.Fatal("Couldn't delete commands: ", err)
		}
    return
	} else if *register {
		if err = c.RegisterCommands(); err != nil {
      c.Close()
			log.Fatal("Couldn't register commands: ", err)
		}
	}
	for {
	}
}
