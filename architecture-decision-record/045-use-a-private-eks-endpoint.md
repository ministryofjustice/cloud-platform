# 45. Use a Private EKS Endpoint

Date: 2026-04-02

## Status

✅ Accepted

## Context

The cluster currently uses a public EKS endpoint, although there is robust authentication in place, there is always a risk if something is exposed on the public internet.  AWS best practice recommends to keep EKS endpoints private.

## Decision

The EKS endpoint will be private. This increases cluster security and follows best practice. It also aligns with other changes we want to make, such as using GitOps and holding all container images in ECR on the platform.

## Consequences

 - The EKS endpoint will no longer be publicly accessible.
 - We will need to provide a method of access for engineers connecting to the endpoint, this will be via VPN in the first instance, with a bastion accessed via SSM as a break glass measure.
 - We will need traffic to our EKS endpoints to be routed via the Global Protect VPN.
 - We will need to set up a Bastion, likely using the [Modernisation Platform Bastion module](https://github.com/ministryofjustice/modernisation-platform-terraform-bastion-linux)
 - We will no longer allow external connections to the cluster endpoint.
 - We will manage our deployments via GitOps. (Potentially ArgoCD)
