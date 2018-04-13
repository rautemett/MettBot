package plugins

import (
	"../ircclient"
	"math/rand"
)

type DecisionPlugin struct {
	ic *ircclient.IRCClient
}

func (x *DecisionPlugin) Register(cl *ircclient.IRCClient) {
	x.ic = cl
	x.ic.RegisterCommandHandler("decide", 0, 0, x)
}

func (x *DecisionPlugin) String() string {
	return "decide"
}

func (x *DecisionPlugin) Info() string {
	return "make an important decision easier"
}

func (x *DecisionPlugin) Usage(cmd string) string {
	switch cmd {
	case "decide":
		return "decide <yes/no question you need help with>: returns the definitive answer to your question"
	}
	return ""
}

func (x *DecisionPlugin) ProcessLine(msg *ircclient.IRCMessage) {
}

func (x *DecisionPlugin) ProcessCommand(cmd *ircclient.IRCCommand) {
	switch cmd.Command {
	case "decide":
		r := rand.Intn(7)
		switch r {
			case 0:
				x.ic.Reply(cmd, "Joah, klingt nach ner guten Idee.")
			case 1:
				x.ic.Reply(cmd, "Klar, was sollte dagegen sprechen?")
			case 2:
				x.ic.Reply(cmd, "Definitiv Ja! Super Idee!")
			case 3:
				x.ic.Reply(cmd, "Nein, würde Ich dir nicht empfehlen.")
			case 4:
				x.ic.Reply(cmd, "Wie kommst du auf so eine bescheuerte Idee?")
			case 5:
				x.ic.Reply(cmd, ":')... Nein.")
			case 6:
				x.ic.Reply(cmd, "Warum fragst du immer mich? Frag doch mal jemand anderen...")
		}
	}
}

func (x *DecisionPlugin) Unregister() {
}
