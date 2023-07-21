package cfg

// ConfigT
// Parent structure for all configuration structs. provides config separation into
// ShortenerT, ServerT and StorageT for the ease of use
type ConfigT struct { // change from localhost to something else
	RSAPath           string `json:"rsa_path" default:"./GophKeeper.keys"`
	ServerAddressHTTP string `json:"server_address_http" default:"localhost:3000"`
	ServerAddressGRPC string `json:"server_address_grpc" default:"localhost:9999"`
}
