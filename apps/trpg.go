package apps

import (
	"GU/refs"
	"GU/utils"
	"fmt"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4/pkg/media/oggwriter"
)

var RoomList = []*Room{
	{
		ConsoleChannel: &discordgo.Channel{
			Type: discordgo.ChannelTypeGuildText,
			ID:   "",
			Name: "dummy",
		},
	},
}
var UsingRecode uint8 = 0
var ContinueRecodeChannel = make(chan uint8)

type Answer struct {
	CustomID string
	Value    string
}

type recorder struct {
	Writers map[uint32]*oggwriter.OggWriter
	Mutex   sync.Mutex
}

type Room struct {
	VC             *discordgo.Channel
	messages       chan *discordgo.MessageCreate
	ConsoleChannel *discordgo.Channel
	IsRecording    bool
	MainChannelID  string
	Role           *discordgo.Role
	GM             *discordgo.Member
	Mute           bool
	Pls            []*discordgo.Member
	PlVotes        map[string]bool
	voteButtonIDs  []string
	fin            bool
}

func (r *recorder) GetWriter(
	userID string,
	ssrc uint32,
) (*oggwriter.OggWriter, error) {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if w, ok := r.Writers[ssrc]; ok {
		return w, nil
	}

	filename := fmt.Sprintf("%s.ogg", userID)

	w, err := oggwriter.New(
		filename,
		48000,
		2,
	)
	if err != nil {
		return nil, err
	}

	r.Writers[ssrc] = w

	return w, nil
}

func (r *Room) GetVCMember(s *discordgo.Session) []*discordgo.Member {
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

func NewRoom(s *discordgo.Session, i *discordgo.InteractionCreate, vc *discordgo.Channel) {
	cc, err := s.GuildChannelCreateComplex(
		refs.Config.GuildID,
		discordgo.GuildChannelCreateData{
			Name: fmt.Sprintf("console : %s", vc.Name),
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
	RoomList = append(RoomList, &Room{
		VC:             vc,
		ConsoleChannel: cc,
		Role:           role,
		Mute:           false,
		GM:             i.Member,
		fin:            false,
		IsRecording:    false,
	})
	utils.SendMessage(cc.ID, "TRPGセッションが開始されました！\n\nコマンドの使用方法を以下に示します\ng!yuetu(g!yuetsu) : 現在VCに参加中で、名前に観戦、愉悦とある人以外にPLロールを付与します\ng!q(g!quit) : TRPGセッションを終了します(このコンソールが閉じます)\ng!m(g!mute) : PLロールを持っている人全員をミュートします。もう一度実行すると解除されます。このミュートはこのコマンドによってしか解除できません。", s)
}

func (r *Room) Vote(s *discordgo.Session) {
	if r.MainChannelID == "" {
		return
	}
	for _, member := range r.GetVCMember(s) {
		for _, role := range member.Roles {
			if role == r.Role.ID {
				m, err := s.ChannelMessageSendComplex(
					r.VC.ID,
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
				r.voteButtonIDs = append(r.voteButtonIDs, m.ID)
			}
		}
	}
}

func (r *Room) ShowAnswer(s *discordgo.Session, res Answer) {
	var id string
	var name string
	for _, pl := range r.Pls {
		if strings.Contains(res.CustomID, pl.User.ID) {
			id = pl.User.ID
			name = pl.User.Username
			break
		}
	}
	if r.PlVotes[id] {
		utils.SendMessage(r.ConsoleChannel.ID, fmt.Sprintf("%s(2回目以降の解答) : %s", name, res.Value), s)
	} else {
		utils.SendMessage(r.ConsoleChannel.ID, fmt.Sprintf("%s : %s", name, res.Value), s)
	}
	i := 0
	for _, v := range r.PlVotes {
		if v {
			i++
		}
	}
	if i == len(r.Pls) {
		utils.SendMessage(r.ConsoleChannel.ID, "全PLが解答しました", s)
		r.closeVote(s)
	}
}

func (r *Room) closeVote(s *discordgo.Session) {
	for key, _ := range r.PlVotes {
		r.PlVotes[key] = false
	}
	for _, id := range r.voteButtonIDs {
		err := s.ChannelMessageDelete(r.MainChannelID, id)
		if err != nil {
			utils.Log(err, "", "closeVote")
			return
		}
	}
}

func (r *Room) Recode(s *discordgo.Session) {
	vc, e := s.ChannelVoiceJoin(
		refs.Config.GuildID,
		r.VC.ID,
		false,
		false,
	)
	UsingRecode++
	r.IsRecording = true
	if e != nil {
		utils.Log(e, "", "Recode")
	}
	rec := &recorder{
		Writers: make(map[uint32]*oggwriter.OggWriter),
	}
	go func() {
		for {
		check:
			select {
			case h := <-ContinueRecodeChannel:
				if h == refs.StopRecode {
					break check
				}
			default:
				break check
			}
			packet, ok := <-vc.OpusRecv
			if !ok {
				break
			}
			userID := vc.UserID
			if userID == "" {
				continue
			}
			writer, err := rec.GetWriter(
				userID,
				packet.SSRC,
			)
			if err != nil {
				utils.Log(err, "", "Recode")
				continue
			}

			err = writer.WriteRTP(&rtp.Packet{
				Header: rtp.Header{
					Version:        2,
					PayloadType:    111,
					SequenceNumber: packet.Sequence,
					Timestamp:      packet.Timestamp,
					SSRC:           packet.SSRC,
				},
				Payload: packet.Opus,
			})

			if err != nil {
				utils.Log(err, "", "Recode")
			}
		}
		for _, writer := range rec.Writers {
			err := writer.Close()
			if err != nil {
				utils.Log(err, "", "Recode")
			}
		}
	}()
}

func StopRecode(room *Room) {
	if room.IsRecording {
		room.IsRecording = false
		UsingRecode--
		if UsingRecode == 0 {
			ContinueRecodeChannel <- refs.StopRecode
			return
		}
		utils.OrderSendMessage(room.ConsoleChannel.ID, "現在は録音中ではありません！")
	}
}
