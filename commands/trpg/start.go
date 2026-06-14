package trpg

import (
	"GU/apps"

	"github.com/bwmarrin/discordgo"
)

func (c *TrpgStartCommand) CreateCommand() []*discordgo.ApplicationCommand {

	dc := []*discordgo.ApplicationCommand{
		{
			Name:        "start",
			Description: "TRPGセッションを開始します",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "ボイスチャンネル",
					Description: "このセッションで使うメインとなるボイスチャンネルを指定してください",
					Required:    true,
				},
			},
		},
	}

	return dc
}

type TrpgStartCommand struct{}

func (c *TrpgStartCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) string {
	vc := i.ApplicationCommandData().Options[0].ChannelValue(s)
	apps.NewRoom(s, i, vc)
	return "Success"
}
