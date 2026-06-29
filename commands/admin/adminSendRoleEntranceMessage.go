package admin

import (
	"GU/refs"
	"GU/utils"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (c *AdminSendRoleEntranceMessageCommand) CreateCommand() []*discordgo.ApplicationCommand {

	dc := []*discordgo.ApplicationCommand{
		{
			Name:        "a-send-roll-entrance",
			Description: "AdminTestMessage",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "old-message",
					Description: "Enter old entrance message(if any)",
					Required:    false,
				},
			},
		},
	}

	return dc
}

type AdminSendRoleEntranceMessageCommand struct{}

func (c *AdminSendRoleEntranceMessageCommand) Execute(s *discordgo.Session, _ *discordgo.InteractionCreate) string {
	embed := &discordgo.MessageEmbed{
		Title:       "その他",
		Description: "プライベートカテゴリー",
		Fields:      []*discordgo.MessageEmbedField{},
		Color:       0x500000,
	}
	if refs.Config.RoleEntranceMessageID != "" {
		err := s.ChannelMessageDelete(refs.Config.RoleEntranceChannelID, refs.Config.RoleEntranceMessageID)
		if err != nil {
			utils.Log(err, "", "adminRollEntranceMessage")
		}
	}
	msg, err := s.ChannelMessageSendEmbed(refs.Config.RoleEntranceChannelID, embed)
	if err != nil {
		utils.Log(err, "", "adminSendRollEntranceMessageCommand")
		return "Failed"
	}
	utils.Log(nil, fmt.Sprintf("MessageID : %s\nChannelID : %s", msg.ID, msg.ChannelID), "AdminSendRoleEntranceMessage")
	refs.Config.RoleEntranceMessageID = msg.ID
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
		err = s.MessageReactionAdd(refs.Config.RoleEntranceChannelID, msg.ID, refs.PrivateCategories[key])
		if err != nil {
			utils.Log(err, "", "adminSendRollEntranceMessageCommand")
		}
	}
	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:      msg.ID,
		Channel: refs.Config.RoleEntranceChannelID,
		Embeds: &[]*discordgo.MessageEmbed{
			embed,
		},
	})
	return "Success"
}
