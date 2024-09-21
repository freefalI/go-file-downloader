package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
)

func setupTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "URL", "Name", "Started at", "Finished at", "Size", "Progress", "Done"})
	table.SetBorder(true)

	return table
}

func printEntries(table *tablewriter.Table) {
	table.ClearRows()

	mu.Lock() // Lock around reading files
	defer mu.Unlock()

	for _, file := range files {
		file.mu.Lock()
		var finishedAt string
		if !file.FinishedAt.IsZero() {
			finishedAt = file.FinishedAt.Format(time.DateTime)
		}

		table.Append([]string{
			strconv.Itoa(file.ID),
			file.URL,
			file.Name,
			file.StartedAt.Format(time.DateTime),
			finishedAt,
			fmt.Sprintf("%dMb", file.Size),
			fmt.Sprintf("%d%%", file.Progress),
			fmt.Sprintf("%t", file.IsDone),
		})
		file.mu.Unlock()
	}

	clearTerminal()
	table.Render()
}

func clearTerminal() {
	os.Stdout.WriteString("\033[H\033[2J")
}
