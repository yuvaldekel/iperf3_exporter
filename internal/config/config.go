// Copyright 2026 Yuval Dekel
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
	"github.com/go-playground/validator/v10"
)

// Config represents the configuration file for the iperf3_exporter.
type configFile struct {
	ListenAddress string 		  		   `yaml:"listenAddress" json:"listen_address"`
	MetricsPath   string		  		   `yaml:"metricsPath" json:"metrics_path"`
	ProbePath     string		  		   `yaml:"probePath" json:"probe_path"`
	TLSCrt		  string				   `yaml:"tlsCrt" json:"tls_crt"`
	TLSKey  	  string				   `yaml:"tlsKey" json:"tls_key"`
    Interval      time.Duration   		   `yaml:"interval" json:"interval" validate:"gt=0"`
	Timeout       time.Duration	  		   `yaml:"timeout" json:"timeout"`

	// Logging configuration for the exporter
	Logging	struct {
		Level 	  string				   `yaml:"level" json:"level"`
		Format	  string				   `yaml:"format" json:"format"`
	} 									   `yaml:"logging"`

	Targets 	  []collector.TargetConfig `yaml:"targets" json:"targets" validate:"dive" default:"[]"` 
}

type argsConfig struct {
	listenAddress  string 		  
	metricsPath    string		  	
	probePath      string
	timeout        time.Duration	  	
	loggingLevel   string
	loggingFormat  string
}

// Config represents the runtime configuration for the iperf3_exporter.
type Config struct {
	ListenAddress string 		  
	MetricsPath   string		  	
	ProbePath     string
	TLSCrt		  string
	TLSKey  	  string
	Timeout       time.Duration	  	
	Targets 	  []collector.TargetConfig 
	Logger        *slog.Logger
}

func validateBitrate(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	return iperf.ValidateBitrate(val)
}

// newConfig creates a new Config with default values.
func newConfig() *configFile {
	return &configFile{
		ListenAddress: "9579",
		MetricsPath:   "/metrics",
		ProbePath:     "/probe",
		TLSCrt: 	   "",
		TLSKey: 	   "",
		Timeout:       30 * time.Second,
		Targets: 	  []collector.TargetConfig{},
		Interval:	  3600 * time.Second,
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

	configFilePath, argsConfig := parseFlags()

	// Load configuration from file if specified
	if err := loadConfigFromFile(configFilePath, configFile, argsConfig); err != nil {
		log.Fatalf("Error loading configuration from file %s: %v", configFilePath, err)
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
	}
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		cfg.Logger.Error("Invalid configuration", "err", err)
		os.Exit(1)
	}

	return cfg
}

// ParseFlags parses the command line flags and returns a Config.
func parseFlags() (string, *argsConfig){
	argsConfig := new(argsConfig)

	// Define command-line flags
	configFilePath := kingpin.Flag("config", "Path to the configuration file").
        Envar("IPERF3_EXPORTER_CONFIG_FILE").
        Default("config.yaml").
		String()
	
	kingpin.Flag("listen-address", "Port to listen on").
        Envar("IPERF3_EXPORTER_PORT").
        Default("").StringVar(&argsConfig.listenAddress)

	kingpin.Flag("metrics-path", "Path under which to expose metrics.").
		Default("").StringVar(&argsConfig.metricsPath)

	kingpin.Flag("probe-path", "Path under which to expose the probe endpoint.").
		Default("").StringVar(&argsConfig.probePath)

	kingpin.Flag("iperf3-timeout", "Timeout for each iperf3 run, in seconds.").
	    Envar("IPERF3_EXPORTER_TIMEOUT").
		Default("").DurationVar(&argsConfig.timeout)

	kingpin.Flag("log-level", "Only log messages with the given severity or above. One of: [debug, info, warn, error]").
        Envar("IPERF3_EXPORTER_LOG_LEVEL").
		Default("").StringVar(&argsConfig.loggingLevel)

	kingpin.Flag("log-format", "Output format of log messages. One of: [logfmt, json]").
		Envar("IPERF3_EXPORTER_LOG_FORMAT").
		Default("").StringVar(&argsConfig.loggingFormat)

	// Version information
	kingpin.Version(version.Print("iperf3_exporter"))
	kingpin.HelpFlag.Short('h')

	kingpin.Parse()
	return *configFilePath, argsConfig
}

// loadConfigFromFile loads the configuration from the specified file path into the provided Config struct.
func loadConfigFromFile(path string, cfg *configFile, argsCfg *argsConfig) error {
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

	// load env and args values if set
	if argsCfg.listenAddress != "" {
		cfg.ListenAddress = argsCfg.listenAddress
	}
	if argsCfg.metricsPath != "" {
		cfg.MetricsPath = argsCfg.metricsPath
	}
	if argsCfg.probePath != "" {
		cfg.ProbePath = argsCfg.probePath
	}
	if argsCfg.timeout != 0 {
		cfg.Timeout = argsCfg.timeout
	}
	if argsCfg.loggingFormat != "" {
		cfg.Logging.Format = argsCfg.loggingFormat
	}
	if argsCfg.loggingLevel != "" {
		cfg.Logging.Level = argsCfg.loggingLevel
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
            cfg.Targets[i].Interval = cfg.Interval
        }
        if cfg.Targets[i].Timeout == 0 {
            cfg.Targets[i].Timeout = cfg.Timeout
        }
	}

	var validate = validator.New()

	if err := validate.RegisterValidation("bitrate", validateBitrate); err != nil {
		return errors.New("config validation failed: " + err.Error())

	}
	
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

	return nil
}
