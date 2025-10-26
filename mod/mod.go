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

func Check(modpath string, lax bool) error {
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
	stable := versions[0]
	oldstable := versions[1]
	if lax {
		// In lax mode, error only if the version is stable (latest version)
		splittedStable := strings.Split(stable, ".")
		stableMinor := strings.Join(splittedStable[:2], ".")
		if strings.HasPrefix(goDirectiveVersion, stableMinor) {
			return fmt.Errorf("version of go directive in go.mod is stable version (stable: %s, current: %s)", stable, goDirectiveVersion)
		}
		// Allow oldstable version or older
		return nil
	}

	// Check if go directive has patch version
	goDirectiveParts := strings.Split(goDirectiveVersion, ".")
	oldstableParts := strings.Split(oldstable, ".")

	// If go directive does not have patch version (e.g., "1.21"), compare only major.minor
	if len(goDirectiveParts) == 2 {
		oldstableMinor := strings.Join(oldstableParts[:2], ".")
		if goDirectiveVersion != oldstableMinor {
			return fmt.Errorf("version of go directive in go.mod is not oldstable (oldstable: %s, current: %s)", oldstable, goDirectiveVersion)
		}
		return nil
	}

	// If go directive has patch version, compare full version
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
