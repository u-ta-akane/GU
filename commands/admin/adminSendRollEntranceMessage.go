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
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "old message",
					Description: "Enter old entrance message(if any)",
					Required:    false,
				},
			},
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
	if refs.Config.RollEntranceMessageID != "" {
		err := s.ChannelMessageDelete(refs.Config.RollEntranceChannelID, refs.Config.RollEntranceMessageID)
		if err != nil {
			utils.Log(err, "", "adminRollEntranceMessage")
		}
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
	chs, _ := s.GuildChannels(refs.Config.GuildID)
	for key, _ := range refs.PrivateCategories {
		for _, ch := range chs {
			if ch.ID == key {
				embed.Description += "\n:" + refs.PrivateCategories[key] + ": : " + ch.Name
			}
		}
		err = s.MessageReactionAdd(refs.Config.RollEntranceChannelID, msg.ID, refs.PrivateCategories[key])
		if err != nil {
			utils.Log(err, "", "adminSendRollEntranceMessageCommand")
		}
	}
	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:      msg.ID,
		Channel: refs.Config.RollEntranceChannelID,
		Embeds: &[]*discordgo.MessageEmbed{
			embed,
		},
	})
	return "Success"
}
