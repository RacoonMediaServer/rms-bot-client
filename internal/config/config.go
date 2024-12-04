package config

import (
	"github.com/RacoonMediaServer/rms-bot-client/pkg/session"
	"github.com/RacoonMediaServer/rms-packages/pkg/configuration"
)

// Remote is a settings for connection to rms-bot-server service

// Configuration represents entire service configuration
type Configuration struct {
	Device           string
	Remote           session.Endpoint
	ContentDirectory string `json:"content-directory"`
	VoiceRecognition bool   `json:"voice-recognition"`
}

var config Configuration

// Load open and parses configuration file
func Load(configFilePath string) error {
	return configuration.Load(configFilePath, &config)
}

// Config returns loaded configuration
func Config() Configuration {
	return config
}
