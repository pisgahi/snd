package sndcfg

import (
	"encoding/json"
	"log"
	"os"

	sndCli "github.com/concernum/snd/cli"
)

type Config struct {
	ServerAddr  string `json:"addr"`
	ReceivedDir string `json:"dir"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func SaveConfig(filename string, config *Config) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

func SetServerAddr(filename, serverAddr string) error {
	config, err := LoadConfig(filename)
	if err != nil {
		return err
	}

	config.ServerAddr = serverAddr
	return SaveConfig(filename, config)
}

func SetReceivedDir(filename, receivedDir string) error {
	config, err := LoadConfig(filename)
	if err != nil {
		return err
	}

	config.ReceivedDir = receivedDir
	return SaveConfig(filename, config)
}

func LoadOrCreateConfig(flags *sndCli.Flags) *Config {
	configFile := ".config.json"

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if flags.ServerAddr == "" || flags.ReceivedDir == "" {
			log.Fatal("Config file does not exist.")
		}

		config := &Config{
			ServerAddr:  flags.ServerAddr,
			ReceivedDir: flags.ReceivedDir,
		}

		if err := SaveConfig(configFile, config); err != nil {
			log.Fatalf("Failed to create configuration file: %v", err)
		}
		log.Println("Configuration file created successfully.")
		os.Exit(0)
	}

	config, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if flags.ServerAddr != "" {
		if err := SetServerAddr(configFile, flags.ServerAddr); err != nil {
			log.Fatalf("Failed to update server address: %v", err)
		}
	}
	if flags.ReceivedDir != "" {
		if err := SetReceivedDir(configFile, flags.ReceivedDir); err != nil {
			log.Fatalf("Failed to update received directory: %v", err)
		}
	}

	return config
}
