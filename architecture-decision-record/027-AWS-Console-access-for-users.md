# 27 AWS Console access for users

Date: 2021-11-03

## Status

🤔 Proposed

## Proposal

Give users (developers in service teams) access to the AWS Console of the Cloud Platform AWS account, so they could view their team's resources.

Implementation:

* AWS resources are tagged by the k8s namespace name, to which they are attached (simply added into our Terraform modules)
* Users would sign into the console using AWS SSO
* CP would set-up authorization in [ABAC](https://docs.aws.amazon.com/singlesignon/latest/userguide/abac-checklist.html), to give users access to their team's resources. Specifically it connects the user's GitHub team(s) with resource tags, using the mappings in the environment repo.

This would give teams console access to their backing resources, including RDS database, ECR, S3, etc.

TODO decide if they are somehow allowed write access - see discussion below.

Note: they would not have access to platform infrastructure, including the cluster which runs their containers, or the NLB where their user traffic ingresses. These are shared resources.

## Context

A key value of CP is to support teams to use **infrastructure as code**, rather than **ClickOps**.
IaC benefits:

* Infrastructure that you tested in the dev environment is going to be exactly the same in production. Whereas ClickOps risks deploying something different
* You can recreate infrastructure from scratch, confident their config is correct. Otherwise you gradually find yourself with 'configuration drift', meaning it is 'snowflake infrastructure', which you cannot risk recreating or recover from failure

With this in mind, CP has been designed so that users define their AWS resources in terraform. This pushes them to use the best practice of 'infrastructure as code', and they are supported in that practice, because CP team writes and maintains most of the terraform modules that they call, and maintain the CI/CD to apply it.

## User needs

"AWS Console Access" can be provided in chunks, outlined below. Each are described with the user needs that have been gathered from chats with users and TA thoughts, pending the proper user research.

### **Read-only** access to **user's AWS resources**

Users could see their RDS database etc. (This doesn't include the machine/networking where their app runs, nor the ingress)

* **Learning and understanding** - The console is much more user-friendly than terraform (or the CLI/API). This is good for understanding what is configured, learning about configuration options, and generally learning AWS/cloud. The console communicates a lot about how the options fit together, current state and show  guidance and warnings, so are great for building a mental model of a service. Tutorials from AWS, CloudGuru, etc use the console, for even the most advanced services, so using the console to learn is a strong convention.
* **Visual feedback** - During usage, some services have state that is particularly useful to view in the console - RDS Performance Insights and Step Functions provide information in a very visually intensive way, and it is unreasonable to try and get the same experience via CLI or API.

### **Read-only** access to **platform infrastructure**

Users could see config of the machine/networking where their app runs, and ingress (load balancer, WAF, DNS)

* **CP understanding** - Users benefit from seeing behind the CP interfaces, to see how the platform works. Reasoning:
    * for interests sake and their learning of cloud
    * knowing how something works enables you to use it better
    * so that they can evaluate it better from a security point of view, seeing the config of VPC, security groups, etc, and the console is how they are used to viewing this things

### **Write** access

Users could change the to:

* **user's AWS resources** - their RDS database etc.
* **ingress** - DNS, WAF

Use cases:

* **Quick iteration** when you're still trying stuff out or getting something to work, e.g. changing a port number, DNS setup
* **One-off operations**, such as upgrading a database version, triggering a database snapshot or restore, DNS cut-overs

If giving users write access to the console, the concern would be not to allow long-lasting configuration changes.

Perhaps a solution would be:

* quick iterations would be allowed, as long as there is an understanding that at the end of the day or week it will be automatically reverted to what the terraform has configured
* one-off actions could maybe be identified and allowed

## Design - Read-only access to user's AWS resources

One way to do it:

* Users login via SSO, with a role for their team
* AWS resources should already be tagged by the owning team.
* IAM created to allow each team role to access resources tagged against that team - https://docs.aws.amazon.com/IAM/latest/UserGuide/access_tags.html
* Concourse pipeline for creating the roles and permissions

There may be a more sophisticated way - maybe AWS folks can suggest.

## Progress

21/1/22 Team is confident that read-only access is a user priority, so focussing on this. User research is in progress to evidence and flesh out the user needs. Work is underway to design how it could work.
