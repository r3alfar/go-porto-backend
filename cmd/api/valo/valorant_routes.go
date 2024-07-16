package valo

import (
	apiHelpers "backend/internal/helpers"
	"backend/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Result1 struct {
	endpoint string
	data     map[string]interface{}
	err      error
}

type Result struct {
	endpoint string
	data     []byte
	err      error
}

type FuncRes struct {
	Status   string `json:"status"`
	ScsCount int    `json:"success_count"`
	ErrCount int    `json:"error_count"`
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
	queryParams.Set("size", "3")
	queryParams.Set("mode", "competitive")

	// Append the encoded query parameters to the URL
	u.RawQuery = queryParams.Encode()

	data, err := apiHelpers.MakeRequest("GET", u.String())
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

	data, err := apiHelpers.MakeRequest("GET", url+endpoint+puuid)
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

	data, err := apiHelpers.MakeRequest("GET", url+endpoint+puuid)
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
		FetchMatches,
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
	var mmr models.MMR
	var mmrHistory models.MMRHistory
	var matchesData models.Matchlist
	var r FuncRes
	// collect results from the channel
	fmt.Println("All goroutines finished.")
	for result := range ch {
		if result.err != nil {
			errCounter++
			fmt.Printf("?????Error caling %s: %v\n", result.endpoint, result.err)
			continue
		} else {
			scsCounter++
		}

		// decode matchList
		// much heavier then others
		if result.endpoint == "v3/by-puuid/matches/ap/" {
			// convert []byte to io.Reader
			reader := bytes.NewReader(result.data)
			decoder := json.NewDecoder(reader)

			if err := decoder.Decode(&matchesData); err != nil {
				fmt.Printf("Failed to decode JSON: %s\n", err)
				errCounter++
				continue
			}

			fmt.Printf("-------Decode to struct Matches: %+v\n", matchesData.Data[0].Metadata.Map)

		} else if result.endpoint == "v1/by-puuid/mmr-history/ap/" {
			// convert []byte to io.Reader
			reader := bytes.NewReader(result.data)
			decoder := json.NewDecoder(reader)

			if err := decoder.Decode(&mmrHistory); err != nil {
				fmt.Printf("Failed to decode JSON: %s\n", err)
				errCounter++
				continue
			}

			fmt.Printf("-------Decode to struct MMRHISTORY FULL: %+v\n", mmrHistory.Data[0])
		} else if result.endpoint == "v2/by-puuid/mmr/ap/" {

			// convert []byte to io.Reader
			reader := bytes.NewReader(result.data)
			decoder := json.NewDecoder(reader)

			if err := decoder.Decode(&mmr); err != nil {
				fmt.Printf("Failed to decode JSON: %s\n", err)
				errCounter++
				continue
			}

			fmt.Printf("-------Decode to struct MMR: %v\n", mmr.Data.HighestRank.PatchedTier)
			fmt.Printf("-------Decode to struct MMRFULL: %+v\n", mmr.Data)
		}

		fmt.Printf("-----Data from %s: SUCCESS\n", result.endpoint)
	}

	r.Status = "Finished"
	r.ScsCount = scsCounter
	r.ErrCount = errCounter

	return r
}

func gatherAndCreate(mmr models.MMR, mmrHistory models.MMRHistory, matchesData models.Matchlist) {
	fmt.Printf("Begin create at %v \n", os.Getenv("DB_LOCALHOST"))

	// load aws config
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Unable to load sdk config: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String("http://localhost:8080")
	})
}

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

	// fmt.Println("---------type of data: %T", body)
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
