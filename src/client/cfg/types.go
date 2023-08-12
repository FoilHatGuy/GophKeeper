package cfg

// ConfigT
// Parent structure for all configuration structs. provides config separation into
// ShortenerT, ServerT and StorageT for the ease of use
type ConfigT struct { // change from localhost to something else
	SecretPath    string `json:"secret_path" default:"./GophKeeper.keys"`
	ServerAddress string `json:"server_address_grpc" default:"localhost:9999"`
	Build         *BuildT
}

// BuildT contains build info and
type BuildT struct {
	BuildVersion string
	BuildDate    string
	BuildCommit  string
}
