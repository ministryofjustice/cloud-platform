# Use kubernetes for running containerised applications

Date: 10/04/18

## Status

✅ Accepted

## Context

MOJ Digital's approach to infrastructure management and ownership has evolved over time, and has led to the following outcomes:

- Unclear boundaries on ownership and responsibilities between service teams and the cloud platforms team
- Significant variation in deployment, monitoring and lifecycle management across products
- Inefficient use of AWS resources due to the use of virtual machine-centric architecture, despite our standardisation on Docker containers

The last few years has seen the advent of several products specifically focused on the problem of running and managing containers in production:

- Kubernetes
- Mesos / Mesosphere / DC/OS
- Docker Swarm
- AWS ECS
- CloudFoundry

Given the technology landscape within MOJ, we require a container management platform that can support a wide range of applications, from "modern" cloud-native 12-factor applications through to "legacy" stateful monolithic applications, potentially encompassing both Linux- and Windows-based applications; this removes CloudFoundry from consideration, given its focus on modern 12-factor applications and reliance on buildpacks to support particular runtimes.

From the remaining list of major container platforms, Kubernetes is the clear market leader:

- Rapid industry adoption during 2017 establishing it as the emerging defacto industry standard
- Managed Kubernetes services from all major cloud vendors
- Broad ecosystem of supporting tools and technologies
- Increasing support for Kubernetes as a deployment target for commercial and open-source software projects

There is also precedent for Kubernetes use within MOJ, as the Analytical Platform team has been building on top of Kubernetes for around 18 months.

## Decision

Use Kubernetes as the container management component and core technology for our new hosting platform.

## Consequences

1. Several technical spikes looking at Kubernetes itself and associated components - logging, monitoring, application deployment, etc.
2. Requirement to use an identity broker for Kubernetes user authentication, given the decision to use Github as an identity provider and the lack of a common auth protocol between Kubernetes and Github
