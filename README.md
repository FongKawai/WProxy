
# WProxy

WProxy is an authenticated SOCKS5 and HTTP proxy tool that supports port multiplexing. It helps you browse the internet securely and efficiently.

[English](README.md) | [中文文档](README_zh-CN.md)

# Features

- [x] Supports authentication to prevent unauthorized access
- [x] Supports port multiplexing, allowing multiple proxy services on a single port
- [x] Lightweight and cross-platform, runs on Windows, Linux and macOS
- [x] Supports SOCKS5 and HTTP/HTTPS proxy protocols
- [x] Supports forwarding to target domains or IPs specified in HTTP/HTTPS headers

# Quick Installation

## One-click Installation for Linux

Use the following command to quickly install WProxy:

```bash
curl -s https://raw.githubusercontent.com/Wenpiner/WProxy/main/install.sh | sudo bash
```

The installation script will automatically:
1. Detect system architecture (supports amd64 and arm64)
2. Download the latest version from GitHub
3. Install to the system directory
4. Create a configuration file
5. Set up a system service (with auto-start on boot)

After installation, you can check the service status with:
```bash
systemctl status wproxy
```

## Uninstallation

Use the following command to quickly uninstall WProxy:

```bash
curl -s https://raw.githubusercontent.com/Wenpiner/WProxy/main/uninstall.sh | sudo bash
```

The uninstallation script will automatically:
1. Stop and disable the WProxy service
2. Remove the system service file
3. Delete the program files
4. Remove the configuration files

## Manual Installation

For Windows, macOS, or Linux without the install script:

1. Download the binary for your platform from [GitHub Releases](https://github.com/Wenpiner/WProxy/releases) (`wproxy.exe` on Windows)
2. Extract to any directory
3. Start with a config file or command-line flags (see **Configuration** below)

# Configuration

## Config file location

| Platform | Default config path | Notes |
|----------|---------------------|-------|
| Linux (`install.sh`) | `/etc/wproxy/config.yaml` | Created by the installer; systemd starts with `-c` |
| Windows / macOS / manual install | None | Create your own YAML and pass `-c`, or use CLI flags only |

Windows example (config next to the binary):

```powershell
.\wproxy.exe -c .\config.yaml
```

## Config file options (YAML)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `listen_addr` | string | `0.0.0.0:1080` | Listen address as `host:port` |
| `username` | string | empty | Proxy auth username; auth is enabled only when both `username` and `password` are set |
| `password` | string | empty | Proxy auth password |
| `certificate.key` | string | empty | TLS private key path (optional) |
| `certificate.cert` | string | empty | TLS certificate path (optional) |

Full example:

```yaml
listen_addr: "0.0.0.0:1080"
username: "admin"
password: "your_strong_password"
certificate:
  key: "/path/to/key.pem"
  cert: "/path/to/cert.pem"
```

Default after Linux one-click install (no `certificate`):

```yaml
listen_addr: "0.0.0.0:1080"
username: "admin"
password: "16-character random password"  # Automatically generated during installation
```

Notes:

1. The Linux installer generates a random 16-character password and prints it when finished
2. Restart the process or service after changing config (`systemctl restart wproxy`)
3. With `-c`, set `certificate` in YAML; `-certificate-key` / `-certificate-cert` apply only when **not** using `-c`

## Command-line flags

| Flag | Default | Description |
|------|---------|-------------|
| `-c` | empty | Path to config file |
| `-host` | `0.0.0.0` | Listen host (combined with `-port` when not using `-c`) |
| `-port` | `1080` | Listen port (1–65535) |
| `-username` | empty | Auth username |
| `-password` | empty | Auth password |
| `-certificate-key` | empty | TLS key path (without `-c` only) |
| `-certificate-cert` | empty | TLS cert path (without `-c` only) |

With `-c`, non-default `-host`, `-port`, `-username`, and `-password` override the config file.

Examples:

```bash
# Config file (Linux install path)
wproxy -c /etc/wproxy/config.yaml

# No config file
wproxy -host 127.0.0.1 -port 7890 -username admin -password secret

# TLS via CLI (no config file)
wproxy -host 0.0.0.0 -port 1080 -certificate-cert cert.pem -certificate-key key.pem
```

# Usage

- Configure your applications or browsers to use WProxy as a proxy server
- Enter the proxy server's address, port, and authentication credentials (if enabled)
- Start browsing the internet through the proxy server
- To forward to a target domain or IP via HTTP/HTTPS headers, set `X-Proxy-Host` and `X-Proxy-Scheme`. If proxy authentication is enabled, also send `Proxy-Authorization` (Basic auth). For example:
  ### Unauthenticated HTTP forwarding
  ```http
  GET /xxx/xxx HTTP/1.1
  Host: example.com
  X-Proxy-Host: target-domain.com
  X-Proxy-Scheme: http
  ```
  ### Unauthenticated HTTP forwarding with custom port (non-TLS)
  ```http
  GET /xxx/xxx HTTP/1.1
  Host: example.com
  X-Proxy-Host: target-domain.com:8080
  X-Proxy-Scheme: http
  ```

  ### Unauthenticated HTTPS forwarding with custom port (TLS)
  ```http
  GET /xxx/xxx HTTP/1.1
  Host: example.com
  X-Proxy-Host: target-domain.com:8443
  X-Proxy-Scheme: https
  ```

  ### Authenticated HTTPS forwarding
  ```http
  GET /xxx/xxx HTTP/1.1
  Host: example.com
  X-Proxy-Host: target-domain.com:8443
  X-Proxy-Scheme: https
  Proxy-Authorization: your_password
  ```

  The proxy server will forward the request to the specified target address and port based on these fields.

### ⚠️ Notes
1. `Proxy-Authorization` authenticates access to **WProxy**, not the host in `X-Proxy-Host`
2. If `username` / `password` are unset (or either is empty), the proxy has no auth and `Proxy-Authorization` is not required
3. The proxy sets `X-Proxy-Loop` internally to prevent loops; clients usually do not need to set it

# Service Management

After installation, you can manage the WProxy service with the following commands:

```bash
# Start the service
sudo systemctl start wproxy

# Stop the service
sudo systemctl stop wproxy

# Restart the service
sudo systemctl restart wproxy

# Check service status
sudo systemctl status wproxy

# Enable auto-start on boot
sudo systemctl enable wproxy

# Disable auto-start on boot
sudo systemctl disable wproxy
```

# Security Best Practices

When using WProxy in production environments, please follow these security recommendations:

## Authentication
- Always use strong, randomly-generated passwords (the installer generates one automatically)
- Change the default username from "admin" to something unique
- Rotate passwords regularly
- Never share credentials over insecure channels

## Network Security
- Use firewall rules to restrict access to the proxy port
- Consider using TLS/SSL certificates for HTTPS proxy connections
- Run the proxy behind a reverse proxy with additional security features when possible
- Monitor access logs for suspicious activity

## Configuration
- Restrict config file permissions (Linux install path: `chmod 600 /etc/wproxy/config.yaml`)
- Store certificates securely with appropriate file permissions
- Regularly update to the latest version to receive security patches

## Production Deployment
- Run the service with minimal required privileges
- Enable system logging and monitor for errors
- Set up automatic restarts on failure (already configured in systemd service)
- Use connection limits and rate limiting for public-facing deployments

# Contributing

If you find any issues or have suggestions for improvements, feel free to submit an issue or pull request. We're happy to improve WProxy together with the community.

# License

WProxy is released under the MIT License, allowing you to freely use, modify, and distribute this project.

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=Wenpiner/WProxy&type=Date)](https://star-history.com/#Wenpiner/WProxy&Date)
