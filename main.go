package main

import (
	"GU/apps"
	"GU/commands"
	"GU/commands/admin"
	"GU/commands/trpg"
	"GU/refs"
	"GU/utils"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
)

var (
	dgs *discordgo.Session
)

// DiscordSessionManager Discordセッションを管理する構造体
type DiscordSessionManager struct{}

func (d *DiscordSessionManager) InitializeSession(token string) *discordgo.Session {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Discordセッションの作成に失敗: %v", err)
	}
	return dg
}

func checkInfo(data interface{}, opt *string, mode string) interface{} {
	isTarget := false
	txt1 := ""
	txt2 := ""
	var secrets refs.SecretData
	var guildStr refs.GuildStructure
	switch mode {
	case "botToken":
		secrets, isTarget = data.(refs.SecretData)
		txt1 = "Bot Token"
		txt2 = "set-bot-token"
		break
	case "moderator":
		guildStr, isTarget = data.(refs.GuildStructure)
		txt1 = "ModeratorChannel ID"
		txt2 = "set-moderator-chan"
		break
	case "YURUBO":
		guildStr, isTarget = data.(refs.GuildStructure)
		txt1 = "YuruboChannel ID"
		txt2 = "set-yurubo-chan"
		break
	case "roll-entrance-chan":
		guildStr, isTarget = data.(refs.GuildStructure)
		txt1 = "Roll-Entrance Channel ID"
		txt2 = "set-roll-entrance-chan"
	case "defaultAuth":
		guildStr, isTarget = data.(refs.GuildStructure)
		txt1 = "Default Authority ID"
		txt2 = "set-default-authority"
		break
	case "guild":
		guildStr, isTarget = data.(refs.GuildStructure)
		txt1 = "Guild ID"
		txt2 = "set-guild-id"
	}
	if !isTarget {
		if *opt == "" {
			fmt.Printf(txt1 + "is empty. Please Use --" + txt2 + ".\n")
			return nil
		}
		switch mode {
		case "botToken":
			secrets = refs.SecretData{
				BotToken: *opt,
			}
			return secrets
		case "moderator":
			guildStr = refs.GuildStructure{
				ModeratorChannelID: *opt,
			}
		case "YURUBO":
			guildStr = refs.GuildStructure{
				YURUBOChannelID: *opt,
			}
		case "defaultAuth":
			guildStr = refs.GuildStructure{
				DefaultAuthorityID: *opt,
			}
		case "roll-entrance-chan":
			guildStr = refs.GuildStructure{
				RollEntranceChannelID: *opt,
			}
		case "guild":
			guildStr = refs.GuildStructure{
				GuildID: *opt,
			}
		}
		return guildStr
	}
	switch mode {
	case "botToken":
		if *opt == "" {
			break
		}
		secrets.BotToken = *opt
		return secrets
	case "moderator":
		if *opt == "" {
			break
		}
		guildStr.ModeratorChannelID = *opt
	case "YURUBO":
		if *opt == "" {
			break
		}
		guildStr.YURUBOChannelID = *opt
	case "roll-entrance-chan":
		if *opt == "" {
			break
		}
		guildStr.RollEntranceChannelID = *opt
	case "guild":
		if *opt == "" {
			break
		}
		guildStr.GuildID = *opt
	default:
		if *opt == "" {
			break
		}
		guildStr.DefaultAuthorityID = *opt
	}
	return guildStr
}

func main() {
	newToken := flag.String("set-bot-token", "", "Enter Bot Token")
	setModerateChannel := flag.String("set-moderator-chan", "", "Enter Moderate Channel ID")
	setYURUBOChannel := flag.String("set-yurubo-chan", "", "Enter YURUBO Channel ID")
	setGuildID := flag.String("set-guild-id", "", "Enter Guild ID")
	setDefaultAuthority := flag.String("set-default-authority", "", "Enter Default Authority ID")
	setRollEntranceChannelID := flag.String("set-roll-entrance-chan", "", "Enter Roll-Entrance Channel ID")
	flag.Parse()
	secrets := utils.JSONFM.Read("secrets.json")
	guildStr := utils.JSONFM.Read("guildStructure.json")
	secrets = checkInfo(secrets, newToken, "botToken")
	guildStr = checkInfo(guildStr, setModerateChannel, "moderator")
	guildStr = checkInfo(guildStr, setYURUBOChannel, "YURUBO")
	guildStr = checkInfo(guildStr, setDefaultAuthority, "defaultAuth")
	guildStr = checkInfo(guildStr, setRollEntranceChannelID, "roll-entrance-chan")
	guildStr = checkInfo(guildStr, setGuildID, "guild")
	if secrets == nil || guildStr == nil {
		fmt.Printf("Some parameters are missing. Please fill all blank parameters.\n")
		return
	}
	refs.Secrets = secrets.(refs.SecretData)
	refs.Config = guildStr.(refs.GuildStructure)
	if refs.Secrets.BotToken == "" || refs.Config.ModeratorChannelID == "" || refs.Config.YURUBOChannelID == "" || refs.Config.GuildID == "" || refs.Config.DefaultAuthorityID == "" || refs.Config.RollEntranceChannelID == "" {
		fmt.Printf("Some parameters are missing. Please fill all blank parameters.\n")
		return
	}
	sessionManager := &DiscordSessionManager{}
	dgs = sessionManager.InitializeSession(refs.Secrets.BotToken)
	dgs.AddHandler(onMemberAdd)
	dgs.AddHandler(onInteraction)
	dgs.AddHandler(TrpgTextHandler)
	if err := dgs.Open(); err != nil {
		var restErr *discordgo.RESTError
		switch {
		case errors.As(err, &restErr):
			switch restErr.Message.Code {
			case discordgo.ErrCodeInvalidAuthenticationToken:
				fmt.Println("Invalid Authentication Token")
				fmt.Printf("Saved Bot Token is incorrect. Please Use --set-bot-token and Enter the Token. ")
				return
			case discordgo.ErrCodeUnauthorized:
				fmt.Printf("%v\n", err)
				fmt.Printf("Please Check Internet, Use --set-bot-token and Enter the Token. ")
				return
			case discordgo.ErrCodeUnknownToken:
				fmt.Printf("%v\n", err)
				fmt.Printf("Saved Bot Token is incorrect. Please Use --set-bot-token and Enter the Token. ")
				return
			}
		case strings.Contains(err.Error(), "4004"):
			fmt.Println("Gateway Authentication Failed")
			fmt.Printf("%v\n", err)
			fmt.Printf("Please Check Internet, Use --set-bot-token and Enter the Token. ")
			return
		default:
			fmt.Printf("%T\n", err)
			fmt.Println(err)
			return
		}
	}
	//cmdsの順番は変えないこと！！！
	//コマンドを追加する際はdataにインデックス番号の追記を適切な位置にすること！
	var cmds = [refs.NumberOfCommands]Command{
		&admin.AdminTestMessageCommand{},
		&commands.AddYURUBOCommand{},
		&commands.DeleteYURUBOCommand{},
		&admin.AdminDeleteMessagesCommands{},
		&admin.AdminStopBotCommand{},
		&admin.AdminReflashRoleDataCommand{},
		&trpg.TrpgStartCommand{},
	}
	createdCommands := func() []*discordgo.ApplicationCommand {
		apps.Ns.InitializeSchedule()
		jd := utils.JSONFM.Read("jobData.json")
		if _, ok := jd.([]refs.JobData); ok {
			utils.JobDataSlice = jd.([]refs.JobData)
		} else {
			utils.JobDataSlice = make([]refs.JobData, 0)
		}
		utils.SendMessage(refs.Config.ModeratorChannelID, "Finished Scheduler Initialization", dgs)
		createdCommands := make([]*discordgo.ApplicationCommand, 0, len(cmds))
		for _, v := range cmds {
			createdCommands = append(createdCommands, v.CreateCommand()...)
		}
		return createdCommands
	}()
	i := 0
	for _, cmd := range createdCommands {
		_, err := dgs.ApplicationCommandCreate(dgs.State.User.ID, refs.Config.GuildID, cmd)
		if err != nil {
			e := fmt.Sprintf("Command Registration Error : %s, %v", cmd.Name, err)
			utils.SendMessage(refs.Config.ModeratorChannelID, e, dgs)
			i++
		} else {
			utils.SendMessage(refs.Config.ModeratorChannelID, cmd.Name+" was registered successfully", dgs)
		}
	}
	if i == 0 {
		log.Printf("Info : All commands have been registered")
		utils.SendMessage(refs.Config.ModeratorChannelID, "Info : All commands have been registered", dgs)
	}
	if refs.Config.RollEntranceChannelID == "" {
		utils.SendMessage(refs.Config.ModeratorChannelID, "Roll-Entrance Channel ID is empty", dgs)
		log.Println("Roll-Entrance Channel ID is empty")
	}
	if refs.Config.RollEntranceMessageID == "" {
		utils.SendMessage(refs.Config.ModeratorChannelID, "Roll-Entrance Message ID is empty", dgs)
		log.Println("Roll-Entrance Message ID is empty")
	}
	SetupCommands(dgs, &cmds)
	defer func(dgs *discordgo.Session) {
		err := dgs.Close()
		if err != nil {
			log.Printf(err.Error())
			utils.SendMessage(refs.Config.ModeratorChannelID, err.Error(), dgs)
		}
	}(dgs)
	channels, _ := dgs.GuildChannels(refs.Config.GuildID)
	for _, channel := range channels {
		if strings.Contains(channel.Name, "priv") {
			refs.PrivateCategories = append(refs.PrivateCategories, channel.ID)
		}
	}
	refs.ReflashRoleData(dgs)
	utils.SendMessage(refs.Config.ModeratorChannelID, "Reflashed RoleData successfully", dgs)
	utils.SendMessage(refs.Config.ModeratorChannelID, "Bot started successfully", dgs)
	log.Println("Bot started. Press CTRL-C to exit")
	node, err := snowflake.NewNode(1)
	if err != nil {
		utils.SendMessage(refs.Config.ModeratorChannelID, err.Error(), dgs)
		log.Fatal(err)
	}
	waitForSignal(refs.Secrets, refs.Config, node)
}

func waitForSignal(secrets refs.SecretData, guildStr refs.GuildStructure, node *snowflake.Node) {
	signal.Notify(admin.SignalChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	utils.IsCreatedChannel = true
	go func() {
		time.Sleep(5000 * time.Millisecond)
		err := utils.JSONFM.Write("secrets.json", secrets)
		fmt.Printf("err=%v type=%T value=%+v\n", err, secrets, secrets)
		err = utils.JSONFM.Write("config.json", guildStr)
		fmt.Printf("err=%v type=%T value=%+v\n", err, guildStr, guildStr)
		if len(utils.JobDataSlice) == 0 {
			err := utils.JSONFM.Write("jobData.json", utils.JobDataSlice)
			if err != nil {
				fmt.Printf("err=%v type=%T value=%+v\n", err, utils.JobDataSlice, utils.JobDataSlice)
			}
		}
	}()
Completed:
	for {
		select {
		case <-admin.SignalChannel:
			break Completed
		case item := <-utils.YURUBOItemChannel:
			utils.SendYURUBOItem(dgs, item)
		case en := <-utils.GeneralMessageChannel:
			utils.SendMessage(en.Channel, en.Message, dgs)
		case <-utils.IDChannel:
			utils.IDChannel <- utils.GenerateID(node)
		case e := <-utils.ErrorChannel:
			utils.Log(e, "", "onMemberAdd")
		}
	}
}
