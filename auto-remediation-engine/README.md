# Auto-Remediation Engine - Phase 2

A Go-based Kubernetes controller that automatically detects cluster issues and remediates them intelligently. Safe fixes execute immediately, risky ones need Teams approval. Devops-DSU channel most likely for now.

## What Does It Do?

Think of it as a **24/7 on-call engineer** that watches your EKS cluster and fixes common problems before they become incidents.

### Problems It Solves

| Issue | Detection | Action | Approval |
|-------|-----------|--------|----------|
| **Pod Crash** | Pod exits repeatedly | Restart pod | Auto (immediate) |
| **OOMKilled** | Pod memory exceeded | Increase memory request + update Git | Auto (immediate) |
| **Missing Security** | Pod runs as root | Add `runAsNonRoot: true` + commit to Git | Auto (immediate) |
| **Cert Expiring** | Certificate <7 days to expiry | Alert Teams immediately | Manual review |
| **Liveness Probe Missing** | Pod has no liveness probe | Alert Teams with template | Manual review |
| **Readiness Probe Missing** | Pod has no readiness probe | Alert Teams with template | Manual review |

## Real-World Impact

**Before Auto-Remediation:**
- Pod crashes → Pager goes off at 3am
- On-call engineer wakes up
- 10 minutes to notice the problem
- 5 minutes to SSH and restart
- **MTTR: 15 minutes**

**After Auto-Remediation:**
- Pod crashes → Auto-restarted in seconds
- Teams alert sent (FYI only)
- On-call engineer sees it in morning standup
- **MTTR: <1 minute**

**Cost Impact:**
- Auto-fixed OOMKilled pods before they spiral
- Security mutations prevent compliance issues
- Proactive alerting on cert expiration prevents outages

## Architecture

```
┌─────────────────────────────────────────────┐
│         Kubernetes Cluster (EKS)            │
├─────────────────────────────────────────────┤
│                                             │
│  ┌──────────────────────────────────────┐   │
│  │  Auto-Remediation Controller Pod     │   │
│  │  (This project)                      │   │
│  └──────────────────────────────────────┘   │
│         ↓              ↓              ↓      │
│    [Detector]    [Remediator]   [Alerter]   │
│         ↓              ↓              ↓      │
│   Watch pods     Execute fixes   Send Teams │
│   Check certs    Mutate specs    messages   │
│   Monitor quotas Update Git               │
│                                             │
└─────────────────────────────────────────────┘
          ↓
    [Git Repo]
  (Mutations committed)
```

## Project Structure

```
auto-remediation-engine/
├── cmd/
│   └── controller/
│       └── main.go                 # Entry point
├── pkg/
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── detector/
│   │   ├── detector.go            # Detects issues in cluster
│   │   ├── crash_detector.go      # Detects pod crashes
│   │   ├── memory_detector.go     # Detects OOMKilled
│   │   ├── security_detector.go   # Detects security issues
│   │   └── cert_detector.go       # Detects expiring certs
│   ├── remediators/
│   │   ├── remediator.go          # Base remediator interface
│   │   ├── pod_restart.go         # Restarts crashed pods
│   │   ├── memory_increment.go    # Increases memory requests
│   │   ├── security_mutator.go    # Adds security settings
│   │   └── git_committer.go       # Commits changes to Git
│   ├── logger/
│   │   └── logger.go              # Structured logging
│   ├── teams/
│   │   └── teams.go               # Teams webhook integration
│   ├── models/
│   │   └── models.go              # Data structures
│   └── git/
│       └── git_client.go           # Git operations (commit/push)
├── k8s/
│   ├── crds/
│   │   └── remediationpolicy.yaml # CRD (future phase)
│   └── manifests/
│       ├── deployment.yaml        # Controller deployment
│       ├── rbac.yaml              # Permissions
│       ├── configmap.yaml         # Configuration
│       └── service-account.yaml   # Service account
├── config/
│   └── remediation-rules.yaml    # Rules to enable/disable
├── docs/
│   ├── ARCHITECTURE.md           # Deep dive
│   ├── REMEDIATION_RULES.md      # All available rules
│   └── DEPLOYMENT.md             # How to deploy
├── Dockerfile
├── go.mod
├── README.md
└── .gitignore
```

## Key Features

### 1. Detector Package
Watches cluster continuously for issues:
- Pod crash patterns (CrashLoopBackOff)
- Memory exhaustion (OOMKilled)
- Missing security context (runAsRoot)
- Certificate expiration (<7 days)
- Missing probes (liveness/readiness)
- Resource quota violations

### 2. Remediator Package
Executes fixes with intelligence:
- **Auto-execute** - Safe actions (restart, add probes, security)
- **Require approval** - Risky actions (cert renewal, quota changes)
- **Git integration** - All mutations committed with reason
- **Rollback** - Each change tracked and reversible

### 3. Teams Integration
Smart alerting:
- Immediate: Pod crashes, OOMKilled, missing security
- Approval needed: Cert expiration, quota changes
- Summary: Daily digest of all remediations
- Metrics: "Fixed X pods, saved Y minutes MTTR"

### 4. Git Integration
All mutations tracked:
```
commit a1b2c3d
Author: auto-remediation <auto@remediation.local>
Date:   Mon Oct 27 14:20:00 2025

    [AUTO] Fix OOMKilled pod: payment-api-xyz

    Pod: payment-api-xyz123
    Namespace: production
    Issue: OOMKilled after 2 restarts
    Action: Increased memory from 512Mi → 1024Mi
    
    Reasoning: Pod consistently hitting memory limit
    Remediation Policy: payment-api-oom-fix-v1
```

## Remediation Rules

### Rule 1: Pod Crash Remediation
**Trigger:** Pod in CrashLoopBackOff
**Action:** Restart pod
**Approval:** Auto (immediate)
**Alert:**
```
✅ AUTO-FIXED: Pod Restarted
Service: payment-api
Pod: payment-api-xyz123
Reason: CrashLoopBackOff detected
Action: Pod restarted
Restarts in last hour: 5
Next check: 5 minutes
```

### Rule 2: OOMKilled Remediation
**Trigger:** Pod status shows OOMKilled
**Action:** 
  1. Increase memory request by 50%
  2. Update deployment in Git
  3. Commit and push changes
**Approval:** Auto (immediate)
**Alert:**
```
✅ AUTO-FIXED: Memory Increased
Service: data-processor
Pod: data-processor-abc789
Memory: 512Mi → 768Mi
Quota: 768Mi / 2Gi
Update committed to Git: devops-repo/deployments/production
Next review: 2 hours
```

### Rule 3: Missing Security Context
**Trigger:** Pod runs with `runAsRoot: true` or missing `runAsNonRoot`
**Action:**
  1. Add security context: `runAsNonRoot: true`
  2. Update deployment spec
  3. Commit to Git with security audit trail
**Approval:** Auto (immediate)
**Alert:**
```
🔒 AUTO-FIXED: Security Hardened
Service: legacy-api
Pod: legacy-api-old123
Issue: Running as root
Action: Added runAsNonRoot: true
Commit: devops-repo/deployments/legacy-api
Compliance: Now aligned with security policy
```

### Rule 4: Certificate Expiring
**Trigger:** Certificate expires in <7 days
**Action:** Alert Teams (let cert-manager handle renewal)
**Approval:** Manual review
**Alert:**
```
⚠️  ALERT: Certificate Expiring Soon
Service: api-gateway
Certificate: api.example.com
Expires in: 5 days (Oct 31, 2025)
Status: cert-manager watching
Action: Please verify cert-manager is running
Alert Frequency: Daily until renewed
```

### Rule 5: Missing Liveness Probe
**Trigger:** Pod has no liveness probe defined
**Action:** Add default liveness probe, commit to Git
**Approval:** Auto (immediate)
**Alert:**
```
✅ AUTO-FIXED: Liveness Probe Added
Service: worker-job
Pod: worker-job-123
Probe Type: TCP on port 8080
Check Interval: 10s
Failure Threshold: 3
Commit: devops-repo/deployments/worker-job
```

### Rule 6: Missing Readiness Probe
**Trigger:** Pod has no readiness probe defined
**Action:** Add default readiness probe, commit to Git
**Approval:** Auto (immediate)
**Alert:**
```
✅ AUTO-FIXED: Readiness Probe Added
Service: web-app
Pod: web-app-456
Probe Type: HTTP /health on port 3000
Check Interval: 5s
Failure Threshold: 2
Commit: devops-repo/deployments/web-app
```

## Configuration

Environment variables (set in Kubernetes ConfigMap or deployment):

```bash
# Teams webhook for alerts
TEAMS_WEBHOOK_URL=https://outlook.webhook.office.com/webhookb2/...

# Git repository for commits
GIT_REPO_URL=git@bitbucket.org:stchome/devops.git
GIT_BRANCH=main
GIT_COMMIT_AUTHOR=auto-remediation@platform.local

# Remediation settings
ENABLE_POD_RESTART=true
ENABLE_MEMORY_INCREMENT=true
ENABLE_SECURITY_MUTATIONS=true
ENABLE_PROBE_INJECTION=true
ENABLE_CERT_ALERTS=true

# Memory increment percentage
MEMORY_INCREMENT_PERCENT=50

# Cert expiry warning threshold (days)
CERT_EXPIRY_THRESHOLD_DAYS=7

# Cluster info
CLUSTER_NAME=production
CLUSTER_REGION=us-east-1

# Logging
LOG_LEVEL=info
```

## Deployment to EKS

### Prerequisites
- EKS cluster running (Phase 1 deployment done)
- `kubectl` configured
- Git SSH keys in cluster
- Teams webhook URL ready

### Step 1: Build and Push Docker Image
```bash
docker build -t auto-remediation:v1.0.0 .
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com
docker tag auto-remediation:v1.0.0 <ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/auto-remediation:v1.0.0
docker push <ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/auto-remediation:v1.0.0
```

### Step 2: Deploy to EKS
```bash
# Create namespace
kubectl create namespace auto-remediation

# Create configmap with rules
kubectl apply -f k8s/manifests/configmap.yaml

# Create RBAC
kubectl apply -f k8s/manifests/rbac.yaml

# Deploy controller
kubectl apply -f k8s/manifests/deployment.yaml

# Verify
kubectl get pods -n auto-remediation
kubectl logs -n auto-remediation -l app=auto-remediation -f
```

## Remediation Flow

```
1. DETECT
   ├─ Watch pod events (create, update, delete)
   ├─ Check pod status (crashes, OOMKilled, etc)
   ├─ Check certificates (expiry dates)
   └─ Check security context (runAsRoot, etc)

2. ASSESS
   ├─ Severity level (critical, warning, info)
   ├─ Is this auto-fixable or needs approval?
   ├─ Have we already tried fixing this?
   └─ Cost/benefit of remediation

3. REMEDIATE
   ├─ If auto-fix:
   │  ├─ Execute fix (restart, mutate spec, etc)
   │  ├─ Commit to Git
   │  └─ Alert Teams (FYI)
   │
   └─ If needs approval:
      ├─ Send Teams message with action buttons
      ├─ Wait for approval
      └─ Execute when approved

4. ALERT
   ├─ Immediate: What was fixed
   ├─ Details: Pod name, resource, action taken
   ├─ Git link: Where changes were committed
   └─ Metrics: Impact (MTTR saved, etc)

5. TRACK
   ├─ Store remediation history
   ├─ Prevent infinite loops (don't fix the same thing repeatedly)
   └─ Generate reports for SRE team
```

## Success Metrics

Track these to measure impact:

1. **MTTR (Mean Time To Recovery)**
   - Before: 15-45 minutes
   - After: <1 minute
   - Target: 99% of issues auto-fixed

2. **Pod Restarts**
   - Manual restarts per week (target: 0)
   - Auto-remediation actions per week (target: >10)

3. **Resource Optimization**
   - OOMKilled incidents prevented
   - Memory waste reduced
   - Cost savings from better resource allocation

4. **Security Compliance**
   - Pods running as root (target: 0)
   - Missing security context (target: 0)
   - Compliance violations caught and fixed

5. **Certificate Issues**
   - Expiry surprises (target: 0)
   - Outages due to cert expiry (target: 0)

## Next Steps (Phase 3)

Once Phase 2 is solid, Phase 3 adds:
- **CRDs** - Define remediation policies as Kubernetes resources
- **Operators** - Full operator pattern with Kubebuilder
- **Multi-cluster** - Manage multiple EKS clusters
- **Advanced GitOps** - ArgoCD integration
- **Audit trails** - Full compliance tracking
- **ML predictions** - Predict and prevent issues before they happen

## Contributing

When adding new remediation rules:

1. Create detector in `pkg/detector/`
2. Create remediator in `pkg/remediators/`
3. Add tests
4. Update this README
5. Create Git commit with `[AUTO]` prefix

## Support

For issues or questions:
- Check `docs/REMEDIATION_RULES.md` for detailed rule docs
- See `docs/ARCHITECTURE.md` for design decisions
- Review logs: `kubectl logs -n auto-remediation -f`

---

**Status:** Phase 2 Development 🚀
**Target Release:** 2 weeks
**Platform:** EKS on AWS
**Language:** Go 1.21+