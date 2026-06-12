package admin

import (
	"GU/refs"
	"GU/utils"

	"github.com/bwmarrin/discordgo"
)

type AdminDeleteMessagesCommands struct{}

func (c *AdminDeleteMessagesCommands) CreateCommand() []*discordgo.ApplicationCommand {
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
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "number",
					Description: "Enter Number of Delete Target Messages",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "Enter Channel of Target",
					Required:    false,
				},
			},
		},
	}

	return dc
}

func (c *AdminDeleteMessagesCommands) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) string {
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
