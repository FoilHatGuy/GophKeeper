package cfg

import (
	"encoding/json"
	"fmt"
	"os"

	defaults "github.com/mcuadros/go-defaults"
	"github.com/sakirsensoy/genv"
	_ "github.com/sakirsensoy/genv/dotenv/autoload" // import for automatic loading of .env config
)

// ConfigOption
// Various options that can be used in New() to set up configs
type ConfigOption func(*ConfigT) *ConfigT

// New
// Accepts config creation options from package.
// Returns the basic config with default values of ConfigT.
func New(opts ...ConfigOption) *ConfigT {
	cfg := &ConfigT{
		Server: &ServerT{},
		Data: &DataStorageT{
			PostgesDSN: "",
		},
	}

	for _, o := range opts {
		cfg = o(cfg)
	}

	return cfg
}

// FromDefaults
// Initializes default values of type ConfigT
func FromDefaults() ConfigOption {
	return func(c *ConfigT) *ConfigT {
		defaults.SetDefaults(c.Server)
		defaults.SetDefaults(c.Data)
		return c
	}
}

// FromJSON
// Overwrites existing values with values from environment (if present)
func FromJSON() ConfigOption {
	return func(c *ConfigT) *ConfigT {
		configPath := genv.Key("CONFIG").String()
		if configPath == "" {
			return c
		}
		data, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Printf("opening JSON failed. Details: %v\n", err)
			return nil
		}
		c2 := *c
		err = json.Unmarshal(data, &c2)
		if err != nil {
			fmt.Println(err)
			return c
		}

		return &c2
	}
}

// FromEnv
// Overwrites existing values with values from environment (if present)
func FromEnv() ConfigOption {
	return func(c *ConfigT) *ConfigT {
		c = &ConfigT{
			Server: &ServerT{
				Address:      os.ExpandEnv(genv.Key("SERVER_ADDRESS").Default(c.Server.Address).String()),
				HTTPS:        genv.Key("HTTPS").Default(c.Server.HTTPS).Bool(),
				LoggingLevel: os.ExpandEnv(genv.Key("LOGGING_LEVEL").Default(c.Server.LoggingLevel).String()),
				SessionLife:  genv.Key("SESSION_LIFE").Default(c.Server.SessionLife).Int(),
			},
			Data: &DataStorageT{
				FileSavePath: os.ExpandEnv(genv.Key("FILE_SAVE_PATH").Default(c.Data.FileSavePath).String()),
				PostgesDSN:   os.ExpandEnv(genv.Key("POSTGES_DSN").Default(c.Data.PostgesDSN).String()),
			},
		}

		return c
	}
}

// WithBuild
// Initializes default values of type ConfigT
func WithBuild(t *BuildT) ConfigOption {
	return func(c *ConfigT) *ConfigT {
		c.Build = t
		return c
	}
}
