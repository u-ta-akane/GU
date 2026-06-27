package commands

import (
	"GU/apps"
	"GU/refs"
	"GU/utils"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (c *AddYURUBOCommand) CreateCommand() []*discordgo.ApplicationCommand {

	dc := []*discordgo.ApplicationCommand{
		{
			Name:        "ゆるぼ",
			Description: "ゆるぼを追加",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "ゲームタイトル",
					Description: "募集するゲームのタイトル",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "募集人数",
					Description: "n(nは任意の自然数)人",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "募集する日",
					Description: "今日から6日後まで作成できます (0:　今日, 1: 明日, 2: 明後日, 3: 明々後日, 4: 4日後, 5: 5日後, 6: 6日後)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "時",
					Description: "0~23",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "分",
					Description: "0~59",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "ロール",
					Description: "対象とするロール",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "ランク",
					Description: "自分のランク(ランクマッチなどで必要であれば)",
					Required:    false,
				},
			},
		},
	}

	return dc
}

// AddYURUBOCommand ゆるぼを追加するコマンド
type AddYURUBOCommand struct{}

func (c *AddYURUBOCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) string {
	party := make([]string, 1, 99)
	party = append(party, i.Member.User.Username)
	var YURUBO = refs.JobData{
		Title:  i.ApplicationCommandData().Options[0].StringValue(),
		Number: int64(int(i.ApplicationCommandData().Options[1].IntValue())),
		Role:   i.ApplicationCommandData().Options[5].RoleValue(s, i.GuildID).Name,
		Party:  party,
	}
	day := i.ApplicationCommandData().Options[2].IntValue()
	hour := i.ApplicationCommandData().Options[3].IntValue()
	minute := i.ApplicationCommandData().Options[4].IntValue()
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "ランク" {
			YURUBO.Rank = opt.StringValue()
		}
	}
	tz, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(tz)
	date := time.Date(now.Year(), now.Month(), now.Day(), int(hour), int(minute), 0, 0, tz).Add(time.Hour * time.Duration(24*day))
	check := (hour-int64(time.Now().Hour()))*60 + minute - int64(time.Now().Minute())
	YURUBO.Cron = date.Format(time.DateTime)
	jobID, e := apps.Ns.RegisterJob(date, utils.OrderSendYURUBOItem, YURUBO, "")
	if day >= 1 || check > 30 {
		date = date.Add(time.Minute * time.Duration(-30))
		YURUBO.Gap = 30
		jobID, e = apps.Ns.RegisterJob(date, utils.OrderSendYURUBOItem, YURUBO, jobID)
	} else if day == 0 && check > 15 {
		date = date.Add(time.Minute * time.Duration(-10))
		YURUBO.Gap = 10
		jobID, e = apps.Ns.RegisterJob(date, utils.OrderSendYURUBOItem, YURUBO, jobID)
	} else if day == 0 && check < 0 {
		return "エラー : 過去の時間は指定できません"
	}
	if e == 0 {
		YURUBO.Id = jobID
		utils.OrderSendYURUBOItem(YURUBO)
	} else {
		return "ゆるぼの追加に失敗しました"
	}
	utils.JobDataSlice = append(utils.JobDataSlice, YURUBO)
	err := utils.JSONFM.Write("jobData.json", utils.JobDataSlice)
	if err != nil {
		utils.Log(err, "", "AddYURUBO")
	}
	return fmt.Sprintf("新しい募集があります！ @%s", YURUBO.Role)
}

type DeleteYURUBOCommand struct{}

func (c *DeleteYURUBOCommand) CreateCommand() []*discordgo.ApplicationCommand {

	dc := []*discordgo.ApplicationCommand{
		{
			Name:        "delete-ゆるぼ",
			Description: "指定したゆるぼを削除します",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "id",
					Description: "削除したいゆるぼのID",
					Required:    true,
				},
			},
		},
	}

	return dc
}

func (c *DeleteYURUBOCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) string {
	status, res, jt := apps.Ns.RemoveJob(i.ApplicationCommandData().Options[0].StringValue())
	if status == 16 {
		err := utils.JSONFM.Write(time.Now().Format(time.DateOnly)+"-backup-job.json", utils.JobDataSlice)
		if err != nil {
			utils.Log(err, "", "DeleteYURUBO")
		}
		return fmt.Sprintf("エラー: このIDのジョブは見つかりません")
	} else if res != "" {
		err := utils.JSONFM.Write(time.Now().Format(time.DateOnly)+"-backup-job.json", utils.JobDataSlice)
		if err != nil {
			utils.Log(err, "", "DeleteYURUBO")
		}
		return fmt.Sprintf("一つ以上のゆるぼを削除できませんでした\n詳細 : \n%s", res)
	}
	res = ""
	j := jt[0]
	for _, msg := range utils.GetMessages(s, refs.Config.YURUBOChannelID, j.Id) {
		utils.DeleteMessages(refs.Config.YURUBOChannelID, msg.ID, 1, "DeleteYURUBOCommand", s)
	}
	date, _ := time.Parse(time.DateTime, j.Cron)
	if j.Rank != "" {
		res += fmt.Sprintf("このゆるぼを削除しました\n(id: %s)\nタイトル: %s\nランク: %s\n人数: %d\n日付: %s\nロール: %s\n", j.Id, j.Title, j.Rank, j.Number, date.Format(time.DateTime), j.Role)
	} else {
		res += fmt.Sprintf("このゆるぼを削除しました\n(id: %s)\nタイトル: %s\nランク: 非ランク\n人数: %d\n日付: %s\nロール: %s\n", j.Id, j.Title, j.Number, date.Format(time.DateTime), j.Role)
	}
	utils.Log(nil, res, "DeleteYURUBOCommand")
	return res
}
