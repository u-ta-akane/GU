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
			Name:        "a-send-role-entrance",
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

func (c *AdminSendRoleEntranceMessageCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) string {
	embed := &discordgo.MessageEmbed{
		Title:       "プライベートカテゴリー",
		Description: "",
		Fields:      []*discordgo.MessageEmbedField{},
		Color:       0x500000,
	}
	if refs.Config.RoleEntranceMessageID != "" {
		err := s.ChannelMessageDelete(refs.Config.RoleEntranceChannelID, refs.Config.RoleEntranceMessageID)
		if err != nil {
			utils.Log(err, "", "adminRollEntranceMessage")
		}
	}
	for _, data := range i.ApplicationCommandData().Options {
		if data.Name == "old-message" {
			if data.Value != "" {
				err := s.ChannelMessageDelete(refs.Config.RoleEntranceChannelID, data.Value.(string))
				if err != nil {
					utils.Log(err, "", "adminRollEntranceMessage")
				}
			}
		}
	}
	msg, err := s.ChannelMessageSendEmbed(refs.Config.RoleEntranceChannelID, embed)
	if err != nil {
		utils.Log(err, "", "adminSendRollEntranceMessageCommand")
		return "Failed"
	}
	utils.Log(nil, fmt.Sprintf("MessageID : %s\nChannelID : %s\nChannelName : %s", msg.ID, msg.ChannelID), "AdminSendRoleEntranceMessage")
	refs.Config.RoleEntranceMessageID = msg.ID
	err = utils.JSONFM.Write("config.json", refs.Config)
	if err != nil {
		utils.Log(err, "", "adminSendRollEntranceMessageCommand")
		return "IO Error"
	}
	chs, _ := s.GuildChannels(refs.Config.GuildID)
	for _, cat := range refs.PrivateCategories {
		for _, ch := range chs {
			if ch.ID == cat.CategoryID {
				embed.Description += "\n" + cat.Emoji + " : " + ch.Name
				err = s.MessageReactionAdd(refs.Config.RoleEntranceChannelID, msg.ID, cat.Emoji)
				if err != nil {
					utils.Log(err, "", "adminSendRollEntranceMessageCommand")
				}
			}
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
