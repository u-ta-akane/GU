package utils

import (
	"GU/refs"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Envelope struct {
	Channel string
	Message string
}

var (
	IsCreatedChannel      = false
	GeneralMessageChannel = make(chan Envelope)
	YURUBOItemChannel     = make(chan refs.JobData)
	ErrorChannel          = make(chan error)
)

func Log(err error, message string, from string) {
	if err != nil {
		msg := time.Now().Format(time.DateTime)
		msg += fmt.Sprintf(", at %s | %v", from, err)
		log.Printf(msg)
		OrderSendMessage(refs.Config.ModeratorChannelID, msg)
	}
	if message != "" {
		msg := time.Now().Format(time.DateTime)
		msg += fmt.Sprintf(" | " + message)
		log.Printf(msg)
		OrderSendMessage(refs.Config.ModeratorChannelID, msg)
	}
}

func SendMessage(channelId string, message string, dgs *discordgo.Session) {
	_, err := dgs.ChannelMessageSend(channelId, message)
	if err != nil {
		log.Printf(time.Now().Format(time.DateTime)+" | Send Message Error : %v", err)
	}
}

func OrderSendMessage(channelId string, message string) {
	GeneralMessageChannel <- Envelope{channelId, message}
}

func DeleteMessages(channelID string, from string, num int, reason string, dgs *discordgo.Session) int {
	msgs, err := dgs.ChannelMessages(channelID, num-1, "", from, "", discordgo.WithAuditLogReason(reason))
	if err != nil {
		Log(err, "", "DeleteMessages")
		return 1
	}
	if len(msgs) == 0 {
		return 0 // メッセージがなくなった
	}

	// メッセージIDの配列
	ids := make([]string, len(msgs)+1)
	ids[0] = from
	for i, msg := range msgs {
		ids[i+1] = msg.ID
	}

	// 一括削除
	err = dgs.ChannelMessagesBulkDelete(channelID, ids)
	if err != nil {
		Log(err, "", "DeleteMessages")
		return 1
	}
	return 0
}

func MakeModal(s *discordgo.Session, i *discordgo.InteractionCreate, customID string) {
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "modal-" + customID,
				Title:    "入力してください",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID: "ans-" + customID,
								Label:    "解答",
								Style:    discordgo.TextInputShort,
								Required: true,
							},
						},
					},
				},
			},
		},
	)
	if err != nil {
		Log(err, "", "MakeModal")
		return
	}
}
