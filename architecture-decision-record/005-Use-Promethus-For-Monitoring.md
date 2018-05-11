# Use Concourse CI

Date: 06/04/2018

## Status

Pending

## Context

We want to understand how monitoring works on a Kubernetes cluster, in particular how we can get metrics from an application running on Kubernetes to somewhere where a team can consume and graph it (dashboard) and we can alert on issues.

We will focus on Prometheus based solutions to begin with – Prometheus is the leading Kubernetes monitoring solution at the moment and this is a pretty big space.

We will also want to understand how this might work for monitoring the Kubernetes cluster itself, rather than applications running on it.

Reasons behind this move were:

* Need to adopt best practice for monitoring in Kubernetes
* ....


## Decision

Use Prometheus Operator

## Consequences

1. Spike into monitoring on Kubernetes [CPT-371](https://dsdmoj.atlassian.net/browse/CPT-371)
2. Decision on where to deploy [CPT-](https://dsdmoj.atlassian.net/browse/CPT-)
3. Automate deployment of Prometheus using Helm.
4. WIP 


~~Cloud Platforms can monitor all resources in all Kubernetes clusters.~~