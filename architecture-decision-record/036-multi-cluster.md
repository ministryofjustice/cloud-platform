# 36. Multi-cluster

Date: 2022-09-16

## Status

✅ Accepted

## Context

At the time of writing (16th September, 2022); Cloud Platform hosts user applications in a single Kubernetes cluster. On 12th April, 2022, we decommissioned the [old kOps cluster](https://github.com/ministryofjustice/cloud-platform-environments/commit/2a36ba40562ff77e64dd18d1e1fca3f253e38c80) in favour of utilising EKS (see [ADR-022 - EKS](022-EKS.md)).

At the time of decommissioning the kOps cluster, Cloud Platform [hosted 409 namespaces](https://github.com/ministryofjustice/cloud-platform-environments/tree/2a36ba40562ff77e64dd18d1e1fca3f253e38c80/namespaces/live.cloud-platform.service.justice.gov.uk) in the new EKS-based Kubernetes cluster.

The key reasoning behind utilising a single Kubernetes cluster for service teams is described in [ADR-021](021-Multi-cluster.md), namely:

- Strong isolation is already required between apps from different teams (via namespaces, network policies), so there is no difference for isolating environments
- Maintaining clusters for each environment is a cost in effort
- You risk the clusters diverging. So you might miss problems when testing on the dev/staging clusters, because they aren't the same as prod.

However, outside of service teams applications; the Cloud Platform does already use multiple clusters for 'management', and ephemeral 'test' purposes, and further information regarding multi-cluster and its challenges are described in [ADR-021](021-Multi-cluster.md).

The Cloud Platform team's use of EKS has matured and usage of the EKS-based cluster has grown since April from 409 namespaces, to 457 at the time of writing and will continue to grow in the future. Due to this, we believe now is a good time to start separating out workloads into a multi-cluster environment, as the nuances and costs of EKS are more widely understood.

It is worth noting that the current cluster (at the time of writing) is not considered a ["large cluster" by Kubernetes](https://kubernetes.io/docs/setup/best-practices/cluster-large/), which stipulates these considerations for a "large cluster":

- No more than 110 pods per node
- No more than 5000 nodes
- No more than 150000 total pods
- No more than 300000 total containers

The Cloud Platform is not near these limits.

## Decision

We have decided to:

- adopt a multi-cluster approach
- start adopting multi-cluster by splitting out workloads by _environment type_, i.e. _production_ and _non-production_
- treat the _non-production_ cluster as a second _production_ cluster within the team, though there may be slight differences such as:
  - compute power, where _non-production_ workloads don't require as intensive power as _production_ workloads
  - hours of operation, for e.g. we may turn off the non-production cluster out of hours or on weekends if there is appetite to do so with our users to save on cost
- we will allow users to self-select their preferred environment initially, though this may change or automate in the future

## Consequences

By running multiple clusters, we can:

- derisk upgrading Kubernetes versions as non-production workloads will be separate from production workloads
- reduce the blast radius by creating a hard separation between clusters
- reduce the cluster single-points of failure, such as ingress controllers
- increase our separation outside of Kubernetes between workloads
- improve our resilience
- we can customise the non-production cluster to meet non-production needs, and the production cluster to meet production needs

However, this does come at both a monetary and team cost:

- we expect the cost of running the Cloud Platform to increase, not significantly but not insignificantly
- there will be less efficient resource usage as some nodes may not fill up
- there will be an increased complexity of cluster management
