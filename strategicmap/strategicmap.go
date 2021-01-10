package strategicmap

import "math"

type Territory struct {
	Name   string
	LeftPX int
	TopPX  int
}

type Map struct {
	Territories []Territory
}

const MAPWIDTHPX = 1536
const MAPHEIGHTPX = 723

var StaticMap = Map{Territories: []Territory{
	{"A", 415, 95},
	{"B", 580, 40},
	{"C", 705, 100},
	{"D", 815, 55},
	{"E", 865, 145},
	{"F", 985, 115},
	{"G", 1100, 130},
	{"H", 1170, 60},
	{"I", 470, 240},
	{"J", 650, 190},
	{"K", 780, 255},
	{"L", 1020, 270},
	{"M", 625, 335},
	{"N", 800, 450},
	{"O", 970, 445},
	{"P", 560, 490},
	{"Q", 580, 630},
	{"R", 840, 645},
}}

func divideAndRoundToNearestInteger(numerator int, denominator int) int {
	return int(math.Round(float64(numerator) / float64(denominator) * 100))
}

func (territory Territory) LeftAsPercent() int {
	return divideAndRoundToNearestInteger(territory.LeftPX, MAPWIDTHPX)
}

func (territory Territory) TopAsPercent() int {
	return divideAndRoundToNearestInteger(territory.TopPX, MAPHEIGHTPX)
}
