package wings

type Config struct {
	Prefix           string
	Title            string
	RegisterDefaults bool
	MaxSuggestions   uint16
}

func DefaultConfig() Config {
	return Config{Prefix: ">>> ", Title: "Dragonfly: A safe and fast MCBE server software", RegisterDefaults: true, MaxSuggestions: 10}
}
