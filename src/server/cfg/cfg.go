package cfg

import (
	"encoding/json"
	"flag"
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
		Data:   &DataStorageT{},
	}

	if !flag.Parsed() {
		flag.Parse()
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
			return nil
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
				AddressHTTP:  genv.Key("SERVER_ADDRESS").Default(c.Server.AddressHTTP).String(),
				AddressGRPC:  genv.Key("SERVER_ADDRESS_GRPC").Default(c.Server.AddressGRPC).String(),
				HTTPS:        genv.Key("HTTPS").Default(c.Server.HTTPS).Bool(),
				LoggingLevel: genv.Key("LOGGING_LEVEL").Default(c.Server.AddressGRPC).String(),
			},
			Data: &DataStorageT{
				FileSavePath: genv.Key("FILE_SAVE_PATH").Default(c.Data.PostgesDSN).String(),
				PostgesDSN:   genv.Key("POSTGES_DSN").Default(c.Data.FileSavePath).String(),
			},
		}

		return c
	}
}
