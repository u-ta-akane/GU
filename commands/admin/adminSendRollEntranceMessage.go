package admin

import (
	"GU/refs"
	"GU/utils"

	"github.com/bwmarrin/discordgo"
)

func (c *AdminSendRollEntranceMessageCommand) CreateCommand() []*discordgo.ApplicationCommand {

	dc := []*discordgo.ApplicationCommand{
		{
			Name:        "a-resend-roll-entrance",
			Description: "AdminTestMessage",
			Options:     []*discordgo.ApplicationCommandOption{},
		},
	}

	return dc
}

type AdminSendRollEntranceMessageCommand struct{}

func (c *AdminSendRollEntranceMessageCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) string {
	embed := &discordgo.MessageEmbed{
		Title:       "その他",
		Description: "プライベートカテゴリー",
		Fields:      []*discordgo.MessageEmbedField{},
		Color:       0x500000,
	}
	msg, err := s.ChannelMessageSendEmbed(refs.Config.RollEntranceChannelID, embed)
	if err != nil {
		utils.Log(err, "", "adminSendRollEntranceMessageCommand")
		return "Failed"
	}
	refs.Config.RollEntranceMessageID = msg.ID
	err = utils.JSONFM.Write("config.json", refs.Config)
	if err != nil {
		utils.Log(err, "", "adminSendRollEntranceMessageCommand")
		return "IO Error"
	}
	return "Success"
}
