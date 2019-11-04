package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)
var bHitRateLimit bool

// CallSingle is a single URL call
// it returns true, if the API Rate Limit was hit
// the error & result object itself are modified in the Request struct
func CallSingle(token string, request *ThousandeyesRequest) bool {

	client := &http.Client{}
	bHitAPILimit := false

	req, err := http.NewRequest("GET", request.URL, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	log.Println(fmt.Sprintf("CALL >>> Url: %s", request.URL))
	resp, err := client.Do(req)
	if resp.StatusCode == 429 {
		bHitAPILimit = true
		request.Error = fmt.Errorf("ThousandEyes API Rate Limit hit (\"Too many requests\") - http code: %d", resp.StatusCode)
		return bHitAPILimit
	} else if err != nil || resp.StatusCode != 200 {
		request.Error = fmt.Errorf("ThousandEyes API Request failed: %s / http code: %d", err, resp.StatusCode)
		return bHitAPILimit
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		request.Error = err
		return bHitAPILimit
	}
	//log.Println(fmt.Sprintf("\nCALL <<< Url: %s |\n%s", request.URL, string(responseData)))
	err = json.Unmarshal(responseData, request.ResponseObject)
	if err != nil {
		log.Println(err.Error())
		request.Error = fmt.Errorf("parse configuration: %s", err.Error())
	}

	return bHitAPILimit
}

// CallSequence does CallSingle calls one after the other
// it returns true, if the API Rate Limit was hit
// the error & result object itself are modified in the Request struct
func CallSequence(token string, requests []ThousandeyesRequest) bool {

	for c := range requests {
		if CallSingle(token, &requests[c]) {
			return true
		}
	}
	return false
}

// CallParallel does CallSingle calls in parallel - can hit thousandeyes api restrictions easily
// it returns true, if the API Rate Limit was hit
// the error & result object itself are modified in the Request struct
func CallParallel(token string, requests []ThousandeyesRequest) error {

	var waitGroup sync.WaitGroup
	var m sync.Mutex

	httpChan := make(chan ThousandeyesRequest, len(requests))
	bHitRateLimit = false;

	for c, request := range requests {
		log.Println(fmt.Sprintf("Count [%d] - URL: %s", c, request.URL))

		waitGroup.Add(1)

		go func(token string, request ThousandeyesRequest, httpChan chan (ThousandeyesRequest), m *sync.Mutex) {
			defer waitGroup.Done()

			if (bHitRateLimit){
				log.Fatal(fmt.Sprintf("Skip Detail request (%s), bcz we hit the API Request Limit.", request.URL))
				return
			}
			//log.Println(fmt.Sprintf("URL: %s | API-Request-Limit-Hit ?: %t", request.URL, b))

			b := CallSingle(token, &request)
			m.Lock()
			bHitRateLimit = b
			m.Unlock()
			httpChan <- request

		}(token, request, httpChan, &m)
	}

	waitGroup.Wait()
	close(httpChan)

	return nil
}
