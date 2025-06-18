<div align="center">

# 🛡️ Hardshell

**Automated hardening tool for Linux systems**

[![Go Report Card](https://goreportcard.com/badge/github.com/mairinkdev/Hardshell)](https://goreportcard.com/report/github.com/mairinkdev/Hardshell)
[![Go Version](https://img.shields.io/github/go-mod/go-version/mairinkdev/Hardshell)](https://github.com/mairinkdev/Hardshell)
[![License](https://img.shields.io/github/license/mairinkdev/Hardshell)](https://github.com/mairinkdev/Hardshell/blob/main/LICENSE)

<img src="https://raw.githubusercontent.com/mairinkdev/Hardshell/logo.png" alt="Hardshell Logo" width="250">

</div>

## 📋 About

Hardshell is a Go CLI tool for automated Linux system hardening, aimed at DevOps, SREs, and security professionals. It analyzes security configurations, identifies issues, and suggests corrections based on best practices.

### ✅ Features

- **Automated analysis of:**
  - SSH configuration (e.g., PermitRootLogin, Protocol, etc.)
  - Security-sensitive sysctl settings
  - Dangerous or insecure services active in the system

- **Report generation:**
  - Output in text, JSON, or HTML
  - Classification of issues as CRITICAL, WARNING, and INFO

- **Automatic fixes:**
  - Generation of shell script with suggestions
  - `--apply` flag to execute corrections (with automatic backup)

- **Container-aware mode:**
  - Capable of analyzing rootfs mounted in a specific directory

## 🚀 Installation

### Via Go

```bash
go install github.com/mairinkdev/Hardshell/cmd/hardshell@latest
```

### Download binary

Download the latest version from [GitHub Releases](https://github.com/mairinkdev/Hardshell/releases)

```bash
wget https://github.com/mairinkdev/Hardshell/releases/latest/download/hardshell_linux_amd64.tar.gz
tar -xzf hardshell_linux_amd64.tar.gz
sudo mv hardshell /usr/local/bin/
```

### Via Docker

```bash
docker pull mairinkdev/hardshell:latest
docker run --rm -v /:/data mairinkdev/hardshell scan
```

## 💻 Usage

### Main command

```bash
# Scan the complete system
hardshell scan

# Scan only SSH configurations
hardshell ssh

# Scan only sysctl configurations
hardshell sysctl

# Scan only services
hardshell services
```

### Options

```bash
# Apply corrections automatically (with backup)
hardshell scan --apply

# Generate report in JSON format
hardshell scan --output json > report.json

# Generate report in HTML format
hardshell scan --output html > report.html

# Analyze a system mounted at /mnt
hardshell scan --mount /mnt

# Use a custom configuration file
hardshell scan --config /path/to/config.yaml
```

## 🔧 Configuration

Hardening rules can be customized through a YAML file:

```yaml
# Example SSH rule
ssh:
  - key: "PermitRootLogin"
    recommended_value: "no"
    severity: "CRITICAL"
    description: "Root direct login should be disabled"

# Example sysctl rule
sysctl:
  - key: "net.ipv4.tcp_syncookies"
    recommended_value: "1"
    severity: "CRITICAL"
    description: "SYN flood protection should be enabled"
```

## 📋 Example output

### Text

```
=== HARDSHELL SECURITY REPORT ===

== SSH ==
1. [CRITICAL] Root direct login should be disabled
   Current value: yes
   Recommended value: no
   Fix: sed -i 's/^PermitRootLogin.*/PermitRootLogin no/' /etc/ssh/sshd_config

2. [WARNING] Password authentication should be disabled, prefer SSH keys
   Current value: yes
   Recommended value: no
   Fix: sed -i 's/^PasswordAuthentication.*/PasswordAuthentication no/' /etc/ssh/sshd_config

== SYSCTL ==
1. [CRITICAL] SYN flood protection should be enabled
   Current value: 0
   Recommended value: 1
   Fix: echo 'net.ipv4.tcp_syncookies = 1' >> /etc/sysctl.conf && sysctl -p

Scan summary:
  Critical issues: 2
  Warnings: 1
  Information: 0
  Total: 3
```

### HTML

<img src="https://raw.githubusercontent.com/mairinkdev/Hardshell/main/docs/report-example.png" alt="Example HTML report" width="600">

## 🛡️ Security checks

- **SSH:**
  - PermitRootLogin, Protocol, PasswordAuthentication, etc.

- **Sysctl:**
  - tcp_syncookies, accept_redirects, accept_source_route, etc.
  - randomize_va_space, protected_hardlinks, protected_symlinks, etc.

- **Services:**
  - telnet, rsh, rlogin, ftp, tftp, etc.

## 🧪 Tests

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit PRs, report bugs, or suggest new features.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

# 🛡️ Hardshell

**Linux 系统自动化加固工具**

## 📋 关于

Hardshell 是一个 Go 语言开发的命令行工具，用于 Linux 系统的自动化加固，面向 DevOps、SRE 和安全专业人员。它能分析安全配置，识别问题，并根据最佳实践提出修正建议。

### ✅ 功能特点

- **自动分析：**
  - SSH 配置（如 PermitRootLogin、Protocol 等）
  - 安全敏感的 sysctl 设置
  - 系统中活跃的危险或不安全服务

- **报告生成：**
  - 支持文本、JSON 或 HTML 输出
  - 将问题分类为严重（CRITICAL）、警告（WARNING）和信息（INFO）

- **自动修复：**
  - 生成带有建议的 shell 脚本
  - 使用 `--apply` 标志执行修复（自动备份）

- **容器感知模式：**
  - 能够分析挂载在特定目录中的 rootfs

## 🚀 安装

### 通过 Go 安装

```bash
go install github.com/mairinkdev/Hardshell/cmd/hardshell@latest
```

### 下载二进制文件

从 [GitHub Releases](https://github.com/mairinkdev/Hardshell/releases) 下载最新版本

```bash
wget https://github.com/mairinkdev/Hardshell/releases/latest/download/hardshell_linux_amd64.tar.gz
tar -xzf hardshell_linux_amd64.tar.gz
sudo mv hardshell /usr/local/bin/
```

### 通过 Docker 使用

```bash
docker pull mairinkdev/hardshell:latest
docker run --rm -v /:/data mairinkdev/hardshell scan
```

## 💻 使用方法

### 主要命令

```bash
# 扫描整个系统
hardshell scan

# 仅扫描 SSH 配置
hardshell ssh

# 仅扫描 sysctl 配置
hardshell sysctl

# 仅扫描服务
hardshell services
```

### 选项

```bash
# 自动应用修复（带备份）
hardshell scan --apply

# 以 JSON 格式生成报告
hardshell scan --output json > report.json

# 以 HTML 格式生成报告
hardshell scan --output html > report.html

# 分析挂载在 /mnt 的系统
hardshell scan --mount /mnt

# 使用自定义配置文件
hardshell scan --config /path/to/config.yaml
```

## 🔧 配置

加固规则可以通过 YAML 文件自定义：

```yaml
# SSH 规则示例
ssh:
  - key: "PermitRootLogin"
    recommended_value: "no"
    severity: "CRITICAL"
    description: "应禁用 root 直接登录"

# sysctl 规则示例
sysctl:
  - key: "net.ipv4.tcp_syncookies"
    recommended_value: "1"
    severity: "CRITICAL"
    description: "应启用 SYN flood 保护"
```

## 📋 输出示例

### 文本输出

```
=== HARDSHELL 安全报告 ===

== SSH ==
1. [严重] 应禁用 root 直接登录
   当前值: yes
   推荐值: no
   修复: sed -i 's/^PermitRootLogin.*/PermitRootLogin no/' /etc/ssh/sshd_config

2. [警告] 应禁用密码认证，首选 SSH 密钥
   当前值: yes
   推荐值: no
   修复: sed -i 's/^PasswordAuthentication.*/PasswordAuthentication no/' /etc/ssh/sshd_config

== SYSCTL ==
1. [严重] 应启用 SYN flood 保护
   当前值: 0
   推荐值: 1
   修复: echo 'net.ipv4.tcp_syncookies = 1' >> /etc/sysctl.conf && sysctl -p

扫描摘要:
  严重问题: 2
  警告: 1
  信息: 0
  总计: 3
```

## 🛡️ 安全检查

- **SSH:**
  - PermitRootLogin, Protocol, PasswordAuthentication 等

- **Sysctl:**
  - tcp_syncookies, accept_redirects, accept_source_route 等
  - randomize_va_space, protected_hardlinks, protected_symlinks 等

- **服务:**
  - telnet, rsh, rlogin, ftp, tftp 等

## 🧪 测试

```bash
# 运行测试
go test ./...

# 运行带覆盖率的测试
go test -cover ./...
```

## 🤝 贡献

欢迎贡献！请随时提交 PR，报告 bug 或建议新功能。

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m '添加某个惊人功能'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 有关详细信息，请参阅 [LICENSE](LICENSE) 文件。
