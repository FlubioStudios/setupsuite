package main

import (
	"fmt"
	"os/exec"
)

// ServiceManager handles service operations across different init systems
type ServiceManager struct {
	manager string
}

// NewServiceManager creates a service manager based on detected init system
func NewServiceManager() (*ServiceManager, error) {
	manager, err := GetServiceManager()
	if err != nil {
		return nil, err
	}
	return &ServiceManager{manager: manager}, nil
}

// Enable enables a service
func (sm *ServiceManager) Enable(serviceName string) error {
	fmt.Printf("Enabling service %s using %s\n", serviceName, sm.manager)

	switch sm.manager {
	case "systemctl":
		return exec.Command("systemctl", "enable", serviceName).Run()
	case "service":
		// For SysV init, enable usually means adding to runlevels
		return exec.Command("chkconfig", serviceName, "on").Run()
	case "rc-service":
		// For OpenRC (Alpine)
		return exec.Command("rc-update", "add", serviceName, "default").Run()
	default:
		return fmt.Errorf("unsupported service manager: %s", sm.manager)
	}
}

// Start starts a service
func (sm *ServiceManager) Start(serviceName string) error {
	fmt.Printf("Starting service %s using %s\n", serviceName, sm.manager)

	switch sm.manager {
	case "systemctl":
		return exec.Command("systemctl", "start", serviceName).Run()
	case "service":
		return exec.Command("service", serviceName, "start").Run()
	case "rc-service":
		return exec.Command("rc-service", serviceName, "start").Run()
	default:
		return fmt.Errorf("unsupported service manager: %s", sm.manager)
	}
}

// Restart restarts a service
func (sm *ServiceManager) Restart(serviceName string) error {
	fmt.Printf("Restarting service %s using %s\n", serviceName, sm.manager)

	switch sm.manager {
	case "systemctl":
		return exec.Command("systemctl", "restart", serviceName).Run()
	case "service":
		return exec.Command("service", serviceName, "restart").Run()
	case "rc-service":
		return exec.Command("rc-service", serviceName, "restart").Run()
	default:
		return fmt.Errorf("unsupported service manager: %s", sm.manager)
	}
}

// Stop stops a service
func (sm *ServiceManager) Stop(serviceName string) error {
	fmt.Printf("Stopping service %s using %s\n", serviceName, sm.manager)

	switch sm.manager {
	case "systemctl":
		return exec.Command("systemctl", "stop", serviceName).Run()
	case "service":
		return exec.Command("service", serviceName, "stop").Run()
	case "rc-service":
		return exec.Command("rc-service", serviceName, "stop").Run()
	default:
		return fmt.Errorf("unsupported service manager: %s", sm.manager)
	}
}

// Reload reloads a service configuration
func (sm *ServiceManager) Reload(serviceName string) error {
	fmt.Printf("Reloading service %s using %s\n", serviceName, sm.manager)

	switch sm.manager {
	case "systemctl":
		return exec.Command("systemctl", "reload", serviceName).Run()
	case "service":
		return exec.Command("service", serviceName, "reload").Run()
	case "rc-service":
		return exec.Command("rc-service", serviceName, "reload").Run()
	default:
		return fmt.Errorf("unsupported service manager: %s", sm.manager)
	}
}

// IsActive checks if a service is active
func (sm *ServiceManager) IsActive(serviceName string) bool {
	switch sm.manager {
	case "systemctl":
		err := exec.Command("systemctl", "is-active", "--quiet", serviceName).Run()
		return err == nil
	case "service":
		err := exec.Command("service", serviceName, "status").Run()
		return err == nil
	case "rc-service":
		err := exec.Command("rc-service", serviceName, "status").Run()
		return err == nil
	default:
		return false
	}
}
