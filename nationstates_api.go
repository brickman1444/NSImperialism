package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Nation struct {
	Name string
	Flag string
}

func GetStandardData(nationName string) (*Nation, error) {

	url := fmt.Sprintf("https://www.nationstates.net/cgi-bin/api.cgi?nation=%s", url.QueryEscape(nationName))
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

	fmt.Println(string(body))

	return nil, nil
}
