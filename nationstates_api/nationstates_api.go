package nationstates_api

import (
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const CENSUSSCALEDEFENSEFORCES int = 46

var rateLimitDuration, _ = time.ParseDuration("30s")
var limiter = NewRateLimiter(40, rateLimitDuration) // API Docs say 50 requests in 30 seconds so I'm being a little conservative so we don't get locked out https://www.nationstates.net/pages/api.html#ratelimits
var cacheExpirationDuration, _ = time.ParseDuration("12h")
var cache = NewCache(cacheExpirationDuration)

type CensusScale struct {
	Id             int `xml:"id,attr"`
	PercentageRank int `xml:"PRANK"`
}

type Nation struct {
	Id           string        `xml:"id,attr"`
	Name         string        `xml:"FULLNAME"`
	FlagURL      string        `xml:"FLAG"`
	Demonym      string        `xml:"DEMONYM"`
	CensusScales []CensusScale `xml:"CENSUS>SCALE"`
}

func (nation *Nation) GetDefenseForces() int {
	for _, censusScale := range nation.CensusScales {
		if censusScale.Id == CENSUSSCALEDEFENSEFORCES {
			return censusScale.PercentageRank
		}
	}
	return 0
}

func (nation *Nation) SetDefenseForces(percentageRank int) {
	for censusIndex, censusScale := range nation.CensusScales {
		if censusScale.Id == CENSUSSCALEDEFENSEFORCES {
			nation.CensusScales[censusIndex].PercentageRank = percentageRank
			return
		}
	}
	nation.CensusScales = append(nation.CensusScales, CensusScale{CENSUSSCALEDEFENSEFORCES, percentageRank})
}

func (nation *Nation) GetURL() string {
	return fmt.Sprintf("https://www.nationstates.net/nation=%s", nation.Id)
}

func (nation Nation) FlagThumbnailURL() string {
	return strings.ReplaceAll(nation.FlagURL, ".png", "t2.png")
}

func (nation *Nation) FlagAndName() template.HTML {
	return template.HTML(fmt.Sprintf("<a href=\"%s\" title=\"%s\"><img src=\"%s\" class=\"flag-thumb\"/>%s</a>", nation.GetURL(), nation.Name, nation.FlagThumbnailURL(), nation.Name))
}

func (nation *Nation) FlagThumbnail() template.HTML {
	return template.HTML(fmt.Sprintf("<a href=\"%s\" title=\"%s\"><img src=\"%s\" class=\"flag-thumb\"/></a>", nation.GetURL(), nation.Name, nation.FlagThumbnailURL()))
}

func ParseNation(xmlData []byte) (*Nation, error) {
	nation := &Nation{}
	err := xml.Unmarshal(xmlData, nation)
	if err != nil {
		return nil, err
	}
	return nation, nil
}

func GetCanonicalName(inNationName string) string {
	return strings.ReplaceAll(strings.ToLower(inNationName), " ", "_")
}

func GetNationData(nationName string) (*Nation, error) {

	if nationName == "" {
		return nil, errors.New("Empty nation name")
	}

	nationName = GetCanonicalName(nationName)

	cachedNation := cache.GetNation(nationName, time.Now())
	if cachedNation != nil {
		return cachedNation, nil
	}

	if limiter.IsAtRateLimit(time.Now()) {
		return nil, errors.New("Hit internal nationstates API rate limit")
	}

	url := fmt.Sprintf("https://www.nationstates.net/cgi-bin/api.cgi?nation=%s;q=census+fullname+flag+demonym;scale=46;mode=prank", url.QueryEscape(nationName))
	log.Println("Pulling down nation data for", nationName)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", "NSImperialism")

	httpClient := &http.Client{Timeout: 10 * time.Second}

	limiter.AddRequestTime(time.Now())

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusTooManyRequests {
		log.Println("Too many requests to NationStates api. Wait", response.Header.Get("X-Retry-After"), "seconds.")
		return nil, errors.New("Too many requests to NationStates api")
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("NationStates API Response Error. StatusCode: " + strconv.Itoa(response.StatusCode))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	parsedNation, err := ParseNation(body)
	if err != nil {
		return nil, err
	}

	cache.AddNation(parsedNation.Id, *parsedNation, time.Now())

	return parsedNation, nil
}

func IsCorrectVerificationCode(nationName string, verificationCode string) (bool, error) {

	if limiter.IsAtRateLimit(time.Now()) {
		return false, errors.New("Hit internal nationstates API rate limit")
	}

	url := fmt.Sprintf("https://www.nationstates.net/cgi-bin/api.cgi?a=verify&nation=%s&checksum=%s", url.QueryEscape(nationName), url.QueryEscape(verificationCode))
	log.Println("Verifying nation", nationName)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	request.Header.Set("User-Agent", "NSImperialism")

	httpClient := &http.Client{Timeout: 10 * time.Second}

	limiter.AddRequestTime(time.Now())

	response, err := httpClient.Do(request)
	if err != nil {
		return false, err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusTooManyRequests {
		log.Println("Too many requests to NationStates api. Wait", response.Header.Get("X-Retry-After"), "seconds.")
		return false, errors.New("Too many requests to NationStates api")
	}

	if response.StatusCode != http.StatusOK {
		return false, errors.New("NationStates API Response Error. StatusCode: " + strconv.Itoa(response.StatusCode))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	bodyString := string(body)

	return strings.HasPrefix(bodyString, "1"), nil
}
