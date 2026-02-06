package config

import (
	"testing"
)

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
		check   func(*testing.T, *ServerConfig)
	}{
		{
			name: "basic web server config",
			content: `.setup_secure{
	ssh_user: "admin",
	user_ssh_rsa: "ssh-rsa AAAAB3...",
	ssh_port: 22022,
	.configuration{
		type: "web",
		domain: "example.com",
		email: "admin@example.com"
	},
	.firewall{
		open_ports: [
			22022,
			80,
			443
		]
	}
}

.install_tools{
	tools: [
		"nginx",
		"certbot",
		"htop"
	]
}`,
			wantErr: false,
			check: func(t *testing.T, cfg *ServerConfig) {
				if cfg.SetupSecure == nil {
					t.Fatal("SetupSecure is nil")
				}
				if cfg.SetupSecure.SSHUser != "admin" {
					t.Errorf("SSHUser = %s, want admin", cfg.SetupSecure.SSHUser)
				}
				if cfg.SetupSecure.SSHPort != 22022 {
					t.Errorf("SSHPort = %d, want 22022", cfg.SetupSecure.SSHPort)
				}
				if cfg.SetupSecure.Config == nil {
					t.Fatal("Config is nil")
				}
				if cfg.SetupSecure.Config.Type != "web" {
					t.Errorf("Type = %s, want web", cfg.SetupSecure.Config.Type)
				}
				if cfg.SetupSecure.Config.Domain != "example.com" {
					t.Errorf("Domain = %s, want example.com", cfg.SetupSecure.Config.Domain)
				}
				if cfg.SetupSecure.Firewall == nil {
					t.Fatal("Firewall is nil")
				}
				if len(cfg.SetupSecure.Firewall.OpenPorts) != 3 {
					t.Errorf("OpenPorts count = %d, want 3", len(cfg.SetupSecure.Firewall.OpenPorts))
				}
				if cfg.InstallTools == nil {
					t.Fatal("InstallTools is nil")
				}
				if len(cfg.InstallTools.Tools) != 3 {
					t.Errorf("Tools count = %d, want 3", len(cfg.InstallTools.Tools))
				}
			},
		},
		{
			name: "database server config",
			content: `.setup_secure{
	ssh_user: "dbadmin",
	user_ssh_rsa: "ssh-rsa AAAAB3...",
	ssh_port: 22022,
	.configuration{
		type: "database",
		db_engine: "mysql",
		root_password: "secret123"
	},
	.firewall{
		open_ports: [
			22022,
			3306
		]
	}
}

.install_tools{
	tools: [
		"mysql-server"
	]
}`,
			wantErr: false,
			check: func(t *testing.T, cfg *ServerConfig) {
				if cfg.SetupSecure == nil {
					t.Fatal("SetupSecure is nil")
				}
				if cfg.SetupSecure.Config.Type != "database" {
					t.Errorf("Type = %s, want database", cfg.SetupSecure.Config.Type)
				}
				if cfg.SetupSecure.Config.Options["db_engine"] != "mysql" {
					t.Errorf("db_engine = %s, want mysql", cfg.SetupSecure.Config.Options["db_engine"])
				}
			},
		},
		{
			name: "minimal config",
			content: `.setup_secure{
	ssh_user: "user",
	ssh_port: 22,
	.configuration{
		type: "basic"
	}
}`,
			wantErr: false,
			check: func(t *testing.T, cfg *ServerConfig) {
				if cfg.SetupSecure == nil {
					t.Fatal("SetupSecure is nil")
				}
				if cfg.SetupSecure.SSHUser != "user" {
					t.Errorf("SSHUser = %s, want user", cfg.SetupSecure.SSHUser)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := ParseConfig(tt.content)
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

func TestParseFirewall(t *testing.T) {
	tests := []struct {
		name      string
		lines     []string
		startIdx  int
		wantPorts []int
	}{
		{
			name: "basic firewall config",
			lines: []string{
				".firewall{",
				"	open_ports: [",
				"		80,",
				"		443,",
				"		22",
				"	]",
				"}",
			},
			startIdx:  0,
			wantPorts: []int{80, 443, 22},
		},
		{
			name: "single port",
			lines: []string{
				".firewall{",
				"	open_ports: [",
				"		8080",
				"	]",
				"}",
			},
			startIdx:  0,
			wantPorts: []int{8080},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			firewall, _ := parseFirewall(tt.lines, tt.startIdx)
			if len(firewall.OpenPorts) != len(tt.wantPorts) {
				t.Errorf("OpenPorts count = %d, want %d", len(firewall.OpenPorts), len(tt.wantPorts))
				return
			}
			for i, port := range tt.wantPorts {
				if firewall.OpenPorts[i] != port {
					t.Errorf("OpenPorts[%d] = %d, want %d", i, firewall.OpenPorts[i], port)
				}
			}
		})
	}
}

func TestParseInstallTools(t *testing.T) {
	tests := []struct {
		name      string
		lines     []string
		startIdx  int
		wantTools []string
	}{
		{
			name: "basic tools config",
			lines: []string{
				".install_tools{",
				"	tools: [",
				`		"git",`,
				`		"htop",`,
				`		"nginx"`,
				"	]",
				"}",
			},
			startIdx:  0,
			wantTools: []string{"git", "htop", "nginx"},
		},
		{
			name: "single tool",
			lines: []string{
				".install_tools{",
				"	tools: [",
				`		"docker.io"`,
				"	]",
				"}",
			},
			startIdx:  0,
			wantTools: []string{"docker.io"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installTools, _ := parseInstallTools(tt.lines, tt.startIdx)
			if len(installTools.Tools) != len(tt.wantTools) {
				t.Errorf("Tools count = %d, want %d", len(installTools.Tools), len(tt.wantTools))
				return
			}
			for i, tool := range tt.wantTools {
				if installTools.Tools[i] != tool {
					t.Errorf("Tools[%d] = %s, want %s", i, installTools.Tools[i], tool)
				}
			}
		})
	}
}
