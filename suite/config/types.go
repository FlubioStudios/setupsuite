package config

// ServerConfig represents the main configuration structure
type ServerConfig struct {
	SetupSecure  *SetupSecure  `json:"setup_secure"`
	InstallTools *InstallTools `json:"install_tools"`
}

// SetupSecure contains security and basic setup configuration
type SetupSecure struct {
	SSHUser    string    `json:"ssh_user"`
	UserSSHRSA string    `json:"user_ssh_rsa"`
	SSHPort    int       `json:"ssh_port"`
	Config     *Config   `json:"configuration"`
	Firewall   *Firewall `json:"firewall"`
}

// Config contains server type and specific configuration
type Config struct {
	Type     string            `json:"type"`
	Domain   string            `json:"domain,omitempty"`
	Email    string            `json:"email,omitempty"`
	Database *DatabaseConfig   `json:"database,omitempty"`
	Docker   *DockerConfig     `json:"docker,omitempty"`
	Proxy    *ProxyConfig      `json:"proxy,omitempty"`
	Options  map[string]string `json:"options,omitempty"`
}

// DatabaseConfig contains database-specific settings
type DatabaseConfig struct {
	Engine   string `json:"engine"` // mysql, postgresql, mongodb
	RootPass string `json:"root_pass"`
	DBName   string `json:"db_name,omitempty"`
	DBUser   string `json:"db_user,omitempty"`
	DBPass   string `json:"db_pass,omitempty"`
}

// DockerConfig contains docker-specific settings
type DockerConfig struct {
	LogDriver  string            `json:"log_driver"`
	LogOptions map[string]string `json:"log_options"`
	Compose    bool              `json:"compose"`
}

// ProxyConfig contains proxy-specific settings
type ProxyConfig struct {
	Upstreams []Upstream `json:"upstreams"`
	SSL       bool       `json:"ssl"`
}

// Upstream represents a proxy upstream server
type Upstream struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Port int    `json:"port"`
}

// Firewall contains firewall configuration
type Firewall struct {
	OpenPorts []int `json:"open_ports"`
}

// InstallTools contains tools to be installed
type InstallTools struct {
	Tools []string `json:"tools"`
}

// ServerType constants
const (
	ServerTypeWeb      = "web"
	ServerTypeDatabase = "database"
	ServerTypeDocker   = "docker"
	ServerTypeProxy    = "proxy"
	ServerTypeBuild    = "build"
)
