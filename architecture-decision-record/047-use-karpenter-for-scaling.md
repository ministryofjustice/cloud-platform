# 47. Use Karpenter For Scaling

Date: 2026-04-21

## Status

✅ Accepted

## Context

The Cloud Platform Team [investigated Karpenter as the technology solution for scaling](https://github.com/ministryofjustice/cloud-platform/issues/8009) during the build of `Container Platform`. We installed and tested Karpetner. We [compared Karpenter](https://justiceuk-my.sharepoint.com/:w:/r/personal/dan_criddle_justice_gov_uk/Documents/CP3%20Auto%20Scaling%20Karpenter%20%20Cluster%20Auto%20Scaler.docx?d=w75429c2098ca4980ae8a5151ff211fd1&csf=1&web=1&e=r7FJZF) with Cluster Auto Scaler and EKS Auto Mode.

## Decision

We will use Karpenter as the primary scaling solution for `Container Platform`.

## Consequences

- Karpenter runs outside of node groups, picling instance types dynamically
- Although it's not our primary use-case Cluster Auto Scaler can run alongside Karpenter. If we need to do that at any point could add in Cluster Auto Scaler with fixed node groups
