package main

import (
	"GU/apps"
	"GU/refs"
	"GU/utils"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var removeOnReactionAddHandler func()

// サーバーにユーザーが参加した時のハンドラー
func onMemberAdd(dgs *discordgo.Session, m *discordgo.GuildMemberAdd) {
	err := dgs.GuildMemberRoleAdd(
		m.GuildID,
		m.User.ID,
		refs.Config.DefaultAuthorityID,
		discordgo.WithAuditLogReason("Set Default Authority"),
	)
	if err != nil {
		utils.ErrorChannel <- err
	}
}

func onReactionAdd(dgs *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.MessageID == refs.Config.RollEntranceMessageID {
		for key, value := range refs.PrivateCategories {
			if m.Emoji.Name == value {
				apps.IOPrivateCategoryMember(dgs, m.Member.User.ID, key)
			}
		}
		err := dgs.MessageReactionRemove(refs.Config.RollEntranceChannelID, refs.Config.RollEntranceMessageID, m.Emoji.ID, m.Member.User.ID)
		if err != nil {
			utils.Log(err, "", "onReactionAdd")
		}
	}
}

func setupOnInteractionHandler(dgs *discordgo.Session, cmds *[refs.NumberOfCommands]Command) {
	dgs.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
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
			case "add-priv-category":
				response := (cmds[refs.IndexAddPrivateCategory]).Execute(s, i)
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: response,
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					utils.Log(err, "", "SetupCommands")
					return
				}
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
			}
		case discordgo.InteractionModalSubmit:
			data := i.ModalSubmitData()
			if strings.Contains(data.CustomID, "mainChannel") {
				for _, room := range apps.RoomList {
					if room.GM.User.ID == i.User.ID {
						for _, row := range data.Components {
							actionRow := row.(*discordgo.ActionsRow)
							for _, comp := range actionRow.Components {
								input := comp.(*discordgo.TextInput)
								for _, r := range apps.RoomList {
									if strings.Contains(input.CustomID, r.VC.Name) {
										r.ShowAnswer(dgs, apps.Answer{
											CustomID: input.CustomID,
											Value:    input.Value,
										})
									}
								}
								if strings.Contains(input.CustomID, "MainChannel") {
									room.MainChannelID = input.Value
								}
							}
						}
					}
				}
				return
			}
			if strings.Contains(data.CustomID, "vote") {
				utils.MakeModal(dgs, i, data.CustomID)
			}
		}
	})
}
