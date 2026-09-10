package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestGetFromLocalMissingFile(t *testing.T) {
	localConfigPath = filepath.Join(t.TempDir(), "missing.yaml")
	_, err := getFromLocal()
	if err == nil {
		t.Fatal("expected error when local config file is missing")
	}
}

func TestGetFromLocalReadsYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mall_config.yaml")
	content := []byte("server:\n  http_port: 9090\n  log_level: info\n")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}
	localConfigPath = path
	conf, err := getFromLocal()
	if err != nil {
		t.Fatal(err)
	}
	if conf.Server.HttpPort != 9090 {
		t.Fatalf("HttpPort=%d want 9090", conf.Server.HttpPort)
	}
}

func TestPrepareRemoteConfigRejectsYAMLPath(t *testing.T) {
	err := prepareRemoteConfig(viper.New(), "config.yaml")
	if err == nil {
		t.Fatal("expected error when -r is given a yaml file path")
	}
	if !strings.Contains(err.Error(), "-c") {
		t.Fatalf("error should tell user to use -c: %v", err)
	}
}

func TestPrepareRemoteConfigAcceptsEtcdEndpoint(t *testing.T) {
	if err := prepareRemoteConfig(viper.New(), "127.0.0.1:2379"); err != nil {
		t.Fatal(err)
	}
}

func TestLooksLikeConfigFile(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"config.yaml", true},
		{"mall_config.yaml", true},
		{"http://127.0.0.1:2379", false},
		{"127.0.0.1:2379", false},
	}
	for _, tt := range cases {
		if got := looksLikeConfigFile(tt.in); got != tt.want {
			t.Fatalf("looksLikeConfigFile(%q)=%v want %v", tt.in, got, tt.want)
		}
	}
}

func TestNormalizeEtcdEndpoint(t *testing.T) {
	if got := normalizeEtcdEndpoint("127.0.0.1:2379"); got != "http://127.0.0.1:2379" {
		t.Fatalf("got %s", got)
	}
	if got := normalizeEtcdEndpoint("http://127.0.0.1:2379"); got != "http://127.0.0.1:2379" {
		t.Fatalf("got %s", got)
	}
}
