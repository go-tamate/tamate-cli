package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-tamate/tamate"
	"github.com/go-tamate/tamate/driver"
	"github.com/olekukonko/tablewriter"
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

	if err := showSchema(schema); err != nil {
		return err
	}
	if err := showRows(rows); err != nil {
		return err
	}
	return nil
}

func showSchema(sc *driver.Schema) error {
	if len(sc.Columns) == 0 {
		return nil
	}

	columnHeaderData := []string{"Name", "Type", "NotNull", "AutoIncrement"}

	columnBodyData := make([][]string, 0)
	for _, col := range sc.Columns {
		notNull := "FALSE"
		if col.NotNull {
			notNull = "TRUE"
		}
		autoIncrement := "FALSE"
		if col.AutoIncrement {
			autoIncrement = "TRUE"
		}
		columnBodyData = append(columnBodyData, []string{col.Name, col.Type.String(), notNull, autoIncrement})
	}

	fmt.Printf("[Column]\n")
	return showTable(columnHeaderData, columnBodyData)
}

func showRows(rows []*driver.Row) error {
	if len(rows) == 0 {
		return nil
	}

	headerData := rows[0].Values.ColumnNames()

	bodyData := make([][]string, 0)
	for _, row := range rows {
		data := make([]string, 0)
		for _, rowVal := range row.Values {
			data = append(data, rowVal.String())
		}
		bodyData = append(bodyData, data)
	}

	fmt.Printf("[Rows]\n")
	return showTable(headerData, bodyData)
}

func showTable(header []string, body [][]string) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetHeader(header)
	table.AppendBulk(body)
	table.Render()
	return nil
}
