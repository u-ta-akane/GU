package admin

import (
	"GU/refs"

	"github.com/bwmarrin/discordgo"
)

type AdminReflashRoleDataCommand struct{}

func (c *AdminReflashRoleDataCommand) CreateCommand() []*discordgo.ApplicationCommand {
	dc := []*discordgo.ApplicationCommand{
		{
			Name:        "a-reflash-roledata",
			Description: "Reflash RoleData Map",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "text",
					Description: "Enter Text",
					Required:    false,
				},
			},
		},
	}

	return dc
}

func (c *AdminReflashRoleDataCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) string {
	txt := "Success"
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "text" {
			txt = opt.StringValue()
		}
	}
	refs.ReflashRoleData(s)
	return txt
}
