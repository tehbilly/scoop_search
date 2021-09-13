package main

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tehbilly/scoop_search"
	"os"
)

type renderer interface {
	Render(results []*scoop_search.AppManifest) error
}

type jsonRenderer struct{}

func (r *jsonRenderer) Render(results []*scoop_search.AppManifest) error {
	out := struct {
		Results []*scoop_search.AppManifest `json:"results"`
	}{
		Results: results,
	}

	ob, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to marshal results: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(string(ob))

	return nil
}

type tableRenderer struct {
	table *tablewriter.Table
}

func newTableRenderer() renderer {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Bucket", "Version", "Description"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	return &tableRenderer{table: table}
}

func (r *tableRenderer) Render(results []*scoop_search.AppManifest) error {
	for _, result := range results {
		r.table.Append([]string{result.Name, result.Bucket, result.Version, result.Description})
	}

	r.table.Render()
	return nil
}
