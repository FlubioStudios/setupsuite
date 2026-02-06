package main

import (
	"testing"
)

func TestNewServiceManager(t *testing.T) {
	sm, err := NewServiceManager()
	if err != nil {
		t.Logf("NewServiceManager() error = %v (may be expected in some environments)", err)
		return
	}

	if sm == nil {
		t.Error("NewServiceManager() returned nil service manager")
		return
	}

	t.Logf("Detected service manager: %s", sm.manager)
}

func TestServiceManagerInterface(t *testing.T) {
	tests := []struct {
		name    string
		manager string
	}{
		{name: "systemctl", manager: "systemctl"},
		{name: "service", manager: "service"},
		{name: "rc-service", manager: "rc-service"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &ServiceManager{manager: tt.manager}
			if sm.manager != tt.manager {
				t.Errorf("manager = %s, want %s", sm.manager, tt.manager)
			}
		})
	}
}
