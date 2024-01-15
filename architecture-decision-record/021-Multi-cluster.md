# 21. Multi-cluster

Date: 2021-05-11

## Status

⌛️ Superseded by [036. Multi-cluster](036-multi-cluster.md)

## What’s proposed

We host user apps across _more than one_ Kubernetes cluster. Apps could be moved between clusters without too much disruption. Each cluster _may_ be further isolated by placing them in separate VPCs or separate AWS accounts.

## Context

Service teams' apps currently run on [one Kubernetes cluster](https://github.com/ministryofjustice/cloud-platform/blob/main/architecture-decision-record/012-One-cluster-for-dev-staging-prod.md). That includes their dev/staging/prod environments - they are not split off. The key reasoning was:

- Strong isolation is already required between apps from different teams (via namespaces, network policies), so there is no difference for isolating environments
- Maintaining clusters for each environment is a cost in effort
- You risk the clusters diverging. So you might miss problems when testing on the dev/staging clusters, because they aren't the same as prod.

(We also have clusters for other purposes: a 'management' cluster for Cloud Platform team's CI/CD and ephemeral 'test' clusters for the Cloud Platform team to test changes to the cluster.)

However we have seen some problems with using one cluster, and advantages to moving to multi-cluster:

- Scaling limits
- Single point of failure
- Derisk upgrading of k8s
- Reduce blast radius for security
- Reduce blast radius of accidental deletion
- Pre-prod cluster
- Cattle not pets

### Scaling limits

Multi-cluster helps us we encounter a scale limitation. For example, we've found ourselves unexpectantly hitting an AWS limit during k8s upgrade. In this situation we could off-load some apps to another cluster. It would be advantageous to put each cluster in its own AWS account, to avoid limits, which are imposed per-account.

### Single point of failure

Running everything on a single cluster is a 'single point of failure', which is a growing concern as more services use CP. Multi-cluster would allow us to quickly move apps off a broken cluster to another cluster.

Several elements in the cluster are a single point of failure:

- ingress (incidents: [1](https://runbooks.cloud-platform.service.justice.gov.uk/incident-log.html#incident-on-2020-10-06-09-07-intermittent-quot-micro-downtimes-quot-on-various-services-using-dedicated-ingress-controllers) [2](https://runbooks.cloud-platform.service.justice.gov.uk/incident-log.html#incident-on-2020-04-15-10-58-nginx-tls))
- external-dns
- cert manager
- kiam
- OPA ([incident](https://runbooks.cloud-platform.service.justice.gov.uk/incident-log.html#incident-on-2020-02-25-10-58))

### Derisk upgrading of k8s

Once you start a Kubernetes version upgrade, rolling back becomes infeasible (incidents: [1](https://runbooks.cloud-platform.service.justice.gov.uk/incident-log.html#q1-2020-january-march)).

With multi-cluster we could do a "blue-green" upgrade - spin up an cluster at the newer k8s version and then carefully move the apps across to it.

### Reduce blast radius for security

Taking a layered approach to security, extra isolation is beneficial. It resists lateral movement and minimizes the impact of a breach.

Isolation is added when you split the workloads across multiple clusters, even if they are in the same VPC. And further isolation is gained with separate VPCs or separate AWS accounts.

More sensitive apps may require this isolation.

Pre-prod environments are likely to often be running new code that has not have been through all the reviews, quality and security checks yet, so there may be a case for keeping these more isolated from environments with access to production data.

### Reduce blast radius of accidental deletion

In the case of accidental deletion being run by a Cloud Platform team member, the ability to run administrative commands on only one cluster at a time would reduce the impact of this event, such as [this incident](https://runbooks.cloud-platform.service.justice.gov.uk/incident-log.html#incident-on-2020-09-21-18-27-some-cloud-platform-components-destroyed). Disaster Recovery procedures are in [good shape now](https://runbooks.cloud-platform.service.justice.gov.uk/disaster-recovery-scenarios.html#cloud-platform-disaster-recovery-scenarios), but it's worth minimizing blast radius all the same.

### Pre-prod cluster

Test clusters are used by the Cloud Platform team to test changes. These are made as realistic as possible, inheriting IaC and config from the main cluster, and the availability of test apps. However test clusters have a different lifecycle, don't run realistic loads and don't have traffic / loads. Ideally we test CP changes in the most realistic way possible, before changes are deployed to the part of the platform where production workloads are running.

Multi-cluster will allow us to put pre-prod environments on a separate cluster to prod environments. Then changes to the platform, once tested on a test cluster, we can roll them out to the pre-prod cluster, as a more realistic test, ahead of rolling out to the prod clusters. This is particularly beneficial for more fundamental changes like upgrading k8s, and single points of failure such as ingress or OPA.

### Challenge of moving apps

If we were to create a fresh cluster, and an app is moved onto it, then there are a lot of impacts:

- **Kubecfg** - a fresh cluster will have a fresh kubernetes key, which invalidates everyone's kubecfg. This means that service teams will need to obtain a fresh token and add it to their app's CI/CD config and every dev will need to refresh their command-line kubecfg for running kubectl.
- **IP Addresses** - unless the load balancer instance and elastic IPs are reused, it'll have fresh IP addresses. This will particularly affect devices on mobile networks that accessing our CP-hosted apps, because they often cache the DNS longer than the TTL. And if CP-hosted apps access third party systems and have arranged for our egress IP to be allow-listed in their firewall, then they will not work until that's updated.

## Steps to achieve it

More ideas for steps towards this are here:
https://docs.google.com/document/d/1-1dXD0a50MZiQUcm4fIIO1ySVQHA_D2-UMIti-oBVNc/edit#heading=h.4hza65nlnk7p

## Decision

TBD
