package apps

import (
	"GU/refs"
	"GU/utils"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func IOPrivateCategoryMember(s *discordgo.Session, user *discordgo.User, target string) {
	chs, _ := s.GuildChannels(refs.Config.GuildID)
	name := ""
	for _, ch := range chs {
		if ch.ID == target {
			name = ch.Name
			_, err := s.ChannelEdit(
				target,
				&discordgo.ChannelEdit{
					PermissionOverwrites: func() []*discordgo.PermissionOverwrite {
						accept := true
						po := make([]*discordgo.PermissionOverwrite, 0)
						for _, oldPO := range ch.PermissionOverwrites {
							if oldPO.ID == user.ID {
								accept = !accept
								continue
							}
							po = append(po, oldPO)
						}
						if accept {
							po = append(po, &discordgo.PermissionOverwrite{
								Type:  discordgo.PermissionOverwriteTypeMember,
								ID:    user.ID,
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
			utils.Log(nil, fmt.Sprintf("User : %s\nTarget : %s", user.GlobalName, name), "IOPrivateCategoryMember")
		}
	}
}
