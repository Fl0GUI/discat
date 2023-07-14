package discord

import (
	"log"
	"mime"
	"strings"

	"github.com/Fl0GUI/discat/catapi"
	"github.com/Fl0GUI/discat/custom"
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
		breed := options[0].Value.(string)
		if brd, ok := breedMap[breed]; ok {
			catRequest.Breed = brd
		}
		if custom.IsCustom(breed) {
			catRequest.Breed = breed
		}
	}

	if custom.IsCustom(catRequest.Breed) {
		customCat(s, m, &catRequest)
	} else {
		apiCat(s, m, &catRequest)
	}
}

func customCat(s *discordgo.Session, m *discordgo.InteractionCreate, catRequest *catapi.Request) {
	fileInfo := custom.GetCat(catRequest.Breed)
	file := custom.Open(catRequest.Breed, fileInfo)
	if file == nil {
		apiCat(s, m, catRequest)
		return
	}
	defer file.Close()
	dot := strings.LastIndex(fileInfo.Name(), ".")
	resp := discordgo.WebhookEdit{
		Files: []*discordgo.File{&discordgo.File{
			Name:        fileInfo.Name(),
			ContentType: mime.TypeByExtension(fileInfo.Name()[dot:]),
			Reader:      file,
		}},
	}

	s.InteractionResponseEdit(m.Interaction, &resp)
}

func apiCat(s *discordgo.Session, m *discordgo.InteractionCreate, catRequest *catapi.Request) {
	cats, err := catRequest.Execute()

	var resp discordgo.WebhookEdit
	var respData string
	if err != nil || len(cats) == 0 {
		respData = "Something went wrong :3"
		log.Printf("Could not get some cats: %s or %v\n", err, len(cats))
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

	addIfmatch := func(option string) {
		if fuzzy.Match(input, strings.ToLower(option)) {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  option,
				Value: option,
			})
		}
	}

	for _, v := range custom.Breeds() {
		if len(choices) >= choiceLimit {
			break
		}
		addIfmatch(v)
	}

	for v, _ := range breedMap {
		if len(choices) >= choiceLimit {
			break
		}
		addIfmatch(v)
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
