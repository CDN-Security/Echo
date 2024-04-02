package config

import (
	"encoding/json"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type ChallengeConfig struct {
	SecretKey  string `json:"secret_key" yaml:"secret"`
	QueryName  string `json:"query_name" yaml:"query_name"`
	CookieName string `json:"cookie_name" yaml:"cookie_name"`
	HeaderName string `json:"header_name" yaml:"header_name"`
}

type ServerConfig struct {
	Host            string `json:"host" yaml:"host"`
	Port            int    `json:"port" yaml:"port"`
	Enable          bool   `json:"enable" yaml:"enable"`
	CertificatePath string `json:"certificate_path" yaml:"certificate_path"`
	PrivateKeyPath  string `json:"private_key_path" yaml:"private_key_path"`
}

type Config struct {
	ChallengeConfig ChallengeConfig `json:"challenge_config" yaml:"challenge_config"`
	ServerConfigs   []ServerConfig  `json:"server_configs" yaml:"server_configs"`
}

var DefaultConfig = NewConfig()

func init() {
	configFilePath := "config.yaml"
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(content, DefaultConfig)
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(DefaultConfig)
	if err != nil {
		panic(err)
	}
	slog.Info("config loaded", slog.String("file", configFilePath))
	slog.Info(string(data))
}

func NewConfig() *Config {
	return &Config{
		ChallengeConfig: ChallengeConfig{
			SecretKey:  "00000000-0000-000000000-000000000000",
			QueryName:  "echo",
			CookieName: "echo",
			HeaderName: "Echo",
		},
		ServerConfigs: []ServerConfig{
			{
				Host: "localhost",
				Port: 80,
			},
			{
				Host:            "localhost",
				Port:            443,
				CertificatePath: "assets/certificates/www.example.com/fullchain.pem",
				PrivateKeyPath:  "assets/certificates/www.example.com/privkey.pem",
			},
		},
	}
}
