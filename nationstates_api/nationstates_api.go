package nationstates_api

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type CensusScale struct {
	Id             int `xml:"id,attr"`
	PercentageRank int `xml:"PRANK"`
}

type Nation struct {
	Id           string        `xml:"id,attr"`
	Name         string        `xml:"NAME"`
	FlagURL      string        `xml:"FLAG"`
	CensusScales []CensusScale `xml:"CENSUS>SCALE"`
}

func (nation *Nation) GetDefenseForces() int {
	for _, censusScale := range nation.CensusScales {
		if censusScale.Id == 46 {
			return censusScale.PercentageRank
		}
	}
	return 0
}

func (nation *Nation) GetURL() string {
	return fmt.Sprintf("https://www.nationstates.net/nation=%s", nation.Id)
}

func ParseNation(xmlData []byte) (*Nation, error) {
	nation := &Nation{}
	err := xml.Unmarshal(xmlData, nation)
	if err != nil {
		return nil, err
	}
	return nation, nil
}

func GetNationData(nationName string) (*Nation, error) {

	url := fmt.Sprintf("https://www.nationstates.net/cgi-bin/api.cgi?nation=%s;q=census+name+flag;scale=46;mode=prank", url.QueryEscape(nationName))
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", "NSImperialism")

	httpClient := &http.Client{Timeout: 10 * time.Second}
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("resp.StatusCode: " + strconv.Itoa(response.StatusCode))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return ParseNation(body)
}
