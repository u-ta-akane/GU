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
		switch m.Emoji.Name {
		}
	}
}

func onInteraction(dgs *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionModalSubmit {
		data := i.ModalSubmitData()
		for _, comp := range data.Components {
			for _, c := range comp.(*discordgo.ActionsRow).Components {
				input := c.(*discordgo.TextInput)
				for _, room := range apps.RoomList {
					if strings.Contains(input.CustomID, room.VC.Name) {
						room.ShowAnswer(dgs, apps.Answer{
							CustomID: input.CustomID,
							Value:    input.Value,
						})
					}
				}
			}
		}

	}

	/*err := dgs.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "操作を実行しました",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})*/
	data := i.MessageComponentData()
	if strings.Contains(data.CustomID, "mainChannel") {
		for _, room := range apps.RoomList {
			if room.GM.User.ID == i.User.ID {
				room.MainChannelID = data.Values[0]
			}
		}
		return
	}
	if strings.Contains(data.CustomID, "vote") {
		utils.MakeModal(dgs, i, data.CustomID)
	}

	switch data.CustomID {

	default:
		utils.YURUBOPartyEdit(dgs, i, data.CustomID)
	}
	/*err := dgs.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		return
	}

	go func() {
		switch data.CustomID {
		default:
			utils.YURUBOPartyEdit(dgs, i, data.CustomID)
		}
		_, e := dgs.FollowupMessageCreate(
			i.Interaction,
			false,
			&discordgo.WebhookParams{
				Content: "操作を実行しました",
			},
		)
		if e != nil {
			utils.Log(e, "", "onInteraction")
		}
	}()
	*/
}
