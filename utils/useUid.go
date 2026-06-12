package utils

import "github.com/bwmarrin/snowflake"

var IDChannel = make(chan string)

func GenerateID(node *snowflake.Node) string {
	return node.Generate().String()
}
