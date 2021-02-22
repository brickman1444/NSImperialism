package strategicmap

import (
	"fmt"
	"html/template"
	"math"

	"github.com/brickman1444/NSImperialism/databasemap"
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

func getTextForTerritory(territoryName string, residents databasemap.DatabaseMap, wars []databasemap.DatabaseWar, nationStatesProvider nationstates_api.NationStatesProvider) (string, error) {
	residentNationID, err := residents.GetResident(territoryName)
	if err != nil {
		return "", err
	}

	if residentNationID == "" {
		return territoryName + " ❓", nil
	}

	residentNation, err := nationStatesProvider.GetNationData(residentNationID)
	if err != nil {
		return "", err
	}

	war := war.FindOngoingWarAt(wars, territoryName)
	if war == nil {
		return territoryName + " " + string(residentNation.FlagThumbnail()), nil
	}

	attacker, err := nationStatesProvider.GetNationData(war.Attacker)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(territoryName, " ", residentNation.FlagThumbnail(), "⚔️", attacker.FlagThumbnail()), nil
}

func Render(strategicMap Map, residents databasemap.DatabaseMap, wars []databasemap.DatabaseWar, nationStatesProvider nationstates_api.NationStatesProvider) (RenderedMap, error) {
	renderedMap := RenderedMap{}

	for _, territory := range strategicMap.Territories {

		text, err := getTextForTerritory(territory.Name, residents, wars, nationStatesProvider)
		if err != nil {
			return RenderedMap{}, err
		}

		renderedTerritory := RenderedTerritory{
			LeftPercent: territory.LeftPercent(),
			TopPercent:  territory.TopPercent(),
			Text:        template.HTML(text),
		}

		renderedMap.Territories = append(renderedMap.Territories, renderedTerritory)
	}

	return renderedMap, nil
}

func DoesTerritoryExist(strategicMap Map, name string) bool {
	for _, territory := range strategicMap.Territories {
		if territory.Name == name {
			return true
		}
	}
	return false
}
