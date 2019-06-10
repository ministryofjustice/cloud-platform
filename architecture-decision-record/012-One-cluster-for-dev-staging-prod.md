# One cluster for dev/staging/prod

Date: 01/06/2018

## Status

Accepted

## Context

The Cloud Platform needs to host both citizen-facing, production services, and development environments for service teams to iterate on their code, or just set up sandboxes for experimentation and learning.

To support this, should we have separate clusters for production, development, and staging? Or, should we run a single cluster hosting all these different types of workload?

## Decision

After consideration of the pros and cons of each approach we went with one cluster, using namespaces to partition different workloads.

Some important reasons behind this move were:

* A single k8s cluster can be made powerful enough to run all of our workloads
* Managing a single cluster keeps our operational overhead and costs to a minimum.
* Namespaces and RBAC keep different workloads isolated from each other.
* It would be very hard to keep multiple clusters (dev/staging/prod) from becoming too different to be representative environments

To clarify the last point; to be useful, a development cluster must be as similar as possible to the production cluster. However, given multiple clusters, with different security and other constraints, some 'drift' is inevitable - e.g. the development cluster might be upgraded to a newer kubernetes version before staging and production, or it could have different connectivity into private networks, or different performance constraints from the production cluster.

Based on our past experience, these differences tend to increase over time, to the point where the development cluster is too far away from production to be representative. The extra work required to maintain multiple environments becomes wasted effort.

If namespace segregation is sufficient to isolate one production service from another, then it is enough to isolate a team's development environment from a production service.

If namespace segregation is not sufficient for this, then the whole cloud platform idea doesn't work.

## Consequences

Having a single cluster to maintain works well for us.

* Service teams know that their development environments accurately reflect the production environments they will eventually create
* There is no duplication of effort, maintaining multiple, slightly different clusters
* All services are managed centrally (e.g. ingress controller, centralised log collection via fluentd, centralised monitoring with Prometheus, cluster security policies)

