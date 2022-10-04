package wings

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/sirupsen/logrus"
)

func RegisterDefaults() {
	for _, c := range []cmd.Command{
		cmd.New("exit", "Closes the server.", []string{}, Exit{}),
	} {
		cmd.Register(c)
	}
}

type Exit struct{}

func (e Exit) Allow(src cmd.Source) bool {
	_, o := src.(*ConsoleSource)
	return o
}

func (e Exit) Run(src cmd.Source, o *cmd.Output) {
	if c, ok := src.(*ConsoleSource); ok {
		o.Printf("Shutting down.")
		err := c.GetServer().Close()
		if err != nil {
			o.Errorf("Error shutting down: %s", err)
			logrus.Exit(1)
			return
		}
		logrus.Exit(0)
	}
}
