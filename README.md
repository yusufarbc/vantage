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

## 🌐 Platform Architecture (Phishing & Scanning Workflows)

Vantage operates from a **VPS with a Static IP bound to a domain via DNS records** (essential for phishing oltalama delivery). It supports three main workflows:
1. **Phishing Simulations (VPS-hosted)**: Orchestrated on the VPS, sending oltalama emails to public targets via the local Postfix SMTP container and receiving landing page interactions via Caddy.
2. **External Scanning (Direct from VPS)**: ProjectDiscovery scanner tools scan public internet assets directly from the VPS.
3. **Internal Scanning (Via Reverse L3 Tunnel Agent)**: Scan packets are routed from the VPS through a secure Chisel tunnel (`tun0` interface) to an agent running inside the corporate network to scan internal assets.

```mermaid
graph TD
    subgraph Cloud / VPS - Static IP bound to DNS Domain
        Caddy["Caddy Reverse Proxy - HTTPS / Port 443"]
        Core["Vantage Core - Docker Container"]
        
        subgraph Vantage Engines
            PhishEngine["Phishing Simulation Engine - Gophish"]
            ScanEngine["ProjectDiscovery Scanner Engine"]
        end

        Postfix["Postfix SMTP Container"]
        ChiselServer["Chisel Server - Port 9090"]
    end

    subgraph Public Internet - Phishing & External Scan Targets
        Victims["Target Recipients and Phishing Victims"]
        PublicAssets["External Assets - Domains, Public IPs, Web Apps"]
    end

    subgraph Internal Corporate Network - Scanning Targets
        Agent["Vantage Chisel Agent - Endpoint"]
        InternalAssets["Internal Corporate Assets - Active Directory, Local Servers"]
    end

    %% Infrastructure Routing
    Caddy -->|Proxy Admin and Phishing Landing Pages| Core
    Core --- PhishEngine
    Core --- ScanEngine
    
    %% Phishing Workflow
    PhishEngine -->|1a. Triggers email delivery| Postfix
    Postfix -->|1b. Sends phishing emails via VPS Static IP| Victims
    Victims -->|1c. Accesses landing pages and inputs credentials via Domain DNS A Record| Caddy

    %% External Scanning Workflow
    ScanEngine -->|2a. Direct external scan requests| PublicAssets
    PublicAssets -->|2b. Returns scan responses directly to VPS| ScanEngine

    %% Internal Scanning Workflow - Reverse Tunnel
    Agent -->|3a. Initiates outbound Chisel connection| Caddy
    Caddy -->|3b. Proxy tunnel connection| ChiselServer
    ChiselServer <-->|3c. Establishes Secure L3 Tunnel tun0| Agent
    
    ScanEngine -->|3d. Directs internal scan to local subnet| ChiselServer
    ChiselServer -->|3e. Routes packets via tun0| Agent
    Agent -->|3f. Performs scan requests locally| InternalAssets
    InternalAssets -->|3g. Returns responses| Agent
    Agent -->|3h. Forwards results back through tunnel| ScanEngine
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
