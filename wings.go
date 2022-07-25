package wings

import (
	"fmt"
	"strings"
	"sync"

	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server"
)
func New(server server.Server) wings {
    return wings{*atomic.NewBool(false),sync.WaitGroup{}, server}
}

type wings struct {
    startedRountinue atomic.Bool
    wg sync.WaitGroup
    server server.Server
}

func (w wings) Executor(s string) {
    s = strings.TrimSpace(s)
	if s == "" {
		return
	}
    fmt.Println(s)
}

func (w wings) Complete(d prompt.Document) []prompt.Suggest {
    s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
} 

func (w wings) Start() error {
    if !w.startedRountinue.CAS(false, true) {
		panic("wings is already processing commands")
	}
    p := prompt.New(
        w.Executor,
        w.Complete,
        prompt.OptionTitle("Dragonfly: A safe and fast MCBE server software"),
        prompt.OptionPrefix("$ "),
        prompt.OptionPrefixTextColor(prompt.Yellow),
        prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
    )
    go p.Run()
    return nil
}












