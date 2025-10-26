package mod

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheck(t *testing.T) {
	versions, err := fetchVersions()
	if err != nil {
		t.Fatal(err)
	}
	stable := versions[0]
	oldstable := versions[1]

	splitted := strings.Split(oldstable, ".")
	minor := strings.Join(splitted[:2], ".")

	tests := []struct {
		name    string
		goVer   string
		lax     bool
		wantErr bool
	}{
		{
			name:    "oldestversion",
			goVer:   "1.17",
			lax:     false,
			wantErr: true,
		},
		{
			name:    "stable",
			goVer:   stable,
			lax:     false,
			wantErr: true,
		},
		{
			name:    "oldstable",
			goVer:   oldstable,
			lax:     false,
			wantErr: false,
		},
		{
			name:    "lax with oldstable",
			goVer:   oldstable,
			lax:     true,
			wantErr: false,
		},
		{
			name:    "lax with oldstable minor version",
			goVer:   minor,
			lax:     true,
			wantErr: false,
		},
		{
			name:    "lax with oldstable minor.0 version",
			goVer:   minor + ".0",
			lax:     true,
			wantErr: false,
		},
		{
			name:    "not lax with oldstable minor version",
			goVer:   minor,
			lax:     false,
			wantErr: true,
		},
		{
			name:    "lax with older version",
			goVer:   "1.17",
			lax:     true,
			wantErr: false,
		},
		{
			name:    "lax with stable version",
			goVer:   stable,
			lax:     true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modpath := filepath.Join(t.TempDir(), "go.mod")
			if err := os.WriteFile(modpath, []byte(fmt.Sprintf("module %s\n\ngo %s\n", t.TempDir(), tt.goVer)), 0600); err != nil {
				t.Fatal(err)
			}
			if err := Check(modpath, tt.lax); (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
