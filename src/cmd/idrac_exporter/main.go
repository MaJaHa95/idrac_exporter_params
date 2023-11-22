package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/mrlhansen/idrac_exporter/internal/config"
	"github.com/mrlhansen/idrac_exporter/internal/logging"
	"github.com/mrlhansen/idrac_exporter/internal/version"
)

func main() {
	var verbose bool
	var configFile string

	flag.BoolVar(&verbose, "verbose", false, "Set verbose logging")
	flag.StringVar(&configFile, "config", "/etc/prometheus/idrac.yml", "Path to idrac exporter configuration file")
	flag.Parse()

	config.ReadConfigFile(configFile)

	if !verbose {
		verbose = config.Config.Verbose
	}

	if verbose {
		logging.SetVerbose(true)

		logging.Info("Verbose mode enabled")
	}

	http.HandleFunc("/metrics", MetricsHandler)
	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("/reset", ResetHandler)
	bind := fmt.Sprintf("%s:%d", config.Config.Address, config.Config.Port)

	logging.Infof("Build information: version=%s revision=%s", version.Version, version.Revision)
	logging.Infof("Server listening on %s", bind)

	if config.Config.SingleHost != "" {
		logging.Infof("Running in single host mode. Only responding to requests for '%s'", config.Config.SingleHost)
	}

	err := http.ListenAndServe(bind, nil)
	if err != nil {
		logging.Fatal(err)
	}
}
