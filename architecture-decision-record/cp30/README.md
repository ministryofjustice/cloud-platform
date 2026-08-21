# Container Platform 3.0 Architecture Decisions

This is the record of architectural decisions for **Container Platform 3.0
(CP30)**, the EKS-based multi-cluster, multi-account platform.

CP30 ADRs are kept separate from the original Cloud Platform (CP2) ADRs in the
[parent directory](../). The two logs use independent numbering:

- **CP2** ADRs are named `NNN-Title.md` (for example `047-use-karpenter-for-scaling.md`).
- **CP30** ADRs are named `ADR-NNN-kebab-case.md` (for example `ADR-002-argocd-gitops-fleet-management.md`).

The `ADR-NNN` numbers are the canonical CP30 identifiers and are referenced
directly from the platform Terraform and user documentation, so they are stable
and are not renumbered to avoid collisions with the CP2 sequence.

These decisions were originally authored during the CP30 design phase and are
being migrated into this repository as the permanent home. Additional CP30 ADRs
will be added here over time.

## Table of contents

- 🤔 [ADR-002: GitOps Fleet Management — EKS Capability for Argo CD](ADR-002-argocd-gitops-fleet-management.md)
- 🤔 [ADR-018: Deployment Model Flexibility — GitOps as Default, Push-Based CD as Accommodation](ADR-018-deployment-model-flexibility-gitops-vs-push-cd.md)

### Statuses

- Proposed: 🤔
- Accepted: ✅
- Rejected: ❌
- Superseded: ⌛️
- Amended: ♻️
