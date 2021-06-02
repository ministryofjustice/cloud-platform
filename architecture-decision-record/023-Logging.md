# Logging

Date: 02/06/2021

## Status

✅ Accepted

## Context

Cloud Platform's existing strategy for logs has been to **centralize** them in an ElasticSearch instance. This allows [service teams](https://user-guide.cloud-platform.service.justice.gov.uk/documentation/logging-an-app/access-logs.html#accessing-application-log-data) and Cloud Platform team to use Kibana's search and browse functionality, for the purpose of debug and resolving incidents. All pods' stdout get [shipped using Fluentbit](https://user-guide.cloud-platform.service.justice.gov.uk/documentation/logging-an-app/log-collection-and-storage.html#application-log-collection-and-storage) and ElasticSearch stored them for 30 days.

Concerns with existing ElasticSearch logging:

* ElasticSearch costs a lot to run - it holds indexes in memory
* CP team doesn't need the power of ElasticSearch very often - rather than use Kibana to look at logs, the CP team mostly uses `kubectl logs`
* Service teams have access to other teams' logs, which is a concern should personal information be inadvertantly logged

With these concerns in mind, and the [migration to EKS](022-EKS.html) meaning we'd need to reimplement log shipping, we reevaluate this strategy.

## User needs

This is what we need logging for:

**Observing apps** - service teams need to see the logs from their apps and ingress, to be able to test new app features and debug issues (e.g. app crashes on start-up). They need to browse and search. Some have asked about JSON logs and using them for dashboards.

**Observing the platform** - CP team need to see the logs from the k8s control plane, k8s components, pipeline and other platform services, from the past few days, to be able to test new platform features and debug issues.

**KPI dashboards** - CP team needs to monitor KPIs, which are machine-readable. It's likely that some of these KPIs will include counts of logged events, such as Auth0 log-ins. (KPIs will probably need to cover other things too, such as grafana monitoring, number of GitHub repos, cloud costs.)

**Security** - service teams, CP team and Security & Privacy team need to monitor and investigate attacks. To quickly react to threats, we need to be able to do a quick search across the platform's logs for an IP or URL pattern, rather than logs de-centralized. We work to provide the Secuity & Privacy to have this on MLAP (MOJ Logging Aggregation Platform), their centralized SIEM (Security Information and Event Management) platform. (SIEMs do threat detection, using constantly updated lists of threats, such as known bad actor IP addresses, hosts, OWASP URL patterns, etc.) For investigation, logs need retaining for 13 months.

## Proposal

As we migrate to EKS, rather than centralized logging, we propose several logging solutions, to fit the varied needs identified.

**Observing apps** - Loki. For occasional searches, a disk-based index seems more appropriate - higher latency than memory, but much lower cost to run. Could setup an instance per team.

**Observing the platform** - CloudWatch is where both EKS and AWS components send platform logs. There may be a driver (such as KPIs) that say we need to ship these to Elastic or Loki in the future. A CloudWatch retention policy can be set to match needs.

**KPI dashboard** - TBD

**Security** - MLAP is designed for this.
