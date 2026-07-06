package trpg

import (
	"GU/apps"
	"GU/utils"
	"fmt"

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
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildVoice,
					},
				},
			},
		},
	}

	return dc
}

type TrpgStartCommand struct{}

func (c *TrpgStartCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) string {
	vc := i.ApplicationCommandData().Options[0].ChannelValue(s)
	_, err := s.ChannelEdit(vc.ID, &discordgo.ChannelEdit{
		Name: fmt.Sprintf("セッション中 - %s - ", vc.Name),
	})
	if err != nil {
		utils.Log(err, "", "TrpgStart")
		return "Failed"
	}
	apps.NewRoom(s, i, vc)
	return "Success"
}
