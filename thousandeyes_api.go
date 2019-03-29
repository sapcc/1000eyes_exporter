package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// ThousandAlerts describes the JSON returned by a request active alerts to ThousandEyes
type ThousandAlerts struct {
	From  string `json:"from"`
	Alert []struct {
		Active    int    `json:"active"`
		AlertID   int    `json:"alertId"`
		DateEnd   string `json:"dateEnd,omitempty"`
		DateStart string `json:"dateStart"`
		Monitors  []struct {
			Active         int    `json:"active"`
			MetricsAtStart string `json:"metricsAtStart"`
			MetricsAtEnd   string `json:"metricsAtEnd"`
			MonitorID      int    `json:"monitorId"`
			MonitorName    string `json:"monitorName"`
			PrefixID       int    `json:"prefixId"`
			Prefix         string `json:"prefix"`
			DateStart      string `json:"dateStart"`
			DateEnd        string `json:"dateEnd"`
			Permalink      string `json:"permalink"`
			Network        string `json:"network"`
		} `json:"monitors,omitempty"`
		Permalink      string `json:"permalink"`
		RuleExpression string `json:"ruleExpression"`
		RuleID         int    `json:"ruleId"`
		RuleName       string `json:"ruleName"`
		TestID         int    `json:"testId"`
		TestName       string `json:"testName"`
		ViolationCount int    `json:"violationCount"`
		Type           string `json:"type"`
		APILinks       []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"apiLinks,omitempty"`
		Agents []struct {
			Active         int    `json:"active"`
			MetricsAtStart string `json:"metricsAtStart"`
			MetricsAtEnd   string `json:"metricsAtEnd"`
			AgentID        int    `json:"agentId"`
			AgentName      string `json:"agentName"`
			DateStart      string `json:"dateStart"`
			DateEnd        string `json:"dateEnd"`
			Permalink      string `json:"permalink"`
		} `json:"agents,omitempty"`
	} `json:"alert"`
	Pages struct {
		Current int `json:"current"`
	} `json:"pages"`
}

func thousandEyesDateTime() string {
	// Go back a bit to have some alerts to parse
	t := time.Now().UTC().Add(-*retrospectionPeriod)
	// 2006-01-02T15:04:05 is a magic date to format dates using example based layouts
	f := t.Format("2006-01-02T15:04:05")
	return string(f)
}

func (t *thousandEyes) getAlerts() (ThousandAlerts, error) {
	client := &http.Client{}
	url := string("https://api.thousandeyes.com/v6/alerts?format=json&from=" + thousandEyesDateTime())
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+t.token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	// TODO return error if http errors occur
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err.Error())
	}

	var a ThousandAlerts
	json.Unmarshal(responseData, &a)
	return a, nil

}
