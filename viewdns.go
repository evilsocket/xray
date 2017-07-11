package xray

import (
	"github.com/PuerkitoBio/goquery"
)

type HistoryEntry struct {
	Address  string
	Location string
	ISP      string
	Updated  string
}

// TODO: Replace this with proper ViewDNS API
func ViewDNS_GetHstory(string domain) map[string]HistoryEntry {
	url := "http://viewdns.info/iphistory/?domain=" + domain
	history := make(map[string]HistoryEntry)

	return history
}
