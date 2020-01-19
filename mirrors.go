package archapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

const (
	MirrorListURL = "https://www.archlinux.org/mirrors/status/json/"
)

func ListMirrors() (*MirrorList, error) {
	return defaultClient.ListMirrors()
}

func (c Client) ListMirrors() (*MirrorList, error) {
	res, err := c.c.Get(MirrorListURL)
	if err != nil {
		return nil, fmt.Errorf("archapi: %w", err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("archapi: %w", err)
	}
	var ml MirrorList
	err = json.Unmarshal(b, &ml)
	if err != nil {
		return nil, fmt.Errorf("archapi: %w", err)
	}
	return &ml, nil
}

type MirrorList struct {
	Cutoff         int       `json:"cutoff"`
	LastCheck      time.Time `json:"last_check"`
	NumChecks      int       `json:"num_checks"`
	CheckFrequency int       `json:"check_frequency"`
	Urls           []Mirror  `json:"urls"`
	Version        int       `json:"version"`
}

type Mirror struct {
	URL            string    `json:"url"`
	Protocol       string    `json:"protocol"`
	LastSync       time.Time `json:"last_sync"`
	CompletionPct  float64   `json:"completion_pct"`
	Delay          int       `json:"delay"`
	DurationAvg    float64   `json:"duration_avg"`
	DurationStddev float64   `json:"duration_stddev"`
	Score          float64   `json:"score"`
	Active         bool      `json:"active"`
	Country        string    `json:"country"`
	CountryCode    string    `json:"country_code"`
	Isos           bool      `json:"isos"`
	Ipv4           bool      `json:"ipv4"`
	Ipv6           bool      `json:"ipv6"`
	Details        string    `json:"details"`
}
