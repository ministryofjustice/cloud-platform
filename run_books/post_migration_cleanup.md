# Post migration cleanup
The intention of this document is to provide a run book for all activities required post migration to the Cloud Platform. This is an instructional document and will not provide a detailed description of each task. 

## Table of contents
  - When to use this doument?
  - Turn off monitoring and alerting
  - Archive the deployment repository
  - Backup data
  - Delete the old Cloudformation stack
  - Ensure removal of resources
  - Set lifecycle on Buckets
  - Ensure live service is still functioning!

## When to use this doument?
After a service has migrated to the Cloud Platform and is considered live. In this instance live refers to a service being used...
## Turn off monitoring and alerting
Cloud watch, how to remove all alerts from cloudwatch, maybe pics...
## Archive the deployment repository
This is the *-deploy GH repositories. How to archive and what messages to place...
## Backup data
Backup RDS and potentially S3 data
## Delete the old Cloudformation stack
Literal deletion of the CF stack using the console. No need to use saltStack
## Ensure removal of resources
Ensure the resources have been removed. Maybe by searching tags...
## Set lifecycle on Buckets
If the service used S3 buckets, set a lifecycle policy that'll remove old data on the bucket. How do you do this?
## Ensure live service is still functioning!
What it says on the tin.

