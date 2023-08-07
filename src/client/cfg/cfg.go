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

var (
	configPath        string
	keyRSAPath        string
	serverAddressHTTP string
	serverAddressGRPC string
)

func init() {
	flag.StringVar(&configPath, "c", "", "path to JSON config. If not specified, ignores json option")
	flag.StringVar(&keyRSAPath, "k", "", "path to RSA key")
	flag.StringVar(&serverAddressHTTP, "a", "", "server's HTTP address")
	flag.StringVar(&serverAddressGRPC, "g", "", "server's GRPC address")
}

// ConfigOption
// Various options that can be used in New() to set up configs
type ConfigOption func(*ConfigT) *ConfigT

// New
// Accepts config creation options from package.
// Returns the basic config with default values of ConfigT.
func New(opts ...ConfigOption) *ConfigT {
	cfg := &ConfigT{}

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
		defaults.SetDefaults(c)
		return c
	}
}

// FromFlags
// Initializes default values of type ConfigT
func FromFlags() ConfigOption {
	return func(c *ConfigT) *ConfigT {
		if keyRSAPath != "" {
			c.RSAPath = keyRSAPath
		}
		if serverAddressHTTP != "" {
			c.ServerAddressHTTP = serverAddressHTTP
		}
		if serverAddressGRPC != "" {
			c.ServerAddressGRPC = serverAddressGRPC
		}
		return c
	}
}

// FromJSON
// Overwrites existing values with values from environment (if present)
func FromJSON() ConfigOption {
	return func(c *ConfigT) *ConfigT {
		if configPath == "" {
			configPath = genv.Key("GKEEPER_CONFIG").String()
			if configPath == "" {
				configPath = "./GophKeeperConfig.json"
				err := os.Setenv("GKEEPER_CONFIG", configPath)
				if err != nil {
					return c
				}
			}
		}
		data, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Printf("opening JSON failed. Details: %v\n", err)
			return c
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
