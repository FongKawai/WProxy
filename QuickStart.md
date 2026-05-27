# WProxy 使用说明

同一端口支持 **SOCKS5** 与 **HTTP/HTTPS**，自动识别协议。

## 启动

```bash
wproxy -c config.yaml
# 或
wproxy -host 0.0.0.0 -port 1080 -username admin -password 密钥
```

`config.yaml` 示例：

```yaml
listen_addr: "0.0.0.0:1080"
username: "admin"
password: "your_password"
```

`username`、`password` 均非空才启用鉴权；改配置后需重启进程。

## 客户端

| 项 | 说明 |
|----|------|
| 类型 | SOCKS5 或 HTTP |
| 地址 / 端口 | 与 `listen_addr` 一致（默认 `1080`） |
| 账号 | 与配置中一致；未启用鉴权则留空 |

## HTTP 转发（可选）

在请求中设置目标，而非 URL 中的 Host：

| 请求头 | 说明 |
|--------|------|
| `X-Proxy-Host` | 目标主机，可加端口 |
| `X-Proxy-Scheme` | `http` 或 `https` |
| `Proxy-Authorization` | 代理已启用鉴权时填写（Basic） |

```http
GET /path HTTP/1.1
Host: proxy.example.com
X-Proxy-Host: target.example.com
X-Proxy-Scheme: http
```
