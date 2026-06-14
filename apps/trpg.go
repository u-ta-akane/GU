package apps

import (
	"GU/refs"
	"GU/utils"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var RoomList = make([]*Room, 0)
var RemoveMessageHandler func()

type Answer struct {
	CustomID string
	Value    string
}

type Room struct {
	VC             *discordgo.Channel
	messages       chan *discordgo.MessageCreate
	consoleChannel *discordgo.Channel
	mainChannel    *discordgo.Channel
	role           *discordgo.Role
	mute           bool
	answer         []Answer
}

func (r *Room) getVCMember(s *discordgo.Session) []*discordgo.Member {
	guild, err := s.State.Guild(refs.Config.GuildID)
	var res []*discordgo.Member
	if err != nil {
		utils.Log(err, "", "getVCMember")
	}
	for _, vs := range guild.VoiceStates {
		if vs.ChannelID != r.VC.ID {
			continue
		}

		member, e := s.GuildMember(refs.Config.GuildID, vs.UserID)
		if e != nil {
			continue
		}
		res = append(res, member)
	}
	return res
}

func trpgMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	for idx, room := range RoomList {
		if room.consoleChannel.ID == m.ChannelID {
			req := strings.Split(m.Content, "!")
			if len(req) == 2 && req[0] == "g" {
				switch req[1] {
				case "quit":
					fallthrough
				case "q":
					err := s.GuildRoleDelete(refs.Config.GuildID, room.role.ID)
					if err != nil {
						utils.Log(err, "", "trpgMessageHandler : quit/q")
						return
					}
					_, err = s.ChannelDelete(room.consoleChannel.ID)
					if err != nil {
						utils.Log(err, "", "trpgMessageHandler : quit/q")
						return
					}
					RoomList[idx] = &Room{}
					break
				case "yuetsu":
					fallthrough
				case "yuetu":
					for _, mem := range room.getVCMember(s) {
						if strings.Contains(m.Member.User.Username, "愉悦") || strings.Contains(m.Member.User.Username, "観戦") {
							err := s.GuildMemberRoleAdd(
								refs.Config.GuildID,
								mem.User.ID,
								room.role.ID,
							)
							if err != nil {
								utils.Log(err, "", "trpgMessageHandler : yuetsu/yuetu")
							}
						}
					}
				case "m":
					fallthrough
				case "mute":
					for _, mem := range room.getVCMember(s) {
						for _, role := range mem.Roles {
							if role == room.role.ID {
								err := s.GuildMemberMute(refs.Config.GuildID, mem.User.ID, !room.mute)
								if err != nil {
									utils.Log(err, "", "trpgMessageHandler : m/mute")
								}
								break
							}
						}
					}
					room.mute = !room.mute
					break
				case "vote":

				}
			}
		}
	}
}

func NewRoom(s *discordgo.Session, i *discordgo.InteractionCreate, vc *discordgo.Channel) {
	cc, err := s.GuildChannelCreateComplex(
		refs.Config.GuildID,
		discordgo.GuildChannelCreateData{
			Name: fmt.Sprintf("monitor-room : %s", vc.Name),
			Type: discordgo.ChannelTypeGuildText,
			PermissionOverwrites: []*discordgo.PermissionOverwrite{
				{
					ID:   refs.Config.GuildID, // @everyone
					Type: discordgo.PermissionOverwriteTypeRole,
					Deny: discordgo.PermissionViewChannel,
				},
				{
					ID:   i.Member.User.ID,
					Type: discordgo.PermissionOverwriteTypeMember,
					Allow: discordgo.PermissionViewChannel |
						discordgo.PermissionSendMessages,
				},
			},
		},
	)
	if err != nil {
		utils.Log(err, "", "NewRoom")
		return
	}
	role, err := s.GuildRoleCreate(refs.Config.GuildID, &discordgo.RoleParams{
		Name:        fmt.Sprintf("PL(%s)", vc.Name),
		Color:       new(refs.ColorTrpg),
		Hoist:       new(true),
		Permissions: new(int64(0)),
		Mentionable: new(false)})
	if err != nil {
		utils.Log(err, "", "NewRoom")
		return
	}
	if RemoveMessageHandler == nil {
		RemoveMessageHandler = s.AddHandler(trpgMessageHandler)
	}
	RoomList = append(RoomList, &Room{
		VC:             vc,
		consoleChannel: cc,
		role:           role,
		mute:           false,
	})
	utils.SendMessage(cc.ID, "TRPGセッションが開始されました！\n\nコマンドの使用方法を以下に示します\ng!yuetu(g!yuetsu) : 現在VCに参加中で、名前に観戦、愉悦とある人以外にPLロールを付与します\ng!q(g!quit) : TRPGセッションを終了します(このコンソールが閉じます)\ng!m(g!mute) : PLロールを持っている人全員をミュートします。もう一度実行すると解除されます。このミュートはこのコマンドによってしか解除できません。", s)
}

func (r *Room) vote(s *discordgo.Session) {
	if r.mainChannel == nil {
	}
	for _, member := range r.getVCMember(s) {
		for _, role := range member.Roles {
			if role == r.role.ID {
				_, err := s.ChannelMessageSendComplex(
					refs.Config.GuildID,
					&discordgo.MessageSend{
						Content: "入力してください",
						Components: []discordgo.MessageComponent{
							discordgo.ActionsRow{
								Components: []discordgo.MessageComponent{
									discordgo.Button{
										Label:    "入力",
										Style:    discordgo.PrimaryButton,
										CustomID: fmt.Sprintf("vote-%s-%s", r.VC.Name, member.User.ID),
									},
								},
							},
						},
					},
				)
				if err != nil {
					utils.Log(err, "", "vote")
				}
			}
		}
	}
}

func (r *Room) ShowAnswer(s *discordgo.Session, res Answer) {

}
