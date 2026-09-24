# ADR-002: GitOps Fleet Management — EKS Capability for Argo CD vs. Self-Managed Argo CD

**Status:** Accepted  
**Date:** 2026-05-07  
**Updated:** 2026-08-27 — recorded the two-hub (one per environment tier) model as the baseline, superseding the earlier single-hub design.  
**Updated:** 2026-09-22 — added the BU-user Argo CD UI access model (Identity Center role mappings plus per-BU AppProject roles), resolving cloud-platform#8548. The original ADR specified deployment-side tenant isolation but never defined how BU engineers reach the Argo CD UI.  
**Decision Maker:** AWS ProServe / MoJ Principal Technical Architect  
**Category:** Compute / Integration

---

## Context

**Problem Statement:** The multi-cluster, multi-account architecture requires centralised GitOps deployment orchestration. The platform team cannot manually push deployments to 20+ private-endpoint EKS clusters in separate accounts. A hub-and-spoke Argo CD model is required to provide unified visibility and declarative state management.

**Requirements:** FR-014 (GitOps deployment model), FR-010 (automated cluster provisioning), NFR-006 (multi-cluster), NFR-007 (private cluster endpoints), NFR-002 (GitHub Actions integration)

**Constraints:**
- Cluster API endpoints are private (NFR-007) — Argo CD must reach spoke clusters without public endpoint or complex network configuration
- Cross-account access is required (spoke clusters in BU accounts, hub in management account)
- No Lambda (NFR-003)
- GitHub is the source control platform (NFR-002)
- EKS Capability for Argo CD was GA in November 2025 — not available when the candidate baseline PDF was written

---

## Decision

**We will use the EKS Capability for Argo CD as the managed GitOps controller, deployed as two hubs — one per environment tier (non-live and live) — with per-BU AppProject-based separation enforcing tenant isolation within each tier.**

This provides AWS-managed HA, upgrades, and native cross-account private cluster connectivity via EKS ARN + Access Entries — without VPC peering or TGW routing between accounts. Each hub runs a single Argo CD Capability instance (EKS limits one Argo CD Capability per cluster) and manages only the spokes in its own tier: the `cloud-platform-nonlive` hub manages non-live spokes, and the `cloud-platform-live` hub manages live spokes. Environment isolation is therefore structural — live and non-live credentials never co-reside on the same hub compute. Within each tier, a per-BU AppProject model provides tenant isolation, where each project's destination allowlist restricts deployments to a single cluster ARN.

### AppProject Hierarchy

The AppProject model provides three layers of isolation: environment separation (live vs. non-live), BU-to-BU isolation, and platform vs. workload separation.

**Platform projects (2 total):**

| AppProject | Destinations | Sources | Sync Policy | Purpose |
|------------|-------------|---------|-------------|---------|
| `platform-live` | All live cluster ARNs (platform namespaces only) | Platform infrastructure repos | Manual | Platform add-ons (Gatekeeper, ADOT, Fluent Bit, ExternalDNS, LBC) on live clusters |
| `platform-nonlive` | All non-live cluster ARNs (platform namespaces only) | Platform infrastructure repos | Auto-sync | Same platform stack on non-live clusters |

**BU workload projects (2 per BU = ~16-18 total):**

| AppProject | Destinations | Sources | Sync Policy | Purpose |
|------------|-------------|---------|-------------|---------|
| `bu1-live` | BU1 live cluster ARN only, workload namespaces | `github.com/ministryofjustice/bu1-*` | Manual | BU1 application deployments to live |
| `bu1-nonlive` | BU1 non-live cluster ARN only, workload namespaces | `github.com/ministryofjustice/bu1-*` | Auto-sync | BU1 application deployments to non-live |
| `bu2-live` | BU2 live cluster ARN only | `github.com/ministryofjustice/bu2-*` | Manual | BU2 live deployments |
| `bu2-nonlive` | BU2 non-live cluster ARN only | `github.com/ministryofjustice/bu2-*` | Auto-sync | BU2 non-live deployments |
| ... | ... | ... | ... | Same pattern for all 8-9 BUs |

**Total: ~20 AppProjects** (2 platform + 8-9 BUs × 2 environments)

### BU Isolation Guardrails

Each BU AppProject enforces:

| Guardrail | Mechanism | Effect |
|-----------|-----------|--------|
| Cluster isolation | `spec.destinations` restricted to single cluster ARN | BU1 project cannot deploy to BU2's cluster |
| Source isolation | `spec.sourceRepos` restricted to BU's repos | BU1 cannot reference BU2's manifests |
| Privilege restriction | `spec.namespaceResourceWhitelist` | BU teams can only create workload resources (Deployment, Service, Ingress, ConfigMap, Secret, HPA) |
| Escalation prevention | `spec.clusterResourceBlacklist` | BU teams cannot create ClusterRole, ClusterRoleBinding, Namespace, or CRDs |

### Multi-Namespace Support Within a BU

Each BU has multiple application teams, each with one or more namespaces within the BU's cluster. Argo CD uses a Git directory-based ApplicationSet to generate one Application per namespace:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: bu1-live-apps
  namespace: argocd
spec:
  generators:
    - git:
        repoURL: https://github.com/ministryofjustice/bu1-deployments
        revision: main
        directories:
          - path: "apps/*"
  template:
    metadata:
      name: "bu1-live-{{path.basename}}"
    spec:
      project: bu1-live
      source:
        repoURL: https://github.com/ministryofjustice/bu1-deployments
        targetRevision: main
        path: "{{path}}"
      destination:
        server: "arn:aws:eks:eu-west-2:111111111111:cluster/bu1-live"
        namespace: "{{path.basename}}"
      syncPolicy:
        automated: false  # Manual sync for live
```

Adding a new app = adding a directory to the BU's deployment repo. The ApplicationSet generates the Argo CD Application automatically.

**Namespace isolation within a BU** is enforced by Kubernetes (not Argo CD AppProjects):
- EKS Access Entries map team identity to namespace-scoped RoleBindings
- Gatekeeper admission policies block cross-namespace resource creation
- Git-level access control (CODEOWNERS, branch protection) prevents Team A from modifying Team B's manifests

### Isolation Layers (Defence in Depth)

```
Layer 1: AWS Account boundary     BU1 live account ≠ BU2 live account (hard isolation)
Layer 2: AppProject destinations  bu1-live can ONLY reach BU1's live cluster ARN
Layer 3: AppProject sources       bu1-live can ONLY pull from BU1's repos
Layer 4: Kubernetes RBAC          Team A can only modify resources in their namespace
Layer 5: Gatekeeper admission     Blocks cross-namespace resource creation
Layer 6: Git access control       CODEOWNERS prevents Team A editing Team B's manifests
```

### BU-User Argo CD UI Access (cloud-platform#8548)

The isolation layers above govern the **deployment path** — what a BU's Argo CD Applications may pull and where they may deploy. They say nothing about the **human path**: which people can open the Argo CD UI, and what they see once inside. The original ADR assumed BU teams interact only through Git and `kubectl`, so no UI access model was defined. In practice BU engineers need read-only visibility of their own deployments in the UI. Access was only ever exercised as `cloud-platform-engineers` and `container-platform-aws`, which both hold default ADMIN, so the gap went unnoticed until a BU user reported an empty Applications list.

The EKS managed Argo CD capability exposes access as two independent layers (see [Configure Argo CD permissions](https://docs.aws.amazon.com/eks/latest/userguide/argocd-permissions.html)). Both are required — either alone yields no useful access:

| Layer | Control | Grants | Where configured |
|-------|---------|--------|------------------|
| UI entry | Capability `rbacRoleMappings` (global `VIEWER`) | Log in to the Argo CD UI. On its own shows **zero** Applications — global roles gate Argo CD itself, not project-scoped resources. | `cluster` component (`argocd_rbac_role_mappings`) |
| App visibility | Per-BU AppProject `roles` (`applications, get, <project>/*`) | See that BU's Applications, read-only, scoped to the BU's own project. | `cluster-core` component (generated per BU AppProject) |

Because visibility is granted per AppProject, a BU only ever sees its own Applications — cross-BU isolation in the shared hub UI is a consequence of the layer-2 grant, not a separate control. Read-only (`VIEWER` + `applications, get`) is deliberate: BU teams observe deployments but do not sync, create, or delete from the UI (GitOps drives changes through Git).

Verified on the non-live hub (2026-09-22): a BU group mapped to global `VIEWER` with no project role sees an empty UI; adding the project `roles` grant makes exactly that BU's Applications visible and no others.

### BU-to-Identity-Group Mapping

Both layers key on a **stable per-BU IAM Identity Center group**, not the per-squad groups that `product.yaml` uses for AWS/EKS access. This is the deliberate asymmetry that keeps the Argo CD configuration static:

- **AWS account / EKS access** is driven by `product.yaml access[].group` — squad-level groups (e.g. `hmpps-sre`), scoped per-cluster and per-namespace, churning as teams onboard.
- **Argo CD UI access** is driven by the **BU parent group** — one per BU, read-only, BU-wide.

The BU parent groups are the children of the `business-units` GitHub team. GitHub surfaces a nested child team's members in the parent team's membership (the team members API marks them `inherited=true`, distinct from directly-added members), and it is the parent group's full membership that syncs to Identity Center. So no per-team churn reaches the Argo CD config: a squad team nested under its BU parent contributes its members to the parent automatically, and the only Argo CD change is onboarding a brand-new BU, which already requires a `bu_configs` edit.

**Operational dependency — squad teams must be nested under their BU parent.** This model relies on that nesting being done. It is the established org convention (verified: `hmpps-developers` includes both directly-added members and `inherited=true` members surfaced purely through child teams), but it is not automatic. A squad team created as a standalone team, not nested under its BU parent, still gets AWS/EKS access through `product.yaml` (which keys on the squad group directly) but its members will **not** appear in the BU parent group, so they get no Argo CD UI access until the team is nested. This is the one recurring manual step: when creating a new squad team, nest it under its BU parent. It is easy to forget precisely because the AWS-access path does not depend on it.

#### Groups in use (onboarded BUs)

These four BUs have `bu_configs` entries and clusters today. The `VIEWER` mapping and AppProject role are wired for each:

| BU code | Clusters | BU parent group (GitHub team → IdC) | Evidence |
|---------|----------|-------------------------------------|----------|
| `octo` | `container-platform-octo-{nonlive,live}` | `office-of-the-cto` | Child teams are `octo-*` (`octo-engineering`, `octo-architecture`, `octo-hosting`). OCTO = Office of the CTO. |
| `laa` | `container-platform-laa-{nonlive,live}` | `laa` | Direct name match; child teams `laa-*`. |
| `hmpps` | `container-platform-hmpps-{nonlive,live}` | `hmpps-developers` | Child teams `hmpps-*` (verified `hmpps-sre` members ⊆ `hmpps-developers`). |
| `cd` | `container-platform-cd-{nonlive,live}` | `central-digital` | Only `central-digital-*` BU team space; no other `cd-*` candidate. |

**Reviewers: please confirm the `octo → office-of-the-cto` and `cd → central-digital` mappings.** `hmpps` and `laa` are unambiguous; the other two are inferred from child-team naming and should be validated against the actual team a BU's engineers belong to.

#### Anticipated groups (not yet onboarded)

The remaining `business-units` children are candidate BU parent groups for future onboarding. They are recorded here so a wrong pick is caught at review time, but are **not** wired until the BU has `bu_configs` clusters:

| BU parent group | Members (2026-09-22) | Notes |
|-----------------|----------------------|-------|
| `opg` | 44 | Office of the Public Guardian — likely a future BU. |
| `technology-services` | 79 | Central technology services. |
| `platforms` | 30 | Platforms group; overlaps the platform team — confirm it is a tenant BU, not platform staff, before mapping. |

Group IdC IDs are not hardcoded in the ADR; the `cluster-core` layer-2 grant resolves BU parent group names to IdC IDs at plan time (Identity Center `ListGroups`), and the `cluster` layer-1 mapping carries the resolved IDs. Names, not IDs, are the reviewable source of truth.

### Onboarding a New BU

When a new BU is onboarded (US-011), the following is automated via Terraform:
1. Provision BU's live and non-live AWS accounts + EKS clusters
2. Create two AppProjects (`buN-live`, `buN-nonlive`) with correct cluster ARN destinations and source repo restrictions
3. Register new cluster secrets in the hub cluster
4. Platform ApplicationSet (cluster generator) automatically deploys platform add-ons to new clusters
5. BU-specific ApplicationSet created to generate per-namespace Applications from the BU's deployment repo
6. Add the BU's parent group to the `VIEWER` role mapping (layer 1, `cluster` component) and add the read-only `roles` grant to the BU's AppProject (layer 2, `cluster-core` component), so BU engineers get read-only UI visibility of their own Applications

Adding the BU is the only Argo CD change; no per-team or per-namespace Argo CD configuration is required as squads onboard.

### EKS Capability Constraint

**Hard limit:** EKS allows only one Argo CD Capability instance per cluster ([source](https://docs.aws.amazon.com/eks/latest/userguide/capabilities.html)). Because a single instance cannot keep live and non-live IAM credentials on separate compute, the platform runs two hubs — one per environment tier — rather than one hub spanning both. This gives each tier its own Argo CD Capability instance and guarantees that live and non-live credentials are never co-located. The two hubs are the baseline design (see Implementation Notes → "Two-Hub Model").

---

## Research Conducted

### Option A: EKS Capability for Argo CD ✓ SELECTED

**Research Confidence:** High

| Source | URL | Key Finding |
|--------|-----|-------------|
| AWS Deep Dive Blog | https://aws.amazon.com/blogs/containers/deep-dive-streamlining-gitops-with-amazon-eks-capability-for-argo-cd/ | AWS-managed; cross-account private cluster access native via EKS ARN |
| EKS Capability User Guide | https://docs.aws.amazon.com/eks/latest/userguide/argocd.html | GA November 2025; all standard regions including eu-west-2 |
| EKS Comparison Guide | https://docs.aws.amazon.com/eks/latest/userguide/argocd-comparison.html | Comparison with self-managed; trade-offs documented |
| AWS What's New | https://aws.amazon.com/about-aws/whats-new/2025/11/amazon-eks-capabilities/ | GA announcement: EKS Capabilities (Argo CD, KRO, ACK) |

**Capabilities Verified:**
- Runs in AWS-managed infrastructure outside customer worker nodes — verified
- Cross-account private cluster access via EKS Access Entries (no VPC peering) — verified
- ApplicationSets supported — verified
- AppProjects supported — verified
- IAM Identity Center required for authentication — verified (constraint documented)
- CodeConnections for GitHub integration — verified
- eu-west-2 available — verified

### Option B: Self-Managed Argo CD on Hub Cluster

**Research Confidence:** High

| Source | URL | Key Finding |
|--------|-----|-------------|
| Candidate Baseline PDF | using-amazon-eks-capabilities...pdf | Hub cluster pattern; self-managed Argo CD; Platform + Workloads project separation |
| Argo CD Docs | https://argo-cd.readthedocs.io/ | HA requires 3-replica setup; manual upgrade lifecycle |
| EKS Fleet Management Research | eks-fleet-management-research.md (§5.1) | Self-managed requires VPC peering or TGW for cross-account private cluster access |

**Capabilities Verified:**
- Full Argo CD AppProject flexibility (multi-namespace) — verified
- Platform team manages upgrades, HA, SSO — confirmed (operational overhead)
- Cross-account private endpoint access: requires VPC peering or TGW + IAM role chaining — confirmed (complex)
- Flexible authentication (OIDC, LDAP, Dex) — verified

---

## Capability Mapping

| Requirement | Option A (EKS Capability) | Evidence | Option B (Self-Managed) | Evidence |
|-------------|--------------------------|----------|------------------------|----------|
| FR-014: Centralised GitOps across private clusters | Native cross-account via EKS ARN + Access Entries | AWS Deep Dive Blog | Requires TGW/VPC peering + role chaining | EKS Fleet research §5.1 |
| NFR-007: Private cluster endpoint support | Native — AWS manages connectivity | EKS User Guide | Complex — VPC peering or TGW required | Research §5.1 |
| NFR-006: Multi-account fleet management | ApplicationSets per cluster; cluster secrets registered by ARN | AWS Blog | Same ApplicationSets; additional network config | Candidate baseline |
| FR-009: Security baseline (AppProject guardrails) | AppProjects supported (single-namespace constraint) | EKS User Guide | Full AppProject flexibility, multi-namespace | Argo CD docs |
| NFR-002: GitHub Actions integration | CodeConnections (GitHub App); no runner access to spokes | AWS Blog | Custom SSH/token secrets; runners need spoke access for CD | Research §5.4 |

---

## Unknowns and Assumptions

| Item | Type | Impact | Mitigation |
|------|------|--------|------------|
| Single-namespace constraint for Argo CD CRs | Constraint | Medium — candidate baseline uses multi-namespace pattern | All Argo CD CRs in one namespace is sufficient for MoJ; AppProjects provide separation |
| One Argo CD Capability per cluster (EKS hard limit) | Constraint | Low — resolved by the two-hub model | Two hubs (one per tier) each run their own Capability instance, so live and non-live credentials are structurally isolated; AppProject destination allowlists provide defence in depth within a tier |
| GitHub Actions integration with private hub cluster | Assumption | Medium — hub Argo CD API must be accessible | Self-hosted runners in management VPC reach hub Argo CD API; spoke clusters not accessed directly |
| EKS Capability Argo CD pricing at scale | Unknown | Low-Medium — per Application pricing | Confirm with AWS account team; ~20-30 Applications per cluster × 22 clusters = budgetable |
| AppProject misconfiguration risk | Risk | Medium — a misconfigured AppProject could allow cross-environment deployment | Terraform-managed AppProject definitions (not manually edited); PR review for AppProject changes; PoC validates guardrails |

---

## Counter-Argument Analysis

### Q1: What specific evidence would make me choose the alternative instead?

If MoJ's AppProject requirements specifically required multiple Argo CD namespaces (multi-namespace Argo CD instance), self-managed would be necessary. The EKS Capability single-namespace constraint limits some advanced multi-tenancy patterns. However, AppProjects with source-repo guardrails satisfy MoJ's separation requirement within the single-namespace model.

### Q2: Is there a managed/integrated AWS service that does this better?

EKS Capability for Argo CD IS the managed AWS service for this capability. It was specifically designed for the hub-and-spoke multi-account use case.

### Q3: What am I NOT seeing about the alternative that might make it superior?

Self-managed Argo CD provides full Argo CD feature parity without the single-namespace constraint and allows custom Dex-based authentication if IAM Identity Center is not the preferred authentication mechanism. If MoJ wanted to use a non-IAM Identity Center OIDC provider for Argo CD, self-managed would be required.

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

### Option B: Self-Managed Argo CD — REJECTED

**Why Rejected:** Requires complex cross-account networking (VPC peering or TGW) to reach private spoke cluster endpoints. The platform team would own Argo CD upgrades, HA, and SSO configuration — significant operational overhead for a small team. EKS Capability eliminates this complexity at the cost of IAM Identity Center dependency, which is required anyway for other platform services.

**Pros:** Full Argo CD feature flexibility; multi-namespace AppProject support; authentication flexibility
**Cons:** VPC peering or TGW required per spoke account for private endpoints; platform team manages upgrades, HA, and SSO; additional operational burden for small team

---

## Rationale

In the context of managing GitOps deployments to 20+ private-endpoint EKS clusters across multiple AWS accounts, facing the constraint that cluster API endpoints must be private (NFR-007) and operational overhead must be minimised for a small platform team, we decided for EKS Capability for Argo CD and rejected self-managed Argo CD, to achieve native cross-account private cluster connectivity without VPC peering and eliminate GitOps infrastructure management overhead, accepting the single-namespace constraint and IAM Identity Center dependency, because the operational complexity of self-managed Argo CD with cross-account private endpoint networking is disproportionate to the team's capacity and the managed capability directly solves the private endpoint connectivity problem.

---

## Consequences

### Positive
- No VPC peering or TGW route between hub and spoke accounts for GitOps traffic
- AWS manages Argo CD HA, upgrades, and scaling
- Native CodeConnections GitHub integration (no webhook credential management)
- Private spoke cluster access from day one without network configuration
- AppProject destination allowlists enforce live/non-live separation at the Argo CD level
- Unified deployment visibility across all environments from a single control plane

### Negative
- IAM Identity Center required (creates dependency on Identity Center being operational)
- Single-namespace constraint for Argo CD CRs (deviation from candidate baseline multi-namespace pattern)
- Per-Application pricing, across two hubs (predictable but adds cost at scale)
- Two hubs to operate and upgrade instead of one (marginal, since the Capability is AWS-managed)
- Within a tier, a misconfigured AppProject could still cross BU boundaries (mitigated by Terraform-managed definitions and PR review); it cannot cross the live/non-live boundary because that is enforced by separate hubs
- BU-user UI access is governed by BU parent GitHub team membership, a control surface outside the `product.yaml` review process that gates AWS/EKS access. A person in a BU parent group but in no `product.yaml` squad gets read-only UI visibility of that BU's Applications without any cluster access. Accepted for a read-only dashboard; the two access decisions are intentionally made at different group tiers (see BU-to-Identity-Group Mapping)

### Neutral
- Platform and Workloads AppProject separation still achievable within single namespace
- ApplicationSets for fleet-level automation supported natively
- Live and non-live are structurally isolated by running one hub per tier (the two-hub model)
- BU-user Argo CD access is read-only (`VIEWER` + per-project `applications, get`); write actions remain GitOps-driven through Git

---

## PoC Validation Approach (Updated — Iteration 3, 3 July 2026)

The ArgoCD hub-and-spoke model was validated using an ephemeral hub/spoke pair before the permanent per-tier hubs (`cloud-platform-nonlive`, `cloud-platform-live`) were provisioned:
- **Hub cluster:** Development cluster in the `cloud-platform-development` account
- **Spoke cluster:** `container-platform-octo-nonlive` cluster

This validates:
1. Cross-account private cluster access via EKS Access Entries
2. ApplicationSet generation from the environments repo (ADR-015)
3. AppProject boundaries enforcing BU isolation
4. Platform add-on deployment to spoke clusters

The dev cluster is purpose-built for validation and can be destroyed/recreated via the `development-cluster-deploy.yml` GitHub Actions workflow in `cloud-platform-github-workflows`.

---

## Implementation Notes (US-015a — 5 July 2026)

The following design refinements were made during implementation:

### Spoke Access Model — Spoke-Driven Registration

The hub cluster creates a `<cluster-name>-argocd-spoke-access` IAM role with a trust policy scoped to the EKS service and the hub cluster ARN. The hub does NOT maintain a list of spoke account IDs or cluster ARNs, and does NOT proactively grant `sts:AssumeRole` or `eks:DescribeCluster` to spoke accounts.

Instead, spoke registration is entirely spoke-driven: when a spoke cluster is provisioned, it adds the hub's spoke access role ARN as an EKS Access Entry with `AmazonEKSClusterAdminPolicy`. This means:
- The hub cannot reach any cluster it hasn't been granted access to
- A random cluster cannot register itself — someone must deliberately add the hub role as an Access Entry on the spoke side (controlled via Terraform PRs)
- No spoke account list to maintain on the hub

### CodeConnections — Optional for Public Repos

MoJ repositories are public. ArgoCD can pull manifests from public repos over HTTPS without authentication. CodeConnections is included as an optional enhancement (defaults to empty/disabled) for production use, providing:
- Webhook-based notifications (immediate sync vs. polling)
- Higher GitHub API rate limits (5,000/hr authenticated vs. 60/hr unauthenticated)
- Commit status write-back

CodeConnections is only required in the hub account.

### IDC Instance Configuration

The IAM Identity Center instance ARN is passed as a GitHub Actions secret (`ARGOCD_IDC_INSTANCE_ARN`) rather than discovered dynamically via a Terraform data source. The `aws_ssoadmin_instances` data source requires `sso:ListInstances` permission which is only available in the management account — the workflow role in `cloud-platform-development` does not have cross-account access to query this. The IDC instance ARN is static infrastructure that rarely changes.

### Two-Hub Model (Decision — supersedes the earlier single-hub baseline)

The team originally considered a single hub managing both live and non-live environments. That raised a specific concern: platform engineers testing ArgoCD configuration changes (ApplicationSets, AppProjects, RBAC) could accidentally affect live clusters, and the mitigations (AppProject destination allowlists, manual sync for live) were procedural (PR review) rather than structural.

**Decision:** the platform runs two hubs, one per environment tier, giving a structural guarantee of isolation. This is now the implemented baseline, not an escalation path.

- **`cloud-platform-nonlive`** — the non-live hub, in the `cloud-platform-nonlive` account. Registers only non-live spokes.
- **`cloud-platform-live`** — the live hub, in the `cloud-platform-live` account. Registers only live spokes.

Each spoke resolves its hub by tier and a hub only registers spokes in its own tier, so live and non-live credentials are never co-located on the same hub compute. The two hubs are declared in `local.argocd_hubs` (Terraform, `cluster/locals.tf`), and each spoke derives its tier hub's Argo CD Capability role ARN by convention — no per-spoke input is required. The infrastructure cost of the second hub is minimal since the Capability runs in AWS-managed infrastructure.

### Hub Enablement Mechanism

ArgoCD is enabled on a cluster via `TF_VAR_enable_argocd=true` passed at deploy time through the `development-cluster-deploy` workflow. The module uses `count = var.enable_argocd ? 1 : 0` — when disabled (default), zero ArgoCD resources are created. This allows the same generic cluster module to provision both hub and spoke clusters.

---

## Related Decisions

- **Depends On:** ADR-001 (EKS multi-cluster), ADR-004 (IAM Identity Center)
- **Influences:** ADR-003 (KRO/ACK), ADR-015 (Environments repo structure)
- **Related:** ADR-007 (DNS per cluster)

---

## Research Sources

1. AWS Deep Dive: EKS Capability for Argo CD — https://aws.amazon.com/blogs/containers/deep-dive-streamlining-gitops-with-amazon-eks-capability-for-argo-cd/
2. EKS Argo CD User Guide — https://docs.aws.amazon.com/eks/latest/userguide/argocd.html
3. EKS Argo CD Comparison — https://docs.aws.amazon.com/eks/latest/userguide/argocd-comparison.html
4. AWS What's New: EKS Capabilities GA — https://aws.amazon.com/about-aws/whats-new/2025/11/amazon-eks-capabilities/
