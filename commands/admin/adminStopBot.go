package admin

import (
	"GU/refs"
	"GU/utils"
	"os"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var SignalChannel = make(chan os.Signal, 1)

func (c *AdminStopBotCommand) CreateCommand() []*discordgo.ApplicationCommand {

	dc := []*discordgo.ApplicationCommand{
		{
			Name:        "a-stop-bot",
			Description: "AdminTestMessage",
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

type AdminStopBotCommand struct{}

func (c *AdminStopBotCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) string {
	txt := ""
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "text" {
			txt = opt.StringValue()
		}
	}
	res, status := utils.HasAuthority(s, i, refs.AuthorityBotManagement)
	if status != 0 {
		utils.Log(nil, "Authorize Error", "AdminStopBotCommand")
		return "Failed"
	}
	if res {
		utils.Log(nil, " Executing /a-shutdown : "+txt, "AdminStopBotCommand")
		SignalChannel <- os.Signal(syscall.SIGTERM)
	}
	return "Success"
}
