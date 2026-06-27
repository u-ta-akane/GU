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
			Name:        "add-priv-category",
			Description: "プライベートカテゴリーを作成します",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "作成したいカテゴリーのタイトルを設定してください",
					Required:    true,
				},
			},
		},
	}

	return dc
}

func (c *AddPrivateCategoryCommands) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) string {
	cat, err := s.GuildChannelCreateComplex(
		refs.Config.GuildID,
		discordgo.GuildChannelCreateData{
			Name: "priv-" + i.ApplicationCommandData().Options[0].StringValue(),
			Type: discordgo.ChannelTypeGuildCategory,
			PermissionOverwrites: []*discordgo.PermissionOverwrite{
				{
					ID:   refs.Config.GuildID, // @everyone
					Type: discordgo.PermissionOverwriteTypeRole,
					Deny: discordgo.PermissionViewChannel,
				},
				{
					ID:    i.Member.User.ID,
					Type:  discordgo.PermissionOverwriteTypeMember,
					Allow: refs.PrivateCategoryMemberPermission,
				},
			},
		},
	)
	if err != nil {
		utils.Log(err, "", "addPrivateCategory")
		return "Error occurred when creating private category"
	}
	refs.PrivateCategories[cat.ID] = ""
	return "Success"
}
