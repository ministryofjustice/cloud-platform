# 21. Multi-cluster

Date: 2021-05-11

## Status

🤔 Proposed

## What’s proposed

We host user apps across *more than one* Kubernetes cluster. Apps could be moved between clusters without too much disruption. Each cluster *may* be further isolated by placing them in separate VPCs or separate AWS accounts.

## Context

Service teams' apps currently run on [one Kubernetes cluster](012-One-cluster-for-dev-staging-prod.html). That includes their dev/staging/prod environments - they are not split off. The key reasoning was:

* Strong isolation is already required between apps from different teams (via namespaces, network policies), so there is no difference for isolating environments
* Maintaining clusters for each environment is a cost in effort
* You risk the clusters diverging. So you might miss problems when testing on the dev/staging clusters, because they aren't the same as prod.

(We also have clusters for other purposes: a 'manager' cluster for Cloud Platform team's CI/CD and ephemeral 'test' clusters for the Cloud Platform team to test changes to the cluster.)

However we have seen some problems with using one cluster, and advantages to moving to multi-cluster:

* Derisk upgrading of k8s
* Isolate environments and sensitive workloads
* Scaling limits
* Single point of failure
* Cattle not pets

### Derisk upgrading of k8s

Once you start a Kubernetes version upgrade, rolling back becomes infeasible (incidents: [1](https://runbooks.cloud-platform.service.justice.gov.uk/incident-log.html#q1-2020-january-march)).

With multi-cluster we could do a "blue-green" upgrade - spin up an cluster at the newer k8s version and then carefully move the apps across to it.

### Isolate environments and sensitive workloads

* **Blast radius** - with only one cluster, the blast radius could be catastrophic, in the case of a breach or accidental deletion ([example](https://runbooks.cloud-platform.service.justice.gov.uk/incident-log)). Multiple clusters reduces that blast radius.

* **Security** can be increased, by adding layers of isolation, such as separate VPCs or AWS accounts. This is a potential requirement for more sensitive apps. It might also be for security of less tested code.

### Scaling limits

Multi-cluster helps us when we encounter a scale limitation. For example, we've found ourselves unexpectantly hitting an AWS limit during k8s upgrade. In this situation we could off-load some apps to another cluster. It would be advantageous to put each cluster in its own VPC and eventually AWS account, to avoid limits which are imposed per-VPC and per-account.

### Single point of failure

Running everything on a single cluster is a 'single point of failure', which is a growing concern as more services use CP. Multi-cluster would allow us to quickly move apps from a degraded or broken cluster to another cluster.

Several elements in the cluster are a single point of failure:

* ingress (incidents: [1](https://runbooks.cloud-platform.service.justice.gov.uk/incident-log.html#incident-on-2020-10-06-09-07-intermittent-quot-micro-downtimes-quot-on-various-services-using-dedicated-ingress-controllers) [2](https://runbooks.cloud-platform.service.justice.gov.uk/incident-log.html#incident-on-2020-04-15-10-58-nginx-tls))
* external-dns
* cert manager
* kiam
* OPA ([incident](https://runbooks.cloud-platform.service.justice.gov.uk/incident-log.html#incident-on-2020-02-25-10-58))

### Challenge of moving apps

If we were to create a fresh cluster, and an app is moved onto it, then there are a lot of impacts:

* **IP Addresses** - unless the load balancer instance and elastic IPs are reused, it'll have fresh IP addresses. This will particularly affect devices on mobile networks that accessing our CP-hosted apps, because they often cache the DNS longer than the TTL. And if CP-hosted apps access third party systems and have arranged for our egress IP to be allow-listed in their firewall, then they will not work until that's updated.

* **Kubecfg** - a fresh cluster will have a fresh kubernetes key, which invalidates everyone's kubecfg. This means that service teams will need to obtain a fresh token and add it to their app's CI/CD config and every dev will need to refresh their command-line kubecfg for running kubectl.

## Decision

TBD