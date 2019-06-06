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
After a service has migrated to the Cloud Platform and is considered live. In this instance live refers to a service being used...

Throughout this process, we recommend you make a note of resources you delete including S3 buckets. This will allow you to reference your actions to service teams.

## Turn off Template Deploy monitoring and alerting

### Step 1: Identify the Cloudformation stack name

### Step 2: Delete Cloudformation stack

Cloud watch, how to remove all alerts from cloudwatch, maybe pics...

## Archive the deployment repository
Traditionally there will be a Template Deploy GitHub repository usually named in the format <service>-deploy (e.g. `graphite-deploy`). This repository is redundant and can be archived.

### Step 1: Change the repository name

### Step 2: Change the repository description

### Step 3: Add a note in the README.md

### Step 4: Actually archive the repository


This is the *-deploy GH repositories. How to archive and what messages to place...
## Backup data

### Step 1: Identify RDS name in AWS console

### Step 2: Take snap shot of database

Backup RDS and potentially S3 data
## Delete the old Cloudformation stack

### Step 1: Identify Cloudformation stack name

### Step 2: Delete Cloudformation stack

Literal deletion of the CF stack using the console. No need to use saltStack
## Ensure removal of resources

### Step 1: Search for resources

### Step 2: Clean resources

Ensure the resources have been removed. Maybe by searching tags...

## Set lifecycle on Buckets

### Step 1: Identify buckets

### Step 2: Manually add lifecycle

If the service used S3 buckets, set a lifecycle policy that'll remove old data on the bucket. How do you do this?
## Ensure live service is still functioning!

### Step 1: Is the service accessible by URL

### Step 2: Check with service team
What it says on the tin.

