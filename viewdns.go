package xray

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ViewDNS struct {
	apikey string
}

type Query struct {
	tool   string
	domain string
}

type Response struct {
	Records []HistoryEntry `json:"records"`
}

type Result struct {
	Query    Query    `json:"query"`
	Response Response `json:"response"`
}

func NewViewDNS(apikey string) *ViewDNS {
	return &ViewDNS{apikey: apikey}
}

func (d *ViewDNS) GetHistory(domain string) []HistoryEntry {
	url := fmt.Sprintf("https://api.viewdns.info/iphistory/?domain=%s&apikey=%s&output=json", domain, d.apikey)
	history := make([]HistoryEntry, 0)

	if d.apikey != "" {
		if res, err := http.Get(url); err == nil {
			defer res.Body.Close()

			decoder := json.NewDecoder(res.Body)
			r := Result{}

			if err = decoder.Decode(&r); err == nil {
				history = r.Response.Records
			}
		}
	}

	return history
}
