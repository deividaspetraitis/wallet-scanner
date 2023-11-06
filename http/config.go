package http

// Config represents HTTP server configuration.
type Config struct {
	Address    string `mapstructure:"address"` // HTTP server address
	Middleware struct {
		RateLimit int `mapstructure:"ratelimit"`
	} `mapstructure:"middleware"`
}
