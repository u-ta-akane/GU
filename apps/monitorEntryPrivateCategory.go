package apps

import (
	"GU/refs"
	"GU/utils"

	"github.com/bwmarrin/discordgo"
)

func IOPrivateCategoryMember(s *discordgo.Session, user string, target string) {
	chs, _ := s.GuildChannels(refs.Config.GuildID)
	for _, ch := range chs {
		if ch.ID == target {
			_, err := s.ChannelEdit(
				target,
				&discordgo.ChannelEdit{
					PermissionOverwrites: func() []*discordgo.PermissionOverwrite {
						accept := true
						po := make([]*discordgo.PermissionOverwrite, 0)
						for _, oldPO := range ch.PermissionOverwrites {
							if oldPO.ID == user {
								accept = !accept
								continue
							}
							po = append(po, oldPO)
						}
						if accept {
							po = append(po, &discordgo.PermissionOverwrite{
								Type:  discordgo.PermissionOverwriteTypeMember,
								ID:    user,
								Allow: refs.PrivateCategoryMemberPermission,
							})
						}
						return po
					}(),
				},
			)
			if err != nil {
				utils.Log(err, "", "monitorPrivateCategory")
			}
		}
	}
}
