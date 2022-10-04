package wings

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"

	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server"
)

func New(server *server.Server, log *logrus.Logger) *wings {
	return &wings{
		*atomic.NewBool(false),
		sync.WaitGroup{},
		server,
		log,
		make([]prompt.Suggest, 0),
		&CmdlineSource{
			server: server,
			log:    log,
		},
	}
}

type wings struct {
	startedRountinue atomic.Bool
	wg               sync.WaitGroup
	server           *server.Server
	log              *logrus.Logger
	cmdsuggestions   []prompt.Suggest
	source           *CmdlineSource
}

type CmdlineSource struct {
	server *server.Server
	log    *logrus.Logger
}

func (c *CmdlineSource) Name() string {
	return "Console"
}

func (c *CmdlineSource) Position() mgl64.Vec3 {
	return mgl64.Vec3{0, 0, 0}
}

func (c *CmdlineSource) SendCommandOutput(o *cmd.Output) {
	for _, err := range o.Errors() {
		c.log.Error(err)
	}
	for _, m := range o.Messages() {
		c.log.Info(m)
	}
}

func (c *CmdlineSource) World() *world.World {
	return c.server.World()
}

func (w *wings) Executor(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}
	sp := strings.Split(s, " ")
	name := sp[0]
	args := sp[1:]
	c, o := cmd.ByAlias(name)
	if o {
		c.Execute(strings.Join(args, " "), w.source)
	} else {
		w.log.Errorf("Unkown command: %s.", name)
	}
}

func (w *wings) UpdateSuggestions() {
	for name, c := range cmd.Commands() {
		if !w.CommandSuggested(name) {
			w.cmdsuggestions = append(w.cmdsuggestions, prompt.Suggest{Text: name, Description: c.Description()})
		}
	}
}

func (w *wings) CommandSuggested(name string) bool {
	for _, suggest := range w.cmdsuggestions {
		if suggest.Text == name {
			return true
		}
	}
	return false
}

func (w *wings) Complete(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return w.cmdsuggestions
	}
	//TODO: THIS IS PRESET UP FOR ARG SUGGESTIONS.
	//args := strings.Split(d.TextBeforeCursor(), " ")
	s := d.GetWordBeforeCursor()
	if strings.HasPrefix(s, "-") {
		//return optionCompleter(args, strings.HasPrefix(s, "--"))
	}
	return prompt.FilterHasPrefix(w.cmdsuggestions, d.GetWordBeforeCursor(), true)
}

func (w *wings) Start() error {
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
	w.UpdateSuggestions()
	go p.Run()
	return nil
}
