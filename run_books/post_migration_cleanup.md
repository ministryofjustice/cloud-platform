# Post migration cleanup
The intention of this document is to provide a run book for all activities required post migration to the Cloud Platform. This is an instructional document and will not provide a detailed description of each task. 

## Table of contents
  - When to use this document?
  - Turn off Template Deploy monitoring and alerting
  - Archive the <service>-deploy repository
  - Backup old Template Deploy data
  - Delete the old Cloudformation stack
  - Ensure removal of resources
  - Set lifecycle on Buckets
  - Ensure live service is still functioning!

## When to use this doument?
Once a service has been migrated to the Cloud Platform from Template Deploy there will be a number of cleanup tasks required to ensure legacy resources have been removed. These tasks will be listed in this document (below). Thoughout this process we recommend making a note/record of deleted resources, including things like S3 buckets and RDS databases. This will allow you to reference your actions to service teams, if required.

## 0 Pre-requisites
There are a couple of suggested pre-requisites before performing the tasks in this document.

#### 0.1 Identify AWS account
Traditionally, production Template Deploy applications lived in the [mojdsd AWS account](). However, this isn't always the case. You'll need to identify which AWS account your application lives in. The recommended approach for this task is to speak to service teams Product Managers.

#### 0.2 Ensure newly migrated application works
Alot of tasks in this document will involve destroying resources, it is vitally important you are sure 100% of production traffic is being routed to the Cloud Platform. 

## 1 Turn off Template Deploy monitoring and alerting
Monitoring and alerting for Template Deploy was achived via a seperate Cloudformation stack using a template from [another project](https://github.com/ministryofjustice/MOJ-Service-Catalog/blob/master/submodules/cloudwatch-legacy-monitoring.template). This stack mainly consisted of Cloudwatch alarms and dashboards.

#### 1.1 Identify the Cloudformation stack name
Access the AWS console and open Cloudformation. Identify the stack name, which has a suffix of `-monitoring`, for example `graphite-monitoring`.

#### 1.2 Delete Cloudformation stack
Cloudformation stack deletion is fairly trivial. In the AWS console, select your stack, click actions and then delete. This will take a few minutes to complete, but will disappear from your list of available stacks. 


## 2 Archive the deployment repository
Traditionally there will be a Template Deploy GitHub repository usually named in the format <service>-deploy (e.g. `graphite-deploy`). This repository is redundant and can be archived.

All actions in this task are performed in the `ministryofjustice` GitHub account. 

#### 2.1 Identify the appropriate repostory name
As stated above, perform a search in the `ministryofjustice` GitHub account for your service name with the `-deploy` suffix.

#### 2.2 Change the repository description
Simply add a `**DEPRECATED**` message at the front of your description and perhaps a message explaining why the repository has been archived. 

#### 2.3 Add a note in the README.md
It is also appropriate to add a note at the top of the `README.md` file signalling the reason for deprecation. 

#### 2.4 Actually archive the repository
GitHub have written a handy article on how to archive repositories. Please read through [this](https://help.github.com/en/articles/archiving-repositories) document.

## 3 Backup old Template Deploy data
Not all Template Deploy applications have an RDS or data storage of some kind. If 

#### 3.1 Identify RDS name in AWS console
This can either by identified by the `DB Instance` name, or via a tag. Either way, you can use the search functionality within the AWS console to identify the correct database. 

#### 3.2 Take snap shot of database
It is vitally important that we keep a snapshot of the database before it's removed. This allows us to perform an emergancy recovery if needed.

To complete this task in the AWS RDS console, select your database > actions > take snapshot. Give the snap shot a meaningful name and select `take snapshot`.

## 4 Delete the old Cloudformation stack

#### 4.1 Identify Cloudformation stack name
Similar to the step above, the Cloudformation for the application stack will follow a logical naming structure, for example, `correspondence-staff-demo`. Find the CloudFormation stack that belongs to your migrated service.

#### 4.2 Delete Cloudformation stack
This step requires to to delete the appropriate stack via the AWS console. Ensure you have the correct stack name (this is vitally important). 

## Ensure removal of resources

#### 5 Search for resources

#### Step 2: Clean resources

Ensure the resources have been removed. Maybe by searching tags...

## Set lifecycle on Buckets

#### Step 1: Identify buckets

#### Step 2: Manually add lifecycle

If the service used S3 buckets, set a lifecycle policy that'll remove old data on the bucket. How do you do this?
## Ensure live service is still functioning!

#### Step 1: Is the service accessible by URL

#### Step 2: Check with service team
What it says on the tin.

