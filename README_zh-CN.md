# WProxy

WProxy 是一个带鉴权的 SOCKS5 和 HTTP 代理工具，支持端口复用。它可以帮助您在安全和高效的环境下上网。

[English](README.md) | [中文文档](README_zh-CN.md)

# 特性

- [x] 支持身份验证，可以防止未经授权的访问
- [x] 支持端口复用，可以复用同一个端口提供多种代理服务
- [x] 轻量级和跨平台，可在 Windows、Linux 和 macOS 上运行
- [x] 支持 SOCKS5 和 HTTP/HTTPS 代理协议
- [x] 支持通过 HTTP/HTTPS 请求头指定目标域名或 IP 进行转发

# 快速安装

## Linux 系统一键安装

使用以下命令快速安装 WProxy：

```bash
curl -s https://raw.githubusercontent.com/Wenpiner/WProxy/main/install.sh | sudo bash
```

安装脚本会自动完成以下操作：
1. 检测系统架构（支持 amd64 和 arm64）
2. 从 GitHub 下载最新版本
3. 安装到系统目录
4. 创建配置文件
5. 设置系统服务（支持开机自启）

安装完成后，您可以通过以下命令查看服务状态：
```bash
systemctl status wproxy
```

## 卸载

使用以下命令快速卸载 WProxy：

```bash
curl -s https://raw.githubusercontent.com/Wenpiner/WProxy/main/uninstall.sh | sudo bash
```

卸载脚本会自动完成以下操作：
1. 停止并禁用 WProxy 服务
2. 删除系统服务文件
3. 删除程序文件
4. 删除配置文件

## 手动安装

适用于 Windows、macOS，或未使用一键脚本的 Linux：

1. 从 [GitHub Releases](https://github.com/Wenpiner/WProxy/releases) 下载对应平台的二进制文件（Windows 为 `wproxy.exe`）
2. 解压到任意目录
3. 使用配置文件或命令行参数启动（见下方「配置」）

# 配置

## 配置文件位置

| 平台 | 默认配置文件路径 | 说明 |
|------|------------------|------|
| Linux（`install.sh` 安装） | `/etc/wproxy/config.yaml` | 安装脚本自动创建，systemd 以 `-c` 加载 |
| Windows / macOS / 手动安装 | 无内置默认路径 | 需自行创建 YAML，启动时用 `-c` 指定；或完全不使用配置文件 |

Windows 示例（配置文件与程序同目录）：

```powershell
.\wproxy.exe -c .\config.yaml
```

## 配置文件选项（YAML）

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `listen_addr` | 字符串 | `0.0.0.0:1080` | 监听地址，格式为 `host:port` |
| `username` | 字符串 | 空 | 代理鉴权用户名；须与 `password` 同时非空才启用鉴权 |
| `password` | 字符串 | 空 | 代理鉴权密码 |
| `certificate.key` | 字符串 | 空 | TLS 私钥文件路径（可选） |
| `certificate.cert` | 字符串 | 空 | TLS 证书文件路径（可选） |

完整配置示例：

```yaml
listen_addr: "0.0.0.0:1080"
username: "admin"
password: "your_strong_password"
certificate:
  key: "/path/to/key.pem"
  cert: "/path/to/cert.pem"
```

Linux 一键安装后的默认配置（不含 `certificate`）：

```yaml
listen_addr: "0.0.0.0:1080"
username: "admin"
password: "16位随机密码"  # 安装时自动生成
```

注意事项：

1. Linux 一键安装时会自动生成 16 位随机密码，安装结束会在终端显示，请妥善保管
2. 修改配置后需重启进程或服务（`systemctl restart wproxy`）
3. 使用 `-c` 时，`certificate` 须在 YAML 中配置；命令行 `-certificate-key` / `-certificate-cert` 仅在**未**使用 `-c` 时生效

## 命令行参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-c` | 空 | 配置文件路径 |
| `-host` | `0.0.0.0` | 监听主机（未使用 `-c` 时与 `-port` 组成监听地址） |
| `-port` | `1080` | 监听端口（1–65535） |
| `-username` | 空 | 鉴权用户名 |
| `-password` | 空 | 鉴权密码 |
| `-certificate-key` | 空 | TLS 私钥路径（仅在不使用 `-c` 时） |
| `-certificate-cert` | 空 | TLS 证书路径（仅在不使用 `-c` 时） |

使用 `-c` 时，命令行中的 `-host`、`-port`、`-username`、`-password` 可在非默认值时覆盖配置文件对应项。

启动示例：

```bash
# 使用配置文件（Linux 安装路径）
wproxy -c /etc/wproxy/config.yaml

# 不使用配置文件，直接指定参数
wproxy -host 127.0.0.1 -port 7890 -username admin -password secret

# 指定 TLS 证书（无配置文件）
wproxy -host 0.0.0.0 -port 1080 -certificate-cert cert.pem -certificate-key key.pem
```

# 使用

- 将您的应用程序或浏览器配置为使用 WProxy 作为代理服务器
- 输入代理服务器的地址和端口，以及身份验证所需的用户名和密码（如果已启用）
- 开始通过代理服务器访问互联网
- 如果需要通过 HTTP/HTTPS 请求头指定目标域名或 IP 进行转发，请在请求头中添加 `X-Proxy-Host`、`X-Proxy-Scheme`；访问代理本身若已启用鉴权，还需 `Proxy-Authorization`（Basic 认证）。例如：
  ### 无鉴权、HTTP 转发
  ```http
  GET /xxx/xxx HTTP/1.1
  Host: example.com
  X-Proxy-Host: target-domain.com
  X-Proxy-Scheme: http
  ```
  ### 无鉴权、HTTP 转发、自定义端口(非TLS)
  ```http
  GET /xxx/xxx HTTP/1.1
  Host: example.com
  X-Proxy-Host: target-domain.com:8080
  X-Proxy-Scheme: http
  ```

  ### 无鉴权、HTTPS 转发、自定义端口(TLS)
  ```http
  GET /xxx/xxx HTTP/1.1
  Host: example.com
  X-Proxy-Host: target-domain.com:8443
  X-Proxy-Scheme: https
  ```

  ### 鉴权、HTTPS 转发 
  ```http
  GET /xxx/xxx HTTP/1.1
  Host: example.com
  X-Proxy-Host: target-domain.com:8443
  X-Proxy-Scheme: https
  Proxy-Authorization: your_password
  ```


  代理服务器会根据这些字段将请求转发到指定的目标地址和端口。
### ⚠️ 注意事项
1. `Proxy-Authorization` 为访问 **WProxy 代理** 的 Basic 认证，与 `X-Proxy-Host` 指定的转发目标无关
2. 未配置 `username` / `password`（或二者任一为空）时，代理无鉴权，无需设置 `Proxy-Authorization`
3. 代理内部会自动设置 `X-Proxy-Loop` 防止环路，客户端一般无需手动添加

# 服务管理

安装后，您可以使用以下命令管理 WProxy 服务：

```bash
# 启动服务
sudo systemctl start wproxy

# 停止服务
sudo systemctl stop wproxy

# 重启服务
sudo systemctl restart wproxy

# 查看服务状态
sudo systemctl status wproxy

# 设置开机自启
sudo systemctl enable wproxy

# 禁用开机自启
sudo systemctl disable wproxy
```

# 安全最佳实践

在生产环境中使用 WProxy 时，请遵循以下安全建议：

## 认证安全
- 始终使用强随机密码（安装程序会自动生成）
- 将默认用户名从 "admin" 更改为唯一的用户名
- 定期更换密码
- 切勿通过不安全的渠道共享凭据

## 网络安全
- 使用防火墙规则限制对代理端口的访问
- 考虑为 HTTPS 代理连接使用 TLS/SSL 证书
- 在可能的情况下，将代理放在具有额外安全功能的反向代理后面
- 监控访问日志以发现可疑活动

## 配置安全
- 保持配置文件权限受限（Linux 安装路径建议 `chmod 600 /etc/wproxy/config.yaml`）
- 安全存储证书并设置适当的文件权限
- 定期更新到最新版本以获取安全补丁

## 生产部署
- 使用最小必需权限运行服务
- 启用系统日志并监控错误
- 设置失败时自动重启（systemd 服务中已配置）
- 对于面向公众的部署，使用连接限制和速率限制

# 贡献

如果您发现任何问题或有任何改进建议，欢迎提交 issue 或 pull request。我们很高兴能与社区一起改进 WProxy。

# 许可证

WProxy 基于 MIT 许可证发布，您可以自由使用、修改和分发本项目。

## Star 历史

[![Star History Chart](https://api.star-history.com/svg?repos=Wenpiner/WProxy&type=Date)](https://star-history.com/#Wenpiner/WProxy&Date) 
`