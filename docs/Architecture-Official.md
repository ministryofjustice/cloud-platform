# The Cloud Platform - Architectural Overview

## Purpose

This document is to provide an overview and breakdown of the individual components of the Cloud Platform.

## Principles

### Platforms scale, ops teams don't

We want to minimize the need for any 'central web ops team', which does bespoke work, inevitably for only a handful of service teams. Instead we believe in platforms, which offer a standardized approach to hosting. A good platform includes lots of automation and a good service wrapper, to meet the common needs of scores of service teams.

Web ops people should mostly work on a platform, delivering value at scale.

### A high tide raises all boats

We embed ops tools and best practices into a platform. This increases service teams' ability to operate their services and promotes high quality engineering.

Service teams should feel the platform supports them to deliver, without being constrained.

For example in Cloud Platform we offer tools: monitoring, logging, alerts, security scanning, and best practices: Infrastructure as Code, GitOps.

### Self service

Providing service teams with the ability to get started with minimal friction means minimizing waiting for humans to decide or provide. Self-service encourages users if they just want to try things out, to spin something up to test a theory, a different config or different infrastructure.

Self-service requires good documentation, automated processes, and maintaining sufficient oversight.

### Cattle not pets

Apps are cattle. When a service team hosts its software on Cloud Platform, it needs to accept it will be treated like cattle not pets. The containers/pods will be herded between server nodes for lots of reasons, such as nodes being drained due to a planned recycling schedule, or due to failures. This approach is great for scaling, efficiency and high availability (HA). But it needs the software containers to be designed as a [12 Factor App](https://12factor.net/) and data/state held in cloud storage.

Infrastructure is cattle. Cloud Platform team should use a 'delete and replace' process to perform upgrades, fix problems and disaster recovery.

Clusters are cattle. We can't have one cluster, vulnerable to a fault, YOLO changes and upgrades. Clusters need to be disposable. Replacement is automatic and seamless.

### Constant reinvention

As hosting technology and ideas change and improve, we will relentlessly reinvent all aspects of the platform to serve user needs better. We shall not be afraid to replace the whole platform, if that serves users best. A strong Cloud Platform team has more

MOJ from made the leap from Template Deploy to Cloud Platform, with a lengthy transition, but provided a step change in benefits - see [Comparison with Template Deploy](#comparison-with-template-deploy).

## Comparison with Template Deploy and earlier hosting

### 2012-2015 Embedded ops

When MOJ Digital was formed, it started developing services on the public cloud. Each service team had an embedded Web Ops engineer, who was encouraged to develop and adopt common standards across teams, and document their ops processes in Confluence, however there was a fair amount of divergence.

### 2015-2018 Template Deploy

Template Deploy introduced standardization of infrastructure definitions and scripting, so that infrastructure could be managed and supported uniformly. Also the Web Ops engineers were brought together into a central ops team.

Web Ops provided:

* CloudFormation templates that defined a standard infrastructure for a typical web app - EC2s with containers, running in a VPC with ELB, ASG, with S3 for assets and RDS backing.
* Fabric scripts that deployed the config+templates
* A centralized Jenkins instance, for continuous deployment

Service teams:

* recorded config values in a YAML file.
* deployed either from the command-line or using Jenkins.

Web Ops and Service teams shared access to a team's AWS account, or a central one, to operate the service.

Reflections on Template Deploy:

* Shared ownership of the infrastructure makes it hard - both service teams and web ops team have to understand each others' work and coordinate a lot.
* Not enough rigour - developers shouldn't be shelling into VMs or making changes on the AWS console - deployment should be only with CI/CD, but still allow service teams to operate e.g. restart containers, snapshot a database etc
* We rolled our own tool, but now there are established tools
* Cost inefficient - VMs are fixed size, so not used effectively when traffic is low
* The architecture of services is not very flexible from the standard template, doesn't allow for different architectures or hosting COTS s/w
* It's a collection of tools, not a platform - there's no central list or view of deployments, needed to manage them centrally

### 2018-present Cloud Platform

Our vision for cloud native infrastructure was expressed in the [MOJ Hosting Strategy](https://mojdigital.blog.gov.uk/2018/10/15/how-were-making-our-hosting-simpler-more-cost-effective-and-more-modern/): infrastructure that's able to be managed all at once, with clear separation between the applications and the infrastructure (using containers), is resilient to failure, and can easily scale.

**Kubernetes** - commoditizes a lot of the undifferentiated heavy lifting to provide a hosting environment. It is a commodity container platform designed to support multiple services will naturally lead to less time being spent on low-level tasks”. It's also the dominant standard for container orchestration. For more about the choice of Kubernetes for orchestration see: [ADR004 Use kubernetes for running containerised applications](/architecture-decision-record/004-use-kubernetes-for-container-management.md)

**Platform** - hosting services in one place is easier to manage and operate.

**Clear boundary of responsibilities** - The CP team are responsible for running containers with the specified arguments, attached volumes, port bindings and hostnames specified by the service team. They are not responsible for anything inside the container. Services are never so tightly coupled to infrastructure primitives that service teams require web operations assistance to reconfigure things at that level.

See also: [History of Cloud Platform](https://docs.google.com/document/d/1Wvyz7SHipAAfZ2idPArKSs977UtubQz4zzHjByEEOuI/edit#)

## Architecture Diagram

This is a high-level overview of the components and services which comprise the Cloud Platform.

![Cloud Platform Architecture](images/cloud-platform-architecture-diagram.png)

Source for this diagram: [Architecture Diagram](https://docs.google.com/drawings/d/1QQpTN8i2n3QZwIELTTbnxpTNy83eP0T50nVv_2aLx5g/edit?usp=sharing)

## Components

As seen in the diagram, there are a few different components that the architecture is comprised of.
This section will break-down these components and detail the context in which they are used.

### Apply Pipeline in Concourse

The 'Apply Pipeline' is the key part of Cloud Platform that deploys the Kubernetes and AWS resources, that a service team has defined in code.

These jobs are managed by Concourse, which runs on the Manager cluster.

The process starts with the user pushing some code to GitHub, in the . Once merged to master, the pipeline is triggered and it applies it. This results in resources in Kubernetes and AWS being created or changed.

![Pipeline diagram](images/arch-dia-v1.png)

* Watches the main branch of the [environments repo](https://github.com/ministryofjustice/cloud-platform-environments) for merged changes to tenants' environment folders.
* Triggers the pipeline in [Concourse](https://concourse.cloud-platform.service.justice.gov.uk/)
* The Namespace Config YAML gets applied to Kubernetes, typically to create a new namespace.
* The Terraform files are applied to AWS, which creates or changes AWS resources, e.g. to create an RDS instance.
* If the pipeline job fails, then Concourse creates a Slack notification in [#lower-priority-alarms](https://mojdt.slack.com/archives/C8QR5FQRX)

### Terraform

Terraform is our chosen format for defining infrastructure as code (IaC).

AWS resources are defined and maintained in Terraform files, written by developers and approved by the Cloud Platform team.

Concourse is able to interpret Terraform files and apply the changes to the relevant AWS resources.

### GitHub

GitHub is where all of our code repositories that are used by pipeline are stored.

Concourse is configured to watch the 'master' branches of our selected application repos.

The GitHub process flows as follows:

1. A developer will make a change to a Terraform file/s in one of our repos and then raise a pull request (PR).
2. A member of the Cloud Platform team will review the PR. The developer is notified of the approval or otherwise.
3. The developer merges the code change into the 'master' branch.
4. The code change in 'main' is noticed by Concourse, and the pipeline is triggered.

### AWS

AWS is where all of our infrastructure is hosted.

AWS is the end-point of the pipeline and destination where we expect to see the resources defined in Terraform to match what we intended before it was sent through the pipeline.
