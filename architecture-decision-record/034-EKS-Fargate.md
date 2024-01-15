# 34 EKS Fargate

Date: 2022-01-21

## Status

🤔 Incomplete

## Proposal

Move from EKS managed nodes to EKS Fargate.

## Discussion

This is really attractive because:

- to reduce our operational overhead
- improve security isolation between pods (it uses Firecracker, so we can stop worrying about an attacker managing to escape a container).

However there’s plenty of things we’d need to tackle, to achieve this (copied from [ADR022 EKS - Fargate considerations](https://github.com/ministryofjustice/cloud-platform/blob/main/architecture-decision-record/022-EKS.md#future-fargate-considerations)):

**Pod limits** - there is a quota limit of [500 Fargate pods per region per AWS Account](https://aws.amazon.com/about-aws/whats-new/2020/09/aws-fargate-increases-default-resource-count-service-quotas/) which could be an issue, considering we currently run ~2000 pods. We can request AWS raise the limit - not currently sure what scope there is. With Multi-cluster stage 5, the separation of loads into different AWS accounts will settle this issue.

**Daemonset functionality** - needs replacement:

- fluent-bit - currently used for log shipping to ElasticSearch. AWS provides a managed version of [Fluent Bit on Fargate](https://aws.amazon.com/blogs/containers/fluent-bit-for-amazon-eks-on-aws-fargate-is-here/) which can be configured to ship logs to ElasticSearch.
- prometheus-node-exporter - currently used to export node metrics to prometheus. In Fargate the node itself is managed by AWS and therefore hidden. However we can [collect some useful metrics about pods running in Fargate from scraping cAdvisor](https://aws.amazon.com/blogs/containers/monitoring-amazon-eks-on-aws-fargate-using-prometheus-and-grafana/), including on CPU, memory, disk and network

**No EBS support** - Prometheus will run still in a managed node group. Likely other workloads too to consider.

**How people check the status of their deployments** - to be investigated

**Ingress can't be nginx? - just the load balancer in front** - to be investigated. Would be fine with AWS Managed Ingress

If we don't use Fargate then we should take advantage of Spot instances for reduced costs. However Fargate is the priority, because the main driver here is engineer time, not EC2 cost.
