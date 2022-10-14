# 25 Domain names

Date: 2021-07-07

## Status

✅ Accepted

## Introduction

The "cluster-specific" domain names *.apps.live-1.cloud-platform.service.justice.gov.uk need addressing with the new EKS cluster which is not called `live-1` any more. This is also an opportunity to review all our domain names that are part of Cloud Platform.

## Wildcard domain for apps

**Decision:** deprecate *.apps.live-1.cloud-platform.service.justice.gov.uk

**Decision:** use instead *.apps.cloud-platform.service.justice.gov.uk

### Context

Tenants have used the `live-1.` domain pattern for their dev and staging environments, following CP [guidance](https://user-guide.cloud-platform.service.justice.gov.uk/documentation/other-topics/custom-domain-cert.html#using-a-custom-domain). It is easy for teams to configure this sort of domain name, because Cloud Platform provides a wildcard SSL/TLS certificate for them. The team implements this by configuring an ingress, which leads to externel-dns setting the DNS name to point at the NLB.

### Reasoning

Removal of `.live-1` - there is an issue with using a specific cluster name `live-1` with the migration to a new cluster. During the hours of migration, service teams will gradually direct more traffic to the new cluster. This is achieved by having a two DNS entries for the service, created by external-dns, with variable weights. So the cluster name is irrelevant and confusing in this name. So we deprecate it and ask people to remove the `.live-1` from their domain names. During and after This change is a burden for service teams, but it only affects internal users of apps, because production environments [must use a custom domain](https://ministryofjustice.github.io/technical-guidance/documentation/standards/naming-domains.html#cloud-and-modernisation-platform-subdomains), and we will give them several months before removing the name.

`.cloud-platform` - having this in the domain makes it easy to delegate NS for each platform. It can't be abbreviated due to our [Naming Things standard](https://ministryofjustice.github.io/technical-guidance/documentation/standards/naming-things.html

`.service.justice.gov.uk` - is for all pre-beta services and all pre-prod environments, as per [Naming Domains standard](https://ministryofjustice.github.io/technical-guidance/documentation/standards/naming-domains.html)

## Custom app domain e.g. myapp

**Existing decision**: myapp.service.justice.gov.uk for prod environments, pre-beta assessment (no change)

**Existing decision**: myapp.service.gov.uk for prod environments, post-beta assessment (no change)

A custom domain is encouraged for prod environments, rather than the wildcard, because it is shorter and more official looking. If a service passes beta assessment then they move to the GDS-managed `service.gov.uk` domain. This is all to meet the [Naming Domains standard](https://ministryofjustice.github.io/technical-guidance/documentation/standards/naming-domains.html)

Cloud Platform supports for the configuration of custom domains - see guidance: [Using a custom domain](https://user-guide.cloud-platform.service.justice.gov.uk/documentation/other-topics/custom-domain-cert.html#using-a-custom-domain)

## Cloud Platform owned application domains e.g. cloud-platform-reports

### Non-cluster-specific domains

**Existing decision**: reports.cloud-platform.service.justice.gov.uk

### Cluster-specific domains

**Decision:**  
* live:   *.prod.cloud-platform.service.justice.gov.uk
* live-2: *.non-prod.cloud-platform.service.justice.gov.uk

### Reasoning

Removal of `.live` - there is an issue with using a specific cluster name such as `live`. The cluster name is irrelevant and confusing in this name. Changes to these domains due to a change in the name of a cluster creates confusion and unnecessarily changes for users. With an alias (instead if cluster name), cloud platform applications can move clusters without requiring a domain name change. 
We will deprecate and remove the `.live` from domain names. 

Current cluster-specific domains requiring a non-cluster-specific domain will be redirected to a landing page with links to all domains and descriptions of each.

## Cluster-specific domains during migration

**Decision**: *.apps.live.cloud-platform.services.justice.gov.uk for use around migrations

In the period  and during migration, service teams will want to test their service running on a particular cluster. For example they may deploy a copy of their service to a new cluster, and wish to test it, before any public traffic is directed to it. So for this purpose we will provide cluster-specific domains.

A future migration might introduce:

* *.apps.live-blue.cloud-platform.services.justice.gov.uk
* *.apps.live-green.cloud-platform.services.justice.gov.uk

Once the migration is complete we should remove these domains and wildcard cert, because we want service teams to use the non-cluster specific domain (or cluster domain) in general use, so that this traffic can be transferred smoothly to a new cluster in the next migration, and old clusters can be switched off.

## Manager domain

**Existing decision**: *.apps.manager.cloud-platform.services.justice.gov.uk

The manager cluster has a separate pattern of URLs and wildcard cert, keeping it separate from the main clusters.
