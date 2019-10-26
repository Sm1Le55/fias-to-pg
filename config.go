package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Username string `yaml:"user"`
		Password string `yaml:"pass"`
		Address  string `yaml:"address"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
	DBF struct {
		Codepage string `yaml:"codepage"`
		Folder   string `yaml:"folder"`
		Threads  int    `yaml:"threads"`
	} `yaml:"dbf"`
	Program struct {
		Update             bool `yaml:"update"`
		DropTables         bool `yaml:"droptables"`
		CreateTables       bool `yaml:"createtables"`
		DropIrrelevantRows bool `yaml:"dropirrelevantrows"`
	} `yaml:"program"`
	Log struct {
		Path string `yaml:"path"`
	} `yaml:"log"`
}

const configName = "config.yml"

func getConfig(configName string) {
	f, err := os.Open(configName)
	if err != nil {
		panic(fmt.Sprintf("Open config file error: %v\n", err))
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		panic(fmt.Sprintf("Decode config file error: %v\n", err))
	}
}
