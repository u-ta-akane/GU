package main

import (
	"GU/refs"
	"GU/utils"

	"github.com/bwmarrin/discordgo"
)

// Command cmdsに追加するために満たすべきインターフェースです。
type Command interface {
	// CreateCommand discordgo.ApplicationCommandをmain.goに返します。
	CreateCommand() []*discordgo.ApplicationCommand
	// Execute コマンドが実行されたときに呼ばれる、処理の本体です。
	Execute(s *discordgo.Session, i *discordgo.InteractionCreate) string
}

func SetupCommands(dgs *discordgo.Session, cmds *[refs.NumberOfCommands]Command) {
	dgs.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}
		switch i.ApplicationCommandData().Name {
		case "a-test-message":
			AdminHandler(s, i, cmds, refs.IndexAdminTestMessage)
		case "a-delete-messages":
			AdminHandler(s, i, cmds, refs.IndexAdminDeleteMessages)
		case "a-stop-bot":
			AdminHandler(s, i, cmds, refs.IndexAdminStopBot)
		case "a-delete-role-data":
			AdminHandler(s, i, cmds, refs.IndexAdminReflashRoleData)
		case "start":
			TrpgHandler(s, i, cmds, refs.IndexTrpgStart)
		case "ゆるぼ":
			response := (cmds[refs.IndexAddYURUBO]).Execute(s, i)
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource, // 「通常の返答」タイプ
				Data: &discordgo.InteractionResponseData{
					Content: response,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				utils.Log(err, "", "SetupCommands")
				return
			}
			break
		case "delete-ゆるぼ":
			response := (cmds[refs.IndexDeleteYURUBO]).Execute(s, i)
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource, // 「通常の返答」タイプ
				Data: &discordgo.InteractionResponseData{
					Content: response,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				utils.Log(err, "", "SetupCommands")
				return
			}
			break
		}
	})
}
