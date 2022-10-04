package wings

import (
	"github.com/c-bata/go-prompt"
	"github.com/df-mc/dragonfly/server/cmd"
	"go/types"
	"strings"
)

// Argument should be implemented on custom command argument types to allow for better cmdline support.
type Argument interface {
	Description() string
}

func (w *wings) GenerateSubCommandSuggestions(c cmd.Command, args []string) []prompt.Suggest {
	params := c.Params(w.source)
	s := make([]prompt.Suggest, 0)
	for i := len(args) - 1; i >= 0; i-- {
		name := args[i]
		for _, l := range params {
			f := false
			bottomLevel := make([]prompt.Suggest, 0)
			for _, param := range l {
				if _, ok := param.Value.(cmd.SubCommand); ok {
					if param.Name == strings.TrimSpace(name) {
						f = true
						// SubSub command.
					} else if f {
						s = append(s, prompt.Suggest{Text: param.Name, Description: paramToDescription(param)})
						break
					}
				} else {
					bottomLevel = append(bottomLevel, prompt.Suggest{Text: param.Name, Description: paramToDescription(param)})
				}
			}
			if f {
				s = append(s, bottomLevel[:]...)
			}
		}
	}
	return s
}

func (w *wings) GenerateArgSuggestions(c cmd.Command) []prompt.Suggest {
	params := c.Params(w.source)
	s := make([]prompt.Suggest, 0)
	for _, l := range params {
		if len(l) == 0 {
			continue
		}
		param := l[0]
		s = append(s, prompt.Suggest{
			Text:        param.Name,
			Description: paramToDescription(param),
		})
	}
	return s
}

func paramToDescription(param cmd.ParamInfo) string {
	d := ""
	switch param.Value.(type) {
	case cmd.SubCommand:
		d = "Subcommand"
	case int64:
		d = "Number"
	case int32:
		d = "Number"
	case int16:
		d = "Number"
	case int:
		d = "Number"
	case uint64:
		d = "Positive Number"
	case uint32:
		d = "Positive Number"
	case uint16:
		d = "Positive Number"
	case uint:
		d = "Positive Number"
	case uint8:
		d = "One byte"
	case int8:
		d = "One byte"
	case float64:

	case float32:
		d = "Decimal Number"
	case string:
		d = "string"
	case types.Slice:
		x := param.Value.(types.Slice)
		d = "List of " + x.Elem().String()
	case Argument:
		d = param.Value.(Argument).Description()
	default:
		d = "No Description/Type found."
	}
	if param.Optional {
		d = "Optional " + d
	}
	return d
}
