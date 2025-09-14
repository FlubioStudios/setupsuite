package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// PackageManager interface for different package management systems
type PackageManager interface {
	Update() error
	Install(packages []string) error
	GetName() string
}

// DebianPackageManager for apt-based systems (Ubuntu, Debian)
type DebianPackageManager struct{}

func (pm *DebianPackageManager) Update() error {
	fmt.Println("Updating package list (apt)...")
	return exec.Command("apt-get", "update").Run()
}

func (pm *DebianPackageManager) Install(packages []string) error {
	args := append([]string{"install", "-y"}, packages...)
	return exec.Command("apt-get", args...).Run()
}

func (pm *DebianPackageManager) GetName() string {
	return "apt"
}

// RedHatPackageManager for yum/dnf-based systems (RHEL, CentOS, Fedora)
type RedHatPackageManager struct {
	useYum bool
}

func (pm *RedHatPackageManager) Update() error {
	fmt.Println("Updating package list (yum/dnf)...")
	if pm.useYum {
		return exec.Command("yum", "check-update").Run()
	}
	return exec.Command("dnf", "check-update").Run()
}

func (pm *RedHatPackageManager) Install(packages []string) error {
	args := append([]string{"install", "-y"}, packages...)
	if pm.useYum {
		return exec.Command("yum", args...).Run()
	}
	return exec.Command("dnf", args...).Run()
}

func (pm *RedHatPackageManager) GetName() string {
	if pm.useYum {
		return "yum"
	}
	return "dnf"
}

// ArchPackageManager for pacman-based systems (Arch Linux)
type ArchPackageManager struct{}

func (pm *ArchPackageManager) Update() error {
	fmt.Println("Updating package list (pacman)...")
	return exec.Command("pacman", "-Sy").Run()
}

func (pm *ArchPackageManager) Install(packages []string) error {
	args := append([]string{"-S", "--noconfirm"}, packages...)
	return exec.Command("pacman", args...).Run()
}

func (pm *ArchPackageManager) GetName() string {
	return "pacman"
}

// AlpinePackageManager for apk-based systems (Alpine Linux)
type AlpinePackageManager struct{}

func (pm *AlpinePackageManager) Update() error {
	fmt.Println("Updating package list (apk)...")
	return exec.Command("apk", "update").Run()
}

func (pm *AlpinePackageManager) Install(packages []string) error {
	args := append([]string{"add"}, packages...)
	return exec.Command("apk", args...).Run()
}

func (pm *AlpinePackageManager) GetName() string {
	return "apk"
}

// OpenSUSEPackageManager for zypper-based systems (openSUSE)
type OpenSUSEPackageManager struct{}

func (pm *OpenSUSEPackageManager) Update() error {
	fmt.Println("Updating package list (zypper)...")
	return exec.Command("zypper", "refresh").Run()
}

func (pm *OpenSUSEPackageManager) Install(packages []string) error {
	args := append([]string{"install", "-y"}, packages...)
	return exec.Command("zypper", args...).Run()
}

func (pm *OpenSUSEPackageManager) GetName() string {
	return "zypper"
}

// DetectPackageManager detects the package manager based on the system
func DetectPackageManager() (PackageManager, error) {
	// Check for package manager binaries in order of preference
	packageManagers := []struct {
		command string
		pm      PackageManager
	}{
		{"apt-get", &DebianPackageManager{}},
		{"dnf", &RedHatPackageManager{useYum: false}},
		{"yum", &RedHatPackageManager{useYum: true}},
		{"pacman", &ArchPackageManager{}},
		{"zypper", &OpenSUSEPackageManager{}},
		{"apk", &AlpinePackageManager{}},
	}

	for _, pmInfo := range packageManagers {
		if _, err := exec.LookPath(pmInfo.command); err == nil {
			fmt.Printf("Detected package manager: %s\n", pmInfo.pm.GetName())
			return pmInfo.pm, nil
		}
	}

	return nil, fmt.Errorf("no supported package manager found")
}

// DetectDistribution detects the Linux distribution
func DetectDistribution() (string, string, error) {
	// Try to read /etc/os-release first (standard)
	if file, err := os.Open("/etc/os-release"); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		distro := make(map[string]string)

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.Contains(line, "=") {
				parts := strings.SplitN(line, "=", 2)
				key := strings.TrimSpace(parts[0])
				value := strings.Trim(strings.TrimSpace(parts[1]), `"`)
				distro[key] = value
			}
		}

		if id, ok := distro["ID"]; ok {
			version := distro["VERSION_ID"]
			return id, version, nil
		}
	}

	// Fallback: check for specific distribution files
	distroFiles := map[string]string{
		"/etc/debian_version": "debian",
		"/etc/redhat-release": "rhel",
		"/etc/centos-release": "centos",
		"/etc/fedora-release": "fedora",
		"/etc/arch-release":   "arch",
		"/etc/alpine-release": "alpine",
		"/etc/SuSE-release":   "opensuse",
	}

	for file, distro := range distroFiles {
		if _, err := os.Stat(file); err == nil {
			return distro, "unknown", nil
		}
	}

	return "unknown", "unknown", fmt.Errorf("could not detect Linux distribution")
}

// GetFirewallManager returns the appropriate firewall management commands
func GetFirewallManager() (string, error) {
	// Check for firewall managers in order of preference
	firewallManagers := []string{"ufw", "firewalld", "iptables"}

	for _, fw := range firewallManagers {
		if _, err := exec.LookPath(fw); err == nil {
			fmt.Printf("Detected firewall manager: %s\n", fw)
			return fw, nil
		}
	}

	return "", fmt.Errorf("no supported firewall manager found")
}

// GetServiceManager returns the service management system
func GetServiceManager() (string, error) {
	// Check for service managers
	serviceManagers := []string{"systemctl", "service", "rc-service"}

	for _, sm := range serviceManagers {
		if _, err := exec.LookPath(sm); err == nil {
			fmt.Printf("Detected service manager: %s\n", sm)
			return sm, nil
		}
	}

	return "", fmt.Errorf("no supported service manager found")
}
