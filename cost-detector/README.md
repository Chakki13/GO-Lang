# Cost Detector

A Go-based Kubernetes watcher that tracks real-time costs in your EKS cluster and sends Teams alerts to TEAMS Channel.
For testing purpose it would be to "Devops DSU" channel.

## What does it do?

Watches your cluster 24/7 and tells you RIGHT NOW how much money is being burned. Instead of waiting 24 hours to find out in AWS, you get instant alerts on Teams.

## What you get

- Real-time cost tracking per pod, service, team
- Teams alerts: "Team X is burning $50/hr RIGHT NOW"
- Cost breakdown by environment (dev, staging, prod)
- Catch runaway costs in minutes instead of days

## Example

Pod spins up with 16 CPU and 64GB RAM by mistake? You get a Teams message in seconds saying "Yo, this pod is costing $5/minute." You fix it immediately instead of discovering the $10k bill later.

## Project structure

- `cmd/` - The main app that runs
- `pkg/watcher/` - Watches pod creation/deletion
- `pkg/calculator/` - Does the cost math
- `pkg/alerts/` - Sends Teams messages
- `pkg/logger/` - Logging stuff
- `pkg/teams/` - Teams API integration
- `pkg/config/` - Configuration
- `pkg/models/` - Data structures
- `k8s/` - Kubernetes deployment files
- `config/` - Config files
- `docs/` - Guides and notes

## Next steps

1. Set up Go modules
2. Watch pods in the cluster
3. Calculate cost per second
4. Send Teams alerts
5. Deploy to EKS
6. Demo it working

Let's go! 🚀
