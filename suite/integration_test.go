package main

import (
	"os"
	"path/filepath"
	"suite/suite/config"
	"testing"
)

func TestReadConfigIntegration(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
		wantErr    bool
		check      func(*testing.T, *config.ServerConfig)
	}{
		{
			name:       "read test web config",
			configFile: "../testdata/configs/test_web.sscfg",
			wantErr:    false,
			check: func(t *testing.T, cfg *config.ServerConfig) {
				if cfg.SetupSecure == nil {
					t.Fatal("SetupSecure is nil")
				}
				if cfg.SetupSecure.SSHUser != "testuser" {
					t.Errorf("SSHUser = %s, want testuser", cfg.SetupSecure.SSHUser)
				}
				if cfg.SetupSecure.SSHPort != 2222 {
					t.Errorf("SSHPort = %d, want 2222", cfg.SetupSecure.SSHPort)
				}
				if cfg.SetupSecure.Config == nil {
					t.Fatal("Config is nil")
				}
				if cfg.SetupSecure.Config.Type != "web" {
					t.Errorf("Type = %s, want web", cfg.SetupSecure.Config.Type)
				}
				if cfg.SetupSecure.Config.Domain != "test.example.com" {
					t.Errorf("Domain = %s, want test.example.com", cfg.SetupSecure.Config.Domain)
				}
			},
		},
		{
			name:       "read test database config",
			configFile: "../testdata/configs/test_database.sscfg",
			wantErr:    false,
			check: func(t *testing.T, cfg *config.ServerConfig) {
				if cfg.SetupSecure == nil {
					t.Fatal("SetupSecure is nil")
				}
				if cfg.SetupSecure.SSHUser != "dbadmin" {
					t.Errorf("SSHUser = %s, want dbadmin", cfg.SetupSecure.SSHUser)
				}
				if cfg.SetupSecure.Config.Type != "database" {
					t.Errorf("Type = %s, want database", cfg.SetupSecure.Config.Type)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			absPath, err := filepath.Abs(tt.configFile)
			if err != nil {
				t.Fatalf("Failed to get absolute path: %v", err)
			}

			content, err := os.ReadFile(absPath)
			if err != nil {
				t.Fatalf("Failed to read config file: %v", err)
			}

			cfg, err := config.ParseConfig(string(content))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.check != nil {
				tt.check(t, cfg)
			}
		})
	}
}

func TestConfigGenerationIntegration(t *testing.T) {
	serverTypes := []string{
		config.ServerTypeWeb,
		config.ServerTypeDatabase,
		config.ServerTypeDocker,
		config.ServerTypeProxy,
		config.ServerTypeBuild,
	}

	for _, serverType := range serverTypes {
		t.Run("generate_"+serverType, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, serverType+".sscfg")

			err := config.CreateDefaultConfig(serverType, configPath)
			if err != nil {
				t.Errorf("CreateDefaultConfig() error = %v", err)
				return
			}

			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				t.Errorf("Config file was not created at %s", configPath)
				return
			}

			content, err := os.ReadFile(configPath)
			if err != nil {
				t.Errorf("Failed to read generated config: %v", err)
				return
			}

			cfg, err := config.ParseConfig(string(content))
			if err != nil {
				t.Errorf("Failed to parse generated config: %v", err)
				return
			}

			if cfg.SetupSecure == nil {
				t.Error("Generated config has nil SetupSecure")
			}

			if cfg.SetupSecure != nil && cfg.SetupSecure.Config == nil {
				t.Error("Generated config has nil Config")
			}
		})
	}
}

func TestFullConfigLifecycle(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.sscfg")

	err := config.CreateDefaultConfig(config.ServerTypeWeb, configPath)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	cfg, err := config.ParseConfig(string(content))
	if err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	if cfg == nil {
		t.Fatal("Config is nil")
	}

	if cfg.SetupSecure == nil {
		t.Fatal("SetupSecure is nil")
	}

	if cfg.SetupSecure.Config == nil {
		t.Fatal("Config.SetupSecure.Config is nil")
	}

	if cfg.SetupSecure.Config.Type != config.ServerTypeWeb {
		t.Errorf("Config type = %s, want %s", cfg.SetupSecure.Config.Type, config.ServerTypeWeb)
	}
}

func TestInvalidConfigHandling(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "malformed config",
			content: "this is not valid config",
			wantErr: false,
		},
		{
			name:    "empty config",
			content: "",
			wantErr: false,
		},
		{
			name: "partial config",
			content: `.setup_secure{
ssh_user: "test"
}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.ParseConfig(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if cfg == nil && !tt.wantErr {
				t.Error("Expected non-nil config")
			}
		})
	}
}
