# ADR-018: Deployment Model Flexibility — GitOps as Default, Push-Based CD as Accommodation

**Status:** Proposed
**Date:** 2026-07-24
**Decision Maker:** AWS ProServe / MoJ Principal Technical Architect
**Category:** Integration / Compute

---

## Context

**Problem Statement:** ADR-002 selected the EKS Capability for Argo CD as the platform's GitOps controller. Application engineering teams have raised a concern: adopting ArgoCD signals a move away from the push-based CI/CD pipelines they have already built and operate. The platform needs a position on whether GitOps is mandatory for all workloads, or whether teams can continue deploying with their existing pipelines.

**Requirements:** FR-014 (GitOps deployment model), FR-009 (security baseline enforcement), NFR-006 (multi-cluster fleet), NFR-007 (private cluster endpoints), and the implicit requirement of tenant team autonomy and low migration friction.

**Constraints:**
- The platform is multi-tenant across BU accounts; namespace-level guardrails (PSA restricted, default-deny NetworkPolicy, Gatekeeper constraints) must apply regardless of how workloads are deployed.
- Spoke cluster API endpoints are private (NFR-007). Any push-based pipeline needs a network and credential path to reach them.
- The platform team is small; every supported deployment path carries an ongoing operational cost.
- App teams have existing, working CI/CD investment that has organisational value.

---

## Decision

**We will make GitOps via ArgoCD the default and recommended deployment model, while supporting push-based deployment of workloads as a documented accommodation. The platform mandates GitOps ownership only for the namespace security baseline, not for workload deployment.**

The mandate boundary is drawn at the namespace baseline. The platform owns namespace creation, RBAC, NetworkPolicy, and quota via the `app-baseline` chart reconciled by ArgoCD (ADR-015). Within a provisioned namespace, teams choose how to deploy their workloads.

Three workload deployment paths are supported:

| Path | Namespace baseline | Workload deployment | Cluster credentials in team CI |
|------|-------------------|---------------------|-------------------------------|
| **A — GitOps monorepo (default)** | Platform via ArgoCD | ArgoCD from shared `container-platform-environments` repo | None |
| **B — GitOps own-repo** | Platform via ArgoCD | ArgoCD from the team's own repo (ApplicationSet points at it) | None |
| **C — Push to namespace** | Platform via ArgoCD | Team's existing pipeline (`helm`/`kubectl`) into the pre-provisioned namespace | Scoped, namespace-only |

Path A is recommended. Path B suits teams that want GitOps but not the shared-monorepo workflow. Path C is the accommodation for teams with mature push pipelines they do not wish to replace.

---

## GitOps + ArgoCD vs Traditional Push CD

A framing point that addresses the core concern: **GitOps replaces the CD half, not CI.** Teams keep their build, test, lint, scan, and image-publish pipelines. Only the final "apply to cluster" step changes — from a pipeline running `kubectl apply` (push) to ArgoCD reconciling the cluster to match Git (pull).

### Advantages of GitOps/ArgoCD

- **No cluster credentials in CI.** With the EKS-managed ArgoCD reaching spokes via EKS Access Entries, no BU pipeline holds a kubeconfig or cluster-admin credential. This removes a large credential-distribution attack surface across many tenant pipelines.
- **Drift detection and self-heal.** ArgoCD continuously reconciles desired vs live state; manual cluster edits are reverted. Push CD only knows state at deploy time.
- **Git as single source of truth.** Every change to running state is a reviewed commit; rollback is `git revert`.
- **Fleet-scale consistency with guardrails.** One hub manages many spokes with per-BU AppProjects enforcing which repos, clusters, and resource kinds are permitted (ADR-002). Replicating those guardrails across N independent push pipelines is materially harder.

### Costs of GitOps/ArgoCD

- **New tool and mental model.** Pull-based reconciliation, sync waves, health/sync status, and app-of-apps carry a genuine learning curve. The engineers' concern is legitimate.
- **Debugging indirection.** "I merged, why isn't it live?" now spans ArgoCD sync status, reposerver cache, and generator behaviour rather than one pipeline log.
- **Two-repo workflow (Path A).** App code and deployment manifests live in separate repos. Teams used to an integrated single-repo deploy feel the split. Path B mitigates this.
- **Imperative steps need patterns.** DB migrations, one-off jobs, and dynamic secret injection require sync hooks or External Secrets rather than an inline pipeline step.
- **Less imperative release control.** "Deploy exactly now, in this order, with this manual gate" needs extra tooling compared to a scripted pipeline.

### Honest summary

For a platform team operating a multi-tenant fleet, GitOps is close to industry standard and the security and consistency benefits are strong. For an individual app team with a working, well-understood pipeline, the marginal benefit is smaller and the switching cost is front-loaded. The decision balances fleet-level operability against team-level autonomy, which is why the platform mandates only the baseline and leaves workload deployment to the team.

---

## Accommodating Teams That Will Not Adopt ArgoCD

The critical architectural question is whether a non-ArgoCD team can still deploy to a private spoke cluster while the platform retains its guardrails. The answer is yes, via Path C.

**How Path C works:**
1. Platform provisions the namespace and baseline (RBAC, NetworkPolicy, quota) via ArgoCD — same as every other team.
2. The team receives a namespace-scoped credential (EKS Access Entry mapped to a namespace RoleBinding, or IRSA) — not cluster admin.
3. The team's existing pipeline deploys workloads into that namespace using their current `helm`/`kubectl` steps.

**What Path C gives up:** cluster credentials return to the team's CI (the main security benefit of GitOps is lost for that tenant), and drift becomes the team's responsibility since ArgoCD does not reconcile their workloads.

**Network path:** private spoke endpoints (NFR-007) mean a push pipeline needs connectivity. Options are self-hosted runners in a VPC with a route to the cluster, or EKS API access via a bastion/relay. This is an operational prerequisite the team must accept for Path C.

**What cannot be accommodated regardless of path:** the namespace security baseline (PSA restricted, default-deny NetworkPolicy, Gatekeeper constraints) is enforced at the Kubernetes admission layer, not by ArgoCD. A push pipeline is subject to the same admission control as a GitOps sync. Teams needing privileged workloads or baseline exemptions are having a platform-policy conversation, not a CI/CD one.

---

## Alternatives Considered

### Option 1: Mandate GitOps for all workloads — REJECTED

**Why Rejected:** Forces every team to rebuild working CD pipelines with no migration path, maximising switching cost and organisational friction for marginal benefit to teams that already deploy reliably. Risks adoption resistance that could stall the wider platform migration.

**Pros:** Uniform operating model; every workload gets drift detection and no cluster credentials in CI; simplest for the platform team to reason about.
**Cons:** High friction for teams with mature pipelines; ignores existing CD investment; adoption risk.

### Option 2: No deployment standard (teams do whatever) — REJECTED

**Why Rejected:** Without platform ownership of the namespace baseline, multi-tenant guardrails become inconsistent and unenforceable at fleet scale. Loses the central value of the hub/spoke model.

**Pros:** Maximum team autonomy; zero migration friction.
**Cons:** No consistent guardrail enforcement; no fleet-level visibility; undermines the rationale for a managed platform.

---

## Rationale

In the context of onboarding application teams with established push-based CI/CD onto a multi-tenant EKS fleet, facing the constraint that namespace-level security guardrails must hold for every tenant while team autonomy and low migration friction are valued, we decided to mandate GitOps only for the namespace baseline and offer three workload deployment paths (GitOps monorepo, GitOps own-repo, push-to-namespace), and rejected both a full GitOps mandate and a no-standard free-for-all, to achieve consistent fleet guardrails without forcing teams to abandon working pipelines, accepting that push-based tenants reintroduce cluster credentials into their CI and own their own drift, because the guardrails that justify the platform are enforced at the admission layer independent of deployment mechanism, so deployment choice can be delegated safely to teams.

---

## Consequences

### Positive
- Teams with mature push pipelines can onboard without rebuilding CD.
- Namespace security baseline is enforced uniformly regardless of deployment path.
- GitOps remains the recommended path, so most teams still gain drift detection and credential-free CI.
- Reduces adoption resistance and de-risks the wider migration.

### Negative
- The platform team supports more than one deployment path, increasing documentation and support surface.
- Path C tenants reintroduce scoped cluster credentials into their CI and lose ArgoCD drift protection.
- Path C requires a network path to private spoke endpoints, which each such team must provision and maintain.
- Mixed models make fleet-wide "what is deployed where" harder to answer for push-based tenants.

### Neutral
- The `enable_argocd` and per-team ApplicationSet mechanisms already support Paths A and B with no new platform code.
- Path C requires a documented namespace-scoped access pattern (EKS Access Entry or IRSA) and runner connectivity guidance.

---

## Unknowns and Assumptions

| Item | Type | Impact | Mitigation |
|------|------|--------|------------|
| Number of teams expected to choose Path C | Unknown | Medium — drives support cost | Survey app teams during onboarding; revisit if Path C demand is high |
| Namespace-scoped credential pattern for Path C | Assumption | Medium | Validate EKS Access Entry namespace-scope + RoleBinding in a PoC before offering Path C |
| Private-endpoint reachability from team runners | Assumption | Medium | Document supported runner/relay patterns; confirm with a Path C pilot team |
| Long-term drift on Path C namespaces | Risk | Low-Medium | Periodic platform audit of push-managed namespaces against baseline |

---

## Related Decisions

- **Depends On:** ADR-002 (ArgoCD GitOps fleet management), ADR-015 (Environments repo structure)
- **Related:** ADR-001 (EKS multi-cluster), ADR-004 (IAM Identity Center), ADR-009 (NetworkPolicy)

---

## Research Sources

1. ADR-002: GitOps Fleet Management — EKS Capability for Argo CD
2. ADR-015: Environments Repo Structure
3. Argo CD Documentation — https://argo-cd.readthedocs.io/
4. AWS Deep Dive: EKS Capability for Argo CD — https://aws.amazon.com/blogs/containers/deep-dive-streamlining-gitops-with-amazon-eks-capability-for-argo-cd/
5. OpenGitOps Principles — https://opengitops.dev/
