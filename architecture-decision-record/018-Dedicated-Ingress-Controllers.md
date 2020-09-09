# Dedicated Ingress Controllers for teams

Date: 02/09/2020

## Status

❌ Rejected (see [019](019-Shared-Ingress-Controllers.md))

## Context

The Cloud Platform was originally set up to have a single ingress controller to
manage all ingresses in the cluster. So, every new ingress added a config block
to one, large nginx config file, and all traffic to all services hosted on the
Cloud Platform is routed through a single AWS load balancer.

Although this was both easy to manage, and saved us some money on load
balancers (approx. $25/month per ingress), it has become unsustainable. We
usually have 6 replicas of the ingress controller pod, and we have started to
see instances of several of these pods crash-looping (usually because they have
run out of shared memory, which cannot be increased in kubernetes. See [this
issue] for more information).

We believe this is because the nginx config has become so large (over 100K
lines), that sometimes pods fail to reload it when it is changed, or the pod is
moved.

## Decision

We will create a separate AWS load balancer and ingress-controller for every
namespace in the cluster. An "ingress class" annotation will cause traffic for
a particular ingress to be routed through the appropriate AWS load balancer and
ingress-controller. See our [module repository] for more details.

"System" ingresses (e.g. those used for concourse, grafana, etc.) will continue
to use the default ingress-controller. There should only ever be a handful of
these, compared with hundreds of team ingresses, so the load on the default
ingress-controller should stay within acceptable limits.

## Consequences

* The default ingress-controller will not crash (and take down all services in
  the cluster).
* The AWS hosting costs of the Cloud Platform will increase (because we will
  have to pay for more load balancers).
* The number of pods on the cluster will increase dramatically (two new pods
  per namespace, plus an additional pod if it's a production namespace).
* We may need more and/or larger worker nodes to run the additional pods, which
  will further increase our hosting costs.

[module repository]: https://github.com/ministryofjustice/cloud-platform-terraform-teams-ingress-controller
[this issue]: https://github.com/kubernetes/kubernetes/issues/28272
