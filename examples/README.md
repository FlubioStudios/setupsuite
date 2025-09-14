# SetupSuite Examples

This directory contains example configuration files for different server types.

## Usage

1. **Generate a configuration template:**
   ```bash
   sudo ./suite -generate -type web -config /path/to/config.sscfg
   ```

2. **Edit the generated configuration file** to match your requirements

3. **Run the setup:**
   ```bash
   sudo ./suite -config /path/to/config.sscfg
   ```

## Available Server Types

- **web** - Web server with Nginx and SSL
- **database** - Database server (MySQL/PostgreSQL)  
- **docker** - Docker host with optimized settings
- **proxy** - Reverse proxy server
- **build** - Build/CI server with development tools

## Configuration Files

- `web-server.sscfg` - Complete web server setup
- `database-server.sscfg` - Database server configuration
- `docker-host.sscfg` - Docker host setup
- `proxy-server.sscfg` - Reverse proxy configuration
- `build-server.sscfg` - CI/CD build server setup

## Important Notes

- Replace `REPLACE_WITH_YOUR_SSH_KEY` with your actual SSH public key
- Replace `REPLACE_WITH_SECURE_PASSWORD` with a strong password
- Update domain names and email addresses as needed
- Review firewall port configurations for your use case
