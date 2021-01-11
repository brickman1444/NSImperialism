package strategicmap

import (
	"fmt"
	"html/template"
	"math"

	"github.com/brickman1444/NSImperialism/nationstates_api"
	"github.com/brickman1444/NSImperialism/war"
)

type Territory struct {
	Name   string
	LeftPX int
	TopPX  int
}

type Map struct {
	Territories []Territory
}

type RenderedTerritory struct {
	Text        template.HTML
	LeftPercent int
	TopPercent  int
}

type RenderedMap struct {
	Territories []RenderedTerritory
}

type Ownerships map[string]nationstates_api.Nation

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

func (territory Territory) LeftPercent() int {
	return divideAndRoundToNearestInteger(territory.LeftPX, MAPWIDTHPX)
}

func (territory Territory) TopPercent() int {
	return divideAndRoundToNearestInteger(territory.TopPX, MAPHEIGHTPX)
}

func Render(strategicMap Map, ownerships Ownerships, wars []*war.War) RenderedMap {
	renderedMap := RenderedMap{}

	testlandia, _ := nationstates_api.GetNationData("testlandia")
	ownerships["A"] = *testlandia

	for _, territory := range strategicMap.Territories {

		renderedTerritory := RenderedTerritory{}
		renderedTerritory.LeftPercent = territory.LeftPercent()
		renderedTerritory.TopPercent = territory.TopPercent()

		residentNation, residentNationExists := ownerships[territory.Name]
		if residentNationExists {

			war := war.FindOngoingWarAt(wars, territory.Name)
			if war != nil {
				renderedTerritory.Text = template.HTML(fmt.Sprint(territory.Name, " ", residentNation.FlagThumbnail(), "⚔️", war.Attacker.FlagThumbnail()))
			} else {
				renderedTerritory.Text = template.HTML(territory.Name + " " + string(residentNation.FlagThumbnail()))
			}
		} else {
			renderedTerritory.Text = template.HTML(territory.Name + " ❓")
		}

		renderedMap.Territories = append(renderedMap.Territories, renderedTerritory)
	}

	return renderedMap
}

func DoesTerritoryExist(strategicMap Map, name string) bool {
	for _, territory := range strategicMap.Territories {
		if territory.Name == name {
			return true
		}
	}
	return false
}

func Colonize(residentNations *Ownerships, strategicMap Map, colonizer nationstates_api.Nation, target string) error {

	if !DoesTerritoryExist(strategicMap, target) {
		return fmt.Errorf("No territory exists at %s", target)
	}

	_, exists := (*residentNations)[target]
	if exists {
		return fmt.Errorf("A nation is already resident at %s", target)
	}

	(*residentNations)[target] = colonizer

	return nil
}
