package utils

import (
	"GU/refs"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func hasBit(num uint32, bitPosition uint32) bool {
	return (num & (1 << bitPosition)) != 0
}

func caluculateAuthority(params ...uint32) uint32 {
	var ans uint32
	for _, i := range params {
		ans += 2 ^ i
	}
	return ans
}

func HasAuthority(s *discordgo.Session, i *discordgo.InteractionCreate, targetAuthority uint32) (bool, int) {
	member, err := s.GuildMember(i.GuildID, i.Member.User.ID)
	var ans uint32
	if err != nil {
		Log(err, "", "HasAuthority")
		return false, 1
	}
	for _, roleID := range member.Roles {
		a, e := strconv.ParseInt(refs.RoleMap[roleID].Name, 10, 32)
		if e == nil {
			ans = uint32(a)
			break
		}
	}
	return hasBit(ans, targetAuthority-1), 0
}
