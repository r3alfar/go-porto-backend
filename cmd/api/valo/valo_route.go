package valo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// matches
func FetchMatches() (map[string]interface{}, error) {
	//url
	endpoint := "v3/matches/mmr-history/ap/"
	puuid := os.Getenv("VALO_MAIN_ACCOUNT_PUUID")
	url := os.Getenv("VALO_DEFAULT_ENDPOINT")

	data, err := makeRequest("GET", url+endpoint+puuid)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// mmr-history
func FetchMMRHistory() (map[string]interface{}, error) {
	//url
	endpoint := "v1/by-puuid/mmr-history/ap/"
	puuid := os.Getenv("VALO_MAIN_ACCOUNT_PUUID")
	url := os.Getenv("VALO_DEFAULT_ENDPOINT")

	data, err := makeRequest("GET", url+endpoint+puuid)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// mmr
func FetchMMRInfo() (map[string]interface{}, error) {
	//url
	endpoint := "v2/by-puuid/mmr/ap/"
	puuid := os.Getenv("VALO_MAIN_ACCOUNT_PUUID")
	url := os.Getenv("VALO_DEFAULT_ENDPOINT")

	data, err := makeRequest("GET", url+endpoint+puuid)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// account
func FetchAccDetail() (map[string]interface{}, error) {
	//url
	endpoint := "v1/by-puuid/account/"
	puuid := os.Getenv("VALO_MAIN_ACCOUNT_PUUID")
	url := os.Getenv("VALO_DEFAULT_ENDPOINT")

	//create request
	data, err := makeRequest("GET", url+endpoint+puuid)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func makeRequest(reqMethod string, url string) (map[string]interface{}, error) {
	client := &http.Client{}
	//create request
	req, err := http.NewRequest(reqMethod, url, nil)
	if err != nil {

		return nil, fmt.Errorf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("VALO_API_KEY"))

	//make request
	resp, err := client.Do(req)
	if err != nil {

		return nil, fmt.Errorf("could not make request: %v", err)
	}

	// check resp satus code
	if resp.StatusCode != http.StatusOK {

		return nil, fmt.Errorf("received Non 200 response: %d", resp.StatusCode)
	}

	//read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error Read response body: %v", err)
	}
	jsonBody := string(body)
	// convert
	var response map[string]interface{}
	err = json.Unmarshal([]byte(jsonBody), &response)
	if err != nil {

		return nil, fmt.Errorf("jSON Parse Error: %v", err)
	}

	data := response["data"].(map[string]interface{})

	return data, nil
}
