package main

import (
	"GU/refs"
	"GU/utils"

	"github.com/bwmarrin/discordgo"
)

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

func onClickButton(dgs *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}
	err := dgs.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "操作を実行しました",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	data := i.MessageComponentData()
	switch data.CustomID {
	default:
		utils.YURUBOPartyEdit(dgs, i, data.CustomID)
	}
	if err != nil {
		utils.Log(err, "", "onClickButton")
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
			utils.Log(e, "", "onClickButton")
		}
	}()
	*/
}
