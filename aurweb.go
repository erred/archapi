package archapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	AurWebURL = "https://aur.archlinux.org/rpc/"

	// Search By
	ByName         = "name"
	ByNameDesc     = "name-desc"
	ByMaintainer   = "maintainer"
	ByDepends      = "depends"
	ByMakeDepends  = "makedepends"
	ByOptDepends   = "optdepends"
	ByCheckDepends = "checkdepends"

	aurApiVersion = "5"
	typeSearch    = "search"
	typeInfo      = "info"
	typeMultiInfo = "multiinfo"
	typeError     = "string"
)

var (
	defaultClient = &Client{
		c: http.DefaultClient,
	}
)

func Info(pkgs ...string) ([]PackageInfo, error) {
	return defaultClient.Info(pkgs...)
}

func Search(by, keywords string) ([]Package, error) {
	return defaultClient.Search(by, keywords)
}

type Client struct {
	c *http.Client
}

func NewClient(c *http.Client) *Client {
	return &Client{
		c: c,
	}
}

func (c Client) Info(pkgs ...string) ([]PackageInfo, error) {
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("archapi: no packages passed")
	}
	v := url.Values{
		"v":    []string{aurApiVersion},
		"type": []string{typeInfo},
	}
	for _, p := range pkgs {
		v.Add("arg[]", p)
	}

	var ir []PackageInfo
	err := c.aurGet(v, &ir)
	if err != nil {
		return nil, fmt.Errorf("archapi: %w", err)
	}
	return ir, nil
}

func (c Client) Search(by, keywords string) ([]Package, error) {
	v := url.Values{
		"v":    []string{aurApiVersion},
		"type": []string{typeSearch},
	}
	if by != "" {
		v.Set("by", by)
	}
	v.Set("arg", keywords)

	var sr []Package
	err := c.aurGet(v, &sr)
	if err != nil {
		return nil, fmt.Errorf("archapi: %w", err)
	}
	return sr, nil
}

func (c Client) aurGet(values url.Values, data interface{}) error {
	res, err := c.c.Get(AurWebURL + "?" + values.Encode())
	if err != nil {
		return err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var ar aurResult
	err = json.Unmarshal(b, &ar)
	if err != nil {
		return err
	}
	if ar.Type == typeError {
		return ErrAur{ar.Error}
	}
	return json.Unmarshal(ar.Results, data)
}

type aurResult struct {
	Version     int             `json:"version"`
	Type        string          `json:"type"`
	ResultCount int             `json:"resultcount"`
	Results     json.RawMessage `json:"results"`
	Error       string          `json:"error"`
}

type Package struct {
	ID             int     `json:"ID"`
	Name           string  `json:"Name"`
	PackageBaseID  int     `json:"PackageBaseID"`
	PackageBase    string  `json:"PackageBase"`
	Version        string  `json:"Version"`
	Description    string  `json:"Description"`
	URL            string  `json:"URL"`
	NumVotes       int     `json:"NumVotes"`
	Popularity     float64 `json:"Popularity"`
	OutOfDate      *string `json:"OutOfDate"`
	Maintainer     string  `json:"Maintainer"`
	FirstSubmitted int64   `json:"FirstSubmitted"`
	LastModified   int64   `json:"LastModified"`
	URLPath        string  `json:"URLPath"`
}

type PackageInfo struct {
	Package

	Depends      []string `json:"Depends"`
	MakeDepends  []string `json:"MakeDepends"`
	OptDepends   []string `json:"OptDepends"`
	CheckDepends []string `json:"CheckDepends"`
	Conflicts    []string `json:"Conflicts"`
	Provides     []string `json:"Provides"`
	Replaces     []string `json:"Replaces"`
	Groups       []string `json:"Groups"`
	License      []string `json:"License"`
	Keywords     []string `json:"Keywords"`
}

type ErrAur struct {
	Err string
}

func (e ErrAur) Error() string {
	return e.Err
}
