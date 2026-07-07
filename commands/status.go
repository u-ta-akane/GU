package commands

import (
	"GU/refs"
	"GU/utils"

	"github.com/bwmarrin/discordgo"
)

type StatusCommand struct{}

func (c *StatusCommand) CreateCommand() []*discordgo.ApplicationCommand {
	dc := []*discordgo.ApplicationCommand{
		{
			Name:        "play",
			Description: "「乱入歓迎」ロールのオン/オフを切り替えます",
			Options:     []*discordgo.ApplicationCommandOption{},
		},
	}

	return dc
}

func (c *StatusCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) string {
	for _, role := range i.Member.Roles {
		if role == refs.Config.PlayingStatusRoleID {
			err := s.GuildMemberRoleRemove(refs.Config.GuildID, i.Member.User.ID, refs.Config.PlayingStatusRoleID)
			if err != nil {
				utils.Log(err, "", "StatusCommand")
				return "Failed"
			}
			return "Success"
		}
	}
	err := s.GuildMemberRoleAdd(refs.Config.GuildID, i.Member.User.ID, refs.Config.PlayingStatusRoleID)
	if err != nil {
		utils.Log(err, "", "StatusCommand")
		return "Failed"
	}
	return "Success"
}
