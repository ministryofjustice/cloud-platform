# ADR-005: Observability Stack — Metrics, Dashboards, and Alerting

**Status:** Accepted — metrics backend decided 2026-09-30 (Option A, AMP); AMG placement decided 2026-09-07; BU metric isolation implemented 2026-09-16 (cloud-platform#8509); collector decided 2026-09-30 (self-managed ADOT)  
**Date:** 2026-05-07 (updated 2026-05-12, 2026-09-07, 2026-09-16, 2026-09-30)  
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

> **Decided 2026-09-30.** The Phase 1 PoC ran Option A and Option D in parallel on
> separate clusters. The decisions below are final. The original "evaluate both
> during the PoC" framing and the full research for each option are retained
> further down as the evidence trail.

**The metrics platform is AMP + AMG + ADOT (Option A).** Metrics are collected on each cluster and remote-written to Amazon Managed Service for Prometheus (AMP); Amazon Managed Grafana (AMG) is the dashboard and query layer. CloudWatch Container Insights (Option D) is **not** the primary metrics store, but its code is **retained behind a feature flag** (`enable_cloudwatch_observability`, default off) for two narrow uses: AWS-native signals that originate in CloudWatch (e.g. RDS status), and CloudWatch's cross-service "explore related" root-cause view, which can be switched on selectively for the platform team.

**The metrics collector is self-managed ADOT**, not the AMP agentless managed collector and not the CloudWatch agent. ADOT runs in-cluster (AWS Distro for OpenTelemetry) and remote-writes to AMP over SigV4 / EKS Pod Identity. It was chosen because it is a **single component that can satisfy several requirements as they arrive**: scraping platform and application `/metrics` (pull), receiving pushed OpenTelemetry (OTLP) metrics from tenants, and — later, if needed — receiving OTLP traces and deriving span metrics into AMP. See "Collector decision" below for the alternatives and why they were ruled out.

**Alerting is Prometheus + AlertManager, not AMG-managed alerting.** This supersedes the AMG 100-alert-rule framing throughout this ADR; that limit no longer applies to the chosen design. The user-facing alert-configuration experience is designed separately (cloud-platform#7871).

**Unchanged across both options and still in force:** AMG is the dashboard UX (team familiarity, Identity Centre SSO), and logs go via Fluent Bit (log architecture is now owned by [ADR-017](ADR-017-logging-and-log-management.md)).

### Collector decision (2026-09-30): self-managed ADOT

Three ways to get metrics into AMP were considered. ADOT was chosen.

| Option | What it is | Ruled out because |
|--------|-----------|-------------------|
| **Self-managed ADOT (CHOSEN)** | OpenTelemetry collector run in-cluster as code; scrapes Prometheus endpoints and/or receives OTLP, remote-writes to AMP | — |
| **AMP agentless managed collector** | AWS-run scraper, no in-cluster collector; AWS owns its scaling/HA/upgrades | Metrics-only and **pull-only**, so pushed OTLP tenant metrics would need a *second* in-cluster collector anyway. Carries a per-collector charge (~$0.04/collector-hour + $0.03/10M samples ≈ $32/cluster/month, ~$7.8k/yr across 22 clusters, mostly fixed). Does **not** reduce scrape-config effort (same Prometheus config) and does **not** remove the need to run node-exporter and kube-state-metrics in-cluster. Soft limit of 10 scrapers/region/account. |
| **CloudWatch agent** | AWS-managed OpenTelemetry collector with CloudWatch components | Sends metrics to CloudWatch, which is the store we did **not** select as primary. Would still need an AMP collector for the Prometheus path, i.e. two agents per node. |

**Why ADOT wins:** there is a **near-term requirement for tenant application metrics**, and tenants may expose them either as Prometheus `/metrics` (pull) or as pushed OTLP. ADOT is the only one of the three that covers **both** patterns with a single agent, and the same agent can later add a traces pipeline (OTLP → X-Ray or OpenSearch) and derive span metrics into AMP. The managed collector's appeal was offloading collector operations to AWS; that appeal is outweighed once a second collector would be needed anyway for OTLP, and given it saves no scrape-config or exporter work.

**Trade-off accepted:** the team owns the ADOT collector lifecycle — upgrades alongside EKS, occasional collector-config deprecations, HA, scrape de-duplication, and memory tuning (estimated at a few engineer-days a year, plus the collectors' node resources) — in exchange for one extensible collector and no per-collector charge. The historical Prometheus pain was running the Prometheus *server* (storage, retention, HA), which AMP removes; the collector is the small part of the stack.

**Known defect to fix in the production build (cloud-platform#7867):** the PoC ADOT collector runs as a DaemonSet with cluster-wide service discovery on every pod, so each target is scraped once per node (measured ~3× on a 3-node cluster). It is silent (no errors) and ingestion cost scales with node count. The production design fixes this with node-scoped discovery (each DaemonSet pod scrapes only its own node) plus a small replica-set collector for cluster-level targets using AMP's `cluster`/`__replica__` de-duplication; the OpenTelemetry Target Allocator is an alternative that additionally enables tenant ServiceMonitor/PodMonitor self-service. Baseline cardinality measured during the PoC: ~16k active series per cluster before node-exporter, kube-state-metrics, control-plane, and add-on metrics are added.

### PoC Evaluation Criteria (how the decision was reached)

| Criterion | Weight | How to Measure |
|-----------|--------|----------------|
| Operational complexity | High | Number of components to configure, cross-account setup effort, failure modes |
| Query capability | High | Can engineers write the queries they need? PromQL vs. CloudWatch Metrics Insights |
| Alert scalability | High | 100-rule AMG limit vs. 5,000 CloudWatch Alarms — does the limit matter at PoC scale? |
| Namespace-scoped access | Medium | How easily can BU teams see only their data in each approach? |
| Ingestion cost at projected scale | Medium | Estimate ingestion cost for 22 clusters × expected metric cardinality (AMP samples ingested vs. CloudWatch custom metrics) |
| Storage cost at projected scale | Medium | Estimate storage cost at the platform's target retention: AMP GB-month vs. CloudWatch per-metric-month, projected as cardinality × retention; include sensitivity to longer retention |
| Team adoption friction | Medium | Which approach do MoJ engineers find easier to use? |

### Decision Gate (closed)

The PoC ran both options and the decision was taken on 2026-09-30: **Option A (AMP), collected by self-managed ADOT.** The PoC epic and its children are closed (cloud-platform#8414). Production build work is tracked under US-107 (#8235) and US-013 (#8247) — see "Implementation tracking" below. This outcome feeds the Phase 2 go/no-go recommendation (US-008).

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

## Unknowns and Assumptions (resolved by the PoC)

| Item | Resolution |
|------|------------|
| AMG 100 alert rules limit at scale | No longer applies. Alerting is Prometheus + AlertManager, not AMG-managed alerting, so the AMG rule limit is not a constraint on the chosen design. |
| PromQL investment vs. fresh start | Resolved in favour of Prometheus. Familiarity and low migration friction favoured PromQL; this was a deciding factor for Option A. |
| CloudWatch custom metric cost at high cardinality | Measured: ~16k active series per PoC cluster. The ingestion/storage cost difference between AMP and CloudWatch was minor relative to overall observability spend (dominated by logging), so cost did not decide the metrics backend. |
| Team preference for query language | Resolved in favour of PromQL/Prometheus (familiar to engineers, no workflow change). |

### Open items carried into the production build

| Item | Where tracked |
|------|---------------|
| Control-plane metrics on EKS Auto Mode (API server/scheduler/controller-manager available; etcd not exposed) — confirm on Auto Mode specifically | cloud-platform#7867 |
| Fix the ADOT 3× scrape duplication; deploy node-exporter + kube-state-metrics; enable per-add-on metrics | cloud-platform#7867 |
| Platform-team fleet-wide view via an explicit grant (not AMG admin-bypass) | cloud-platform#7867 |
| ADR text and `amg.tf` (`unifiedAlerting`) still reflect the old AMG-alerting assumption and need reconciling with the AlertManager decision | cloud-platform#7871 |
| Tenant application metrics (near-term): define the exposure contract (annotations / ServiceMonitor / OTLP) and per-tenant cardinality guardrails | cloud-platform#7907 |

---

## Phase 1 PoC Validation Plan (completed — retained as the record of what was run)

> This plan was executed; both options ran on separate clusters and the decision
> (Option A, ADOT) was taken on 2026-09-30. Retained for the record.

Both Option A and Option D were deployed on the PoC clusters:

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

The AMG service role and workspace need no additional deploy-role permissions from the main pipeline's apply role (confirmed with the platform team), so AMG does not depend on the `oidc.tf` grant. The AMG workspace has no dependency on any cluster existing, so applying it in the root (first) stage is safe.

### BU metric isolation

> **Updated 2026-09-16 (cloud-platform#8509 implemented):** the isolation model
> below is now realised in Terraform in the root `cloud-platform` component
> (`grafana-objects.tf`), with a dev PoC deployed against two simulated BUs. The
> implementing change is
> [modernisation-platform-environments#19165](https://github.com/ministryofjustice/modernisation-platform-environments/pull/19165).

**Requirement (security team, 2026-09-07):** logical separation of metrics between BUs — if a BU publishes sensitive metrics, other BUs must not be able to view them by default. This is default-deny visibility, not a hard cryptographic boundary.

**Key constraint:** in a single shared AMG workspace, Grafana **folder permissions control dashboard visibility only — they do not restrict which data sources a user can query.** A user with Explore access, or a dashboard exposing a data-source picker, could otherwise query any data source in the workspace. Folders alone therefore do **not** provide the required isolation.

**LBAC ruled out:** Grafana label-based access control (LBAC), which would let a single data source restrict results by label per team, is **not available in AMG** (Grafana Cloud/Enterprise only — verified against AMG docs). Isolation must therefore use AMG's own primitives: per-BU data sources, data-source permissions, Teams, and folders.

Isolation is enforced in three combined layers:

1. **Grafana data source permissions (primary enforcement).** One data source per BU (e.g. `amp-bu1`, `amp-bu2`), each restricted so only that BU's Grafana team can query it. `grafana_data_source_permission` manages the **entire** permission set for a data source, so the default-deny is achieved by granting `Query` to only the BU's team and deliberately **not** granting the built-in Viewer/Editor basic roles. This removes default query access for everyone else — every other BU's users — in dashboards and in Explore.
2. **Per-BU IAM assume-roles (AWS-layer defence in depth).** Each BU data source uses a distinct cross-account IAM role scoped to that BU's AMP/CloudWatch, rather than one broad role over all backends. Delivered by cloud-platform#8517 (production form).
3. **Folders + Viewer role (visibility hygiene).** Per-BU folders and read-only (Viewer) access so users see only their dashboards and cannot author dashboards pointing at other data sources.

The binding chain is: **IAM Identity Center group(s) (per BU) → Grafana team → {folder permission + data source permission}**.

**Many teams per BU.** A BU has many delivery teams, not one. Every identity group granted access to any namespace in a BU should get Grafana access to that BU's metrics, so each BU's Grafana team syncs from a **list** of IdC groups (`grafana_team_external_group.groups` is a list; deduplicated). BUs are configured by group **name** as those appear in `product.yaml`; names are resolved to IdC group IDs (AMG team sync matches on group ID) via the plural `data.aws_identitystore_groups` data source (ListGroups) using the read-only Identity Center provider — no hardcoded IDs. Deriving each BU's group list automatically from the `access[].group` entries across its `product.yaml` files is expected follow-up work. Isolation granularity is the BU; namespace-level isolation within a BU is not achievable with data-source permissions (that would need LBAC).

**Data source type:** AMG v12+ uses the dedicated Amazon Managed Service for Prometheus plugin (`grafana-amazonprometheus-datasource`); SigV4 is built in and signed with the AMG workspace's own IAM role, so no credentials are configured on the data source.

**Implementation dependency:** the AWS provider manages the AMG *workspace*, the service account, and IAM roles, but Grafana-internal objects (teams, folders, data sources, and data-source permissions) are managed through the **Grafana Terraform provider** (`grafana/grafana ~> 4.46`), which authenticates into the workspace with a **short-lived service account token minted per pipeline run** (AMG tokens are capped at 30 days, so nothing is stored). The service account itself is Terraform-managed at root; only the token is per-run. This is the same provider and token mechanism dashboards-as-code (cloud-platform#8510) will use.

**Known caveat (permissions):** the `github-actions-plan` role does not currently hold `grafana:CreateWorkspaceServiceAccountToken` (only the apply role does), so at plan time the Grafana provider is unauthenticated; the token-mint pipeline step is tolerant of this and does not fail the plan job. Granting the plan role that permission is a modernisation-platform-framework change tracked separately.

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

This covers both the **metrics-backend** alternatives (A chosen over B, C, D) and the **collector** alternatives (ADOT chosen over the managed collector and the CloudWatch agent).

### Metrics backend

#### Option D: CloudWatch Container Insights + AMG — REJECTED as primary (retained behind a flag)

**Why Rejected:** It was the challenger and ran in the PoC alongside Option A. Its genuine advantages (fewer moving parts, no AMP to provision, CloudWatch Alarms with no 100-rule limit) did not outweigh Option A's: PromQL is familiar to MoJ engineers so migration friction is low; the cost difference was minor against total observability spend (dominated by logging); and CloudWatch's query/dashboard model is less expressive and makes custom application metrics (EMF/PutMetricData) less natural than Prometheus scraping. Not discarded entirely — kept behind `enable_cloudwatch_observability` for AWS-native signals (e.g. RDS status) and the "explore related" root-cause view, usable selectively by the platform team.

**Pros:** fewest moving parts; no AMP/ADOT; no alert-rule limit; single metrics+logs pipeline in CloudWatch.
**Cons:** no PromQL; less expressive dashboards; EMF/API for custom metrics; per-metric cost high at high cardinality; metrics land in CloudWatch, not the Prometheus-native store the team wanted.

#### Option B: Self-Managed Prometheus + Grafana — REJECTED

**Why Rejected:** The current self-managed stack costs ~£7,000/month and is "compute management problematic." Repeating the self-managed model perpetuates the operational burden; AWS-managed alternatives meet the requirements at comparable cost. (AMP removes the Prometheus-*server* burden specifically — the part that hurt — while ADOT keeps only the lightweight collector in-cluster.)

**Pros:** no alert-rule limits; full feature control; no Identity Centre dependency for Grafana.
**Cons:** ~£7,000/month confirmed pain point; platform team owns upgrades, HA, storage, scaling; multi-cluster needs Thanos/Cortex.

#### Option C: Datadog — NOT APPLICABLE (team-level tool)

**Why Not Applicable:** Datadog is a team-level choice, not a platform observability decision. Teams may keep using it alongside the platform stack. No platform licensing evaluation required.

### Collector (how metrics reach AMP)

#### AMP agentless managed collector — REJECTED

**Why Rejected:** metrics-only and pull-only, so pushed OTLP tenant metrics would need a second in-cluster collector anyway, defeating the point of offloading. Per-collector charge (~$32/cluster/month, ~$7.8k/yr across 22 clusters, mostly fixed). Saves no scrape-config effort (same Prometheus config) and still needs node-exporter and kube-state-metrics in-cluster. Soft limit of 10 scrapers/region/account would bite in the ephemeral-cluster dev account.

**Pros:** AWS owns collector scaling/HA/upgrades; default config includes control-plane jobs; no collector pods on nodes.
**Cons:** per-collector charge; pull-only (no OTLP); metrics-only (no traces); doesn't reduce config or remove in-cluster exporters.

#### CloudWatch agent — REJECTED (for metrics)

**Why Rejected:** sends metrics to CloudWatch, not the selected AMP store; would still need an AMP collector for the Prometheus path, i.e. two agents per node. May re-enter scope later purely for **traces** to CloudWatch/X-Ray, independent of the metrics collector.

**Pros:** AWS-managed; can also handle traces and logs.
**Cons:** wrong destination for primary metrics; its log path is itself Fluent Bit routing only to CloudWatch, so it does not replace the Fluent Bit pipeline in ADR-017.

#### Self-managed ADOT — CHOSEN

Single extensible collector covering pull `/metrics`, pushed OTLP tenant metrics, and (later) traces + span-derived metrics; vendor-neutral (OpenTelemetry), no per-collector charge. Trade-off: the team owns the collector lifecycle (see "Collector decision" above).

---

## Rationale

The PoC ran AMP (Option A) and CloudWatch Container Insights (Option D) in parallel on separate clusters. We chose **AMP + AMG + ADOT, collected by self-managed ADOT**, because:

1. **Familiarity and low migration friction.** Engineers know Prometheus/PromQL; workflows and mental models don't change. This was the decisive factor over Option D.
2. **AMP removes the pain that mattered.** The historical Prometheus burden was running the *server* at scale (storage, retention, HA) — AMP removes exactly that. The collector is the small, remaining in-cluster piece.
3. **Cost was not decisive.** Measured cardinality (~16k active series/cluster) showed the AMP-vs-CloudWatch metrics cost difference is minor against total observability spend, which is dominated by logging.
4. **ADOT is one component that grows with us.** It covers both tenant metric patterns (pull `/metrics` and pushed OTLP) today, and can later add traces and span-derived metrics — avoiding a second collector. The managed collector and the CloudWatch agent each would have forced a second agent for one of these needs.
5. **Self-managed is rejected** — repeating the ~£7,000/month self-managed Grafana/Prometheus model perpetuates the operational burden the platform is trying to shed; managed AMP/AMG meets the requirements at comparable cost.
6. **Datadog is out of platform scope** — a team-level choice, not a platform decision.

CloudWatch Container Insights is retained behind a feature flag for AWS-native signals and selective root-cause analysis, but is not the primary metrics store. Alerting is Prometheus + AlertManager. This decision feeds the Phase 2 go/no-go recommendation (US-008).

---

## Consequences

### Positive
- No Prometheus-server or Grafana infrastructure to manage (AMP + AMG are managed); cross-cluster unified view via AMG with Identity Centre SSO.
- PromQL compatibility — engineers' existing queries and skills port directly.
- ADOT is vendor-neutral (OpenTelemetry) and extensible: one collector can scrape `/metrics`, receive OTLP tenant metrics, and later carry traces + span metrics — avoiding a second agent and avoiding pipeline lock-in.
- No AMG alert-rule limit in play — alerting is Prometheus + AlertManager.
- CloudWatch remains available (behind a flag) for AWS-native signals and selective root-cause analysis, without being the primary store.

### Negative / cost of ownership
- The team owns the ADOT collector lifecycle: upgrades alongside EKS, occasional collector-config deprecations, HA, scrape de-duplication, and memory tuning (estimated a few engineer-days/year, plus the collectors' node CPU/memory).
- The PoC collector's 3× scrape duplication must be fixed before production rollout (cloud-platform#7867); node-exporter and kube-state-metrics must be deployed and per-add-on metrics enabled for full coverage.
- AMG carries many data sources (one per BU AMP), needing clear naming conventions and the per-BU data-source-permission isolation model.
- Alerting UX (how tenants configure alerts, secure webhook storage) is still to be designed (cloud-platform#7871), and the ADR/`amg.tf` AMG-alerting assumption needs reconciling.

### Neutral
- AMG is unchanged from the original plan — the team builds Grafana skills regardless.
- Tenant application metrics are a near-term requirement; the exposure contract and cardinality guardrails are tracked in cloud-platform#7907.

---

## Related Decisions

- **Depends On:** ADR-004 (IAM Identity Center — AMG SSO, OpenSearch SSO)
- **Related:** ADR-001 (EKS), ADR-002 (Argo CD), ADR-017 (Logging and Log Management)

### Implementation tracking (cloud-platform issues)

**PoC — closed (decision taken 2026-09-30):**
- **#8414** — observability metrics PoC parent epic (closed; Option A selected)
- **#8416** — Option A pipeline (closed; selected)
- **#8417** — Option D pipeline (closed; evaluated, not selected, retained behind flag)
- **#8418** — cost analysis (closed; cost not decisive)
- **#8509** — BU metric isolation mechanism (closed; validated in dev with simulated BUs — see the section above and modernisation-platform-environments#19165)
- **#8508, #8523** — AMG workspace / plan-role token permission (closed, enabling work)

**Production build — under US-107 (#8235) Managed Observability Stack:**
- **#7867** — platform metrics capture (production ADOT: fix duplication, node-exporter/KSM, add-on metrics, fleet-wide view)
- **#8517** — per-BU AMP workspaces + cross-account query roles
- **#8558** — consolidated `business-units.json` group mapping
- **#8560** — split AMG into live/non-live workspaces + ephemeral-cluster feature flag

**Production build — under US-013 (#8247) Direct Namespace Observability:**
- **#8561** — productionise BU isolation with real BU parent groups
- **#8510** — dashboards-as-code pipeline

**Related:**
- **#7871** — re-architect user alert configuration (Prometheus + AlertManager)
- **#7907** — tenant service application metrics (near-term requirement)

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
