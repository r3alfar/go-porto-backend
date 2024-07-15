package apiHelpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func ParseJSON(response []byte) (interface{}, error) {
	// parse single json object
	var singleObject map[string]interface{}
	err := json.Unmarshal(response, &singleObject)
	if err == nil {
		return singleObject, nil
	}

	// parse arrayof obejcts
	var arrayOfObjects []map[string]interface{}
	err = json.Unmarshal(response, &arrayOfObjects)
	if err == nil {
		return arrayOfObjects, nil
	}

	// returning error
	return nil, fmt.Errorf("unable to parse JSON data")
}

func MakeRequest(reqMethod string, url string) ([]byte, error) {
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

	return body, nil

	// jsonBody := string(body)
	// // convert
	// var response map[string]interface{}
	// err = json.Unmarshal([]byte(jsonBody), &response)
	// if err != nil {

	// 	return nil, fmt.Errorf("jSON Parse Error: %v", err)
	// }

	// fmt.Println("---------type of data: %T", body)
	// data := response["data"].(map[string]interface{})

	// return data, nil
}
