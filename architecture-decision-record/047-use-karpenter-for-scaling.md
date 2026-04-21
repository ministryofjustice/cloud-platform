# 47. Use Karpenter For Scaling

Date: 2026-04-21

## Status

✅ Accepted

## Context

The Cloud Platform Team [investigated Karpenter as the technology solution for scaling](https://github.com/ministryofjustice/cloud-platform/issues/8009) during the build of `Container Platform`. We installed and tested **Karpenter**. We [compared Karpenter](https://justiceuk-my.sharepoint.com/:w:/r/personal/dan_criddle_justice_gov_uk/Documents/CP3%20Auto%20Scaling%20Karpenter%20%20Cluster%20Auto%20Scaler.docx?d=w75429c2098ca4980ae8a5151ff211fd1&csf=1&web=1&e=r7FJZF) with **Cluster Auto Scaler** and **EKS Auto Mode**. Testing was successful.

Testing was carried out on `Modernisation Platform` running on **AWS**. It is assumed the cloud provider being used for `Container Platform` will continue to be **AWS**.

## Decision

We will use **Karpenter** as the primary scaling solution for `Container Platform`.

## Consequences

- Scaling is event-driven via direct EC2 Fleet API calls
- User workloads will not run in node groups (**Karpenter** runs outside of node groups, picking instance types dynamically)
- Scaling should occur faster than in `Cloud Platform`
- Node consolidation is more aggressive
- We need to come up with a list of acceptable AWS Instance types and limit to that
- We need to perform tuning and testing with **Karpetner**
- We need to ensure security - see [Karpenter Threat Model](https://karpenter.sh/docs/reference/threat-model/)
- Although it's not our primary use-case **Cluster Auto Scaler** could still run alongside **Karpenter** if needed
