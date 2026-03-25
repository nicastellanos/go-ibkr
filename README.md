# E-mini Delta Neutrality Bot

A Go-based automated monitoring tool for tracking and rebalancing portfolio delta across **ES** (S&P 500) and **RTY** (Russell 2000) futures via the Interactive Brokers API.

## 🚀 Quick Start

### 1. Prerequisites
* **Go 1.21+**
* **IB Gateway** or **TWS** installed.
* **Paper Trading Account** (Recommended for testing).

### 2. IBKR Configuration
To allow the bot to connect, configure the Gateway as follows:
1. Open **Settings > API > Settings**.
2. Check **"Enable ActiveX and Socket Clients"**.
3. Uncheck **"Read-Only API"**.
4. Set **Socket Port** to `7497` (TWS Paper) or `4002` (Gateway Paper).
5. Add `127.0.0.1` to **Trusted IPs**.

### 3. Installation & Run
```bash
# Clone the repository
git clone [https://github.com/your-username/go-ibkr-collector.git](https://github.com/your-username/go-ibkr-collector.git)
cd go-ibkr-collector

# Install dependencies
go mod tidy

# Run the bot
go run cmd/bot/main.go
