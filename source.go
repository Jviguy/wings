package wings

import (
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
)

type ConsoleSource struct {
	server *server.Server
	log    *logrus.Logger
}

func (c *ConsoleSource) Name() string {
	return "Console"
}

func (c *ConsoleSource) Position() mgl64.Vec3 {
	return mgl64.Vec3{0, 0, 0}
}

func (c *ConsoleSource) SendCommandOutput(o *cmd.Output) {
	for _, err := range o.Errors() {
		c.log.Error(text.ANSI(err))
	}
	for _, m := range o.Messages() {
		c.log.Info(text.ANSI(m))
	}
}

func (c *ConsoleSource) World() *world.World {
	return c.server.World()
}

func (c *ConsoleSource) GetServer() *server.Server {
	return c.server
}

func (c *ConsoleSource) GetLog() *logrus.Logger {
	return c.log
}
