package discord

import (
	"discat/catapi"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

const choiceLimit = 25

var breedOption = discordgo.ApplicationCommandOption{
	Type:         discordgo.ApplicationCommandOptionString,
	Name:         "breed",
	Description:  "The cat you want",
	Required:     false,
	Autocomplete: true,
}

var catCommand = &discordgo.ApplicationCommand{
	Name:        "cat",
	Description: "get a cat",
	Options:     []*discordgo.ApplicationCommandOption{&breedOption},
}

func getCatCommand() (*discordgo.ApplicationCommand, error) {
	return catCommand, nil
}

func handleCat(s *discordgo.Session, m *discordgo.InteractionCreate) {
	if m.ApplicationCommandData().ID != catCommand.ID {
		return
	}

	switch m.Interaction.Type {
	case discordgo.InteractionApplicationCommand:
		respondWithCat(s, m)
	case discordgo.InteractionApplicationCommandAutocomplete:
		completeBreed(s, m)
	}
}

func respondWithCat(s *discordgo.Session, m *discordgo.InteractionCreate) {
	go deferMessage(s, m)

	var catRequest catapi.Request
	options := m.ApplicationCommandData().Options
	if len(options) == 1 {
		if brd, ok := breedMap[options[0].Value.(string)]; ok {
			catRequest.Breed = brd
		}
	}

	cats, err := catRequest.Execute()

	var resp discordgo.WebhookEdit
	var respData string
	if err != nil {
		respData = "Something went wrong :3"
		log.Println("Could not get some cats:", err)
	} else {
		respData = cats[0].Url
	}
	resp = discordgo.WebhookEdit{
		Content: &respData,
	}
	s.InteractionResponseEdit(m.Interaction, &resp)
}

func completeBreed(s *discordgo.Session, m *discordgo.InteractionCreate) {
	input := strings.ToLower(getBreedOption(m))
	choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, choiceLimit)

	for v, _ := range breedMap {
		if len(choices) >= choiceLimit {
			break
		}
		if fuzzy.Match(input, strings.ToLower(v)) {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  v,
				Value: v,
			})
		}
	}

	data := discordgo.InteractionResponseData{
		Choices: choices,
	}

	resp := discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &data,
	}
	s.InteractionRespond(m.Interaction, &resp)
}

func deferMessage(s *discordgo.Session, m *discordgo.InteractionCreate) {
	data := discordgo.InteractionResponseData{
		Content: ":3",
	}
	resp := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &data,
	}
	s.InteractionRespond(m.Interaction, &resp)
}

func getBreedOption(m *discordgo.InteractionCreate) string {
	options := m.ApplicationCommandData().Options
	if len(options) == 1 {
		return options[0].Value.(string)
	}
	return ""
}
