package main

import (
	"context"
	"fmt"

	"github.com/go-tamate/tamate"
)

type ShowCommand struct {
	config *Config
}

func NewShowCommand(cfg *Config) (*ShowCommand, error) {
	return &ShowCommand{
		config: cfg,
	}, nil
}

func (c *ShowCommand) Run(args []string) int {
	if len(args) < 2 {
		// FIXME comment
		fmt.Println("datasource is not specified")
		return 1
	}

	datasourceName := args[0]
	taregtName := args[1]

	for _, dscfg := range c.config.Datasources {
		if dscfg.Name == datasourceName {
			if err := c.show(dscfg, taregtName); err != nil {
				fmt.Printf("Error: %s", err)
				return 1
			}
			break
		}
	}

	return 0
}

func (c *ShowCommand) Synopsis() string {
	return "Show column and row"
}

func (c *ShowCommand) Help() string {
	return "Usage: tamate show [--help] <datasource_name> <target_name>"
}

func (c *ShowCommand) show(cfg *DatasourceConfig, name string) error {
	ctx := context.Background()

	ds, err := tamate.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return err
	}
	defer ds.Close()

	schema, err := ds.GetSchema(ctx, name)
	if err != nil {
		return err
	}
	rows, err := ds.GetRows(ctx, name)
	if err != nil {
		return err
	}

	fmt.Printf("[Schema]\n")
	fmt.Printf("%v\n", schema)
	fmt.Printf("[Rows]\n")
	fmt.Printf("%v\n", rows)

	return nil
}
