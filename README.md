<div align="center">
  <h1>🌩️ awsx</h1>
  <p><em>A Free, Open-Source AWS Reachability Analyzer Built with Golang</em></p>
  <img src="https://cdn.worldvectorlogo.com/logos/aws-2.svg" width="80" alt="AWS Logo"/>
</div>

---

<div align="center">
  <a href="https://github.com/itzrahulyadav/awsx/stargazers">
    <img src="https://img.shields.io/github/stars/itzrahulyadav/awsx?style=social" alt="Stars"/>
  </a>
  <a href="https://github.com/itzrahulyadav/awsx/issues">
    <img src="https://img.shields.io/github/issues/itzrahulyadav/awsx?style=social" alt="Issues"/>
  </a>
  <a href="https://github.com/itzrahulyadav/awsx/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/itzrahulyadav/awsx?style=social" alt="License"/>
  </a>
</div>

---

### 🚀 What is awsx?
`awsx` is a lightweight, free alternative to the [AWS Reachability Analyzer](https://docs.aws.amazon.com/vpc/latest/reachability/what-is-reachability-analyzer.html), crafted in **Golang** to help you analyze network connectivity within your AWS infrastructure. Tired of AWS pricing tiers or just want a CLI-driven, open-source solution? `awsx` has got you covered—simple, fast, and community-driven!

- **Purpose:** Diagnose network reachability issues between AWS resources (e.g., EC2 instances, VPCs, subnets).
- **Why Golang?** Blazing-fast performance and easy cross-platform builds.
- **Cost:** 100% free—no hidden fees, no AWS bills!

Think of it as your trusty sidekick for untangling VPC mysteries—without the corporate price tag.

---

### ✨ Features
- ✅ Check connectivity between AWS resources (e.g., EC2 to RDS, VPC peering).
- ✅ Detailed path analysis with hop-by-hop insights.
- ✅ CLI-first design for DevOps workflows.
- ✅ Cross-platform support (Linux, macOS, Windows).
- ✅ Lightweight—no heavy dependencies, just pure Go goodness.

*(Note: Add more features as you develop them!)*

---

### 🛠️ Getting Started
#### Prerequisites
- Go 1.18+ installed ([Download here](https://go.dev/dl/)).
- AWS CLI configured with credentials (`aws configure`).
- An AWS account with VPC resources to analyze.

#### Installation
1. Clone the repo:
   ```bash
   git clone https://github.com/itzrahulyadav/awsx.git
   cd awsx
