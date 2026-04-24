# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected changes.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start the daemon with a baseline scan of your current open ports:

```bash
portwatch start
```

Run a one-time scan and print all open ports:

```bash
portwatch scan
```

Watch for changes and alert when a new or closed port is detected:

```bash
portwatch watch --interval 30s --alert email
```

Example output when an unexpected port opens:

```
[ALERT] 2024-01-15 14:32:01 — New port detected: TCP 8080 (PID 3821, nginx)
[INFO]  Baseline ports: 22, 443, 80
[INFO]  Current ports:  22, 443, 80, 8080
```

Configuration is stored in `~/.portwatch/config.yaml`. You can define an allowlist of expected ports to suppress false positives.

```yaml
allowlist:
  - 22
  - 80
  - 443
interval: 30s
```

## Requirements

- Go 1.21+
- Linux or macOS
- Root or `CAP_NET_ADMIN` privileges for full port visibility

## License

MIT © 2024 [yourusername](https://github.com/yourusername)