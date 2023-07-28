package cfg

// ConfigT
// Parent structure for all configuration structs. provides config separation into
// ShortenerT, ServerT and StorageT for the ease of use
type ConfigT struct {
	Server *ServerT      `json:"server"`
	Data   *DataStorageT `json:"data"`
}

// ServerT stores Server side related configuration for both GRPC and HTTP.
// Can be accessed via a structure of type ConfigT
type ServerT struct {
	AddressHTTP  string `default:"localhost:3000" json:"address_http"`
	AddressGRPC  string `default:"localhost:9999" json:"address_grpc"`
	HTTPS        bool   `default:"false" json:"enable_https"`
	LoggingLevel string `default:"Debug" json:"logging_level"` // exactly like in logrus
}

// DataStorageT stores Data storage related configuration.
// DatabaseDSN contains string used for connection to Postgres and Redis.
// Can be accessed via a structure of type ConfigT
type DataStorageT struct {
	FileSavePath string `default:"./fileData" json:"file_save_path"`
	PostgesDSN   string `default:"" json:"postges_dsn"`
}
