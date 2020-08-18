# Migrate from live-0 to live-1 Cluster

Date: 28/02/2019

## Status

✅ Accepted

## Context

Migrating from live-0 to live-1 cluster. The reason behind this is based on the need to move to a dedicated AWS account (moj-cp), which will be much easier to support, and the need to move away from the Ireland (EU) region to the London (UK) region as Cloud Platform requirement to host data in the UK, rather than in Europe.


## Decision

After some long consideration of possible options, the decision has been made to migrate from the live-0 cluster to the new live-1 cluster.

Since we only want to be running a single cluster, we will need to shut down live-0 as soon as it's no longer needed. Also services migrate from live-0 to live-1 sooner will avoid the complexities of running two parallel clusters.

## Consequences

Creation of the live-1 cluster, which includes a move to the moj-cp AWS account, move to the eu-west-2 region and below tasks performed in live-1 cluster.

1. In Live-1, Introduce Pod Security Policies, to tighten the security on production cluster.

2. Live-1 Kubernetes version bump to 1.11.8 using kops, this is for security and feature implementation. 1.11.8 is the latest version kops will allow. live-0 is using 1.10

3. Change tenant to justice-cloud-platform.

4. Change all references to the integration.dsd.io to cloud-platform.service.justice.gov.uk.

5. Change the kops store to the moj-cp AWS account under s3://cloud-platform-kops-state. All S3 buckets need to be moved over to the new account.

6. In live-1 Amend all region information to eu-west-2 (London).

7. Update EBS version to 3.3.10 and Encrypted EBS volumes

8. Create Elasticsearch cluster in new AWS account live and audit

9. Remove nginx-ingress-controllers and made nginx-ingress-acme default, create tls rules for all Cloud Platform applications

10. Add new kuberos image.

11. Move pipelines running in live0 concourse cluster to live1 concourse cluster.
