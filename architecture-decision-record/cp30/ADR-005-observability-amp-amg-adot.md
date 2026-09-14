# ADR-005: Observability Stack — Metrics, Dashboards, and Alerting

**Status:** Proposed — Phase 1 PoC validation required (metrics backend); AMG placement decided 2026-09-07  
**Date:** 2026-05-07 (updated 2026-05-12, 2026-09-07)  
**Decision Maker:** AWS ProServe / MoJ Principal Technical Architect  
**Category:** Compute

---

## Context

**Problem Statement:** The current observability stack (self-managed Grafana ~£7,000/month, OpenSearch used as log dump, AlertManager) has high operational overhead and is not configured for daily use. The new multi-cluster platform requires a managed observability solution with cross-cluster aggregation, namespace-scoped access, direct engineer access without access-request barriers, and UK data residency.

**Requirements:** FR-012 (shared observability, direct engineer access), FR-202 (self-service alerts), NFR-101 (managed observability stack), NFR-004 (UK data residency). Log-related requirements (FR-204, FR-208) are addressed in [ADR-017](ADR-017-logging-and-log-management.md).

**Constraints:**
- UK data residency required (NFR-004)
- No Lambda (NFR-003)
- IAM Identity Center already required (ADR-004) — AMG dependency is already satisfied
- Current self-managed Grafana costs ~£7,000/month
- MoJ team has existing Grafana skills (current self-managed stack)
- 07 May meeting: logs should go to CloudWatch first, then optionally to per-BU OpenSearch
- 20 May meeting: full-text search is important, cannot dispense with OpenSearch; platform team defines clusters, app teams consume without admin access; explore placement in Platform Services account; avoid moving logs around for ingestion (cost concern)

---

## Decision

**Both approaches will be deployed on the Phase 1 PoC cluster for comparative evaluation. The final decision will be made before Phase 2 based on PoC findings.**

**Option A (baseline):** AMP + AMG + ADOT — full Prometheus-native pipeline with Grafana UX  
**Option D (challenger):** CloudWatch Container Insights + AMG (querying CloudWatch as data source) — simpler pipeline, fewer moving parts

Both options use Amazon Managed Grafana for the dashboard UX (team familiarity) and Fluent Bit + CloudWatch Logs for log aggregation. The difference is the metrics backend: AMP (Prometheus-compatible, PromQL) vs. CloudWatch Metrics (native, CloudWatch query language).

### PoC Evaluation Criteria

| Criterion | Weight | How to Measure |
|-----------|--------|----------------|
| Operational complexity | High | Number of components to configure, cross-account setup effort, failure modes |
| Query capability | High | Can engineers write the queries they need? PromQL vs. CloudWatch Metrics Insights |
| Alert scalability | High | 100-rule AMG limit vs. 5,000 CloudWatch Alarms — does the limit matter at PoC scale? |
| Namespace-scoped access | Medium | How easily can BU teams see only their data in each approach? |
| Ingestion cost at projected scale | Medium | Estimate ingestion cost for 22 clusters × expected metric cardinality (AMP samples ingested vs. CloudWatch custom metrics) |
| Storage cost at projected scale | Medium | Estimate storage cost at the platform's target retention: AMP GB-month vs. CloudWatch per-metric-month, projected as cardinality × retention; include sensitivity to longer retention |
| Team adoption friction | Medium | Which approach do MoJ engineers find easier to use? |

### Decision Gate

The PoC team will produce a comparison report by Phase 1 Week 10 covering the criteria above. The final observability architecture decision will be included in the Phase 2 go/no-go recommendation (US-008).

---

## Research Conducted

### Option A: AMP + AMG + ADOT (SELECTED)

**Research Confidence:** High

| Source | URL | Key Finding |
|--------|-----|-------------|
| AWS Observability Accelerator | https://aws.amazon.com/solutions/guidance/monitoring-amazon-eks-workloads-using-amazon-managed-services-for-prometheus-and-grafana/ | Terraform-based; AMP + AMG + ADOT pattern; EKS-specific dashboards |
| Cross-Account AMG | https://aws.amazon.com/blogs/opensource/setting-up-amazon-managed-grafana-cross-account-data-source-using-customer-managed-iam-roles/ | Single AMG workspace; cross-account AMP data sources via IAM roles |
| AMG Quotas | https://docs.aws.amazon.com/grafana/latest/userguide/AMG_quotas.html | 100 alert rules per workspace (hard limit) — must be planned for |
| AMP Quotas | https://docs.aws.amazon.com/prometheus/latest/userguide/AMP_quotas.html | 50M active series per workspace; auto-adjusting |
| Observability Research | observability-research.md | eu-west-2 confirmed available; Identity Center SSO integration |

**Capabilities Verified:**
- AMG requires IAM Identity Center — confirmed (aligned with ADR-004)
- Cross-account AMP data sources via customer-managed IAM roles — verified
- Folder-level access control per BU in AMG — verified
- Terraform provider for AMG workspace, data sources, dashboards — verified
- ADOT EKS add-on available in eu-west-2 — confirmed
- CloudWatch Container Insights + Fluent Bit for logs (FluentD deprecated February 2025) — verified

### Option B: Self-Managed Prometheus + Grafana

**Research Confidence:** High

| Source | URL | Key Finding |
|--------|-----|-------------|
| Current MoJ stack | .apex/customer-context.md | Self-managed Grafana ~£7,000/month; compute management problematic |
| Thanos/Cortex for multi-cluster | General research | Thanos requires sidecar per Prometheus; complex HA |

**Capabilities Verified:**
- Full Prometheus/Grafana feature parity — confirmed
- Current cost: ~£7,000/month confirmed as pain point — confirmed
- Multi-cluster: requires Thanos or Cortex (additional operational complexity) — confirmed
- Platform team manages upgrades, HA, storage, scaling — confirmed operational overhead

### Option D: CloudWatch Container Insights + AMG (querying CloudWatch) — CHALLENGER

**Research Confidence:** High

| Source | URL | Key Finding |
|--------|-----|-------------|
| CloudWatch Container Insights | https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Container-Insights-EKS.html | Native EKS integration; automatic cluster/node/pod metrics; zero-config |
| AMG CloudWatch Data Source | https://docs.aws.amazon.com/grafana/latest/userguide/using-amazon-cloudwatch-in-AMG.html | AMG queries CloudWatch directly; cross-account via IAM |
| CloudWatch Cross-Account | https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Cross-Account-Cross-Region.html | Delegated admin for cross-account metric/log visibility |
| CloudWatch Alarms | https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/AlarmThatSendsEmail.html | 5,000 alarms per region (adjustable); no hard limit concern |

**Capabilities Verified:**
- Container Insights: automatic EKS metrics (cluster, node, pod, container) — verified
- AMG can query CloudWatch as a data source (no AMP needed) — verified
- CloudWatch cross-account observability via delegated admin — verified
- CloudWatch Alarms: 5,000 per region (adjustable) — no 100-rule constraint
- Fluent Bit → CloudWatch Logs: same log pipeline as Option A — verified
- No ADOT collector needed (Container Insights agent handles metrics) — verified
- No AMP workspace provisioning per cluster — simpler setup
- CloudWatch Metrics Insights query language (not PromQL) — confirmed trade-off

**Key advantages over Option A:**
- Fewer moving parts (no ADOT, no AMP, no cross-account AMP IAM roles)
- No 100-alert-rule constraint (CloudWatch Alarms scale to 5,000+)
- No Identity Centre dependency for alerting (only for AMG dashboard access)
- Single metrics+logs pipeline (CloudWatch handles both)

**Key disadvantages vs. Option A:**
- No PromQL — CloudWatch uses its own query language
- Less flexible dashboard building (CloudWatch data source in Grafana is less expressive than Prometheus)
- Custom application metrics require CloudWatch EMF format or custom metric API calls (vs. Prometheus scraping)

### Option C: Datadog

**Research Confidence:** Medium

| Source | URL | Key Finding |
|--------|-----|-------------|
| Datadog EKS Integration | https://docs.datadoghq.com/integrations/amazon_eks/ | Per-host pricing; multi-cluster native; UK region available |
| BRD Recommendation | MoJ Container Platform BRD v1.1 | "Check existing Datadog/Grafana Cloud licensing before finalising" |

**Capabilities Verified:**
- Native multi-cluster observability — verified
- Per-host pricing can be expensive at scale — verified
- UK data residency: Datadog EU region (not AWS-native) — confirmed
- MoJ may have existing Datadog licensing — unknown (BRD open item)

---

## Capability Mapping

| Requirement | Option A (AMP+AMG+ADOT) | Option D (CloudWatch+AMG) | Evidence |
|-------------|-------------------------|---------------------------|----------|
| FR-012: Direct engineer access | AMG folder permissions + Identity Center groups | AMG folder permissions (same UX) + CloudWatch IAM | AMG docs |
| NFR-101: Managed observability | AWS-managed; 3 services to configure | AWS-managed; 1 service (CloudWatch) + AMG | AMP/CW docs |
| NFR-004: UK data residency | AMP and AMG in eu-west-2 | CloudWatch and AMG in eu-west-2 | Regional availability |
| FR-202: Self-service alerts | AMG alert rules (100-rule hard limit) | CloudWatch Alarms (5,000/region, adjustable) + AMG | AMG/CW Quotas |
| FR-012: Platform alerts before BU | AMP Alert Manager + SNS | CloudWatch Alarms + SNS (no limit concern) | CW docs |
| PromQL support | Native (AMP is Prometheus-compatible) | Not available (CloudWatch query language) | AMP/CW docs |
| Custom app metrics | Prometheus scraping via ADOT | CloudWatch EMF or PutMetricData API | ADOT/CW docs |
| Operational complexity | Medium (ADOT + AMP + AMG + cross-account IAM) | Low (Container Insights toggle + AMG) | Architecture assessment |

---

## Unknowns and Assumptions

| Item | Type | Impact | Mitigation |
|------|------|--------|------------|
| AMG 100 alert rules limit at scale (Option A) | Risk | Medium — 8-9 BUs × N rules may approach limit | Platform-critical alerts in AMP Alert Manager; BU alerts in AMG; OR use Option D (CloudWatch Alarms, no limit) |
| PromQL investment vs. fresh start | Unknown | High — determines whether Option A's PromQL advantage matters | PoC evaluation: if MoJ has no existing PromQL to migrate, Option D's simpler pipeline may win |
| CloudWatch custom metric cost at high cardinality | Risk | Medium — CloudWatch per-metric pricing can be expensive with high-cardinality labels | PoC: measure actual metric cardinality on PoC cluster; compare AMP vs. CW cost |
| Team preference for query language | Unknown | Medium — affects daily usability | PoC: have MoJ engineers try both PromQL (AMP) and CloudWatch Metrics Insights; gather feedback |

---

## Phase 1 PoC Validation Plan

Both Option A and Option D will be deployed on the PoC cluster:

**Option A deployment:**
- ADOT Collector DaemonSet (EKS add-on)
- AMP workspace (single, for PoC cluster)
- AMG workspace with AMP data source

**Option D deployment:**
- CloudWatch Container Insights (EKS Auto Mode native)
- AMG workspace with CloudWatch data source (same AMG workspace, additional data source)

**Evaluation activities (Phase 1 Weeks 6–10):**
1. Deploy both pipelines on PoC cluster simultaneously
2. Build equivalent dashboards in AMG using both data sources
3. Create equivalent alert rules (AMP Alert Manager vs. CloudWatch Alarms)
4. Measure: setup time, cross-account configuration effort, query expressiveness
5. Have 2-3 MoJ engineers use both approaches for 2 weeks; gather preference feedback
6. Produce cost impact analysis comparing current observability tooling vs both proposed options at full scale (22 clusters × projected metric cardinality), itemising ingestion and storage costs separately for each option (AMP GB-month vs. CloudWatch per-metric-month, driven by cardinality × retention) — satisfies NFR-106
7. Document findings in comparison report by Week 10

**Decision gate:** Comparison report feeds into US-008 (Phase 2 go/no-go recommendation). Final observability stack decision is made before Phase 2 build begins.

---

## Account Placement

> **Updated 2026-09-07 (tech lead huddle):** AMG is placed in the existing
> **`cloud-platform-live` account**, not a new dedicated Observability account,
> and uses a **single AMG workspace**. The original dedicated-account proposal
> is retained below under "Superseded proposal" for context.

Amazon Managed Grafana (AMG) resides in the existing **`cloud-platform-live` account**. A separate dedicated Observability account was considered but rejected: the overhead of managing an additional account was judged excessive for what is primarily dashboard hosting, and app teams receive only **read-only** dashboard access, so the isolation a separate account would provide is not warranted.

**Rationale:**
- App teams get read-only access to AMG dashboards and never receive IAM access to the account. The security exposure that motivated a separate account (protecting the Argo CD hub) is low under a read-only dashboard access model.
- Avoids the operational overhead of provisioning and maintaining an additional AWS account for a visualisation layer.
- `cloud-platform-live` already exists and is managed by the platform team.

**Single AMG workspace (all BUs, both live and non-live):**
- One workspace, not per-BU and not split live/non-live.
- Driven by AMG licensing: **$5 / viewer / month, $9 / editor / month**, charged only for users who log in that month, minimum one license per workspace per month. Separate workspaces would double-charge any user who spans BUs or environments.
- Live/non-live workspace separation is **deferred** (roughly 80% live / 20% non-live dashboards). Grafana folder permissions separate views within the single workspace if needed.

**Dashboard capacity (confirmed sufficient):** each AMG workspace has a hard limit of 2000 dashboards (default 5 workspaces per account, adjustable on request via quota `L-2C2D5119` — AWS publishes no fixed maximum, increases are case-by-case). MoJ currently has **fewer than 400 dashboards** across all environments, so a single workspace is sufficient for production with substantial headroom. This is a monitor-only item, not a constraint on the single-workspace decision.

**Access model:**

| Actor | AMG Role | Scope | Access method |
|-------|----------|-------|---------------|
| Platform engineer | Admin | All folders, all data sources | Identity Centre SSO → AMG |
| BU app engineer | Viewer (read-only) | Own BU folder | Identity Centre SSO → AMG |

**Data flow:**
- AMP workspaces (per BU) and CloudWatch remain in each BU account (data stays close to source).
- The single AMG workspace in `cloud-platform-live` queries BU accounts cross-account via IAM role assumption. This means AMG carries **many data sources** (one per BU AMP/CloudWatch), which requires clear data-source naming conventions and sample dashboards.
- No metric/log data is stored in `cloud-platform-live` for this purpose — AMG is a query and visualisation layer only.

**Dashboards-as-code:** users author dashboards in the Grafana UI, export JSON, commit to the repo, and a pipeline applies them (mirroring the Argo CD environment-folder pattern with an added dashboards folder). This gives version control, audit trail, and recovery if a workspace is lost or upgraded.

### Deployment: AMG in the cloud-platform root component

AMG is deployed as part of the **root `cloud-platform` Terraform component** (alongside the other account-level singletons such as IAM roles, ECR, and Route53), not a dedicated component. A standalone `observability` component was prototyped but judged overkill for a single workspace plus one IAM role: root already provides the required providers (`aws`, `http`), the per-workspace model, and the related account-level IAM, and it applies as the first pipeline stage. The metrics collectors (AMP workspace + ADOT, or the CloudWatch Observability add-on) remain **per-cluster** in `cluster-components`; only AMG — a single, account-level resource — lives at root.

Which environments host an AMG workspace is declared in root's `locals.tf` (not passed as a runtime flag):

- **`cloud-platform-live`** — the production central AMG serving all BUs and both live and non-live environments.
- **`cloud-platform-development`** — a test AMG in the development account, so ephemeral clusters' dashboards and cross-cluster data sources can be validated.

`enable_amg` is derived from `contains(local.amg_host_workspaces, terraform.workspace)`, so every other workspace (BU spokes, preproduction, nonlive) creates no AMG resources. Designating a new AMG host is a reviewed one-line code change, which prevents an AMG workspace being created in the wrong account by accident.

### BU metric isolation

**Requirement (security team, 2026-09-07):** logical separation of metrics between BUs — if a BU publishes sensitive metrics, other BUs must not be able to view them by default. This is default-deny visibility, not a hard cryptographic boundary.

**Key constraint:** in a single shared AMG workspace, Grafana **folder permissions control dashboard visibility only — they do not restrict which data sources a user can query.** A user with Explore access, or a dashboard exposing a data-source picker, could otherwise query any data source in the workspace. Folders alone therefore do **not** provide the required isolation.

Isolation is enforced in three combined layers:

1. **Grafana data source permissions (primary enforcement).** One data source per BU (e.g. `amp-hmpps`, `amp-laa`), each restricted so only that BU's Grafana team can query it. This is what actually prevents cross-BU querying, including from Explore.
2. **Per-BU IAM assume-roles (AWS-layer defence in depth).** Each BU data source uses a distinct cross-account IAM role scoped to that BU's AMP/CloudWatch, rather than one broad role over all backends.
3. **Folders + Viewer role (visibility hygiene).** Per-BU folders and read-only (Viewer) access so users see only their dashboards and cannot author dashboards pointing at other data sources.

The binding chain is: **IAM Identity Center group (per BU) → Grafana team → {folder permission + data source permission}**.

**Implementation dependency:** the AWS provider manages the AMG *workspace* and IAM roles, but Grafana-internal objects (teams, folders, data sources, and data source permissions) are managed through the Grafana provider / dashboards-as-code pipeline that authenticates into the workspace. The isolation is only realised once that second layer exists; the root `cloud-platform` component provisions the workspace and IAM foundation, and the Grafana-object layer is tracked as follow-up work (see Related Decisions).

**Superseded proposal (original, pre-2026-09-07):** AMG in a dedicated Observability account separate from the Platform Services account, on the grounds that app engineers should never touch the account hosting Argo CD. Rejected in favour of the above because app-team access is read-only and dashboard-only, making the extra account's overhead unjustified.

---

## Counter-Argument Analysis

### Q1: What specific evidence would make me choose self-managed instead?

If MoJ had an existing enterprise Grafana Cloud contract covering the full fleet, the cost model for AMP+AMG might not provide savings. However, the current self-managed Grafana is confirmed as a pain point (~£7k/month, operational burden). No platform-level Datadog evaluation is needed — Datadog is a team-level choice outside platform scope.

### Q2: Is there a managed/integrated AWS service that does this better?

AMP + AMG + ADOT is the AWS-native managed stack for this capability. CloudWatch-only is another AWS option (fully managed, no AMG) but lacks the visual dashboard experience and has higher per-metric cost at large scale.

### Q3: What am I NOT seeing about the alternative that might make it superior?

Self-managed Grafana + Thanos provides full feature control without the AMG 100 alert rule limit. At large scale (50+ BUs), the AMG 100-rule limit becomes a genuine constraint. MoJ's scale (8-9 BUs) makes it manageable but requires careful budget allocation.

---

## Alternative Consideration Checklist

- [x] Searched for managed/integrated AWS services for this capability domain
- [x] Researched minimum 2 viable alternatives with equal depth
- [x] Documented specific AWS documentation URLs for each alternative
- [x] Created capability mapping with evidence sources
- [x] Documented unknowns and assumptions
- [x] Assigned research confidence level to each alternative
- [x] Completed Counter-Argument analysis

---

## Alternatives Considered

### Option B: Self-Managed Prometheus + Grafana — REJECTED

**Why Rejected:** The current self-managed stack costs ~£7,000/month and is described as having "compute management problematic." Moving to the same self-managed model in the new architecture perpetuates the operational burden. AWS-managed alternatives exist that satisfy all requirements at comparable cost.

**Pros:** No alert rule limits; full feature control; no IAM Identity Center dependency for Grafana
**Cons:** ~£7,000/month confirmed pain point; platform team manages upgrades, HA, storage; multi-cluster requires Thanos/Cortex (additional complexity)

### Option C: Datadog — NOT APPLICABLE (team-level tool)

**Why Not Applicable:** Datadog is used by individual application teams at their own discretion. It is not a platform-level observability decision. Teams may continue using Datadog alongside the platform-provided managed stack if they choose. No platform licensing evaluation required.

---

## Rationale

In the context of replacing a self-managed Grafana stack costing ~£7,000/month with high operational overhead, facing the need for multi-cluster cross-account observability with direct engineer access and UK data residency, we decided to validate both AMP+AMG+ADOT (Option A) and CloudWatch Container Insights+AMG (Option D) during Phase 1 PoC, and rejected self-managed Prometheus/Grafana outright, because:

1. Self-managed is the confirmed pain point — any managed alternative is an improvement
2. Option A provides PromQL compatibility and richer cross-account aggregation, but has a 100-alert-rule hard limit and more moving parts
3. Option D provides simpler operations with fewer components and no alert-rule constraint, but loses PromQL and has less expressive custom metrics
4. MoJ is building a new platform (not migrating existing Prometheus configs), so the PromQL advantage may not be decisive
5. The 07 May meeting clarified that logs go to CloudWatch first regardless — both options share the same log pipeline
6. A 2-week comparative evaluation on the PoC cluster is low-cost and eliminates guesswork

The final decision will be evidence-based, made by Phase 1 Week 10, and included in the go/no-go recommendation.

---

## Consequences

### Positive (both options)
- No Grafana/Prometheus infrastructure to manage
- Cross-cluster unified view via AMG
- Identity Center SSO for Grafana (no separate user management)
- CloudWatch Logs as shared log backend (Fluent Bit DaemonSet per cluster)

### Positive (Option A specific)
- PromQL compatibility — existing Prometheus queries port directly
- ADOT is vendor-neutral (OpenTelemetry); avoids metric pipeline lock-in
- Richer cross-account data source model in AMG

### Positive (Option D specific)
- Fewer moving parts (no ADOT, no AMP, no per-cluster workspace provisioning)
- No 100-alert-rule constraint (CloudWatch Alarms scale to 5,000+)
- Single pipeline for metrics and logs (CloudWatch handles both)
- Simpler cross-account setup (CloudWatch delegated admin)

### Negative (Option A specific)
- AMG 100 alert rule hard limit requires hybrid alert management
- 3 services to configure and maintain (ADOT + AMP + AMG)
- Cross-account AMP IAM roles per cluster

### Negative (Option D specific)
- No PromQL — CloudWatch query language is less expressive
- Custom application metrics require EMF format (less natural than Prometheus scraping)
- CloudWatch per-metric cost can be high at high cardinality

### Neutral
- Phase 1 PoC runs both approaches simultaneously (low additional cost)
- Final decision deferred to Phase 1 Week 10 based on evidence
- AMG is common to both options — team builds Grafana skills regardless

---

## Related Decisions

- **Depends On:** ADR-004 (IAM Identity Center — AMG SSO, OpenSearch SSO)
- **Related:** ADR-001 (EKS), ADR-002 (Argo CD), ADR-017 (Logging and Log Management)

---

## OpenSearch and Log Management

**This section has been moved to [ADR-017: Logging and Log Management Architecture](ADR-017-logging-and-log-management.md).**

ADR-017 covers the full logging pipeline architecture (Fluent Bit multi-output), OpenSearch deployment model and account placement evaluation, structured logging standards (FR-208), SOC integration (Cortex XSIAM), tiered storage lifecycle, log volume guardrails, PII handling, and platform component logging. This ADR (005) retains scope over metrics, dashboards, and alerting only.

---

## Research Sources

1. AWS Observability Accelerator for EKS — https://aws.amazon.com/solutions/guidance/monitoring-amazon-eks-workloads-using-amazon-managed-services-for-prometheus-and-grafana/
2. Cross-Account AMG Data Sources — https://aws.amazon.com/blogs/opensource/setting-up-amazon-managed-grafana-cross-account-data-source-using-customer-managed-iam-roles/
3. AMG Quotas — https://docs.aws.amazon.com/grafana/latest/userguide/AMG_quotas.html
4. AMP Quotas — https://docs.aws.amazon.com/prometheus/latest/userguide/AMP_quotas.html
5. ADOT EKS Metrics/Traces — https://aws.amazon.com/blogs/containers/metrics-and-traces-collection-using-amazon-eks-add-ons-for-aws-distro-for-opentelemetry/
