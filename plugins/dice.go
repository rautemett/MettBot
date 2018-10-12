package plugins

import (
	"../ircclient"
	"math/rand"
	"strings"
	"fmt"
)

type DicePlugin struct {
	ic *ircclient.IRCClient
}

func (x *DicePlugin) Register(cl *ircclient.IRCClient) {
	x.ic = cl
}

func (x *DicePlugin) String() string {
	return "dice"
}

func (x *DicePlugin) Info() string {
	return "roll a dice"
}

func (x *DicePlugin) Usage(cmd string) string {
	return ""
}

func (x *DicePlugin) ProcessLine(msg *ircclient.IRCMessage) {
	if msg.Command != "PRIVMSG" || len(msg.Args) == 0 {
		return
	}

	if strings.HasPrefix(msg.Args[0], "!d") {
		var sides int
		n, err := fmt.Sscanf(msg.Args[0], "!d%d", &sides);
		if err == nil && n == 1 && sides > 0 {
			r := rand.Intn(sides) + 1
			s := fmt.Sprintf("Your d%d rolled a %d", sides, r)
			x.ic.ReplyMsg(msg, s)
		}
	}
}

func (x *DicePlugin) ProcessCommand(cmd *ircclient.IRCCommand) {
}

func (x *DicePlugin) Unregister() {
}
