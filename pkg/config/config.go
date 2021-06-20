package config

// Credential defines a single project credential.
type Credential struct {
	Project  string `json:"project" yaml:"project"`
	URL      string `json:"url" yaml:"url"`
	Insecure bool   `json:"insecure" yaml:"insecure"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Org      string `json:"org" yaml:"org"`
	Vdc      string `json:"vdc" yaml:"vdc"`
}

// Server defines the general server configuration.
type Server struct {
	Addr string `json:"addr" yaml:"addr"`
	Path string `json:"path" yaml:"path"`
	Web  string `json:"web_config" yaml:"web_config"`
}

// Logs defines the level and color for log configuration.
type Logs struct {
	Level  string `json:"level" yaml:"level"`
	Pretty bool   `json:"pretty" yaml:"pretty"`
}

// Target defines the target specific configuration.
type Target struct {
	Engine      string       `json:"engine" yaml:"engine"`
	File        string       `json:"file" yaml:"file"`
	Refresh     int          `json:"refresh" yaml:"refresh"`
	Credentials []Credential `json:"credentials" yaml:"credentials"`
}

// Config is a combination of all available configurations.
type Config struct {
	Server Server `json:"server" yaml:"server"`
	Logs   Logs   `json:"logs" yaml:"logs"`
	Target Target `json:"target" yaml:"target"`
}

// Load initializes a default configuration struct.
func Load() *Config {
	return &Config{
		Target: Target{
			Credentials: make([]Credential, 0),
		},
	}
}
