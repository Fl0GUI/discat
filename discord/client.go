package discord

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	tokenEnv = "DISCORD_TOKEN"
	guildEnv = "DISCORD_GUILD"
	appEnv   = "DISCORD_APP"
)

type Client struct {
	*discordgo.Session

	toClose []func()

	appId, guildId string
}

var client *Client = nil

func New(ctx context.Context) (*Client, error) {
	if client != nil {
		return client, nil
	}

	client = &Client{}
	var err error

	token := os.Getenv(tokenEnv)
	if len(token) == 0 {
		return client, envError(tokenEnv)
	}
	client.appId = os.Getenv(appEnv)
	if len(client.appId) == 0 {
		return client, envError(appEnv)
	}
	client.guildId = os.Getenv(guildEnv)

	client.Session, err = discordgo.New("Bot " + token)
	if err != nil {
		return client, err
	}
	err = client.Open()
	if err != nil {
		return client, err
	}

	client.toClose = []func(){
		client.AddHandler(handleCat),
		client.AddHandler(handleBreed),
		func() {
			for client.Session.Close() != nil {
			}
		},
	}

	return client, err
}

func (c *Client) Close() {
	grp := sync.WaitGroup{}
	for _, f := range c.toClose {
		grp.Add(1)
		go func() {
			f()
			grp.Done()
		}()
	}
	grp.Wait()
}

func (c *Client) RegisterCommands() error {
	catCmd, err := getCatCommand()
	if err != nil {
		return err
	}
	catCommand, err = c.ApplicationCommandCreate(c.appId, c.guildId, catCmd)
	if err != nil {
		return err
	}

	breedCmd, err := getBreedCommand()
	if err != nil {
		return err
	}
	breedCommand, err = c.ApplicationCommandCreate(c.appId, c.guildId, breedCmd)
	return err
}

func envError(env string) error {
	return errors.New(fmt.Sprintf("%s not set", env))
}

func (c *Client) ResetCommands() error {
	commands, err := c.ApplicationCommands(c.appId, c.guildId)
	if err != nil {
		return err
	}
	grp := sync.WaitGroup{}
	grp.Add(len(commands))
	for _, comm := range commands {
		go func() {
			defer grp.Done()
			err := c.ApplicationCommandDelete(c.appId, c.guildId, comm.ID)
			if err != nil {
				fmt.Printf("Could not delete command %s: %s", comm.Name, err)
			}
		}()
	}
	grp.Wait()
	return nil
}
