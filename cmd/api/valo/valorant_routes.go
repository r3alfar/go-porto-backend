package valo

import (
	apiHelpers "backend/internal/helpers"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
)

type Result struct {
	endpoint string
	data     map[string]interface{}
	err      error
}

type FuncRes struct {
	status   string
	scsCount int
	errCount int
}

// matches
func FetchMatches() Result {
	//url
	endpoint := "v3/by-puuid/matches/ap/"
	puuid := os.Getenv("VALO_MAIN_ACCOUNT_PUUID")
	defaultUrl := os.Getenv("VALO_DEFAULT_ENDPOINT")
	fullUrl := defaultUrl + endpoint + puuid

	// Create a URL object
	u, err := url.Parse(fullUrl)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return Result{endpoint: endpoint, data: nil, err: err}
	}

	// Create a url.Values object and set query parameters
	queryParams := url.Values{}
	queryParams.Set("size", "1")
	queryParams.Set("mode", "competitive")

	// Append the encoded query parameters to the URL
	u.RawQuery = queryParams.Encode()

	data, err := makeRequest("GET", u.String())
	if err != nil {
		return Result{endpoint: endpoint, data: data, err: err}
	}

	return Result{endpoint: endpoint, data: data, err: err}
}

// mmr-history
func FetchMMRHistory() Result {
	//url
	endpoint := "v1/by-puuid/mmr-history/ap/"
	puuid := os.Getenv("VALO_MAIN_ACCOUNT_PUUID")
	url := os.Getenv("VALO_DEFAULT_ENDPOINT")

	data, err := makeRequest("GET", url+endpoint+puuid)
	if err != nil {
		return Result{endpoint: endpoint, data: nil, err: err}
	}

	return Result{endpoint: endpoint, data: data, err: err}
}

// mmr
func FetchMMRInfo() Result {
	//url
	endpoint := "v2/by-puuid/mmr/ap/"
	puuid := os.Getenv("VALO_MAIN_ACCOUNT_PUUID")
	url := os.Getenv("VALO_DEFAULT_ENDPOINT")

	data, err := makeRequest("GET", url+endpoint+puuid)
	if err != nil {
		return Result{endpoint: endpoint, data: data, err: err}
	}

	return Result{endpoint: endpoint, data: data, err: err}
}

func dailyGrab() (map[string]interface{}, error) {

	return nil, nil
}

func InitialGrab() FuncRes {

	endpointsFunc := []func() Result{
		// FetchMatches,
		FetchMMRHistory,
		FetchMMRInfo,
	}

	var wg sync.WaitGroup
	ch := make(chan Result, len(endpointsFunc))

	// launch go routine for each function
	for _, fn := range endpointsFunc {
		wg.Add(1)
		go func(f func() Result) {
			defer wg.Done()
			ch <- f()
		}(fn)
		// fn is passed to go func as f
	}

	// close the channel after all goroutines are done
	go func() {
		wg.Wait()
		close(ch)
	}()

	// collect the channel after all go routines are done
	var errCounter, scsCounter int
	var r FuncRes
	for result := range ch {
		if result.err != nil {
			errCounter++
			fmt.Printf("?????Error caling %s: %v\n", result.endpoint, result.err)
		} else {
			scsCounter++
		}
		fmt.Printf("-----Data from %s: %v\n", result.endpoint, result.data)
	}

	// collect results from the channel
	fmt.Println("All goroutines finished.")

	r.status = "Finished"
	r.scsCount = scsCounter
	r.errCount = errCounter

	return r
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

	fmt.Println("---------type of data: %T", body)
	data := response["data"].(map[string]interface{})

	return data, nil
}

// account
// func FetchAccDetail() (map[string]interface{}, error) {
// 	//url
// 	endpoint := "v1/by-puuid/account/"
// 	puuid := os.Getenv("VALO_MAIN_ACCOUNT_PUUID")
// 	url := os.Getenv("VALO_DEFAULT_ENDPOINT")

// 	//create request
// 	// data, err := makeRequest("GET", url+endpoint+puuid)
// 	res, err := apiHelpers.MakeRequest("GET", url+endpoint+puuid)
// 	if err != nil {
// 		return nil, err
// 	}

// 	parsed, err := apiHelpers.ParseJSON(res)
// 	if err != nil {
// 		return nil, err
// 	}
// 	//since it a single object
// 	data := parsed.(map[string]interface{})

// 	return data, nil

// }

func FetchAccDetail() ([]byte, error) {
	//url
	endpoint := "v1/by-puuid/account/"
	puuid := os.Getenv("VALO_MAIN_ACCOUNT_PUUID")
	url := os.Getenv("VALO_DEFAULT_ENDPOINT")

	//create request
	// data, err := makeRequest("GET", url+endpoint+puuid)
	res, err := apiHelpers.MakeRequest("GET", url+endpoint+puuid)
	if err != nil {
		return nil, err
	}

	// parsed, err := apiHelpers.ParseJSON(res)
	// if err != nil {
	// 	return nil, err
	// }
	// //since it a single object
	// data := parsed.(map[string]interface{})

	return res, nil

}
