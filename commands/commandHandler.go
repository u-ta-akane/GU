package commands

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

func adminHandler(s *discordgo.Session, i *discordgo.InteractionCreate, cmds *[refs.NumberOfCommands]Command, index int) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	switch index {
	case refs.IndexAdminTestMessage:
		result, e := utils.HasAuthority(s, i, refs.AuthoritySendAdminMessage)
		response := "Authorization Error"
		if e == 0 && result {
			response = (cmds[refs.IndexAdminTestMessage]).Execute(s, i)
		}
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource, // 「通常の返答」タイプ
			Data: &discordgo.InteractionResponseData{
				Content: response,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			utils.Log(err, "", "adminHandler")
			return
		}
		break
	case refs.IndexAdminDeleteMessages:
		result, e := utils.HasAuthority(s, i, refs.AuthorityControlMessages)
		response := "Authorization Error"
		if e == 0 && result {
			response = (cmds[refs.IndexAdminDeleteMessages]).Execute(s, i)
		}
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource, // 「通常の返答」タイプ
			Data: &discordgo.InteractionResponseData{
				Content: response,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			utils.Log(err, "", "adminHandler")
			return
		}
		break
	case refs.IndexAdminStopBot:
		result, e := utils.HasAuthority(s, i, refs.AuthorityBotManagement)
		response := "Authorization Error"
		if e == 0 && result {
			response = (cmds[refs.IndexAdminStopBot]).Execute(s, i)
		}
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource, // 「通常の返答」タイプ
			Data: &discordgo.InteractionResponseData{
				Content: response,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			utils.Log(err, "", "adminHandler")
			return
		}
		break
	case refs.IndexAdminReflashRoleData:
		result, e := utils.HasAuthority(s, i, refs.AuthorityReflashData)
		response := "Authorization Error"
		if e == 0 && result {
			response = (cmds[refs.IndexAdminReflashRoleData]).Execute(s, i)
		}
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource, // 「通常の返答」タイプ
			Data: &discordgo.InteractionResponseData{
				Content: response,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			utils.Log(err, "", "adminHandler")
			return
		}
		break
	}
}

func SetupCommands(dgs *discordgo.Session, cmds *[refs.NumberOfCommands]Command) {
	dgs.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}
		switch i.ApplicationCommandData().Name {
		case "a-test-message":
			adminHandler(s, i, cmds, refs.IndexAdminTestMessage)
			break
		case "a-delete-messages":
			adminHandler(s, i, cmds, refs.IndexAdminDeleteMessages)
			break
		case "a-stop-bot":
			adminHandler(s, i, cmds, refs.IndexAdminStopBot)
			break
		case "a-delete-role-data":
			adminHandler(s, i, cmds, refs.IndexAdminReflashRoleData)
			break
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
