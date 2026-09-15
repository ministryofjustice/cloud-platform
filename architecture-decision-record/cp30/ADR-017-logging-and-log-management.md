# ADR-017: Logging and Log Management Architecture

**Status:** Proposed — Phase 1 PoC validation required  
**Date:** 2026-07-22  
**Decision Maker:** AWS ProServe / MoJ Principal Technical Architect  
**Category:** Observability

---

## Context

**Problem Statement:** The current Cloud Platform logging architecture relies on Fluent Bit shipping application logs to a self-managed OpenSearch cluster (~$54k/month) that serves as both full-text search and a "log dump." Additionally, logs are forwarded to Palo Alto Cortex XSIAM for SOC security monitoring via S3 → SQS and CloudWatch → Firehose patterns. The new multi-cluster, multi-account platform (CP3) requires a logging architecture that:

- Supports multiple log consumers (operations, search, security/SOC) without duplicating log data unnecessarily
- Reduces the current $54k/month OpenSearch cost through right-sizing, tiered storage, and lifecycle management
- Provides structured logging standards to improve queryability
- Maintains the existing Cortex XSIAM integration for SOC
- Scales across 22 clusters and 8-9 BUs with namespace-scoped access control
- Operates within UK data residency requirements (eu-west-2)

**Requirements:** FR-204 (log aggregation and full-text search), FR-206 (observability guardrails — log volume), FR-207 (platform component observability), FR-208 (structured logging standards), NFR-004 (UK data residency), NFR-101 (managed services), NFR-106 (cost impact analysis)

**Constraints:**
- UK data residency required — all log storage in eu-west-2 (NFR-004)
- No Lambda (NFR-003)
- SOC integration via Palo Alto Cortex XSIAM must be preserved (existing pattern)
- EKS Auto Mode — platform-managed components emit logs via CloudWatch Vended Logs
- Current costs: OpenSearch ~$54k/month, self-managed Grafana ~£7k/month
- 20 May meeting: full-text search is important, cannot dispense with OpenSearch; avoid moving logs across accounts for ingestion (cost concern)
- 07 May meeting: logs should go to a central store first, then optionally to per-BU search

---

## Decision

**A multi-destination Fluent Bit pipeline with S3 as the durable log archive, CloudWatch Logs for operational querying, and OpenSearch for full-text search.** SOC forwarding to Cortex XSIAM is maintained via the existing S3 → SQS pull pattern.

The architecture separates log concerns by consumer:

| Consumer | Purpose | Destination | Retention |
|----------|---------|-------------|-----------|
| Operations (platform team) | Real-time troubleshooting, recent log tailing | CloudWatch Logs | 30 days (configurable) |
| Full-text search (all engineers) | Historical investigation, pattern discovery | OpenSearch | Hot: 14 days, Warm: 90 days |
| Security/SOC | Threat detection, incident response | Cortex XSIAM (via S3 → SQS) | Per SOC policy |
| Archive/compliance | Long-term retention, audit | S3 (lifecycle to Glacier) | 1 year+ (per policy) |

### Pipeline Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│  EKS Cluster (per BU account)                                       │
│                                                                     │
│  ┌──────────┐    ┌───────────────────────────────────────────────┐  │
│  │ App Pods │───▶│ Fluent Bit DaemonSet                          │  │
│  │ (stdout) │    │                                               │  │
│  └──────────┘    │  ┌─────────┐  ┌──────────┐  ┌────────────┐   │  │
│                  │  │Output: │  │Output:  │  │Output:    │   │  │
│  ┌──────────┐    │  │   S3    │  │CloudWatch│  │OpenSearch │   │  │
│  │ Platform │───▶│  │(primary)│  │  Logs    │  │(optional) │   │  │
│  │Components│    │  └────┬────┘  └────┬─────┘  └─────┬─────┘   │  │
│  └──────────┘    └───────┼────────────┼───────────────┼─────────┘  │
│                          │            │               │             │
└──────────────────────────┼────────────┼───────────────┼─────────────┘
                           │            │               │
                           ▼            ▼               ▼
                    ┌────────────┐ ┌──────────┐  ┌────────────┐
                    │  S3 Bucket │ │CloudWatch│  │ OpenSearch  │
                    │ (per-BU)   │ │  Logs    │  │  Service   │
                    └──────┬─────┘ └──────────┘  └────────────┘
                           │
                    ┌──────┴──────┐
                    │ SQS Queue   │
                    │(S3 events)  │
                    └──────┬──────┘
                           │
                           ▼
                    ┌─────────────┐       ┌──────────────────┐
                    │Cortex XSIAM │       │ S3 Lifecycle     │
                    │  (SOC)      │       │ → Glacier (1yr+) │
                    └─────────────┘       └──────────────────┘


EKS Control Plane Logs:
    EKS ──▶ CloudWatch Logs ──▶ Firehose ──▶ Cortex XSIAM

EKS Auto Mode Component Logs (Karpenter, EBS CSI, LB Controller, VPC CNI):
    CloudWatch Vended Logs ──▶ CloudWatch Logs (+ optional S3 delivery)
```

### Log Classification

| Log Type | Sources | Destinations | Access |
|----------|---------|--------------|--------|
| Application logs | Pod stdout/stderr | S3, CloudWatch, OpenSearch | App team (own namespace) |
| Platform component logs | Argo CD, Gatekeeper, Gateway API, cert-manager, ExternalDNS | S3, CloudWatch, OpenSearch | Platform team |
| EKS control plane logs | API server, audit, authenticator, controller manager, scheduler | CloudWatch → Firehose → Cortex XSIAM | Platform team + SOC |
| EKS Auto Mode component logs | Karpenter, EBS CSI, LB controller, VPC CNI IPAM | CloudWatch (Vended Logs) | Platform team |
| Kubernetes audit logs | API server audit | CloudWatch → Firehose → Cortex XSIAM | Platform team + SOC |
| VPC Flow Logs | VPC ENIs | S3 → SQS → Cortex XSIAM | Platform team + SOC |
| CloudTrail logs | AWS API calls | S3 → SQS → Cortex XSIAM | Platform team + SOC |
| DNS query logs | Route 53 Resolver | S3 → SQS → Cortex XSIAM | Platform team + SOC |

---

## Options Evaluated

### Option A: Fluent Bit Multi-Output (S3 + CloudWatch + OpenSearch) — SELECTED

Fluent Bit ships logs directly to multiple destinations from its output plugins. S3 is the primary durable store (cheapest per-GB, feeds SOC via existing pattern). CloudWatch provides real-time operational queries. OpenSearch provides full-text search.

**Pros:**
- Single collection agent handles all routing — no intermediate hops for application logs
- S3 as primary archive is the cheapest durable storage (~$0.023/GB/month standard, ~$0.004/GB Glacier)
- Preserves existing Cortex XSIAM pattern (S3 → SQS) with no architectural change
- Fluent Bit back-pressure and retry per output — failure in one destination doesn't block others
- No cross-account log transfer for the primary archive (S3 in same account as cluster)
- CloudWatch retention kept short (30 days) to control cost — long-term in S3
- OpenSearch output can be optional per cluster (non-live clusters may not need it)

**Cons:**
- Fluent Bit manages three output plugins — more configuration complexity per cluster
- Fluent Bit S3 plugin buffers to local disk — requires adequate node storage
- If Fluent Bit pod restarts, buffered-but-undelivered logs may be lost (mitigated by filesystem buffer on host path)

### Option B: Fluent Bit → CloudWatch Only → Subscription Filters to Downstream

All logs go to CloudWatch Logs first, then subscription filters route to S3 (archive), OpenSearch (search), and Firehose (SOC).

**Pros:**
- Single output from Fluent Bit (simpler agent config)
- CloudWatch handles durability and fan-out
- Native CloudWatch features (Logs Insights, metric filters, alarms) available immediately

**Cons:**
- CloudWatch ingestion cost: $0.50/GB (eu-west-2) — at scale this dominates cost
- Subscription filter limitation: max 2 subscription filters per log group — constrains fan-out
- CloudWatch → Firehose → OpenSearch is not supported (documented AWS limitation)
- Subscription filter to OpenSearch requires Lambda (violates NFR-003) or direct subscription (complex cross-account)
- All logs pass through CloudWatch before reaching archive — adds latency and cost for the SOC/archive path
- Cross-account subscription filters add IAM complexity

### Option C: Fluent Bit → Firehose → S3/OpenSearch

Fluent Bit ships to Amazon Data Firehose, which handles buffering and delivery to S3 and OpenSearch.

**Pros:**
- Firehose provides managed buffering and retry (more reliable than Fluent Bit S3 plugin alone)
- Single Firehose stream can deliver to both S3 and OpenSearch simultaneously
- Firehose handles format transformation (e.g., converting to Parquet for S3)

**Cons:**
- Firehose cost: $0.029/GB ingestion + S3 storage — more expensive than direct S3 writes
- Firehose adds latency (60-second minimum buffer interval for S3, configurable)
- Per-stream throughput limits (quota management at scale)
- Requires Firehose delivery stream per cluster or per log group — operational overhead at 22 clusters
- Less flexible filtering — Fluent Bit's native filtering is richer than Firehose transformation

### Option D: Fluent Bit → S3 Only (with Athena for Search)

All logs go to S3 in a queryable format (Parquet/JSON). Use Athena or S3 Select for ad-hoc queries instead of OpenSearch.

**Pros:**
- Cheapest possible architecture — S3 storage only
- No OpenSearch operational cost ($54k/month saving)
- Athena queries are pay-per-query — no idle infrastructure cost

**Cons:**
- Athena query latency (seconds to minutes) — not suitable for real-time troubleshooting
- No live tail or streaming log view
- Engineers lose the familiar OpenSearch Dashboards/Kibana UX for log exploration
- 20 May meeting explicitly confirmed "full-text search is important, cannot dispense with OpenSearch"
- Poor experience for incident response where speed matters

---

## Selected Approach: Option A with Refinements

Option A is selected as the baseline architecture. Key refinements:

1. **S3 as the durable primary** — All logs land in S3 regardless of other destinations. This is the cheapest store, feeds the SOC, and provides the compliance archive.
2. **CloudWatch for operational use** — Short retention (30 days default), used for real-time tailing and Logs Insights queries during active troubleshooting.
3. **OpenSearch for full-text search** — Provides the rich search experience engineers need for historical investigation. Receives logs via Fluent Bit direct output (not via CloudWatch subscription filters, avoiding the Lambda/cost issue).
4. **SOC forwarding unchanged** — S3 bucket notifications → SQS → Cortex XSIAM pulls. Same pattern as today. EKS control plane logs continue via CloudWatch → Firehose → Cortex XSIAM.

### PoC Evaluation Items

The following must be validated during Phase 1 PoC:

| Item | Question | How to Measure |
|------|----------|----------------|
| Fluent Bit S3 reliability | Does the S3 output plugin handle back-pressure and node restarts without log loss? | Deploy on PoC cluster, simulate high-volume logging, kill Fluent Bit pods, verify no gaps in S3 |
| CloudWatch cost at scale | What is the projected CloudWatch ingestion cost for 22 clusters? | Measure log volume on PoC cluster, extrapolate |
| OpenSearch sizing | Can a right-sized OpenSearch domain serve full-text search at projected volume for <$20k/month? | Deploy OpenSearch, ingest PoC logs, measure query performance and cost |
| OpenSearch account placement | Centralised (one cluster) vs per-BU (multiple clusters) — which model works? | See evaluation criteria below |
| Fluent Bit parsing | Do the parsers handle the mix of structured/unstructured logs from existing MoJ apps? | Sample real application logs, test parser configurations |
| OpenSearch Serverless viability | Is Serverless cost-effective and performant for non-live or per-BU use? | Deploy Time series collection alongside provisioned domain, compare ingestion reliability, query latency, Dashboards UX, and projected cost |

---

## SOC Integration — Palo Alto Cortex XSIAM

### Current Architecture (to be preserved)

The existing CP2 platform forwards logs to the SOC team's Cortex XSIAM using the following patterns:

| Log Source | Delivery Path | Mechanism |
|------------|---------------|-----------|
| CloudTrail | S3 bucket → S3 event notification → SQS → Cortex XSIAM pulls | IAM user access keys |
| VPC Flow Logs | S3 bucket → S3 event notification → SQS → Cortex XSIAM pulls | IAM user access keys |
| Route 53 DNS query logs | S3 bucket → S3 event notification → SQS → Cortex XSIAM pulls | IAM user access keys |
| EKS control plane logs | CloudWatch Logs → Firehose delivery stream → Cortex XSIAM HTTP endpoint | Firehose IAM role |
| Application logs | Fluent Bit dual output → S3 bucket → S3 event notification → SQS → Cortex XSIAM pulls | IAM user access keys |

**Source:** [Cloud Platform Runbook — Logs to SOC Cortex XSIAM](https://runbooks.cloud-platform.service.justice.gov.uk/logs-to-soc-cortex-xsiam.html)

### CP3 Architecture (proposed)

The CP3 architecture preserves these patterns with the following adjustments:

| Change | CP2 (current) | CP3 (proposed) | Rationale |
|--------|---------------|----------------|-----------|
| Application log S3 delivery | Fluent Bit → S3 (secondary output alongside OpenSearch) | Fluent Bit → S3 (primary output) | S3 becomes the primary durable store, not a secondary copy |
| IAM authentication | IAM user with long-lived access keys | IAM role with cross-account assume (preferred) or IAM user with automated key rotation | Eliminate long-lived credentials where possible |
| Per-BU S3 buckets | Single S3 bucket (single cluster) | Per-BU S3 bucket (multi-account) | Each BU account has its own log bucket; Cortex XSIAM pulls from each |
| EKS control plane logs | CloudWatch → Firehose → Cortex | Same pattern, per-cluster Firehose delivery stream | One Firehose stream per cluster, same destination |
| EKS Auto Mode component logs | N/A (new in CP3) | CloudWatch Vended Logs → CloudWatch → Firehose → Cortex | New log source, follows EKS control plane pattern |

> **Scope note (PoC status):** SOC forwarding of EKS Auto Mode component logs via Firehose is **deferred** — it is not required at this stage and needs SOC input for scoping (see #8430). The `→ Firehose → Cortex` path for Auto Mode logs above is the intended design, not committed PoC scope. Auto Mode logs currently land in CloudWatch Logs (and optionally S3) only.

### Open Question: Cortex XSIAM Authentication

The current pattern uses IAM user access keys that require periodic rotation. The CP3 architecture should evaluate:
- Whether Cortex XSIAM supports IAM role-based cross-account access (preferred — no key rotation needed)
- If not, implement automated key rotation via Secrets Manager with rotation Lambda (exception to NFR-003 for security automation) or an alternative rotation mechanism
- Confirm with SOC team during Phase 1 engagement

---

## Structured Logging Standards (FR-208)

### Published Standard

The platform publishes a recommended structured logging format. Teams are encouraged (not mandated) to adopt it:

```json
{
  "timestamp": "2026-07-22T14:30:00.123Z",
  "level": "ERROR",
  "service": "payment-api",
  "message": "Failed to process payment",
  "trace_id": "abc123def456",
  "span_id": "789ghi",
  "error": {
    "type": "TimeoutException",
    "message": "Connection timed out after 30s",
    "stack": "..."
  },
  "context": {
    "user_id": "REDACTED",
    "transaction_id": "tx-98765"
  }
}
```

**Required fields** (platform-enriched if missing): `timestamp`, `level`, `service`  
**Recommended fields**: `trace_id`, `span_id`, `message`, `error`  
**Optional fields**: `context` (application-specific key-value pairs)

### Platform-Side Enrichment

Regardless of application log format, Fluent Bit adds Kubernetes metadata to every log record:

| Field | Source | Example |
|-------|--------|---------|
| `kubernetes.namespace` | Pod metadata | `payments-prod` |
| `kubernetes.pod_name` | Pod metadata | `payment-api-7f8d9c-x2k4` |
| `kubernetes.container_name` | Container spec | `payment-api` |
| `kubernetes.node_name` | Node assignment | `ip-10-195-32-45.eu-west-2.compute.internal` |
| `kubernetes.labels.*` | Pod labels | `app=payment-api, team=fintech` |
| `cluster_name` | Fluent Bit env var | `cloud-platform-development` |
| `account_id` | Fluent Bit env var | `123456789012` |

### Log Parsing Pipeline (Fluent Bit)

```
INPUT (tail) → PARSER (try JSON, logfmt, regex) → FILTER (kubernetes metadata)
    → FILTER (redaction rules) → OUTPUT (S3, CloudWatch, OpenSearch)
```

1. **Parser stage**: Attempts JSON parse first. If that fails, tries logfmt. If both fail, treats as unstructured text with the raw message in a `log` field.
2. **Kubernetes filter**: Enriches with pod/namespace/node/labels metadata.
3. **Redaction filter**: Applies regex-based redaction for known PII patterns (see PII Handling below).
4. **Output routing**: Tag-based routing sends logs to appropriate destinations.

### Adoption Tracking

The platform tracks the percentage of logs arriving in structured JSON format vs unstructured text. This metric is reported per namespace and per BU to encourage adoption without enforcement:
- Structured (JSON parsed successfully): target >80% within 6 months of Phase 2 launch
- Unstructured (raw text): still ingested, still searchable, but with reduced field-level query capability

---

## OpenSearch — Full-Text Log Search

### Observability Account — Single Dedicated Account (All Environments)

All observability query and visualisation tooling (AMG, OpenSearch Dashboards) resides in a **single dedicated Observability account** that serves both live and non-live environments. This aligns with ADR-005's account placement decision for AMG.

**Why a single account (not split live/non-live):**

| Factor | Single account | Split (non-live + live) |
|--------|----------------|------------------------|
| Incident correlation | Engineers compare live vs non-live logs in one place | Must switch between two OpenSearch instances |
| Data isolation | Index-level RBAC in OpenSearch; actual log data stays in BU accounts (S3/CloudWatch) | Account-boundary separation (redundant — data is already in source accounts) |
| Operational overhead | 1 OpenSearch domain (or Serverless collection), 1 set of index policies | 2 domains, duplicate ISM policies and access config |
| Access control | Identity Centre SSO, role-scoped to BU/namespace indices | Same, but duplicated across accounts |
| Security boundary | OpenSearch has read-only access to log data (no write to source accounts) | Same posture, but in two places |
| Cost | Single domain sized for total query load | Two domains with idle capacity in non-live |

**Key insight:** The Observability account doesn't store the primary copy of log data. S3 (per-BU account) is the durable archive. OpenSearch receives a copy for search purposes. If the OpenSearch domain is compromised, the attacker gains read access to log data — the same data already available to anyone with Identity Centre SSO access to the relevant BU scope. Splitting accounts doesn't add a meaningful security boundary.

**When splitting would be justified:**
- A compliance requirement mandating production observability tooling in a separate account from non-production (not currently the case)
- Different teams managing live vs non-live observability (unlikely — same platform team)

**What remains per-BU account (not in the Observability account):**
- S3 log buckets (primary archive, SOC feed)
- CloudWatch log groups (operational queries, short retention)
- AMP workspaces or CloudWatch metrics (per ADR-005)

### Requirement Confirmation

The 20 May 2026 meeting confirmed: "full-text search is important, cannot dispense with OpenSearch." The platform team manages the OpenSearch infrastructure; application teams consume it via IAM Identity Centre SSO without administrative access.

### OpenSearch Placement Within the Observability Account — PoC Evaluation

The Observability account is confirmed as the home for query/visualisation tooling. The remaining question is whether OpenSearch lives **centralised in the Observability account** (one domain, all BUs) or **per-BU account** (domain co-located with the log source). Both models use the single Observability account for AMG regardless.

| | Option 1: Centralised (Observability account) | Option 2: Per-BU account |
|--|--|--|
| Operational model | Platform team manages one (or few) OpenSearch domain(s) | Platform team manages N domains across BU accounts |
| Platform team access | Direct — single account, all indices visible | Cross-account — must assume roles into each BU account |
| App team access | Cross-account IAM Identity Centre SSO, scoped to BU indices | Same-account IAM Identity Centre SSO, simpler IAM |
| Log ingestion | Fluent Bit writes cross-account to central OpenSearch | Fluent Bit writes same-account to local OpenSearch |
| Noisy-neighbour risk | Index-level isolation within shared domain; one team's volume could affect performance | Natural isolation — each BU has own domain capacity |
| Cost | Potentially lower (shared infrastructure) | Potentially higher (per-BU domain overhead) |
| Complexity | Simpler infra management; more complex access control (index-level RBAC) | More domains to manage; simpler access control (account boundary) |

### Access Control Model

| Actor | Access scope | Mechanism |
|-------|-------------|-----------|
| Platform team | All logs across all BUs and namespaces | Identity Centre SSO → OpenSearch (admin role) |
| App team engineer | Only logs for namespaces they have access to | Identity Centre SSO → OpenSearch (BU/namespace-scoped role) |
| SOC team | All security-relevant logs | Cortex XSIAM (via S3 pull — not via OpenSearch) |

### Tiered Storage Lifecycle

| Tier | Storage | Retention | Query Latency | Cost (relative) |
|------|---------|-----------|---------------|-----------------|
| Hot | OpenSearch (instance store/EBS) | 14 days | Milliseconds | $$$ |
| Warm | OpenSearch UltraWarm (S3-backed) | 90 days | Seconds | $$ |
| Cold | OpenSearch Cold Storage or S3 | 1 year | Minutes (requires restore) | $ |
| Archive | S3 Glacier Deep Archive | 1 year+ (per compliance) | Hours (requires restore) | ¢ |

### Selective Retrieval from Archive

Engineers must be able to restore logs for a specific time window, namespace, or application without restoring the entire archive. Implementation approach:
- **Bucket model:** one S3 bucket **per business unit** (not per cluster), named `container-platform-{bu}-{env}-eks-logs` (naming confirmed during PoC — see #8430). A single BU bucket is shared across the BU's clusters and log producers.
- **Prefix layout within the bucket:** Fluent Bit application/platform logs land under `logs/`; EKS Auto Mode component logs land under `auto-mode/`. Fluent Bit objects are partitioned by `logs/{cluster}/{namespace}/{YYYY}/{MM}/{DD}/{HH}/`.
- Restore request targets a specific prefix (e.g., one namespace, one day)
- Documented runbook for the restore process, validated during PoC

### Cost Reduction Strategy

Current OpenSearch spend: ~$54k/month. Target: significant reduction through:

| Lever | Mechanism | Expected Impact |
|-------|-----------|-----------------|
| Right-sizing | Match instance types to actual query load (current likely over-provisioned) | 20-40% reduction |
| Tiered storage | UltraWarm for 14-90 day logs (fraction of hot storage cost) | 30-50% reduction on older data |
| Retention enforcement | Automated ISM policies delete indices past retention — no unbounded growth | Prevents cost creep |
| Structured logging | Better indexing reduces storage (fewer catch-all text fields) | 10-20% storage reduction |
| Log filtering | Fluent Bit drops known-noisy, low-value log patterns before ingestion | 10-30% volume reduction |
| Sampling (optional) | For very high-volume debug logs, sample rather than ingest all | Variable |

### OpenSearch Serverless — Evaluation for PoC

OpenSearch Serverless eliminates domain management entirely (no instance sizing, patching, shard rebalancing). It uses OpenSearch Compute Units (OCUs) that auto-scale based on load and can scale to zero when idle.

**Key characteristics:**

| Aspect | Detail |
|--------|--------|
| Pricing | $0.24/OCU/hour (eu-west-2). Production with standby replicas: minimum 4 OCUs (2 indexing + 2 search) = ~$691/month baseline |
| Scale to zero | NextGen collections scale to 0 OCU after 10 minutes of inactivity — no cost when idle |
| Cold start | 10-30 seconds latency on first request after scale-to-zero |
| Collection types | Time series (log analytics), Search (full-text), Vector search |
| Storage tiering | Time series collections automatically tier hot → warm. No manual ISM policy management |
| Access control | Data access policies (different model from provisioned FGAC/RBAC) |
| Dashboards | Supported, but with reduced plugin availability compared to provisioned |

**Where Serverless fits best:**

| Environment | Recommendation | Rationale |
|-------------|----------------|-----------|
| Non-live (development, preproduction) | Serverless (Time series collection) | Scale-to-zero eliminates idle cost overnight/weekends. Cold start acceptable for non-live. |
| Live (production) | Provisioned domain | Predictable sub-second query latency, UltraWarm tiering, full Dashboards plugin support, no cold start during incidents |
| Per-BU model (if selected) | Serverless per-BU | Avoids managing 8-9 separate provisioned domains. Each BU collection scales independently. |

**Trade-offs vs provisioned:**

| Consideration | Provisioned | Serverless |
|---------------|-------------|------------|
| Operational overhead | Platform team manages sizing, patching, shard balance | Zero — fully managed |
| Cost predictability | Fixed (instance hours) + variable (storage) | Variable (OCU hours) — harder to predict at scale |
| UltraWarm/Cold tiers | Full control via ISM policies | Automatic (Time series) but less configurable |
| Fine-grained access | RBAC with backend roles, index-level permissions | Data access policies (principal-based, less granular) |
| Steady-state cost at high volume | Often cheaper at sustained high throughput | Can be more expensive — OCU cost accumulates at continuous load |
| OpenSearch Dashboards | Full plugin ecosystem | Reduced plugin support |
| Ingestion from Fluent Bit | Well-established (direct HTTPS) | Supported via SigV4 auth (less battle-tested at scale) |

**PoC evaluation plan:**

Deploy both models on the PoC cluster:
1. Provisioned domain (right-sized for PoC volume) — baseline
2. Serverless Time series collection — challenger for non-live use

Measure: ingestion reliability from Fluent Bit, query latency, Dashboards functionality, cost projection at full scale, access control expressiveness for namespace-scoped visibility.

**Decision gate:** Include in the Week 10 comparison report. Likely outcome is a hybrid model (Serverless for non-live, provisioned for live) unless Serverless proves cost-effective and performant at sustained production load.

---

## Log Volume Guardrails (FR-206)

### Per-Namespace Controls

| Control | Mechanism | Default | Override |
|---------|-----------|---------|----------|
| Rate alert | CloudWatch metric alarm on log group ingestion rate | Warning at 1 GB/hour per namespace | Platform team adjusts per namespace |
| Hard rate limit | Fluent Bit throttle filter (per-namespace tag) | No hard limit initially — alert only | Can be enabled if a namespace is impacting others |
| Log level recommendation | Documentation + onboarding guidance | INFO in production, DEBUG only in non-live | Team discretion |
| Max log line size | Fluent Bit buffer_max_size | 1 MB per log event | Configurable |

### Platform-Level Controls

- **Total cluster log volume monitoring**: Alert if aggregate cluster output exceeds capacity planning threshold
- **Fluent Bit back-pressure**: If any output destination is slow/down, Fluent Bit pauses reads from that tail input (prevents memory exhaustion)
- **S3 lifecycle enforcement**: Automated — no team can disable archival policies on the shared log bucket

---

## PII and Sensitive Data Handling

### Redaction at Ingestion

Fluent Bit applies regex-based redaction filters before logs leave the cluster:

| Pattern | Description | Action |
|---------|-------------|--------|
| Email addresses | `[\w.]+@[\w.]+\.\w+` | Replace with `[EMAIL_REDACTED]` |
| National Insurance numbers | `[A-Z]{2}\d{6}[A-Z]` | Replace with `[NINO_REDACTED]` |
| Credit card numbers | `\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}` | Replace with `[CARD_REDACTED]` |

### Responsibilities

| Layer | Responsibility |
|-------|---------------|
| Application team | Must not log secrets, tokens, passwords, or keys. Must redact domain-specific PII (case reference numbers, victim names, etc.) |
| Platform (Fluent Bit) | Applies best-effort pattern-based redaction for common PII. Not a guarantee — defence in depth. |
| OpenSearch | Fine-grained access control ensures teams only see their own namespace indices |
| S3 archive | Encrypted at rest (SSE-S3 or SSE-KMS). Bucket policy restricts access to platform team + SOC |

### Audit Logging (Distinct from Application Logs)

Kubernetes API audit logs and CloudTrail logs are security artifacts, not operational logs:
- Forwarded to Cortex XSIAM as priority (SOC requirement)
- Retained in CloudWatch for platform team access (90-day retention for audit logs)
- Not sent to the application-facing OpenSearch domain (separation of concerns)
- Access restricted to platform team and SOC — application teams do not see audit logs

---

## Platform Component Logging (FR-207)

Platform infrastructure components produce their own logs that must be collected separately from application logs:

| Component | Log Source | Key Signals |
|-----------|-----------|-------------|
| Argo CD | Application controller, repo server, server | Sync failures, health check failures, OOM |
| Gatekeeper | Audit controller, controller manager | Policy violations, admission denials |
| Gateway API / Envoy | Envoy access logs, controller logs | 5xx rates, connection errors, certificate issues |
| cert-manager | Controller logs | Certificate renewal failures, ACME errors |
| ExternalDNS | Controller logs | DNS sync failures, throttling |
| EKS Auto Mode | Karpenter, EBS CSI, LB controller, VPC CNI | Scaling events, volume attach failures, ENI exhaustion |

Platform component logs:
- Tagged with `log_type=platform` in Fluent Bit (vs `log_type=application` for tenant pods)
- Sent to a dedicated CloudWatch log group (`/platform/{cluster}/{component}`)
- Indexed in a separate OpenSearch index pattern (`platform-*`)
- Platform team has dashboards and alerts on these logs (separate from app team views)

### EKS Auto Mode Log-Source Enablement Constraint

The four Auto Mode Vended Logs sources (`AUTO_MODE_COMPUTE_LOGS`, `AUTO_MODE_BLOCK_STORAGE_LOGS`, `AUTO_MODE_LOAD_BALANCING_LOGS`, `AUTO_MODE_IPAM_LOGS`) cannot be enabled in a single Terraform apply: each source triggers an EKS cluster update, and EKS allows only one in-flight update at a time, so creating all four together fails with `ConflictException`. Sleep-based sequencing proved unreliable (it failed on idle and Argo CD clusters alike — this is Terraform parallel-processing behaviour, not cluster load).

**Current state:** compute logs only are enabled, to unblock cluster deploys. Reliable sequencing for all four sources is tracked in #8501. Measured Auto Mode log volume is negligible (~3.2 MB/day/cluster, ~2.1 GB/month across 22 clusters), so this is a sequencing problem, not a cost or value one.

---

## Unknowns and Assumptions

| Item | Type | Impact | Mitigation |
|------|------|--------|------------|
| Cortex XSIAM IAM role support | Unknown | Medium — determines whether long-lived access keys can be eliminated | Engage SOC team during Phase 1 to confirm authentication options |
| Actual log volume per cluster | Unknown | High — drives cost estimates for all destinations | Measure on PoC cluster in Phase 1 Weeks 1-4, extrapolate |
| OpenSearch domain sizing | Unknown | High — affects cost reduction target | PoC evaluation with real log volume |
| Fluent Bit S3 plugin reliability under back-pressure | Risk | Medium — log loss during node pressure | PoC testing with simulated load; enable filesystem buffering on host path volume |
| Multi-account Fluent Bit → OpenSearch auth | Unknown | Medium — cross-account writes need IAM or SAML | PoC evaluation of both centralised and per-BU models |
| Application team structured logging adoption rate | Unknown | Low-Medium — affects OpenSearch query experience | Track metrics; provide tooling and examples; do not mandate |

---

## Consequences

### Positive
- S3 as primary archive — cheapest durable storage, feeds SOC without additional infrastructure
- Existing Cortex XSIAM pattern preserved — no SOC team rework required
- CloudWatch provides immediate operational value (Logs Insights, alarms) with short retention to limit cost
- OpenSearch retains the full-text search capability engineers rely on
- Structured logging standard improves query experience over time without blocking unstructured logs
- Platform and application logs separated — different retention, access, and alerting models
- Tiered storage lifecycle automates cost control

### Negative
- Three Fluent Bit outputs per cluster — more configuration to manage (mitigated by Terraform/GitOps templating)
- OpenSearch still required (cannot eliminate the $54k cost entirely) — but right-sizing and tiering should reduce substantially
- PII redaction is best-effort at platform level — application teams retain responsibility for sensitive data
- Fluent Bit is a single point of collection — pod failure means brief log gap (mitigated by DaemonSet, filesystem buffer)

### Neutral
- CloudWatch is used for operational convenience, not as the primary archive — cost-controlled via short retention
- OpenSearch account placement decision deferred to PoC evaluation
- Cortex XSIAM authentication mechanism (IAM role vs access keys) depends on vendor capability

---

## Related Decisions

- **Depends On:** ADR-004 (IAM Identity Center — OpenSearch SSO), ADR-005 (Observability Stack — metrics/dashboards)
- **Supersedes:** ADR-005 OpenSearch section (moved here)
- **Related:** ADR-001 (EKS), ADR-008 (EKS Auto Mode)

## Related Issues

Tracking for this ADR lives in `ministryofjustice/cloud-platform`. The PoC that produces the evidence for the Phase 1 Week 10 decision gate is broken down as follows:

| Issue | Scope |
|-------|-------|
| [ministryofjustice/cloud-platform#8415](https://github.com/ministryofjustice/cloud-platform/issues/8415) | **Parent spike** — Logging and Log Management Architecture Validation (comparison report, decision gate) |
| [ministryofjustice/cloud-platform#8419](https://github.com/ministryofjustice/cloud-platform/issues/8419) | Fluent Bit multi-output pipeline (S3 + CloudWatch + OpenSearch) reliability and parsing |
| [ministryofjustice/cloud-platform#8420](https://github.com/ministryofjustice/cloud-platform/issues/8420) | OpenSearch deployment model — Provisioned vs Serverless, centralised vs per-BU |
| [ministryofjustice/cloud-platform#8421](https://github.com/ministryofjustice/cloud-platform/issues/8421) | SOC log forwarding — Cortex XSIAM in multi-account |
| [ministryofjustice/cloud-platform#8422](https://github.com/ministryofjustice/cloud-platform/issues/8422) | Logging cost analysis at scale (NFR-106) |
| [ministryofjustice/cloud-platform#8430](https://github.com/ministryofjustice/cloud-platform/issues/8430) | EKS Auto Mode enhanced logging — CloudWatch Vended Logs |
| [ministryofjustice/cloud-platform#8501](https://github.com/ministryofjustice/cloud-platform/issues/8501) | EKS Auto Mode logging — reliable sequencing of all four log-type sources |

---

## Research Sources

1. Cloud Platform Runbook — Logs to SOC Cortex XSIAM: https://runbooks.cloud-platform.service.justice.gov.uk/logs-to-soc-cortex-xsiam.html
2. EKS Auto Mode Managed Component Logs: https://docs.aws.amazon.com/eks/latest/userguide/auto-managed-component-logs.html
3. CloudWatch Logs Subscription Filters: https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/SubscriptionFilters.html
4. Fluent Bit S3 Output Plugin: https://docs.fluentbit.io/manual/pipeline/outputs/s3
5. Fluent Bit Multi-Destination Configuration: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/firelens-docker-buffer-limit.html
6. CloudWatch Logs Pricing (eu-west-2): https://aws.amazon.com/cloudwatch/pricing/
7. OpenSearch UltraWarm Storage: https://docs.aws.amazon.com/opensearch-service/latest/developerguide/ultrawarm.html
