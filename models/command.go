package models
type Command struct {
	Description string
	Option      map[string]string // it's not used now
}

var Commands = map[string]Command{
	"timeline": Command{
		Description: "Get timeline",
		Option:      make(map[string]string),
	},
	"tweet": Command{
		Description: "Post tweet",
		Option:      make(map[string]string),
	},
	"clear": Command{
		Description: "Clear console",
		Option:      make(map[string]string),
	},
	"exit": Command{
		Description: "Terminate this cli client",
		Option:      make(map[string]string),
	},
	"help": Command{
		Description: "Show how to use",
		Option:      make(map[string]string),
	},
}