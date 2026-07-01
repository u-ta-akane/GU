package main

import (
	"GU/refs"
	"GU/utils"

	"github.com/bwmarrin/discordgo"
)

func AdminHandler(s *discordgo.Session, i *discordgo.InteractionCreate, cmds *[refs.NumberOfCommands]Command, index int) {
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
	case refs.IndexAdminSendRoleEntranceMessage:
		result, e := utils.HasAuthority(s, i, refs.AuthorityRoleEntranceManagement)
		response := "Authorization Error"
		if e == 0 && result {
			response = (cmds[refs.IndexAdminSendRoleEntranceMessage]).Execute(s, i)
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
	}
}
