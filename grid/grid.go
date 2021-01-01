package grid

import (
	"fmt"
	"html/template"
	"strconv"
	"unicode"

	"github.com/brickman1444/NSImperialism/nationstates_api"
	"github.com/brickman1444/NSImperialism/war"
)

const NUMROWS = 5
const NUMCOLUMNS = 5

type Cell struct {
	ResidentNation *nationstates_api.Nation
}

type Row struct {
	Cells [NUMCOLUMNS]Cell
}

type Grid struct {
	Rows [NUMCOLUMNS]Row
	Year int
}

type RenderedCell struct {
	Text template.HTML
}

type RenderedRow struct {
	Cells [NUMCOLUMNS]RenderedCell
}

type RenderedGrid struct {
	Rows [NUMCOLUMNS]RenderedRow
	Year int
}

func toCharStr(i int) string {
	return string(rune('A' - 1 + i))
}

func toIndex(r rune) (int, error) {
	if unicode.IsLetter(r) || unicode.IsUpper(r) {
		return int(r - 'A' + 1), nil
	} else if unicode.IsDigit(r) {
		return strconv.Atoi(string(r))
	} else {
		return -1, fmt.Errorf("characters out of valid range")
	}
}

func (grid *Grid) Render(wars []*war.War) *RenderedGrid {
	renderedGrid := &RenderedGrid{}
	renderedGrid.Year = grid.Year

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

				war := war.FindOngoingWarAt(wars, rowIndex, columnIndex)

				if war != nil {
					cellText = cellText + "⚔️" + war.Attacker.FlagThumbnail()
				}

				renderedGrid.Rows[rowIndex].Cells[columnIndex].Text = template.HTML(cellText)

				continue
			}

			renderedGrid.Rows[rowIndex].Cells[columnIndex].Text = template.HTML("❓")
		}
	}
	return renderedGrid
}

func (grid *Grid) GetCoordinates(coordinatesString string) (int, int, error) {

	if len(coordinatesString) != 2 {
		return 0, 0, fmt.Errorf("Invalid string length %s", coordinatesString)
	}

	firstRune := rune(coordinatesString[0])
	if firstRune < 'A' || firstRune > 'D' {
		return 0, 0, fmt.Errorf("First character isn't an upper case character %s", coordinatesString)
	}

	secondRune := rune(coordinatesString[1])
	if secondRune < '1' || secondRune > '4' {
		return 0, 0, fmt.Errorf("Second character isn't a number %s", coordinatesString)
	}

	columnIndex, err := toIndex(firstRune)
	if err != nil {
		return 0, 0, err
	}

	rowIndex, err := toIndex(secondRune)
	if err != nil {
		return 0, 0, err
	}

	return rowIndex, columnIndex, nil
}

func (grid *Grid) Colonize(colonizer nationstates_api.Nation, target string) error {

	rowIndex, columnIndex, err := grid.GetCoordinates(target)
	if err != nil {
		return err
	}

	if grid.Rows[rowIndex].Cells[columnIndex].ResidentNation != nil {
		return fmt.Errorf("A nation is already resident at %s", target)
	}

	grid.Rows[rowIndex].Cells[columnIndex].ResidentNation = &colonizer

	return nil
}
