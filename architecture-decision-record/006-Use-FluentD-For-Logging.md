# Use FluentD

Date: 09/04/2018

## Status

Pending

## Context

We want to understand how monitoring works on a Kubernetes cluster, in particular how we can get logging from an our Kubernetes cluster to somewhere where a team can consume and graph it (dashboard) and we can alert on issues.

We want to be able to get logs out from a Kubernetes cluster to somewhere that we can consume and act on them – elasticsearch is the primary candidate.

There are 2 categories of logs we are interested in: 
1. Application logs for things running on Kubernetes
2. Logs from Kubernetes itself that might be interesting to us (like API logs).

Reasons behind this move were:

* Greater visibility of what is currently on cluster.
* With the move to Kubernetes a cloud native logging solution is needed

## Decision

Use FluentD 

## Consequences

1. Spike into what are the best solutions for logging on Kubernetes [CPT-372](https://dsdmoj.atlassian.net/browse/CPT-372)
2. Decision made to deploy to Kubernetes [CPT-488](https://dsdmoj.atlassian.net/browse/CPT-488)
3. A diagram of how logs go from running applications to an Elasticsearch cluster [CPT-372](https://dsdmoj.atlassian.net/browse/CPT-372)
3. Automate deployment [CPT-372](https://dsdmoj.atlassian.net/browse/CPT-372)
4. FluentD is logging EVERYTHING.
5. WIP