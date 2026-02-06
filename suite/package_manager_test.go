package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectPackageManager(t *testing.T) {
	pm, err := DetectPackageManager()
	if err != nil {
		t.Logf("DetectPackageManager() error = %v (may be expected in some environments)", err)
		return
	}

	if pm == nil {
		t.Error("DetectPackageManager() returned nil package manager")
		return
	}

	name := pm.GetName()
	validPMs := []string{"apt", "yum", "dnf", "pacman", "apk", "zypper"}
	found := false
	for _, valid := range validPMs {
		if name == valid {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("DetectPackageManager() returned unknown package manager: %s", name)
	}
}

func TestDebianPackageManager(t *testing.T) {
	pm := &DebianPackageManager{}
	if pm.GetName() != "apt" {
		t.Errorf("GetName() = %s, want apt", pm.GetName())
	}
}

func TestRedHatPackageManager(t *testing.T) {
	pm := &RedHatPackageManager{useYum: true}
	if pm.GetName() != "yum" {
		t.Errorf("GetName() = %s, want yum", pm.GetName())
	}

	pm = &RedHatPackageManager{useYum: false}
	if pm.GetName() != "dnf" {
		t.Errorf("GetName() = %s, want dnf", pm.GetName())
	}
}

func TestArchPackageManager(t *testing.T) {
	pm := &ArchPackageManager{}
	if pm.GetName() != "pacman" {
		t.Errorf("GetName() = %s, want pacman", pm.GetName())
	}
}

func TestAlpinePackageManager(t *testing.T) {
	pm := &AlpinePackageManager{}
	if pm.GetName() != "apk" {
		t.Errorf("GetName() = %s, want apk", pm.GetName())
	}
}

func TestOpenSUSEPackageManager(t *testing.T) {
	pm := &OpenSUSEPackageManager{}
	if pm.GetName() != "zypper" {
		t.Errorf("GetName() = %s, want zypper", pm.GetName())
	}
}

func TestDetectDistribution(t *testing.T) {
	distro, version, err := DetectDistribution()
	if err != nil {
		t.Logf("DetectDistribution() error = %v", err)
		return
	}

	if distro == "" {
		t.Error("DetectDistribution() returned empty distro")
	}

	t.Logf("Detected distribution: %s %s", distro, version)
}

func findExecutable(name string) (string, error) {
	paths := []string{
		"/usr/bin",
		"/bin",
		"/usr/sbin",
		"/sbin",
		"/usr/local/bin",
	}

	for _, dir := range paths {
		fullPath := filepath.Join(dir, name)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, nil
		}
	}

	return "", os.ErrNotExist
}
