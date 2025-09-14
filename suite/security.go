package main

import (
	"fmt"
	"os"
	"suite/suite/config"
)

// setupBasicSecurity handles the basic security configuration
func setupBasicSecurity(cfg *config.SetupSecure) error {
	VerboseLogger.LogInfo("Starting basic security setup")

	if cfg == nil {
		return fmt.Errorf("security config is nil")
	}

	// Add user
	if cfg.SSHUser != "" {
		fmt.Printf("Adding user %s\n", cfg.SSHUser)
		VerboseLogger.LogInfo("Creating user: %s", cfg.SSHUser)

		err := VerboseCommandRun("adduser", "--disabled-password", "--gecos", "", cfg.SSHUser)
		if err != nil {
			VerboseLogger.LogWarning("Error adding user: %s", err)
			fmt.Printf("Warning: Error adding user: %s\n", err)
		}

		// Group add
		fmt.Println("Adding group sshuser")
		VerboseCommandRun("groupadd", "sshuser")

		// Usermod
		fmt.Println("Adding user to group")
		VerboseCommandRun("usermod", "-aG", "sshuser", cfg.SSHUser)

		// Add to sudo
		fmt.Println("Configuring sudo")
		sudoCmd := fmt.Sprintf("echo '%s ALL=(ALL:ALL) ALL' | EDITOR='tee -a' visudo", cfg.SSHUser)
		VerboseCommandRun("sh", "-c", sudoCmd)

		// Setup SSH keys
		if cfg.UserSSHRSA != "" && cfg.UserSSHRSA != "REPLACE_WITH_YOUR_SSH_KEY" {
			setupSSHKeys(cfg.SSHUser, cfg.UserSSHRSA)
		}
	}

	// Configure SSH daemon
	if cfg.SSHPort > 0 {
		configureSSHD(cfg.SSHPort)
	}

	// Setup root bashrc
	setupRootBashrc()

	// Update system
	fmt.Println("Updating system")
	VerboseLogger.LogInfo("Starting system update")
	VerboseCommandRun("apt-get", "update")
	VerboseCommandRun("apt-get", "upgrade", "-y")

	VerboseLogger.LogInfo("Basic security setup completed")
	return nil
}

func setupSSHKeys(user, sshKey string) {
	VerboseLogger.LogInfo("Setting up SSH keys for user: %s", user)

	homeDir := "/home/" + user
	sshDir := homeDir + "/.ssh"

	// Create .ssh directory
	VerboseLogger.LogInfo("Creating SSH directory: %s", sshDir)
	VerboseMkdirAll(sshDir, 0700)

	// Create authorized_keys file
	authKeys := sshDir + "/authorized_keys"
	VerboseLogger.LogInfo("Creating authorized_keys file: %s", authKeys)
	file, err := VerboseOpenFile(authKeys, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		VerboseLogger.LogError("Error creating authorized_keys: %s", err)
		fmt.Printf("Error creating authorized_keys: %s\n", err)
		return
	}
	defer file.Close()

	VerboseWriteString(file, sshKey+"\n")

	// Set ownership and permissions
	VerboseLogger.LogInfo("Setting ownership and permissions for SSH files")
	VerboseCommandRun("chown", "-R", user+":"+user, sshDir)
	VerboseCommandRun("chmod", "700", sshDir)
	VerboseCommandRun("chmod", "600", authKeys)

	// Copy to root for security backup
	rootSshDir := "/root/.ssh"
	VerboseLogger.LogInfo("Creating backup in root SSH directory: %s", rootSshDir)
	VerboseMkdirAll(rootSshDir, 0700)
	VerboseCommandRun("cp", authKeys, rootSshDir+"/authorized_keys_copied_due_to_security")
}

func configureSSHD(port int) {
	fmt.Printf("Configuring SSH daemon on port %d\n", port)
	VerboseLogger.LogInfo("Configuring SSH daemon on port %d", port)

	// Backup original sshd_config
	VerboseLogger.LogInfo("Backing up original sshd_config")
	VerboseCommandRun("cp", "/etc/ssh/sshd_config", "/etc/ssh/sshd_config.backup")

	// Create new sshd_config
	sshdConfig := fmt.Sprintf(`# SetupSuite generated SSH configuration
Port %d
PermitRootLogin no
AllowGroups sshuser
PubkeyAuthentication yes
PasswordAuthentication no
PermitEmptyPasswords no
ChallengeResponseAuthentication no
UsePAM yes
X11Forwarding yes
PrintMotd no
AcceptEnv LANG LC_*
Subsystem sftp /usr/lib/openssh/sftp-server

# Include original config for other settings
Include /etc/ssh/sshd_config.backup
`, port)

	VerboseLogger.LogInfo("Writing new SSH configuration")
	file, err := VerboseCreate("/etc/ssh/sshd_config")
	if err == nil {
		VerboseWriteString(file, sshdConfig)
		file.Close()

		// Test configuration and restart SSH
		VerboseLogger.LogInfo("Testing SSH configuration")
		if VerboseCommandRun("sshd", "-t") == nil {
			VerboseLogger.LogInfo("SSH configuration test passed, restarting SSH service")
			VerboseCommandRun("systemctl", "restart", "sshd")
		} else {
			fmt.Println("SSH configuration test failed, reverting...")
			VerboseLogger.LogError("SSH configuration test failed, reverting to backup")
			VerboseCommandRun("mv", "/etc/ssh/sshd_config.backup", "/etc/ssh/sshd_config")
		}
	} else {
		VerboseLogger.LogError("Failed to create SSH config file: %s", err)
	}
}

func setupRootBashrc() {
	fmt.Println("Configuring root bashrc")
	VerboseLogger.LogInfo("Configuring root bashrc")

	rootBashrc := "/root/.bashrc"
	file, err := VerboseOpenFile(rootBashrc, os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		bashrcContent := "\n# SetupSuite additions\nexport LS_OPTIONS='--color=auto'\nalias ls='ls -la $LS_OPTIONS'\nPATH=$PATH:/usr/sbin\n"
		VerboseWriteString(file, bashrcContent)
		file.Close()
		VerboseLogger.LogInfo("Root bashrc configuration completed")
	} else {
		VerboseLogger.LogError("Failed to open root bashrc: %s", err)
	}
}
