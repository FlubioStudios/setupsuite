# Distribution Support

SetupSuite now supports multiple Linux distributions automatically by detecting the package manager, service manager, and firewall system at runtime.

## Supported Distributions

### Package Managers
- **APT** (apt-get) - Ubuntu, Debian, Linux Mint
- **DNF** - Fedora 22+, RHEL 8+, CentOS 8+
- **YUM** - RHEL 7, CentOS 7, Amazon Linux
- **Pacman** - Arch Linux, Manjaro
- **Zypper** - openSUSE, SUSE Linux Enterprise
- **APK** - Alpine Linux

### Service Managers
- **systemd** (systemctl) - Most modern distributions
- **SysV Init** (service) - Older RHEL/CentOS, Debian
- **OpenRC** (rc-service) - Alpine Linux, Gentoo

### Firewall Managers
- **UFW** - Ubuntu, Linux Mint (user-friendly frontend)
- **firewalld** - Fedora, RHEL 7+, CentOS 7+
- **iptables** - Universal fallback for older systems

## Auto-Detection Process

SetupSuite automatically detects your system configuration:

1. **Distribution Detection**: Reads `/etc/os-release` and fallback files
2. **Package Manager**: Checks for available binaries in order of preference
3. **Service Manager**: Detects systemd, SysV, or OpenRC
4. **Firewall**: Finds UFW, firewalld, or falls back to iptables

## Distribution-Specific Handling

### Package Installation Examples

#### Ubuntu/Debian (APT)
```bash
apt-get update
apt-get install -y nginx mysql-server
```

#### Fedora (DNF)
```bash
dnf check-update
dnf install -y nginx mysql-server
```

#### CentOS 7 (YUM)
```bash
yum check-update
yum install -y nginx mysql-server
```

#### Arch Linux (Pacman)
```bash
pacman -Sy
pacman -S --noconfirm nginx mysql
```

#### Alpine Linux (APK)
```bash
apk update
apk add nginx mysql
```

### Service Management Examples

#### systemd (Most Modern Distributions)
```bash
systemctl enable nginx
systemctl start nginx
systemctl restart nginx
```

#### SysV Init (Older Systems)
```bash
chkconfig nginx on
service nginx start
service nginx restart
```

#### OpenRC (Alpine)
```bash
rc-update add nginx default
rc-service nginx start
rc-service nginx restart
```

### Firewall Configuration

#### UFW (Ubuntu/Debian)
```bash
ufw default deny incoming
ufw default allow outgoing
ufw allow 80
ufw allow 443
ufw --force enable
```

#### firewalld (RHEL/Fedora)
```bash
firewall-cmd --permanent --add-port=80/tcp
firewall-cmd --permanent --add-port=443/tcp
firewall-cmd --reload
```

#### iptables (Universal Fallback)
```bash
iptables -A INPUT -p tcp --dport 80 -j ACCEPT
iptables -A INPUT -p tcp --dport 443 -j ACCEPT
iptables-save
```

## Node.js Installation by Distribution

### Ubuntu/Debian
Uses NodeSource repository for latest LTS:
```bash
curl -fsSL https://deb.nodesource.com/setup_lts.x | bash -
apt-get install -y nodejs
```

### RHEL/CentOS/Fedora
Uses NodeSource RPM repository:
```bash
curl -fsSL https://rpm.nodesource.com/setup_lts.x | bash -
dnf install -y nodejs  # or yum
```

### Arch Linux
Uses official repositories:
```bash
pacman -S --noconfirm nodejs npm
```

### Alpine Linux
Uses Alpine repositories:
```bash
apk add nodejs npm
```

## MySQL/PostgreSQL Service Names

Different distributions use different service names:

| Distribution | MySQL Service | PostgreSQL Service |
|-------------|---------------|-------------------|
| Ubuntu/Debian | mysql | postgresql |
| RHEL/CentOS | mysqld | postgresql |
| Fedora | mysqld | postgresql |
| Alpine | mariadb | postgresql |
| Arch | mysqld | postgresql |

## Configuration Examples by Distribution

### Ubuntu 20.04+ Server
```
.install_tools{
    tools: [
        "nginx",
        "mysql-server",
        "php8.1",
        "php8.1-fpm",
        "php8.1-mysql",
        "certbot",
        "python3-certbot-nginx",
        "ufw",
        "fail2ban"
    ]
}
```

### CentOS 8 / RHEL 8
```
.install_tools{
    tools: [
        "nginx",
        "mysql-server",
        "php",
        "php-fpm",
        "php-mysqlnd",
        "certbot",
        "python3-certbot-nginx",
        "firewalld",
        "fail2ban"
    ]
}
```

### Fedora 35+
```
.install_tools{
    tools: [
        "nginx",
        "mysql-community-server",
        "php",
        "php-fpm",
        "php-mysqlnd",
        "certbot",
        "python3-certbot-nginx",
        "firewalld"
    ]
}
```

### Alpine Linux
```
.install_tools{
    tools: [
        "nginx",
        "mariadb",
        "mariadb-client",
        "php81",
        "php81-fpm",
        "php81-mysqli",
        "certbot",
        "iptables"
    ]
}
```

### Arch Linux
```
.install_tools{
    tools: [
        "nginx",
        "mysql",
        "php",
        "php-fpm",
        "certbot",
        "certbot-nginx",
        "ufw"
    ]
}
```

## Troubleshooting

### Package Not Found
If a package is not found, check the distribution-specific package name:
```bash
# Ubuntu/Debian
apt search package-name

# Fedora/RHEL
dnf search package-name

# Arch
pacman -Ss package-name

# Alpine
apk search package-name
```

### Service Failures
Check service status with the appropriate command:
```bash
# systemd
systemctl status service-name

# SysV
service service-name status

# OpenRC
rc-service service-name status
```

### Firewall Issues
Verify firewall rules:
```bash
# UFW
ufw status

# firewalld
firewall-cmd --list-all

# iptables
iptables -L
```

## Adding Support for New Distributions

To add support for a new distribution:

1. **Add package manager** in `package_manager.go`
2. **Add service manager** support in `service_manager.go`
3. **Update detection logic** in `DetectDistribution()`
4. **Add distribution-specific handling** in server setup functions
5. **Update documentation** with new examples
