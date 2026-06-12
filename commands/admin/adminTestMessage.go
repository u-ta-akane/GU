package admin

import (
	"GU/utils"
	"log"

	"github.com/bwmarrin/discordgo"
)

func (c *AdminTestMessageCommand) CreateCommand() []*discordgo.ApplicationCommand {

	dc := []*discordgo.ApplicationCommand{
		{
			Name:        "a-test-message",
			Description: "AdminTestMessage",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "text",
					Description: "Enter Text",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "Enter Channel",
					Required:    true,
				},
			},
		},
	}

	return dc
}

type AdminTestMessageCommand struct{}

func (c *AdminTestMessageCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) string {
	text := i.ApplicationCommandData().Options[0].StringValue()
	channel := i.ApplicationCommandData().Options[1].ChannelValue(s)
	log.Printf(text)
	utils.OrderSendMessage(channel.ID, text)
	return "Success"
}
