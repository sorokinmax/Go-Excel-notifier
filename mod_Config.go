package main

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// Config struct
type Config struct {
	Paths struct {
		ExcelFile string `yaml:"excelfile"`
	} `yaml:"paths"`
	SMTP struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"user"`
		Password string `yaml:"pass"`
		From     string `yaml:"from"`
		To       string `yaml:"to"`
		CC       string `yaml:"cc"`
		Subject  string `yaml:"subject"`
	} `yaml:"smtp"`
}

func readConfigFile(cfg *Config) {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func readConfigEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatal(err)
	}
}
