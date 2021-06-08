package main

import (
	"log"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// Config struct
type Config struct {
	Common struct {
		NotifyForDays             time.Duration `yaml:"notify-for-days"`
		AdminsEmails              []string      `yaml:"admins-emails"`
		TableCaption              string        `yaml:"table-caption"`
		TableHeaderNameColumn     string        `yaml:"table-header-name-column"`
		TableHeaderCheckingColumn string        `yaml:"table-header-checking-column"`
		DateFormat                string        `yaml:"date-format"`
	} `yaml:"common"`
	Excel struct {
		File                        string   `yaml:"file"`
		Sheet                       string   `yaml:"sheet"`
		NameColumn                  string   `yaml:"name-column"`
		CheckingColumn              string   `yaml:"checking-column"`
		CheckingRowStart            int      `yaml:"checking-row-start"`
		CheckingRowEnd              int      `yaml:"checking-row-end"`
		PersonalNotificationBundles []string `yaml:"personal-notification-bundles"`
		PersonalNotificationEmails  []string `yaml:"personal-notification-emails"`
	} `yaml:"excel"`
	SMTP struct {
		Host     string   `yaml:"host"`
		Port     int      `yaml:"port"`
		Username string   `yaml:"user"`
		Password string   `yaml:"pass"`
		From     string   `yaml:"from"`
		To       []string `yaml:"to"`
		CC       string   `yaml:"cc"`
		Subject  string   `yaml:"subject"`
	} `yaml:"smtp"`
}

func readConfigFile(cfg *Config) {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Fatalln(err)
	}
}

func readConfigEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatalln(err)
	}
}
