# driftwatch

A CLI tool that detects configuration drift between live infrastructure and declared state files.

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git
cd driftwatch
go build -o driftwatch .
```

---

## Usage

```bash
# Compare live infrastructure against a local state file
driftwatch check --state ./infra/state.yaml --provider aws

# Output drift report in JSON format
driftwatch check --state ./infra/state.yaml --provider aws --output json

# Watch for drift on an interval (every 5 minutes)
driftwatch watch --state ./infra/state.yaml --interval 5m
```

**Example output:**

```
[DRIFT DETECTED] ec2/instance i-0abc123def456
  expected: t3.micro
  actual:   t3.small

[OK] s3/bucket my-app-bucket
[OK] rds/instance prod-db
```

---

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--state` | Path to declared state file | `./state.yaml` |
| `--provider` | Cloud provider (`aws`, `gcp`, `azure`) | `aws` |
| `--output` | Output format (`text`, `json`) | `text` |
| `--interval` | Polling interval for watch mode | `10m` |

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)