# 37. Serverless (Lambda/Functions-as-a-Service)

Date: 2022-09-22

## Status

❌ Rejected

## Context

### Overview

The Cloud Platform team often gets questions regarding the use of "serverless" within the context of the Cloud Platform. The Cloud Platform is (at the time of writing) hosted on AWS, and their definition of "serverless" products is:

- AWS Lambda _(Functions-as-a-Service)_
- Amazon API Gateway _(APIs-as-a-Service)_
- Amazon DynamoDB _(Serverless key value store)_
- Amazon EventBridge _(Serverless event bus)_
- Amazon Simple Notification Service (SNS) _(Serverless messaging service)_
- Amazon Simple Queue Service (SQS) _(Serverless queuing service)_
- Amazon Simple Storage Service (S3) _(Serverless file store)_
- AWS AppSync _(Serverless GraphQL / pub/sub)_
- AWS Fargate _(Serverless compute for containers)_
- AWS Step Functions _(Serverless state machine)_

The Cloud Platform does support the usage of some of these services already (namely [DynamoDB](https://github.com/ministryofjustice/cloud-platform-terraform-dynamodb-cluster), [SNS](https://github.com/ministryofjustice/cloud-platform-terraform-sns-topic), [SQS](https://github.com/ministryofjustice/cloud-platform-terraform-sqs), [S3](https://github.com/ministryofjustice/cloud-platform-terraform-s3-bucket)).

The most common query in the context of serverless is the ability for teams to use AWS Lambda, which is the focus of this ADR.

#### What is the Cloud Platform?

The Cloud Platform is a centralised platform, utilising Kubernetes, which is running containerised workloads. This is provided as a _platform-as-a-service_ to the wider Ministry of Justice. See [What is the Cloud Platform?](https://user-guide.cloud-platform.service.justice.gov.uk/documentation/concepts/what-is-the-cloud-platform.html#what-is-the-cloud-platform) for more information.

#### What is AWS Lambda?

[AWS Lambda](https://aws.amazon.com/lambda/) provides the (self-described) ability to "run code without provisioning or managing infrastructure". Lambda supports container images (if they implement the [runtime interface client](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-images.html#runtimes-api-client)) or directly as code.

### Cloud Platform **is** serverless

The widely accepted answer for "what is serverless" is that teams _don't need to manage underlying servers for the services they run_. In this context, the Cloud Platform _is_ in itself serverless. The Cloud Platform team is the provider in this case, who manage the "underlying servers" (in the context of cloud), networking infrastructure, etc. This is in the same way that typical AWS Lambda usage is also serverless.

### Our thoughts on using AWS Lambda/Functions-as-a-Service in the Cloud Platform

When people raise serverless usage with the Cloud Platform team, they usually have different reasons for doing so, so it is impossible to list and answer each question separately, so this is limited to the _common_ thoughts on using AWS Lambda/Functions-as-a-Service in the context of Cloud Platform:

1. Cloud Platform is already serverless

    The Cloud Platform itself is already serverless for users, they don't need to manage the underlying infrastructure or common network configuration between their service and elsewhere.

2. We can't (currently) provide a good developer experience when using AWS Lambda

    Cloud Platform provides read-only access to _users_ of the Cloud Platform, to encourage infrastructure-as-code and standardise infrastructure deployment across the MOJ. AWS Lambda functions need to be uploaded often (when changed, for e.g.), and our current pipeline for _infrastructure_ deployment utilises the [cloud-platform-environments](https://github.com/ministryofjustice/cloud-platform-environments) repository, which requires one of the Cloud Platform team to approve for deployment. This can become frustrating for debugging Lambda or initially configuring it.

3. There is little monetary benefit for using serverless in the context of what Cloud Platform currently hosts

    Small services running on the Cloud Platform, even if idle, are very cheap to run anyway in the wider context of the Cloud Platform. This cost is also (as with all services running on Cloud Platform) absorbed by the Cloud Platform team.

4. We don't currently have safeguards or common, reusable functionality in place for using Functions-as-a-Service

    The Cloud Platform currently has out-of-the-box functionality for basic alerting, monitoring, and logging that has been standardised for use across every service that runs on Cloud Platform. However, we would need to start from scratch to implement feature parity for the use of AWS Lambda (and Functions-as-a-Service generally), which can create a derivation in standardisation. This derivation is linked to the risk appetite across several teams, including the service team, Cloud Platform team, and other stakeholders.

5. There is appetite across MOJ teams to explore the use of Functions-as-a-Service

    It's worth noting that there _is_ growing appetite for teams to start exploring the use of Functions-as-a-Service, and we typically revisit FaaS offerings that can be run on Kubernetes, such as [Kubeless (now archived)](https://github.com/vmware-archive/kubeless), [OpenFaaS](https://docs.openfaas.com/) and [Knative](https://knative.dev/docs/).

6. We treat all infrastructure _and_ services as standard, not unique

    Cloud Platform currently treats _all_ provisioned infrastructure for the underlying platforms as standard and replaceable. A typical requirement for a service to run on Cloud Platform is, for example, it to be stateless as per [12FA](https://12factor.net/): so even the _service_ is treated as standard and replaceable. By introducing AWS Lambda (specifically), we start to treat specific services as unique until we hit a tipping point for standardisation of cloud infrastructure and creates divergence both locally (in specific areas of the MOJ) and more generally (across the whole of the MOJ).

## Decision

Based on the above points, we have ❌ rejected the introduction of AWS Lambda/Functions-as-a-Service into the Cloud Platform. On a point-by-point basis:

1. Cloud Platform is already serverless for our users (i.e. they don't need to manage infrastructure behind the scenes)
2. We can't (currently) provide a good developer experience when using AWS Lambda, which would likely lead to frustration of teams trying to deliver services
3. There is little or no monetary benefit to using Functions-as-a-Service in the context of Cloud Platform
4. We don't have safeguards or common, reusable functionality in place for using Functions-as-a-Service, and we don't have the capability or time to currently start exploring this
5. As there is appetite in the wider MOJ to explore the use of Functions-as-a-Service, Cloud Platform will periodically review Kubernetes integrations with both AWS Lambda (e.g via [AWS Controllers for Kubernetes](https://aws-controllers-k8s.github.io/community/docs/community/overview/)) and more generic Functions-as-a-Service offerings, such as [OpenFaaS](https://docs.openfaas.com/) and [Knative](https://knative.dev/docs/).
6. By introducing Functions-as-a-Service (or AWS Lambda) to Cloud Platform, we need to be cognisant and understand how we can treat the underlying infrastructure for it as standard and replaceable; not unique and hand-built.

Further, in March 2022, we conducted user research around the use of serverless and more explicitly, what needs teams across the MOJ are trying to solve with serverless. This research revealed two important points:

- _some_ aspects of serverless are already provided by the team; and adopting a serverless approach won't (by itself) resolve outstanding needs for users, which are already being considered with non-serverless infrastructure
- further research needs to be taken but in-line with more general approaches rather than focussing on serverless as a solution

## Consequences

By rejecting introducing AWS Lambda/Functions-as-a-Service to the Cloud Platform, we (the Cloud Platform team):

- are asking teams to continue using more "traditional" methods for building their services, such as containerised, stateless, and meeting the [12 Factor app methodology](https://12factor.net/)
- will periodicially review Functions-as-a-Service on Kubernetes to see what's happening in the market
- will carry out further research to meet needs using the best technology to meet user needs, rather than focussing on serverless as a solution
