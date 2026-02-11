// Copyright 2019 Edgard Castro
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package config provides configuration handling for the iperf3_exporter.
package config

import (
	"errors"
	"log/slog"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
	"github.com/yuvaldekel/iperf3_exporter/internal/collector"
	"github.com/yuvaldekel/iperf3_exporter/internal/iperf"
	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/go-playground/validator/v10"
)
const (
	// defaultTimeout is the default timeout for iperf3 tests when no timeout is configured.
	defaultTimeout = 30.0
)

// Config represents the configuration file for the iperf3_exporter.
type ConfigFile struct {
	ListenAddress string 		  		   `yaml:"listenAddress" json:"listen_address"`
	MetricsPath   string		  		   `yaml:"metricsPath" json:"metrics_path"`
	ProbePath     string		  		   `yaml:"probePath" json:"probe_path"`
	Timeout       time.Duration	  		   `yaml:"timeout" json:"timeout"`
	Targets 	  []collector.TargetConfig `yaml:"targets" json:"targets" validate:"dive",default:"[]"` 
	// Logging configuration for the exporter
	Logging	struct {
		Level 	  string				   `yaml:"level" json:"level"`
		Format	  string				   `yaml:"format" json:"format"`
	} 									   `yaml:"logging"`
}

// Config represents the runtime configuration for the iperf3_exporter.
type Config struct {
	ListenAddress string 		  
	MetricsPath   string		  	
	ProbePath     string		  	
	Timeout       time.Duration	  	
	Targets 	  []collector.TargetConfig 
	Logger        *slog.Logger
	WebConfig     *web.FlagConfig
}

// newConfig creates a new Config with default values.
func newConfig() *ConfigFile {
	return &ConfigFile{
		ListenAddress: "9579",
		MetricsPath:   "/metrics",
		ProbePath:     "/probe",
		Timeout:       30 * time.Second,
		Targets: 	  []collector.TargetConfig{},
		Logging: struct {
			Level  string `yaml:"level" json:"level"`
			Format string `yaml:"format" json:"format"`
		}{
			Level:  "info",
			Format: "logfmt",
		},
	}
}

// LoadConfig loads the configuration from command-line flags and optionally from a configuration file.
func LoadConfig() *Config {
	configFile := newConfig()
	configPath := parseFlags(configFile)

	// Load configuration from file if specified
	if err := loadConfigFromFile(configPath, configFile); err != nil {
		log.Fatalf("Error loading configuration from file: %v", err)
	}

	log.Println("%s", configFile.ListenAddress)
	// Create web configuration for the exporter
	webConfig := &web.FlagConfig{
        WebListenAddresses: &[]string{":" + configFile.ListenAddress},
        WebSystemdSocket:   new(bool),
        WebConfigFile:      new(string),
    }

	// Initialize logger
	var logLevelSlog slog.Level

	switch configFile.Logging.Level {
	case "debug":
		logLevelSlog = slog.LevelDebug
	case "info":
		logLevelSlog = slog.LevelInfo
	case "warn":
		logLevelSlog = slog.LevelWarn
	case "error":
		logLevelSlog = slog.LevelError
	default:
		logLevelSlog = slog.LevelInfo
	}

	var handler slog.Handler
	if configFile.Logging.Format == "json" {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: logLevelSlog})
	} else {
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevelSlog})
	}

	logger := slog.New(handler)

	cfg := &Config{
		ListenAddress: configFile.ListenAddress,
		MetricsPath:   configFile.MetricsPath,
		ProbePath:     configFile.ProbePath,
		Timeout:       configFile.Timeout,
		Targets: 	   configFile.Targets,
		Logger:        logger,
		WebConfig:     webConfig,
	}
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		cfg.Logger.Error("Invalid configuration", "err", err)
		os.Exit(1)
	}

	return cfg
}

// ParseFlags parses the command line flags and returns a Config.
func parseFlags(cfg *ConfigFile) string {

	// Define command-line flags
	configFilePath := kingpin.Flag("config", "Path to the configuration file").
        Envar("IPERF3_EXPORTER_CONFIG_FILE").
        Default("config.yaml").
		String()

	kingpin.Flag("listen-address", "Port to listen on").
        Envar("IPERF3_EXPORTER_PORT").
        Default(cfg.ListenAddress).StringVar(&cfg.ListenAddress)

	kingpin.Flag("metrics-path", "Path under which to expose metrics.").
		Default(cfg.MetricsPath).StringVar(&cfg.MetricsPath)

	kingpin.Flag("probe-path", "Path under which to expose the probe endpoint.").
		Default(cfg.ProbePath).StringVar(&cfg.ProbePath)

	kingpin.Flag("iperf3-timeout", "Timeout for each iperf3 run, in seconds.").
	    Envar("IPERF3_EXPORTER_TIMEOUT").
		Default(cfg.Timeout.String()).DurationVar(&cfg.Timeout)

	kingpin.Flag("log-level", "Only log messages with the given severity or above. One of: [debug, info, warn, error]").
        Envar("IPERF3_EXPORTER_LOG_LEVEL").
		Default(cfg.Logging.Level).StringVar(&cfg.Logging.Level)

	kingpin.Flag("log-format", "Output format of log messages. One of: [logfmt, json]").
		Envar("IPERF3_EXPORTER_LOG_FORMAT").
		Default(cfg.Logging.Format).StringVar(&cfg.Logging.Format)

	// Version information
	kingpin.Version(version.Print("iperf3_exporter"))
	kingpin.HelpFlag.Short('h')

	kingpin.Parse()

	return *configFilePath
}

// loadConfigFromFile loads the configuration from the specified file path into the provided Config struct.
func loadConfigFromFile(path string, cfg *ConfigFile) error {
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return errors.New("error reading config file: " + err.Error())
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return errors.New("error unmarshaling config file: " + err.Error())
	}

	for i := range cfg.Targets {
        if cfg.Targets[i].Port == 0 {
            cfg.Targets[i].Port = 5201
        }
        if cfg.Targets[i].Protocol == "" {
            cfg.Targets[i].Protocol = "tcp"
        }
        if cfg.Targets[i].Period == 0 {
            cfg.Targets[i].Period = 5 * time.Second
        }
        if cfg.Targets[i].Interval == 0 {
            cfg.Targets[i].Interval = 3600 * time.Second
        }
        if cfg.Targets[i].Timeout == 0 {
            cfg.Targets[i].Timeout = defaultTimeout * time.Second
        }
	}

	var validate = validator.New()
	validate.RegisterValidation("bitrate", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return iperf.ValidateBitrate(val)
    })

	if err := validate.Struct(cfg); err != nil {
        return errors.New("config validation failed: " + err.Error())
    }

	return nil
}


// Validate validates the configuration.
func (c *Config) Validate() error {
	// Validate basic configuration
	if c.MetricsPath == "" {
		return errors.New("metrics path cannot be empty")
	}

	if c.ProbePath == "" {
		return errors.New("probe path cannot be empty")
	}

	if c.Timeout <= 0 {
		return errors.New("timeout must be greater than 0")
	}

	if c.Logger == nil {
		return errors.New("logger cannot be nil")
	}

	// Validate web configuration
	if c.WebConfig == nil {
		return errors.New("web configuration cannot be nil")
	}

	// Note: Additional web config validation is handled by web.ListenAndServe
	// which checks for listen addresses and validates TLS config if provided

	return nil
}
