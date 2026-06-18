package commands

import (
	"GU/refs"
	"GU/utils"

	"github.com/bwmarrin/discordgo"
)

type AddPrivateCategoryCommands struct{}

func (c *AddPrivateCategoryCommands) CreateCommand() []*discordgo.ApplicationCommand {
	dc := []*discordgo.ApplicationCommand{
		{
			Name:        "a-delete-messages",
			Description: "AdminDeleteMessages",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "from",
					Description: "Enter Message ID of Start Point",
					Required:    true,
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
				},
			},
		},
	}

	return dc
}

func (c *AddPrivateCategoryCommands) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) string {
	from := i.ApplicationCommandData().Options[0].StringValue()
	num := i.ApplicationCommandData().Options[1].IntValue()
	channel := refs.Config.ModeratorChannelID
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "channel" {
			channel = opt.ChannelValue(s).ID
		}
	}
	err := utils.DeleteMessages(channel, from, int(num), "Use of AdminDeleteMessages", s)
	if err != 0 {
		return "Error occurred when deleting messages"
	}
	return "Delete Messages Success"
}
