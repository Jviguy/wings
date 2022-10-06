package wings

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"

	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server"
)

func New(server *server.Server, log *logrus.Logger, config Config) *wings {
	if config.RegisterDefaults {
		RegisterDefaults()
	}
	return &wings{
		*atomic.NewBool(false),
		sync.WaitGroup{},
		server,
		log,
		make([]prompt.Suggest, 0),
		&ConsoleSource{
			server: server,
			log:    log,
		},
		config,
	}
}

type wings struct {
	started        atomic.Bool
	wg             sync.WaitGroup
	server         *server.Server
	log            *logrus.Logger
	cmdsuggestions []prompt.Suggest
	source         *ConsoleSource
	config         Config
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
	args := strings.Split(d.TextBeforeCursor(), " ")
	//if there is args
	if len(args) > 1 {
		name := args[0]
		args = args[1:]

		c, o := cmd.ByAlias(name)
		if o {
			as := make([]prompt.Suggest, 0)
			if len(args) > 1 {
				as = w.GenerateSubCommandSuggestions(c, args)
			} else {
				as = w.GenerateArgSuggestions(c)
			}
			s := make([]string, 0)
			for _, x := range as {
				s = append(s, x.Text)
			}
			_, d := FindLevenshtein(args[len(args)-1], s)
			if d < 5 {
				return prompt.FilterFuzzy(as, args[len(args)-1], true)
			} else {
				return as
			}
		}
	}
	return prompt.FilterHasPrefix(w.cmdsuggestions, d.GetWordBeforeCursor(), true)
}

func (w *wings) Start() error {
	if !w.started.CAS(false, true) {
		panic("wings is already processing commands")
	}
	p := prompt.New(
		w.Executor,
		w.Complete,
		prompt.OptionTitle(w.config.Title),
		prompt.OptionPrefix(w.config.Prefix),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionMaxSuggestion(w.config.MaxSuggestions),
		prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
	)
	w.UpdateSuggestions()
	go p.Run()
	return nil
}
