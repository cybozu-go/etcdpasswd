package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		data     string
		wantKeep bool
	}{
		{
			name: "default",
			data: "endpoints:\n  - http://127.0.0.1:2379\n",
		},
		{
			name:     "keep supplementary groups",
			data:     "endpoints:\n  - http://127.0.0.1:2379\nkeep-supplementary-groups: true\n",
			wantKeep: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p := filepath.Join(t.TempDir(), "config.yml")
			if err := os.WriteFile(p, []byte(tc.data), 0o644); err != nil {
				t.Fatal(err)
			}

			cfg, err := loadConfig(p)
			if err != nil {
				t.Fatal(err)
			}
			if cfg.KeepSupplementaryGroups != tc.wantKeep {
				t.Errorf("KeepSupplementaryGroups = %v, want %v", cfg.KeepSupplementaryGroups, tc.wantKeep)
			}
			if len(cfg.Endpoints) != 1 || cfg.Endpoints[0] != "http://127.0.0.1:2379" {
				t.Errorf("Endpoints = %v", cfg.Endpoints)
			}
			if cfg.Prefix != "/passwd/" {
				t.Errorf("Prefix = %q, want the default", cfg.Prefix)
			}
		})
	}
}
