package sndcfg

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	sndCli "github.com/pisgahi/snd/cli"
)

type Config struct {
	ServerAddr  string `json:"addr"`
	ReceivedDir string `json:"dir"`
	CertDir     string `json:"certDir"`
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
	configFile := "sndcfg/.config.json"

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if flags.ServerAddr == "" || flags.ReceivedDir == "" || flags.CertDir == "" {
			log.Fatal("Config file does not exist and required flags are not provided.")
		}

		config := &Config{
			ServerAddr:  flags.ServerAddr,
			ReceivedDir: flags.ReceivedDir,
			CertDir:     flags.CertDir,
		}

		if err := SaveConfig(configFile, config); err != nil {
			log.Fatalf("Failed to create configuration file: %v", err)
		}
		fmt.Println("Configuration file created successfully.")
		os.Exit(0)
	}

	config, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if config.CertDir == "" {
		log.Fatal("Certificate directory is not configured.")
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
	if flags.CertDir != "" {
		if err := SetCertDir(configFile, flags.CertDir); err != nil {
			log.Fatalf("Failed to update certificate directory: %v", err)
		}
	}

	return config
}

func SetCertDir(filename, certDir string) error {
	config, err := LoadConfig(filename)
	if err != nil {
		return err
	}

	config.CertDir = certDir
	return SaveConfig(filename, config)
}
