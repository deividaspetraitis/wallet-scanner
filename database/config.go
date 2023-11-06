package database

// Config represents database configuration.
type Config struct {
	Host     string `mapstructure:"host"`     // server address
	Port     int    `mapstructure:"port"`     // server port
	Username string `mapstructure:"username"` // user
	Password string `mapstructure:"password"` // pass
	Database string `mapstructure:"database"` // database
}
