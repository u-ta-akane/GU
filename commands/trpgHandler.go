package commands

import (
	"GU/refs"
	"GU/utils"

	"github.com/bwmarrin/discordgo"
)

func trpgHandler(s *discordgo.Session, i *discordgo.InteractionCreate, cmds *[refs.NumberOfCommands]Command, index uint8) {
	switch index {
	case refs.IndexTrpgStart:
		response := (cmds[refs.IndexTrpgSetMute]).Execute(s, i)
	case refs.IndexTrpgSetMute:
		result, e := utils.HasAuthority(s, i, refs.AuthorityVoiceChannelManagement)
		response := "Authorization Error"
		if e == 0 && result {
			response = (cmds[refs.IndexTrpgSetMute]).Execute(s, i)
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
