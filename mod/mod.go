package mod

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/mod/modfile"
)

const versionURL = "https://golang.org/dl/?mode=json"

func Check(modpath string) error {
	b, err := os.ReadFile(modpath)
	if err != nil {
		return err
	}
	m, err := modfile.Parse(modpath, b, nil)
	if err != nil {
		return err
	}
	goDirectiveVersion := m.Go.Version
	versions, err := fetchVersions()
	if err != nil {
		return err
	}
	oldstable := versions[1]
	if goDirectiveVersion != oldstable {
		return fmt.Errorf("version of go directive in go.mod is not latest oldstable (oldstable: %s, current: %s)", oldstable, goDirectiveVersion)
	}
	return nil
}

var cached []string

type versions []version

type version struct {
	Version string `json:"version"`
}

func fetchVersions() ([]string, error) {
	if len(cached) != 0 {
		return cached, nil
	}
	res, err := http.Get(versionURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var vs versions
	if err := json.NewDecoder(res.Body).Decode(&vs); err != nil {
		return nil, err
	}
	if len(vs) < 2 {
		cached = nil
		return nil, fmt.Errorf("failed to fetch versions")
	}
	for _, v := range vs {
		cached = append(cached, strings.TrimPrefix(v.Version, "go"))
	}
	return cached, nil
}
