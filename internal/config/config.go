package config

type Config struct {
	Name string
	Env  Env
}

// singleton instance of Config
var instance *Config

// GetConfig returns the singleton instance of Config
func GetConfig() *Config {
	if instance == nil {
		instance = &Config{
			Name: "default",
			Env:  LoadEnv(),
		}
	}
	return instance
}
