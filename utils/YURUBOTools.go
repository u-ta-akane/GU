package utils

import (
	"GU/refs"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func OrderSendYURUBOItem(item refs.JobData) { YURUBOItemChannel <- item }

func SendYURUBOItem(s *discordgo.Session, item refs.JobData) {
	rank := func() string {
		if item.Rank == "" {
			return "非ランク"
		}
		return "ランク(募集帯 : " + item.Rank + ")"
	}()
	date, e := time.Parse(time.DateTime, item.Cron)
	if e != nil {
		Log(e, "", "SendYURUBOItem")
		return
	}
	date = date.Add(time.Minute * time.Duration(item.Gap))
	embed := &discordgo.MessageEmbed{
		Title: item.Title,
		Description: fmt.Sprintf("@%s\n募集日時 : %s\nランク/非ランク : %s\n募集人数 : @%d\n\n参加者 : %s", item.Role, func() string {
			if item.Cron == refs.UndecidedYURUBOCron {
				return "未定"
			}
			return date.Format(time.DateTime)
		}(), rank, item.Number,
			func() string {
				str := ""
				for _, party := range item.Party {
					str += party + ", "
				}
				return str
			}(),
		),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  item.Id,
				Value: "(ID)",
			},
		},
		Color: refs.GetColor(item.Role),
	}
	// ボタンを作成
	button := discordgo.Button{
		Label:    "参加する",
		Style:    discordgo.PrimaryButton,
		CustomID: item.Id,
	}
	// ActionsRow にボタンを追加
	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{button},
	}
	// 送信
	_, err := s.ChannelMessageSendComplex(refs.Config.YURUBOChannelID, &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: []discordgo.MessageComponent{actionRow},
	})
	if err != nil {
		Log(err, "", "SendYURUBOItem")
	}
}

func GetMessages(dgs *discordgo.Session, channelID string, id string) []*discordgo.Message {
	msgs, err := dgs.ChannelMessages(channelID, 100, "", "", "")
	if err != nil {
		Log(err, "", "YURUBOPartyEdit")
		return nil
	}
	var targets []*discordgo.Message
	for _, msg := range msgs {
		if len(msg.Embeds) == 0 {
			continue
		}
	Embed:
		for _, embed := range msg.Embeds {
			// 特定のフィールド名を持つか確認
			for _, field := range embed.Fields {
				if field.Name == id {
					targets = append(targets, msg)
					break Embed
				}
			}
		}
	}
	return targets
}

func YURUBOPartyEdit(dgs *discordgo.Session, i *discordgo.InteractionCreate, id string) []*discordgo.Message {
	targets := GetMessages(dgs, refs.Config.YURUBOChannelID, id)
	if len(targets) == 0 {
		return nil
	}
	jt := JSONFM.SearchJobFromJSON(id)
	if len(jt) == 0 {
		return nil
	}
	updatedJob := func() refs.JobData {
		for _, j := range jt {
			if j.Gap == 0 {
				return j
			}
		}
		return jt[0]
	}()
	rank := func() string {
		if updatedJob.Rank == "" {
			return "非ランク"
		}
		return "ランク(募集帯 : " + updatedJob.Rank + ")"
	}()
	p := ""
	gapIndex := 0
	for idx, party := range jt[0].Party {
		if party == i.Member.User.Username {
			gapIndex++
			updatedJob.Party[idx] = ""
			continue
		}
		p += party
		updatedJob.Party[idx-gapIndex] = party
		if idx != len(jt[0].Party)-1 {
			p += ", "
		}
	}
	if gapIndex == 0 {
		updatedJob.Party = append(updatedJob.Party, i.Member.User.Username)
		p += i.Member.User.Username
	} else {
	}
	var updatedEmbed = discordgo.MessageEmbed{
		Title:       updatedJob.Title,
		Description: fmt.Sprintf("@%s\n募集日時 : %s\nランク/非ランク : %s\n募集人数 : @%d\n\n参加者 : %s", updatedJob.Role, updatedJob.Cron, rank, updatedJob.Number, p),
		Color:       refs.GetColor(updatedJob.Role),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  id,
				Value: "(ID)",
			},
		},
	}
	embeds := []*discordgo.MessageEmbed{&updatedEmbed}
	buttons := []discordgo.MessageComponent{
		discordgo.Button{
			Label:    "参加する",
			Style:    discordgo.PrimaryButton,
			CustomID: id,
		},
	}
	contents := []discordgo.MessageComponent{
		&discordgo.ActionsRow{Components: buttons},
	}
	edit := discordgo.MessageEdit{
		Components: &contents,
		Embeds:     &embeds,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{
				discordgo.AllowedMentionTypeUsers,
			},
		},
		Flags:   discordgo.MessageFlags(0),
		ID:      "",
		Channel: refs.Config.YURUBOChannelID,
	}
	var msgs []*discordgo.Message
	for _, target := range targets {
		edit.ID = target.ID
		m, e := dgs.ChannelMessageEditComplex(&edit)
		if e != nil {
			Log(e, "", "YURUBOPartyEdit")
		}
		msgs = append(msgs, m)
	}
	for idx, job := range JobDataSlice {
		if job.Id == id {
			JobDataSlice[idx].Party = updatedJob.Party
		}
	}
	err := JSONFM.Write("jobData.json", JobDataSlice)
	if err != nil {
		Log(err, "", "YURUBOPartyEdit")
	}
	return msgs
}
