package main

import (
	"github.com/bwmarrin/discordgo"
)

// Command cmdsに追加するために満たすべきインターフェースです。
type Command interface {
	// CreateCommand discordgo.ApplicationCommandをmain.goに返します。
	CreateCommand() []*discordgo.ApplicationCommand
	// Execute コマンドが実行されたときに呼ばれる、処理の本体です。
	Execute(s *discordgo.Session, i *discordgo.InteractionCreate) string
}
