package strategicmap

import "math"

type Territory struct {
	LeftPX int
	TopPX  int
}

type Map struct {
	Territories []Territory
}

const MAPWIDTHPX = 1536
const MAPHEIGHTPX = 723

var StaticMap = Map{Territories: []Territory{
	{415, 95},
	{580, 40},
	{705, 100},
	{815, 55},
	{865, 145},
	{985, 115},
	{1100, 130},
	{1170, 60},
	{470, 240},
	{650, 190},
	{780, 255},
	{1020, 270},
	{625, 335},
	{800, 450},
	{970, 445},
	{560, 490},
	{580, 630},
	{840, 645},
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
