package main

import (
	"GU/apps"
	"GU/refs"
	"GU/utils"
	"fmt"
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
	if m.MessageID == refs.Config.RoleEntranceMessageID {
		if m.Member.User.ID == dgs.State.User.ID {
			return
		}
		for _, cat := range refs.PrivateCategories {
			if m.Emoji.Name == cat.EmojiName {
				apps.IOPrivateCategoryMember(dgs, m.Member.User, cat.CategoryID)
			}
		}
		err := dgs.MessageReactionRemove(refs.Config.RoleEntranceChannelID, refs.Config.RoleEntranceMessageID, m.Emoji.Name, m.Member.User.ID)
		if err != nil {
			utils.Log(err, "", "onReactionAdd")
		}
	}
}

func onMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	for idx, room := range apps.RoomList {
		if m.ChannelID == room.ConsoleChannel.ID {
			req := strings.Split(m.Content, "!")
			if len(req) == 2 && req[0] == "g" {
				switch req[1] {
				case "main":
					_, err := s.ChannelMessageSendComplex(
						room.ConsoleChannel.ID,
						&discordgo.MessageSend{
							Content: "PLに一般公開するチャンネルを選択してください",
							Components: []discordgo.MessageComponent{
								discordgo.ActionsRow{
									Components: []discordgo.MessageComponent{
										discordgo.SelectMenu{
											CustomID:    fmt.Sprintf("MainChannel : %s", room.VC.Name),
											Placeholder: "チャンネルを選択",
											MenuType:    discordgo.ChannelSelectMenu,
											ChannelTypes: []discordgo.ChannelType{
												discordgo.ChannelTypeGuildText,
												discordgo.ChannelTypeGuildVoice,
											},
										},
									},
								},
							},
						},
					)
					if err != nil {
						utils.Log(err, "", "vote")
					}
				case "quit":
					fallthrough
				case "q":
					if room.IsRecording {
						apps.StopRecode(room)
					}
					err := s.GuildRoleDelete(refs.Config.GuildID, room.Role.ID)
					if err != nil {
						utils.Log(err, "", "trpgMessageHandler : quit/q")
						return
					}
					_, err = s.ChannelDelete(room.ConsoleChannel.ID)
					if err != nil {
						utils.Log(err, "", "trpgMessageHandler : quit/q")
						return
					}
					_, err = s.ChannelEdit(room.VC.ID, &discordgo.ChannelEdit{
						Name: room.Name,
					})
					if err != nil {
						utils.Log(err, "", "trpgMessageHandler : quit/q")
						return
					}
					apps.RoomList[idx] = &apps.Room{}
				case "yuetsu":
					fallthrough
				case "yuetu":
					for _, mem := range room.GetVCMember(s) {
						if !strings.Contains(m.Member.User.Username, "愉悦") || !strings.Contains(m.Member.User.Username, "観戦") {
							err := s.GuildMemberRoleAdd(
								refs.Config.GuildID,
								mem.User.ID,
								room.Role.ID,
							)
							room.Pls = append(room.Pls, mem)
							room.PlVotes[mem.User.ID] = false
							if err != nil {
								utils.Log(err, "", "trpgMessageHandler : yuetsu/yuetu")
							}
						}
					}
				case "m":
					fallthrough
				case "mute":
					for _, mem := range room.GetVCMember(s) {
						for _, role := range mem.Roles {
							if role == room.Role.ID {
								err := s.GuildMemberMute(refs.Config.GuildID, mem.User.ID, !room.Mute)
								if err != nil {
									utils.Log(err, "", "trpgMessageHandler : m/mute")
								}
								break
							}
						}
					}
					room.Mute = !room.Mute
					break
				case "vote":
					room.Vote(s)
				case "recode":
					room.Recode(s)
				case "stop-rec":
					fallthrough
				case "stop-recode":
					apps.StopRecode(room)
				}

			}
		}
	}
}

func onPresenceUpdate(s *discordgo.Session, p *discordgo.PresenceUpdate) {
	if p.Status == discordgo.StatusOffline {
		user, err := s.GuildMember(refs.Config.GuildID, p.User.ID)
		if err != nil {
			utils.Log(err, "", "onPresenceUpdate")
			return
		}
		for _, role := range user.Roles {
			if role == refs.Config.PlayingStatusRoleID {
				err = s.GuildMemberRoleRemove(refs.Config.GuildID, p.User.ID, refs.Config.PlayingStatusRoleID)
				if err != nil {
					utils.Log(err, "", "onPresenceUpdate")
					return
				}
			}
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
			case "a-send-role-entrance":
				AdminHandler(s, i, cmds, refs.IndexAdminSendRoleEntranceMessage)
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
					utils.Log(err, "", "setupOnInteractionHandler")
					return
				}
			case "play":
				response := (cmds[refs.IndexStatusCommand]).Execute(s, i)
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
