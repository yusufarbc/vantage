# Vantage (Gophish Security Operations Hub)

**Vantage** is an advanced fork of **Gophish**, transformed into a unified Security Operations Hub. It acts as a powerful, asynchronous Go-based CLI wrapper and orchestrator for **ProjectDiscovery** tools, allowing security teams to manage phishing simulations alongside automated network scanning.

---

## 🖼️ Screenshots

### Unified Security Ops Dashboard
![Vantage Dashboard](vantage-dashboard.JPG)

### Enterprise Docker Stack & Reverse Tunneling
![Vantage Docker Architecture](vantage-docker.JPG)

---

## 🚀 Key Features

*   **Unified Dashboard**: Manage phishing campaigns and security scans from a single modern UI.
*   **ProjectDiscovery CLI Wrapper**: Fully integrates and orchestrates ProjectDiscovery tools (**Nuclei**, **Subfinder**, **HTTPx**, **Naabu**, **DNSx**, **Katana**, **TLSx**, **ASNMap**, **Uncover**) as an asynchronous scan engine.
*   **Reverse L3 Tunneling**: Perform internal network scans through a secure Chisel-based reverse tunnel.
*   **Real-time Insights Discovery**: WebSocket-based live scan logs and campaign performance tracking.
*   **Enterprise-Ready Deployment**: Orchestrated via Docker Compose with Caddy (HTTPS) and Postfix integration.

---

## 📚 Documentation

For detailed setup and usage instructions, please refer to the following guides:

*   📖 **[Vantage Overview & API Reference](doc/README_VANTAGE.md)** - Main project documentation.
*   🚀 **[Deployment & Operations Guide](doc/DEPLOYMENT_GUIDE.md)** - Step-by-step VPS/Server deployment.
*   🌐 **[Reverse L3 Tunnel Guide](doc/REVERSE_TUNNEL_GUIDE.md)** - **[NEW]** Setup and usage for internal scanning.

## 🌐 Platform Architecture (Phishing Simulation & Scanning Engine)

Vantage merges the **Gophish oltalama (phishing) simulation workflow** with an **asynchronous ProjectDiscovery scanning engine** that routes internal traffic through a secure Chisel-based reverse L3 tunnel. Below is the comprehensive platform architecture and network flow:

```mermaid
graph TD
    subgraph Cloud / VPS
        Caddy["Caddy Reverse Proxy - HTTPS / Port 443"]
        Core["Vantage Core - Docker Container"]
        
        subgraph Vantage Engines
            PhishEngine["Phishing Simulation Engine - Gophish"]
            ScanEngine["ProjectDiscovery Scanner Engine - Nuclei, Naabu, etc."]
        end

        Postfix["Postfix SMTP Container"]
        ChiselServer["Chisel Server - Port 9090"]
    end

    subgraph Public Internet - Phishing Targets
        Victims["Target Recipients and Phishing Victims"]
    end

    subgraph Internal Corporate Network - Scanning Targets
        Agent["Vantage Chisel Agent - Endpoint"]
        InternalAssets["Internal Corporate Assets - Active Directory, Servers"]
    end

    %% Infrastructure Routing
    Caddy -->|Proxy Admin and Phishing Landing Pages| Core
    Core --- PhishEngine
    Core --- ScanEngine
    
    %% Phishing Workflow
    PhishEngine -->|1. Triggers email delivery| Postfix
    Postfix -->|2. Sends phishing emails| Victims
    Victims -->|3. Accesses oltalama pages / inputs credentials| Caddy

    %% Scanning Workflow - Reverse Tunnel
    Agent -->|4. Initiates outbound Chisel connection| Caddy
    Caddy -->|5. Proxy tunnel connection| ChiselServer
    ChiselServer ---|6. Establishes Secure L3 Tunnel tun0| Agent
    
    ScanEngine -->|7. Directs scan to internal subnet| ChiselServer
    ChiselServer -->|8. Routes packets via tun0| Agent
    Agent -->|9. Performs local scan requests| InternalAssets
    InternalAssets -->|10. Returns responses| Agent
    Agent -->|11. Returns scan results back| ScanEngine
```

---

## 🏗️ Quick Start

```bash
git clone https://github.com/your-org/gophish-vantage.git
cd gophish-vantage
cp .env.example .env
# Edit .env and start
docker-compose up -d
```

Access your dashboard at `https://yourdomain.com/` (as configured in `.env`).

---

## ⚖️ License

This project extends **Gophish** (MIT License) and integrates **ProjectDiscovery** tools.
See [LICENSE](./LICENSE) for full details.

---

**Built with ❤️ for offensive security teams.**
