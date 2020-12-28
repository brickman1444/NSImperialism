package grid

import "fmt"

type Cell struct {
	Text string
}

type Row struct {
	Cells [4]Cell
}

type Grid struct {
	Rows [4]Row
}

func Get() *Grid {
	grid := &Grid{}

	for rowIndex, _ := range grid.Rows {
		for columnIndex, _ := range grid.Rows[rowIndex].Cells {
			grid.Rows[rowIndex].Cells[columnIndex].Text = fmt.Sprintf("(%d,%d)", rowIndex, columnIndex)
		}
	}
	return grid
}
