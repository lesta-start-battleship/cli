package ui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type Table struct {
	width       int
	colWidths   []int
	rows        [][]string
	header      []string
	headerStyle lipgloss.Style
	rowStyle    lipgloss.Style
}

func NewTable(totalWidth int, colWidths []int) *Table {
	return &Table{
		width:     totalWidth,
		colWidths: colWidths,
		headerStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).Bold(true),
		rowStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0E0E0")).Bold(true),
	}
}

func (t *Table) AddHeader(headers []string) {
	t.header = headers
}

func (t *Table) AddRow(row []string) {
	t.rows = append(t.rows, row)
}

func (t *Table) Render() string {
	var sb strings.Builder

	if len(t.header) > 0 {
		var headerParts []string
		for i, h := range t.header {
			width := t.colWidths[i]
			headerParts = append(headerParts, t.headerStyle.Render(fmt.Sprintf("%-*s", width, h)))
		}
		sb.WriteString(strings.Join(headerParts, " "))
		sb.WriteString("\n")

		sb.WriteString(strings.Repeat("-", t.width))
		sb.WriteString("\n")
	}

	for _, row := range t.rows {
		var rowParts []string
		for i, cell := range row {
			width := t.colWidths[i]
			rowParts = append(rowParts, t.rowStyle.Render(fmt.Sprintf("%-*s", width, cell)))
		}
		sb.WriteString(strings.Join(rowParts, " "))
		sb.WriteString("\n")
	}

	return sb.String()
}
