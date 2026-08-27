# ADR-018: Deployment Model — GitOps Monorepo Mandated as the Target, Push-Based CD as a Transitional Exception

**Status:** Accepted
**Date:** 2026-07-24
**Updated:** 2026-08-27 — GitOps via ArgoCD from the shared monorepo (Path A) is now the mandated target deployment model. GitOps from per-team repos (Path B) is documented for potential future consideration but is not part of the current decision, owing to drift and per-repo operational overhead concerns. Push-based CD (Path C) is retained only as a time-boxed migration exception. To be revisited if teams face strong, well-evidenced opposition.
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

**We will mandate GitOps via ArgoCD from the shared monorepo (Path A) as the target deployment model for all workloads. The platform owns the namespace security baseline via ArgoCD in every case, and GitOps from the monorepo is the expected mechanism for workload deployment. Per-team GitOps repos (Path B) are documented for potential future consideration only. Push-based deployment (Path C) is supported only as a time-boxed migration exception for teams with existing pipelines, not as a permanent co-equal option.**

The platform owns namespace creation, RBAC, NetworkPolicy, and quota via the `app-baseline` chart reconciled by ArgoCD (ADR-015). The target state is that workloads are also deployed through ArgoCD from the shared `container-platform-environments` monorepo. A team may deploy with its own push pipeline only during an agreed migration window, under the exception process below.

Deployment paths:

| Path | Namespace baseline | Workload deployment | Cluster credentials in team CI | Status |
|------|-------------------|---------------------|-------------------------------|--------|
| **A — GitOps monorepo** | Platform via ArgoCD | ArgoCD from shared `container-platform-environments` repo | None | Target (current decision) |
| **B — GitOps own-repo** | Platform via ArgoCD | ArgoCD from the team's own repo (ApplicationSet points at it) | None | Documented for future consideration; not adopted |
| **C — Push to namespace** | Platform via ArgoCD | Team's existing pipeline (`helm`/`kubectl`) into the pre-provisioned namespace | Scoped, namespace-only | Transitional exception |

Path A is the decision: GitOps from the shared monorepo. Path B (per-team repos) is not adopted now — it risks configuration drift across many repos and adds per-repo operational overhead for a small platform team — but is documented below so it can be reconsidered later if there is a clear need. Path C is a transitional exception with a migration expectation, not a free choice a team can settle on indefinitely.

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
- **Two-repo workflow (Path A).** App code and deployment manifests live in separate repos. Teams used to an integrated single-repo deploy feel the split. Path B (per-team repos) could reduce this friction and is noted as a possible future option, but is not adopted now.
- **Imperative steps need patterns.** DB migrations, one-off jobs, and dynamic secret injection require sync hooks or External Secrets rather than an inline pipeline step.
- **Less imperative release control.** "Deploy exactly now, in this order, with this manual gate" needs extra tooling compared to a scripted pipeline.

### Summary

For a platform team operating a multi-tenant fleet, GitOps is the industry standard and the security and consistency benefits are strong. For an individual app team with a working, well-understood pipeline, the marginal benefit is smaller and the switching cost is front-loaded. This cost is a one-off migration expense rather than a reason to run two deployment models permanently: a single, uniform target keeps fleet-level operability simple for a small platform team. GitOps is therefore the mandated target, and the push-based path exists to smooth the transition, not to be a standing alternative.

---

## Transitional Push-Based Deployment (Path C — Exception)

GitOps is the target, but a team migrating from an established push pipeline may need time to get there. Path C exists for that transition. It is an exception, granted per team for a bounded period, not a permanent operating mode.

**Exception process:**
1. The team requests a Path C exception from the platform team, with a migration intent to Path A or B.
2. The platform team agrees a time-boxed window and records the exception (owning team, reason, target GitOps path, review date).
3. The exception is reviewed at the agreed date; it is renewed only if migration is genuinely blocked, otherwise the team is expected to be on GitOps.

**How Path C works during the window:**
1. Platform provisions the namespace and baseline (RBAC, NetworkPolicy, quota) via ArgoCD — same as every other team.
2. The team receives a namespace-scoped credential (EKS Access Entry mapped to a namespace RoleBinding, or IRSA) — not cluster admin.
3. The team's existing pipeline deploys workloads into that namespace using their current `helm`/`kubectl` steps.

**What Path C gives up:** cluster credentials return to the team's CI (the main security benefit of GitOps is lost for that tenant), and drift becomes the team's responsibility since ArgoCD does not reconcile their workloads. These are the reasons it is transitional rather than permanent.

**Network path:** private spoke endpoints (NFR-007) mean a push pipeline needs connectivity. Options are self-hosted runners in a VPC with a route to the cluster, or EKS API access via a bastion/relay. This is an operational prerequisite the team must accept for the duration of the exception.

**What is not negotiable regardless of path:** the namespace security baseline (PSA restricted, default-deny NetworkPolicy, Gatekeeper constraints) is enforced at the Kubernetes admission layer, not by ArgoCD. A push pipeline is subject to the same admission control as a GitOps sync. Teams needing privileged workloads or baseline exemptions are having a platform-policy conversation, not a CI/CD one.

---

## Alternatives Considered

### Option 1: GitOps as one of several co-equal, permanently-supported paths — REJECTED (superseded)

This was the earlier position in this ADR: GitOps recommended, but push-to-namespace offered as a standing accommodation a team could settle on indefinitely.

**Why Rejected:** Running two deployment models permanently doubles the platform team's support and documentation surface and leaves fleet-wide "what is deployed where" partially unanswerable for push-based tenants. Without a clear target, teams have no signal to converge, so the split becomes permanent by default. The platform wants a single, clear target model; the switching cost for teams with mature pipelines is real but is a one-off migration expense, not a reason to carry two models forever.

**Pros:** Lowest short-term friction; no team has to change anything.
**Cons:** Permanent dual operating model; higher ongoing platform cost; weaker fleet consistency and visibility; no convergence signal.

### Option 2: No deployment standard (teams do whatever) — REJECTED

**Why Rejected:** Without platform ownership of the namespace baseline, multi-tenant guardrails become inconsistent and unenforceable at fleet scale. Loses the central value of the hub/spoke model.

**Pros:** Maximum team autonomy; zero migration friction.
**Cons:** No consistent guardrail enforcement; no fleet-level visibility; undermines the rationale for a managed platform.

### Option 3: GitOps from per-team repos (Path B) — DEFERRED (documented for future consideration)

Under Path B, each team keeps its deployment manifests in its own repo and an ApplicationSet points ArgoCD at it, rather than everyone sharing the `container-platform-environments` monorepo. This keeps GitOps' security and drift-detection benefits while giving teams a single-repo feel.

**Why Deferred (not adopted now):** many per-team repos multiply the surface where configuration can drift between what each repo declares and the platform's baseline expectations, and each additional repo and its ApplicationSet wiring is ongoing operational overhead for a small platform team. The monorepo keeps configuration and review in one place. Path B is documented here so it can be reconsidered if a clear need emerges (for example, teams with strong single-repo requirements or scaling limits on the monorepo).

**Pros:** Single-repo developer experience while staying on GitOps; no cluster credentials in CI; keeps drift detection.
**Cons:** Drift risk spread across many repos; per-repo ApplicationSet and review overhead; weaker central visibility than a monorepo.

### Chosen: GitOps monorepo (Path A) mandated as the target, push-based (Path C) as a transitional exception

GitOps from the shared monorepo (Path A) is the single target model. Per-team GitOps repos (Path B) are deferred and documented for future consideration. Push-based deployment (Path C) is retained only as a time-boxed migration exception so teams with existing pipelines are not forced to rebuild everything on day one, but the expectation is convergence on Path A. This gives the clear target the platform wants while keeping a humane migration path.

**Revisit trigger:** the mandate will be reconsidered if teams face strong, well-evidenced opposition or genuine technical blockers to GitOps adoption; Path B may be revisited separately if a clear need for per-team repos emerges. Until then, the monorepo is the target and Path C is an exception, not a default.

---

## Rationale

In the context of onboarding application teams with established push-based CI/CD onto a multi-tenant EKS fleet, facing the constraint that namespace-level security guardrails must hold for every tenant and that a small platform team needs a single tractable operating model, we decided to mandate GitOps via ArgoCD from the shared monorepo (Path A) as the target deployment model, to defer per-team GitOps repos (Path B) to potential future consideration on drift and per-repo overhead grounds, and to retain push-based deployment (Path C) only as a time-boxed migration exception, and rejected both keeping push-based CD as a permanent co-equal path and a no-standard free-for-all, to give teams a clear target while still offering a humane migration route, accepting that transitional Path C tenants reintroduce scoped cluster credentials into their CI and own their own drift until they migrate, because a single target model keeps fleet consistency, security, and support cost manageable, with the mandate to be revisited if teams face strong, well-evidenced opposition.

---

## Consequences

### Positive
- A single, clear target model (GitOps from the shared monorepo) that teams can plan towards, rather than an open-ended choice.
- Keeping all deployment config in one monorepo bounds drift and review surface, versus spreading it across many per-team repos.
- Namespace security baseline is enforced uniformly regardless of deployment path.
- Most teams gain drift detection and credential-free CI; the fleet converges on one operating model over time.
- Teams with mature push pipelines still get a migration route (Path C) rather than a day-one cliff edge.

### Negative
- Teams with mature push pipelines face a one-off migration cost to reach the target model.
- During transition the platform team supports two deployment paths, increasing documentation and support surface until Path C exceptions close out.
- Path C tenants reintroduce scoped cluster credentials into their CI and lose ArgoCD drift protection for the duration of the exception.
- Path C requires a network path to private spoke endpoints, which each such team must provision and maintain while on the exception.
- Exceptions need tracking and periodic review, which is a small ongoing platform overhead.

### Neutral
- The ApplicationSet mechanism could support Path B (per-team repos) with little new platform code if it is adopted in future, but Path B is deferred for now.
- Path C requires a documented namespace-scoped access pattern (EKS Access Entry or IRSA) and runner connectivity guidance for the transition period.

---

## Unknowns and Assumptions

| Item | Type | Impact | Mitigation |
|------|------|--------|------------|
| Number of teams needing a Path C exception, and for how long | Unknown | Medium — drives transitional support cost | Track exceptions with target migration dates; review at each window's end |
| Strength of team opposition to the GitOps mandate | Unknown | Medium — could trigger a revisit of the mandate | Gather evidence during onboarding; escalate to decision maker if opposition is strong and well-evidenced |
| Namespace-scoped credential pattern for Path C | Assumption | Medium | Validate EKS Access Entry namespace-scope + RoleBinding in a PoC before granting Path C exceptions |
| Private-endpoint reachability from team runners | Assumption | Medium | Document supported runner/relay patterns; confirm with a Path C pilot team |
| Drift on Path C namespaces during the exception window | Risk | Low-Medium | Periodic platform audit of push-managed namespaces against baseline |

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
