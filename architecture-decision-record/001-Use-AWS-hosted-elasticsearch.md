# Use AWS hosted Elasticsearch

Date: 05/02/2018

## Status

✅ Accepted

## Context

The cloud platforms team self-host an Elasticsearch cluster with Kibana and Logstash (ELK). This cluster has suffered numerous outages (see [CPT-282](https://dsdmoj.atlassian.net/browse/CPT-282) and [CPT-152](https://dsdmoj.atlassian.net/browse/CPT-152) in Jira) that have been difficult to recover from.

Reasons behind this move were:

* Average of almost one week per month spent on debugging, fixing and reviving ELK
* Lengthy downtimes which made data recovery pointless
* Self hosted ELK stacks cost was significantly higher than AWS ElasticSearch solution
* Not working ELK cluster was also a blocker for product teams as they couldn't see any application logs

## Decision

Replace our self hosted ELK stack with the managed AWS Elasticsearch

## Consequences

    1. Application logs were already shipped directly from the servers to AWS CloudWatch Logs ( thanks to works on Service Catalogue project )
    2. We have created lambda function for shipping logs into newly created AWS ES cluster
    3. We have created a script to set lambda function mentioned above as a target for logs from AWS CloudWatch Logs ( script@opstools )
    4. We have deleted the ELK stacks from all environments
    5. We have modified Sensu alerting to use new ES cluster
    5. Cluster healthchecks set up via Cloudwatch to ensure reliability.

Product teams have confirmed their application logs being shipped and searchable in new ELK cluster with as of now - zero downtime and maintenance from Cloud Platforms side.

We don't have any admin power over the cluster, unable to install either kibana plugins or change ES connectivity settings.
