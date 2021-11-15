# 26 Managed Prometheus

Date: 2021-10-08

## Status

🤔 Proposed

## Proposal

Use [Amazon Managed Service for Prometheus](https://aws.amazon.com/prometheus/) (AMP) for monitoring Cloud Platform, instead of a self-managed instance of Prometheus.

## Context

### Monitoring metrics

It's good operational practice to have good 'observability'. This includes monitoring, achieved by regular checking the metrics, or health numbers, of the containers running. The timeseries data which is collected can be shown as graphs or other indicators in a dashboard, and evaluated against rules which trigger alerts to the operators. Typical use by operators include:

* to become familiar with the typical quantity of resources consumed by their software
* to be alerted to deteriorating health, so that they can fix it, before it becomes an incident
* being alerted to an incident, to be able to react quickly, not just when users flag it
* during an incident getting an at-a-glance overview of where problems exist
* after an incident to understand what went wrong, and help review the actions taken during the response
* reviewing long-term patterns of health

### Choice of Prometheus

Prometheus has been used by Cloud Platform for monitoring metrics since being established in 2017. CP have used Prometheus for a number of years and are happy with its functionality. Prometheus remains a popular choice in industry, is open source with a large community, and is [recommended by CNCF](https://raw.githubusercontent.com/cncf/trailmap/master/CNCF_TrailMap_latest.png).

**Commercial proprietary options**: The CP team has [investigated](https://docs.google.com/document/d/1eG2Fd4G8f_bPBYnNAwE2seF2PdNXsxjggsqck6IXNeI/edit) some commercial monitoring solutions, such as DataDog, Splunk, Honeycomb. These tend to do several related things, including metric monitoring, dashboard, alerting, logging and application performance monitoring. They don't offer any functionality we particular need. But it's nice that monitoring and logging are nicely integrated. Having it as a managed service, would reduce operations, and having to architect for future scaling. There was some concern about off-shoring log data, for circumstances where personal data leaks into logs. The cost varied $10k - $160k. There is some suggestion that even with a managed service, we'd still need to retain some use of Prometheus. There were concerns about users' config as code - whether it was as easy to deploy as it is with open source. And there would be the cost of migration users config would need migrating from their existing Prometheus/AlertManager/Grafana syntax.

So overall we are happy to stick with Prometheus.

### How CP uses Prometheus

Prometheus is setup to monitor the whole of Cloud Platform:

* Tenant containers
* Tenant AWS resources
* Kubernetes cluster. kube-prometheus

Prometheus is configured to store 24h worth of data, which is enough to support most use cases. The data is also sent on to Thanos, which efficiently stores 1 year of metrics data, and makes it available for queries using the same PromQL syntax.

Alertmanager uses the Prometheus data when evaluating its alert rules.

### Concerns about hosting Prometheus (as we currently do)

The Prometheus container has not run smoothly in recent months:

* **Performance (resolved)** - There were some serious performance issues - alert rules were taking too long to evaluate against the Prometheus data, however this was successfully alleviated by increasing the disk iops, so is not a remaining concern.

* **Custom node group** - Being a single Prometheus instance for monitoring the entire platform, it consumes a lot of resources. We've put it on a dedicated node, so it has the full resources. And it needs more memory than other nodes, which means it needs a custom node group, which is a bit of extra management overhead.

* **Scalability** - Scaling in this vertical way is not ideal - scaling up is not smooth and eventually we'll hit a limit of CPU/memory/iops. There are options to shard - see below.

We also need to address:

* **Management overhead** - Managed cloud services are generally preferred to self-managed because the cost tends to be amortized over a large customer base and be far cheaper than in-house staff. And people with ops skills are at a premium. The management overhead is:
  * for each of Prometheus, kube-prometheus

* **High availability** - We have a single instance of Prometheus, simply because we've not got round to choosing and implementing a HA arrangement yet. This risks periods of outage where we don't collect metrics data. Although the impact on the use cases is not likely to be very disruptive, there is some value in fixing this up.

### Options for addressing the concerns

**Thanos take load off Prometheus**: It's been suggested we could reduce load on Prometheus by it only retaining say 2h of logs, and shift as much as possible of the work to Thanos. However since the latest data is what you want to be running most queries and alert rules on, it is not clear how close to real-time Thanos can be kept. And Thanos rules are [unreliable](https://github.com/thanos-io/thanos/blob/main/docs/components/rule.md#risk)). This would likely reduce the load on Prometheus, but it may only be temporary, and it might simply shift the scalability concerns onto Thanos.

**Sharding**: We could split/shard the Prometheus instance: perhaps dividing into two - tenants and platform. Or if we did multi-cluster we could have one Prometheus instance per cluster. This appears relatively straightforward to do. There would be concern that there would be scaling thresholds where it is necessary to change how to divide it into shards, so a bit of planning would be needed.

**High Availability**: The recommended approach would be to run multiple instances of Prometheus configured the same, scraping the same endpoints independently. [Source](https://prometheus-operator.dev/docs/operator/high-availability/#prometheus) There is a `replicas` option to do this. However for HA we would also need to have a load balancer for the PromQL queries to the Prometheus API, to fail-over if the primary is unresponsive. And it's not clear how this works with duplicate alerts being sent to AlertManager. This doesn't feel like a very paved path, with Prometheus Operator [saying](https://prometheus-operator.dev/docs/operator/high-availability/) "We are currently implementing some of the groundwork to make this possible, and figuring out the best approach to do so, but it is definitely on the roadmap!" - Jan 2017, and not updated since.

**Managed Prometheus**: Using a managed service of prometheus, such as AMP, would address most of these concerns, and is evaluated in detail in the next section.

### Evaluation of managed Prometheus

Managed Prometheus, such as AMP, would out-source the operational concerns, including performance and scalability, as well as management of custom node groups. We'll examine these against the at a financial cost. And potentially monitoring from outside the VPC will change the architecture. This ADR aims to understand all these issues. Specifically interested in Amazon Managed Service for Prometheus (AMP).

Scaling: AMP scales automatically, without configuration.

High Availability: AMP is highly available and distributed across multiple AZs automatically. In contrast there are question marks about how to do this with self-hosted - see above.

Resilience: AMP is relatively isolated against cluster issues. The data kept in durable storage, away from the cluster.

Lock-in: the configuration syntax and other interfaces are the same or similar to our existing self-hosted Prometheus, so we maintain low lock-in / migration cost.


### Existing install

The 'monitoring' namespace is configured in [components terraform](https://github.com/ministryofjustice/cloud-platform-infrastructure/blob/main/terraform/aws-accounts/cloud-platform-aws/vpc/eks/components/components.tf#L115-L138) calling the [cloud-platform-terraform-monitoring module](https://github.com/ministryofjustice/cloud-platform-terraform-monitoring). This [installs](https://github.com/ministryofjustice/cloud-platform-terraform-monitoring/blob/main/prometheus.tf#L88) the [kube-prometheus-stack Helm chart](https://github.com/prometheus-community/helm-charts/blob/main/charts/kube-prometheus-stack/README.md) / [kube-prometheus](https://github.com/prometheus-operator/kube-prometheus) (among other things).

[kube-prometheus](https://github.com/prometheus-operator/kube-prometheus) contains a number of things:

* [Prometheus Operator](https://github.com/prometheus-operator/prometheus-operator) - adds kubernetes-native wrappers for managing Prometheus
  * CRDs for install: Prometheus, Alertmanager, Grafana, ThanosRuler
  * CRDs for configuring: ServiceMonitor, PodMonitor, Probe, PrometheusRule, AlertmanagerConfig
    - allows specifying monitoring targets using kubernetes labels 
* Kubernetes manifests
* Grafana dashboards
* Prometheus rules
* example configs for: node_exporter, scrape targets, alerting rules for cluster issues

High Availability - not implemented (yet).

Resilience - In the absence of HA, some resilience is offered by the presence of a back-up node. If the node on which Prometheus instance is running fails, then the back-up node is ready to run a Prometheus instance, so downtime is capped at 2 minutes. The back-up node is the same spec and configuration. It is running in the same AZ, so that it can access the same EBS with the Prometheus data.
https://github.com/ministryofjustice/cloud-platform/issues/1749#issue-587058014

Prometheus config is held in k8s resources:

* ServiceMonitor
* PrometheusRule - alerting

Barriers to adopting AMP:

* AMP is not released in the London region yet (at the time of writing, 3/11/21)
* alertmanager - it doesn't do our alerts yet (at the time of writing, 3/11/21). but alertmanager is due later in the year
* storage - you can throw as much data at it. no need for Thanos then.
* Managed grafana has no terraform support yet so just setup in AWS console.
* AMP is headless. ie no web interface. TODO is this an issue?

## How it would work with AMP

A 'forwarding' Prometheus instance installed in the cluster, which scrapes data, keeps a copy and uses "remote write" to forward all of its data to AMP. This transfer looks pretty robust, using a write ahead log and queue.

The forwarding Prometheus could simply be the existing Prometheus installed by kube-prometheus, but instead of storing it for 24h, just long enough to 'remote write' it AMP.

(Alternatively, if we followed what [what AP do](https://github.com/moj-analytical-services/analytical-platform-flux/blob/main/clusters/development/prometheus/prometheus.yaml)), the staging Prometheus would be installed with the [Prometheus Helm chart](https://github.com/prometheus-community/helm-charts/tree/main/charts/prometheus). The Prometheus Helm chart configures Prometheus to scrape k8s resources that have [prometheus.io annotations](https://github.com/prometheus-community/helm-charts/tree/main/charts/prometheus#scraping-pod-metrics-via-annotations). So all the ServiceMonitor configs etc would need to be changed to annotations, which is a big burden for tenants.)


TODO:  prometheus != kube-prometheus (currently in the cluster) - jason
kube-prometheus compatibility? there are a load of specific things from kube-prometheus that we’ll lose? Prometheus rules for example, there are  currently 3543 defined in live-1. How do end users define a rule in managed prometheus? (as an example) I imagine it would probably be manual in the short term.

AMP supports both recording and alert rules, these are defined in one or more YAML rules files which are the same format as standalone Prometheus, so you you should be able to reuse your existing rules. You can have multiple rules files in an AMP workspace but each must be in a separate namespace. See https://docs.aws.amazon.com/prometheus/latest/userguide/AMP-Ruler.html.
Once you have defined the rules YAML files, you can base64 the files and upload them via the CLI https://docs.aws.amazon.com/prometheus/latest/userguide/AMP-rules-upload.html
or API using the create-rule-groups-namespace command: https://docs.aws.amazon.com/prometheus/latest/userguide/AMP-APIReference.html#AMP-APIReference-CreateRuleGroupsNamespace

AWS Region
  AMP workspace 
    AMP namespace - one rule group, one rules YAML file

TODO: look at scale and costs

