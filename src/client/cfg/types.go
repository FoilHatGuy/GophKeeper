package cfg

// ConfigT
// Parent structure for all configuration structs. provides config separation into
// ShortenerT, ServerT and StorageT for the ease of use
type ConfigT struct { // change from localhost to something else
	SecretPath    string  `json:"secret_path" default:"./GophKeeper.keys"`
	ServerAddress string  `json:"server_address_grpc" default:"localhost:9999"`
	ConfigPath    string  `json:"-" default:"./GophKeeperConfig.json"`
	Build         *BuildT `json:"-"`
}

// BuildT contains build info and
type BuildT struct {
	BuildVersion string `default:"N/A"`
	BuildDate    string `default:"N/A"`
	BuildCommit  string `default:"N/A"`
}
