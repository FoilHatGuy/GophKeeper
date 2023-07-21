package main

import (
	log "github.com/sirupsen/logrus"

	"gophKeeper/src/server/cfg"
	app "gophKeeper/src/server/server"
)

const (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	config := cfg.New(
		cfg.FromDefaults(),
		cfg.FromEnv(),
		cfg.FromJSON(),
	)

	log.SetFormatter(&log.TextFormatter{
		ForceColors:            true,
		ForceQuote:             false,
		DisableLevelTruncation: false,
		QuoteEmptyFields:       true,
	})
	lvl, err := log.ParseLevel(config.Server.LoggingLevel)
	if err != nil {
		defaultConfig := cfg.New(
			cfg.FromDefaults(),
		)
		log.Errorf("level is set wrongly, defaulting to %s", defaultConfig.Server.LoggingLevel)
		config.Server.LoggingLevel = defaultConfig.Server.LoggingLevel
		lvl, err = log.ParseLevel(config.Server.LoggingLevel)
		if err != nil {
			log.Panic("default logging value is incorrect")
		}
	}
	log.SetLevel(lvl)

	log.Infof("buildVersion\t= %q\n", buildVersion)
	log.Infof("buildDate\t= %q\n", buildDate)
	log.Infof("buildCommit\t= %q\n", buildCommit)
	log.Debug("Wow, server is running!")

	app.RunHTTPServer(config)
	app.RunGRPCServer(config)
}
