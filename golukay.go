// Package golukay fetches the list of UK Bank Holidays
package golukay

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// Date When the Bank Holiday occurs
type Date struct {
	time.Time
}

// UnmarshalJSON Converts raw date string from Gov API to time.Time
func (ct *Date) UnmarshalJSON(b []byte) (err error) {
	rawDate := strings.Trim(string(b), "\"")
	layout := "2006-01-02"
	ct.Time, err = time.Parse(layout, rawDate)
	return
}

// BankHoliday - A day off work in some part of the UK, e.g. Christmas
type BankHoliday struct {
	Title string `json:"title"`
	Date  Date   `json:"date"`
}

// Division A part of the UK e.g. Scotland
type Division struct {
	Name   string        `json:"division"`
	Events []BankHoliday `json:"events"`
}

// GovResponse The payload returned from gov.uk
type GovResponse struct {
	EnglandAndWales Division `json:"england-and-wales"`
	Scotland        Division `json:"scotland"`
	NorthernIreland Division `json:"northern-ireland"`
}

// GetHolidays Fetches UK Bank Holidays from the gov.uk website
func GetHolidays() (GovResponse, error) {
	var govResponse GovResponse

	url := "https://www.gov.uk/bank-holidays.json"

	resp, err := http.Get(url)
	if err != nil {
		return govResponse, err
	}
	if resp.StatusCode != 200 {
		fmt.Printf("API request failed with code %d", resp.StatusCode)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return govResponse, err
	}

	err = json.Unmarshal(body, &govResponse)
	if err != nil {
		fmt.Printf("JSON Error %s\n", err)
		return govResponse, err
	}

	return govResponse, nil
}
