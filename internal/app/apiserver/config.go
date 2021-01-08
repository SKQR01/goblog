package apiserver

// Config of the server.
type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
	SessionKey  string `toml:"session_key"`
}

// NewConfig of server (if no bind_addr in apiserver.toml then sets the :8080 value of BindAddr(port) and debug logLevel by default).
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}
