package grid

import (
	"fmt"
	"html/template"

	"github.com/brickman1444/NSImperialism/nationstates_api"
)

const NUMROWS = 5
const NUMCOLUMNS = 5

type Cell struct {
	ResidentNation *nationstates_api.Nation
	AttackerNation *nationstates_api.Nation
}

type Row struct {
	Cells [NUMCOLUMNS]Cell
}

type Grid struct {
	Rows [NUMCOLUMNS]Row
}

type RenderedCell struct {
	Text template.HTML
}

type RenderedRow struct {
	Cells [NUMCOLUMNS]RenderedCell
}

type RenderedGrid struct {
	Rows [NUMCOLUMNS]RenderedRow
}

func toCharStr(i int) string {
	return string(rune('A' - 1 + i))
}

func (grid *Grid) Render() *RenderedGrid {
	renderedGrid := &RenderedGrid{}

	for rowIndex, _ := range grid.Rows {
		for columnIndex, _ := range grid.Rows[rowIndex].Cells {
			if rowIndex == 0 && columnIndex == 0 {
				continue
			}

			if rowIndex == 0 {
				renderedGrid.Rows[rowIndex].Cells[columnIndex].Text = template.HTML(toCharStr(columnIndex))
				continue
			}

			if columnIndex == 0 {
				renderedGrid.Rows[rowIndex].Cells[columnIndex].Text = template.HTML(fmt.Sprintf("%d", rowIndex))
				continue
			}

			if grid.Rows[rowIndex].Cells[columnIndex].ResidentNation != nil {

				cellText := grid.Rows[rowIndex].Cells[columnIndex].ResidentNation.FlagThumbnail()
				if grid.Rows[rowIndex].Cells[columnIndex].AttackerNation != nil {
					cellText = cellText + "⚔️" + grid.Rows[rowIndex].Cells[columnIndex].AttackerNation.FlagThumbnail()
				}

				renderedGrid.Rows[rowIndex].Cells[columnIndex].Text = template.HTML(cellText)

				continue
			}

			renderedGrid.Rows[rowIndex].Cells[columnIndex].Text = template.HTML("❓")
		}
	}
	return renderedGrid
}
