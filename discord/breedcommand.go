package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/Fl0GUI/discat/catapi"
)

var breedMap = make(map[string]string)

var breedCommand = &discordgo.ApplicationCommand{
	Name:        "breeds",
	Description: "get all breeds",
}

func getBreedCommand() (result *discordgo.ApplicationCommand, err error) {
	var breedsReq catapi.BreedRequest
	breeds, err := breedsReq.Execute()
	if err != nil {
		return
	}

	for _, b := range breeds {
		breedMap[b.Name] = b.Id
	}

	return breedCommand, nil
}

func handleBreed(s *discordgo.Session, m *discordgo.InteractionCreate) {
	if m.Interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}

	if m.ApplicationCommandData().ID != breedCommand.ID {
		return
	}

	content := strings.Builder{}
	content.WriteString("We've got: ")
	sep := ""
	for name := range breedMap {
		content.WriteString(fmt.Sprintf("%s%s", sep, name))
		sep = ", "
	}

	data := discordgo.InteractionResponseData{
		Content: content.String(),
	}
	resp := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &data,
	}
	s.InteractionRespond(m.Interaction, &resp)
}
