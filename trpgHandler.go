package main

import (
	"GU/apps"
	"GU/refs"
	"GU/utils"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func TrpgHandler(s *discordgo.Session, i *discordgo.InteractionCreate, cmds *[refs.NumberOfCommands]Command, index uint8) {
	switch index {
	case refs.IndexTrpgStart:

		response := (cmds[refs.IndexTrpgStart]).Execute(s, i)
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource, // 「通常の返答」タイプ
			Data: &discordgo.InteractionResponseData{
				Content: response,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			utils.Log(err, "", "trpgHandler")
			return
		}
	}
}

func TrpgTextHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	for idx, room := range apps.RoomList {
		if room.ConsoleChannel.ID == m.ChannelID {
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
						stopRecode(room)
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
					apps.RoomList[idx] = &apps.Room{}
					break
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
					stopRecode(room)
				}

			}
		}
	}
}

func stopRecode(room *apps.Room) {
	if room.IsRecording {
		room.IsRecording = false
		apps.UsingRecode--
		if apps.UsingRecode == 0 {
			apps.ContinueRecodeChannel <- refs.StopRecode
			return
		}
		utils.SendMessage(room.ConsoleChannel.ID, "現在は録音中ではありません！", dgs)
	}
}
