package mod

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestCheck(t *testing.T) {
	versions, err := fetchVersions()
	if err != nil {
		t.Fatal(err)
	}
	stable := versions[0]
	oldstable := versions[1]

	tests := []struct {
		name    string
		goVer   string
		wantErr bool
	}{
		{
			name:    "oldestversion",
			goVer:   "1.17",
			wantErr: true,
		},
		{
			name:    "stable",
			goVer:   stable,
			wantErr: true,
		},
		{
			name:    "oldstable",
			goVer:   oldstable,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modpath := filepath.Join(t.TempDir(), "go.mod")
			if err := os.WriteFile(modpath, []byte(fmt.Sprintf("module %s\n\ngo %s\n", t.TempDir(), tt.goVer)), 0600); err != nil {
				t.Fatal(err)
			}
			if err := Check(modpath); (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
