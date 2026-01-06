# Cloud Platform Architecture Decisions

This is a record of architectural decisions made during the development of the
Cloud Platform.

To understand why we are recording decisions and how we are doing it, please
see [ADR-000](000-Record-Architecture-Decisions.md)

## Table of contents

- ✅ [0. Record Architecture Decisions](000-Record-Architecture-Decisions.md)
- ✅ [1. Use AWS hosted Elasticsearch](001-Use-AWS-hosted-elasticsearch.md)
- ✅ [2. Use GitHub for architecture decision log](002-Use-github-for-architecture-decision-record.md)
- ✅ [3. Use Concourse CI](003-Use-Concourse-CI.md)
- ✅ [4. Use kubernetes for running containerised applications](004-use-kubernetes-for-container-management.md)
- ✅ [5. Monitoring and Alerting of ECR Image counts](005-ECR-monitoring-and-alerting.md)
- ✅ [6. Use GitHub as our identity provider](006-Use-github-as-user-directory.md)
- ✅ [7. Use ECR As Container Registry](007-Use-ECR-As-Container-Registry.md)
- ✅ [8. Support Deployments from Third Party CI](008-Support-Deployments-from-Third-Party-CI.md)
- ✅ [9. Naming convention for clusters](009-Naming-convention-for-clusters.md)
- ✅ [10. live-0 to live-1 Cluster](010-live-0-to-live-1-Cluster.md)
- ✅ [11. Introduce Open Policy Agent](011-Introduce-Open-Policy-Agent.md)
- ✅ [12. One cluster for dev/staging/prod](012-One-cluster-for-dev-staging-prod.md)
- ✅ [13. Use RSpec for cluster tests](013-Use-RSpec-for-cluster-tests.md)
- ✅ [14. Why we build our own kubernetes cluster](014-Why-we-build-our-own-kubernetes-cluster.md)
- ✅ [15. Keeping Terraform modules up to date](015-Keeping-Terraform-modules-up-to-date.md)
- ✅ [16. Kibana is open to all service teams](016-Kibana-is-open-to-all-service-teams.md)
- 🤔 [17. Variable names are in snake-case](017-Variable-Naming.md)
- ❌ [18. Dedicated ingress controllers](018-Dedicated-Ingress-Controllers.md)
- 🤔 [19. Shared ingress controllers](019-Shared-Ingress-Controllers.md)
- 🤔 [20. Environments and Pipeline](020-Environments-and-Pipeline.md)
- ⌛️ [21. Multi-cluster](021-Multi-cluster.md)
- ✅ [22. EKS](022-EKS.md)
- ✅ [23. Logging](023-Logging.md)
- 🤔 [24. Use SSO for CP Team to access AWS Console](024-SSO-for-CP-Team-to-access-AWS.md)
- ✅ [25. Domain names](025-Domain-names.md)
- 🤔 [26. Managed Prometheus](026-Managed-Prometheus.md)
- 🤔 [27. AWS Console access for users](027-AWS-Console-access-for-users.md)
- 🤔 [28. Repo vulnerability scanning](028-Repo-vulnerability-scanning.md)
- 🤔 [29. Kubernetes security config auditing](029-Kubernetes-security-config-auditing.md)
- 🤔 [30. AWS Root User Security](030-AWS-root-user-security.md)
- 🤔 [31. Image vulnerability scanning](031-Image-vulnerability-scanning.md)
- ✅ [33. Manager cluster](033-Manager-cluster.md)
- 🤔 [34. EKS Fargate](034-EKS-Fargate.md)
- ✅ [35. Deprecated TLS Versions](035-deprecated-tls-versions.md)
- ✅ [36. Multi-cluster](036-multi-cluster.md)
- ❌ [37. Serverless (Lambda/Functions-as-a-Service)](037-serverless.md)
- 🤔 [38. Split Components Infrastructure](038-split-components-infrastructure.md)
- ✅ [39. AWS Network Firewall](039-AWS-Network-Firewall.md)
- ✅ [40. Use GitHub Actions](040-use-github-actions.md)
- ✅ [41. Decommission CLI Tool](041-decommission-cp-cli-tool.md)
- ✅ [42. Use Modernisation Platform accounts](042-use-modernisation-platform-accounts.md)

### Statuses:

- Proposed: 🤔
- Accepted: ✅
- Rejected: ❌
- Superseded: ⌛️
- Amended: ♻️
