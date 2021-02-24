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
	ID     string
	LeftPX int
	TopPX  int
}

type Map struct {
	Territories []Territory
}

type RenderedTerritory struct {
	Text        template.HTML
	ID          string
	LeftPercent int
	TopPercent  int
}

type RenderedMap struct {
	Territories []RenderedTerritory
	Name        string
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

func GetTerritoryDisplayName(territory databasemap.DatabaseCell) string {
	if len(territory.Name) != 0 {
		return territory.Name
	} else {
		return territory.ID
	}
}

func getTerritoryDisplayNameLink(territory databasemap.DatabaseCell, mapID string) string {
	name := GetTerritoryDisplayName(territory)

	url := "/maps/" + mapID + "/territories/" + territory.ID
	return fmt.Sprintf("<a href=\"%s\" title=\"%s\">%s</a>", url, name, name)
}

func getTextForTerritory(territoryDefinition Territory, databaseMap databasemap.DatabaseMap, nationStatesProvider nationstates_api.NationStatesProvider) (string, error) {

	territory, doesTerritoryExist := databaseMap.Cells[territoryDefinition.ID]
	if !doesTerritoryExist {
		return territoryDefinition.ID + " ❓", nil
	}

	territoryDisplayNameLink := getTerritoryDisplayNameLink(territory, databaseMap.ID)

	if territory.Resident == "" {
		return territoryDisplayNameLink + " ❓", nil
	}

	residentNation, err := nationStatesProvider.GetNationData(territory.Resident)
	if err != nil {
		return "", err
	}

	war := war.FindOngoingWarAt(databaseMap.GetWars(), territoryDefinition.ID)
	if war == nil {
		return territoryDisplayNameLink + " " + string(residentNation.FlagThumbnail()), nil
	}

	attacker, err := nationStatesProvider.GetNationData(war.Attacker)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(territoryDisplayNameLink, " ", residentNation.FlagThumbnail(), "⚔️", attacker.FlagThumbnail()), nil
}

func Render(strategicMap Map, databaseMap databasemap.DatabaseMap, nationStatesProvider nationstates_api.NationStatesProvider) (RenderedMap, error) {
	renderedMap := RenderedMap{}
	renderedMap.Name = databasemap.GetDisplayName(databaseMap)

	for _, territoryDefinition := range strategicMap.Territories {

		text, err := getTextForTerritory(territoryDefinition, databaseMap, nationStatesProvider)
		if err != nil {
			return RenderedMap{}, err
		}

		renderedTerritory := RenderedTerritory{
			LeftPercent: territoryDefinition.LeftPercent(),
			TopPercent:  territoryDefinition.TopPercent(),
			Text:        template.HTML(text),
			ID:          territoryDefinition.ID,
		}

		renderedMap.Territories = append(renderedMap.Territories, renderedTerritory)
	}

	return renderedMap, nil
}

func DoesTerritoryExist(strategicMap Map, territoryID string) bool {
	for _, territory := range strategicMap.Territories {
		if territory.ID == territoryID {
			return true
		}
	}
	return false
}
