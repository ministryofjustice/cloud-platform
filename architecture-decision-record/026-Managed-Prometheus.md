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

Prometheus is setup to monitor the whole of Cloud Platform, including:

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

**Sharding**: We could split/shard the Prometheus instance: perhaps dividing into two - tenants and platform. Or if we did multi-cluster we could have one Prometheus instance per cluster. This appears relatively straightforward to do. There would be concern that however we split it, as we scale in the future we'll hit future scaling thresholds, where it will be necessary to change how to divide it into shards, so a bit of planning would be needed.

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

## How it would work with AMP

We've spiked a terraform module to implement this: https://github.com/ministryofjustice/cloud-platform-terraform-amp

A 'forwarding' Prometheus instance installed in the cluster, which scrapes data, keeps a copy and uses "remote write" to forward all of its data to AMP. This transfer looks pretty robust, using a write ahead log and queue.

The forwarding Prometheus could simply be the standard community Helm chart: https://prometheus-community.github.io/helm-charts which is simple and makes upgrades easy. We don't need kube-prometheus any more. It doesn't need to store much - just enough to forward it to AMP. TODO work out if the link to AMP breaks down for a period.

(This arrangement is the same as [what AP do](https://github.com/moj-analytical-services/analytical-platform-flux/blob/main/clusters/development/prometheus/prometheus.yaml)). The Prometheus Helm chart configures Prometheus to scrape k8s resources that have [prometheus.io annotations](https://github.com/prometheus-community/helm-charts/tree/main/charts/prometheus#scraping-pod-metrics-via-annotations). However we'll let users define ServiceMonitor configs as before.)

Storage: - you can throw as much data at it. Instead there is a days limit of 150 days. This is plenty long enough for the use cases, so no need for **Thanos**, so this is very useful.

Alertmanager:

* AMP has an Alertmanager-compatible option, which we'd use with the same rules
* Sending alerts would need to us to configure: create SNS topic that forwards to user Slack channels

Grafana:

* Amazon Managed Grafana has no terraform support yet so just setup in AWS console. So in the meantime we stick with self-managed Grafana, which works fine.

Prometheus web interface - previously AMP was headless, but now it comes with the web interface

Prometheus Rules and Alerts:

* In our existing cluster:
   * we get ~3500 Prometheus rules from: https://github.com/kubernetes-monitoring/kubernetes-mixin 
   * kube-prometheus compiles it to JSON and applies it to the cluster
* So for our new cluster:
   * we need to do the same thing for our new cluster. But let's avoid using kube-prometheus. Just copy what it does.
   * when we upgrade the prometheus version, we'll manually [run the jsonnet config generation](https://github.com/kubernetes-monitoring/kubernetes-mixin#generate-config-files), and paste the resulting rules into our terraform module e.g.: https://github.com/ministryofjustice/cloud-platform-terraform-amp/blob/main/example/rules.tf

### Still to figure out

#### Tenants' rules and alerts

Tenants [create PrometheusRule k8s resources by kubectl apply](https://user-guide.cloud-platform.service.justice.gov.uk/documentation/monitoring-an-app/how-to-create-alarms.html#create-a-prometheusrule) or in environments repo. How are we going to get this config inserted into AMP? Could it just be a cronjob that copies these resources' config into AMP? (We do something similar for Grafana.) And we should check that invalid rules don't get applied and cause problems for other tenants.

#### Costs

Look at scale and costs. Ingestion: $1 for 10m samples

Prices (Ireland):

* EU-AMP:MetricSampleCount - $0.35 per 10M metric samples for the next 250B metric samples
* EU-AMP:MetricStorageByteHrs - $0.03 per GB-Mo for storage above 10GB

#### Region

AMP is not released in the London region yet (at the time of writing, 3/11/21). However we could run it in another region - data ingestion into the regio is free - we would pay only for the Grafana queries.

#### Other components we use

We should check our usage of these related components, and if we still need them in the new cluster:

* CloudWatch exporter
* Node exporter
* ECR exporter
* Pushgateway

#### Showing alerts

We need to show users their alerts - both the config and their firing. Currently they use alertmanager's web interface - so we'll need an equivalent.

Maybe we could run alertmanager itself, purely for this purpose, but rely on AMP's alertmanager for actually sending them? But they would show up as inactive, and might expose differences in functionality.

Or maybe we can give users read-only access to the console, for their team's SNS.

#### Workspace as a service?

We could offer users a Prometheus workspace to themselves - a full monitoring stack that they fully control. Just a terraform module they can run.  Maybe this is better for everyone, than a centralized one, or just for some specialized users - do some comparison?
