package refs

import "github.com/bwmarrin/discordgo"

type SecretData struct {
	BotToken string `json:"bot_token"`
}
type GuildStructure struct {
	GuildID            string `json:"guild_id"`
	YURUBOChannelID    string `json:"yurubo_channel_id"`
	ModeratorChannelID string `json:"moderator_channel_id"`
	DefaultAuthorityID string `json:"default_authority_id"`
}

var (
	Secrets = SecretData{}
	Config  = GuildStructure{}
)

const (
	ColorFPS       = 0x800000
	ColorRPG       = 0x808000
	ColorMineCraft = 0x008000
	ColorHorror    = 0x000080
	ColorTableGame
	ColorDeveloper = 0x378566
	ColorTrpg      = 0x0
)

const (
	_ uint8 = iota
	AuthorityControlMessages
	AuthoritySendAdminMessage
	AuthorityBotManagement
	AuthorityReflashData
)

const (
	IndexAdminTestMessage = iota
	IndexAddYURUBO
	IndexDeleteYURUBO
	IndexAdminDeleteMessages
	IndexAdminStopBot
	IndexAdminReflashRoleData
	IndexTrpgStart
	IndexTrpgSetMute
	// NumberOfCommands この上にコマンドを追加する！
	NumberOfCommands
)

const (
	MakeTrpgTextHandler uint8 = iota
	RemoveTrpgTextHandler
)

const StopRecode uint8 = 0

// JobData は、チーム、cronスケジュール、および役職を持つジョブの構造体
type JobData struct {
	Id     string   `json:"id"`
	Title  string   `json:"title"`
	Rank   string   `json:"rank"`
	Number int64    `json:"number"`
	Cron   string   `json:"cron"`
	Role   string   `json:"role"`
	Gap    int      `json:"gap"`
	Party  []string `json:"party"`
}

var RoleMap = make(map[string]*discordgo.Role)

func ReflashRoleData(s *discordgo.Session) {
	// ギルドのすべてのロール取得
	guildRoles, _ := s.GuildRoles(Config.GuildID)
	for _, role := range guildRoles {
		RoleMap[role.ID] = role
	}

}

func GetColor(role string) int {
	color := 0x000
	switch role {
	case "OW":
		fallthrough
	case "Valorant":
		fallthrough
	case "Apex":
		fallthrough
	case "BF":
		fallthrough
	case "R6S":
		fallthrough
	case "Strinova":
		color = ColorFPS
		break
	case "Identity V":
		fallthrough
	case "DBD":
		color = ColorHorror
		break
	case "Minecraft":
		color = ColorMineCraft
		break
	case "Shadowverse":
		fallthrough
	case "デュエマ":
		fallthrough
	case "遊戯王":
		fallthrough
	case "雀魂":
		color = ColorTableGame
		break
	case "Dev":
		color = ColorDeveloper
		break
	default:
		color = ColorRPG
	}
	return color
}
