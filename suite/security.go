package main

import (
	"fmt"
	"os"
	"os/exec"
	"suite/suite/config"
)

// setupBasicSecurity handles the basic security configuration
func setupBasicSecurity(cfg *config.SetupSecure) error {
	if cfg == nil {
		return fmt.Errorf("security config is nil")
	}

	// Add user
	if cfg.SSHUser != "" {
		fmt.Printf("Adding user %s\n", cfg.SSHUser)
		err := exec.Command("adduser", "--disabled-password", "--gecos", "", cfg.SSHUser).Run()
		if err != nil {
			fmt.Printf("Warning: Error adding user: %s\n", err)
		}

		// Group add
		fmt.Println("Adding group sshuser")
		exec.Command("groupadd", "sshuser").Run()

		// Usermod
		fmt.Println("Adding user to group")
		exec.Command("usermod", "-aG", "sshuser", cfg.SSHUser).Run()

		// Add to sudo
		fmt.Println("Configuring sudo")
		cmd := exec.Command("sh", "-c", "echo '"+cfg.SSHUser+" ALL=(ALL:ALL) ALL' | EDITOR='tee -a' visudo")
		cmd.Run()

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
	exec.Command("apt-get", "update").Run()
	exec.Command("apt-get", "upgrade", "-y").Run()

	return nil
}

func setupSSHKeys(user, sshKey string) {
	homeDir := "/home/" + user
	sshDir := homeDir + "/.ssh"

	// Create .ssh directory
	os.MkdirAll(sshDir, 0700)

	// Create authorized_keys file
	authKeys := sshDir + "/authorized_keys"
	file, err := os.OpenFile(authKeys, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Error creating authorized_keys: %s\n", err)
		return
	}
	defer file.Close()

	file.WriteString(sshKey + "\n")

	// Set ownership and permissions
	exec.Command("chown", "-R", user+":"+user, sshDir).Run()
	exec.Command("chmod", "700", sshDir).Run()
	exec.Command("chmod", "600", authKeys).Run()

	// Copy to root for security backup
	rootSshDir := "/root/.ssh"
	os.MkdirAll(rootSshDir, 0700)
	exec.Command("cp", authKeys, rootSshDir+"/authorized_keys_copied_due_to_security").Run()
}

func configureSSHD(port int) {
	fmt.Printf("Configuring SSH daemon on port %d\n", port)

	// Backup original sshd_config
	exec.Command("cp", "/etc/ssh/sshd_config", "/etc/ssh/sshd_config.backup").Run()

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

	file, err := os.Create("/etc/ssh/sshd_config")
	if err == nil {
		file.WriteString(sshdConfig)
		file.Close()

		// Test configuration and restart SSH
		if exec.Command("sshd", "-t").Run() == nil {
			exec.Command("systemctl", "restart", "sshd").Run()
		} else {
			fmt.Println("SSH configuration test failed, reverting...")
			exec.Command("mv", "/etc/ssh/sshd_config.backup", "/etc/ssh/sshd_config").Run()
		}
	}
}

func setupRootBashrc() {
	fmt.Println("Configuring root bashrc")

	rootBashrc := "/root/.bashrc"
	file, err := os.OpenFile(rootBashrc, os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		bashrcContent := "\n# SetupSuite additions\nexport LS_OPTIONS='--color=auto'\nalias ls='ls -la $LS_OPTIONS'\nPATH=$PATH:/usr/sbin\n"
		file.WriteString(bashrcContent)
		file.Close()
	}
}
