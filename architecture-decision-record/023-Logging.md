# 23 Logging

Date: 02/06/2021

## Status

✅ Accepted

## Context

Cloud Platform's existing strategy for logs has been to **centralize** them in an ElasticSearch instance (Saas hosted by AWS OpenSearch). This allows [service teams](https://user-guide.cloud-platform.service.justice.gov.uk/documentation/logging-an-app/access-logs.html#accessing-application-log-data) and Cloud Platform team to use Kibana's search and browse functionality, for the purpose of debug and resolving incidents. All pods' stdout get [shipped using Fluentbit](https://user-guide.cloud-platform.service.justice.gov.uk/documentation/logging-an-app/log-collection-and-storage.html#application-log-collection-and-storage) and ElasticSearch stored them for 30 days.

Concerns with existing ElasticSearch logging:

- ElasticSearch costs a lot to run - it uses a lot of memory (for lots of things, although it is disk first for the documents and indexes)
- CP team doesn't need the power of ElasticSearch very often - rather than use Kibana to look at logs, the CP team mostly uses `kubectl logs`
- Service teams have access to other teams' logs, which is a concern should personal information be inadvertantly logged
- Fluentd + AWS OpenSearch combination has no flexibility to parse/define the JSON structure of logs, so all our teams right now have to contend with grabbing the contents of a single log field and parsing it outside ES

With these concerns in mind, and the [migration to EKS](https://github.com/ministryofjustice/cloud-platform/blob/main/architecture-decision-record/022-EKS.md) meaning we'd need to reimplement log shipping, we reevaluate this strategy.

## User needs

Logs are used for several purposes, so it's worth considering them separately:

**Observing apps** - service teams need to see the logs from their apps and ingress, to be able to test new app features and debug issues (e.g. app crashes on start-up). They need to browse and search. A popular feature request is for structured logs, where each log line is a JSON object, with each field searchable in the log UI. Some others have asked about using them for dashboards.

**Observing the platform** - CP team need to see the logs from the k8s control plane, k8s components, pipeline and other platform services, from the past few days, to be able to test new platform features and debug issues.

**KPI dashboards** - CP team needs to monitor KPIs, which are machine-readable. It's likely that some of these KPIs will include counts of logged events, such as Auth0 log-ins. (KPIs will probably need to cover other things too, such as grafana monitoring, number of GitHub repos, cloud costs.)

**Security** - service teams, CP team and Security & Privacy team need to monitor and investigate attacks. To quickly react to threats, we need to be able to do a quick search across the platform's logs for an IP or URL pattern, rather than logs de-centralized. We work to provide the Security & Privacy to have this on MLAP (MOJ Logging Aggregation Platform), their centralized SIEM (Security Information and Event Management) platform. (SIEMs do threat detection, using constantly updated lists of threats, such as known bad actor IP addresses, hosts, OWASP URL patterns, etc.) For investigation, logs need retaining for 13 months. CP team needs to be involved when findings come in, to add knowledge of the system to the assessment.

## Design discussion

Rather than centralized logging in ES, we'll evaluate different logging solutions, to fit the varied needs identified.

**AWS services for logging** - with the cluster now in EKS, it wouldn't be too much of a leap to centralizing logs in CloudWatch and make use of the AWS managed tools. One one hand it's proprietary to AWS, so adds cost of switching away. But it might be preferable to the cost of running ES, and related tools like GuardDuty and Security Hub, with use across Modernization Platform, is attractive.

### Observing apps\*\*

- Loki - seems a good fit. For occasional searches, a disk-based index seems more appropriate - higher latency than memory, but much lower cost to run. (In comparison, ES describes itself as primarily disk based indexes, but it _requires_ heavy use of memory.) Could setup an instance per team. Need to evaluate how we'd integrate it, and usability.
- CloudWatch Logs - possible and low operational overhead - needs further evaluation.
- Sentry - Some teams have beeing using Sentry for logs, but [Sentry says themself it is better suited to error management](https://sentry.io/vs/logging/), which is a narrower benefit than full logging.

### Observing the platform

CloudWatch is where both EKS and AWS components send platform logs. There may be a driver (such as KPIs) that say we need to ship these to Elastic or Loki in the future. A CloudWatch retention policy can be set to match needs.

### KPI dashboard

TBD

### Security

- MLAP was designed for this, but it is stalled, so probably best to manage it ourselves.
- ElasticSearch does have open source plugins for SIEM scanning. And it offers quick searching needed during a live incident. Maybe we could reduce the amount of data we put in it. But fundamentally it is an expensive option, to get both live searching and long retention period.
- AWS-native solution using GuardDuty and CloudWatch Logs may provide something analogous.

## Next steps

21/1/22 DR suggests investigating centralizing logs in CloudWatch Logs
