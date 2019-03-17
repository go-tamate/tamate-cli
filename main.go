package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mitchellh/cli"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/go-tamate/tamate-mysql"

	// spreadsheet driver
	_ "github.com/go-tamate/tamate-spreadsheet"
)

type DatasourceConfig struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
	DSN    string `json:"dsn"`
}

type Config struct {
	Datasources []*DatasourceConfig `json:"datasources"`
}

func getConfig(path string) (*Config, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.NewDecoder(r).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	c := cli.NewCLI("tamate", "0.0.1")

	// FIXME: ファイルパスを可変式に
	config, err := getConfig("./config.json")
	if err != nil {
		panic(err)
	}

	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"show": func() (cli.Command, error) {
			return NewShowCommand(config)
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(exitStatus)
}
