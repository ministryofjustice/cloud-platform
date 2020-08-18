# Monitoring/Alerting of ECR Image Counts

Date: 04/07/2019

## Status

✅ Accepted

## Context

We use ECR as the Docker container registry that makes it easy for users to store, manage, and deploy Docker container images.
Due to some applications having a constant rate of images being pushed to their ECR repo, we found that the AWS limit of 1000 images was being hit by some teams. To avoid this we had implemented a lifecycle policy of *100 images* per ECR repo. ECR repositories created for use in the Cloud Platform will have this default lifecycle policy applied.

As lifecycle policy will only keep 100 most recent versions of an image and silently delete images, application users raised an [issue][user-issue] on imposing any limit on number of images is potentially dangerous, unless teams have their own clean-up mechanism.

## Decision

After discussing with application teams and consideration of possible options, the decision has been made to remove the lifecycle policy altogether, but adding monitoring and alerting such that we can take action before an ECR runs out of space.

As it is to do with metrics & alerts, since prometheus is our monitoring solution we decided to use prometheus for metrics & alerts.

## Consequences

* We are gathering metrics about the number of ECRs we have, and the number of images in each of them
* We have an 'amber' alert when any ECR contains more than 500 images
* We have a 'red' alert when any ECR contains more than 800 images
* We have an 'amber' alert when the number of ECRs rises above 500
* We have a 'red' alert when the number of ECRs rises above 800
* When alerts are trigerred for Image count, Inform relavent service teams to manage their repositories by cleaning up some old images.
* When alerts are trigerred for ECR count, raise a request for AWS to increase the count of repositories.

[user-issue]: https://github.com/ministryofjustice/cloud-platform/issues/839
