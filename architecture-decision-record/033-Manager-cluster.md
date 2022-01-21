# 33 Manager cluster

Date: 2022-01-07

## Status

✅ Accepted (documenting an existing thing)

## The existing decision

We have a 'Manager cluster', for running infrastructure CI/CD (Concourse), which is separate to the clusters that run user workloads.

## Context

On the Manager cluster we run:

* Concourse
* Thanos compactor
* Hoodaw report generators

Reasons to run these workloads separate to the main cluster:

1. **Recovering the main cluster** - Concourse deploys infrastructure, including EKS. So if it runs on the same cluster that it creates, we may encounter a "chicken and egg problem". For example if the live cluster is damaged, then we would want to reapply that infrastructure with Concourse, so it makes sense to run it on another cluster. (This arrangement however doesn't cover some other scenarious - it runs in the same AWS account and VPC, so might not recover some issues with those things, or bootstrap a fresh install.)

2. **Test cluster** - It seems genuinely useful to test applying changes to the Manager cluster before they are applied to Live, particularly significant ones like kubernetes upgrades. We already test on a (personal) Test cluster, but the Manager cluster runs more life-like workloads, and there is more pressure to apply changes to it without downtime.

3. **Security isolation of CI/CD** - it's a good security practice to run CI/CD separately to where the app runs. CI/CD tends to need wide permissions to build containers, and to deploy things. Admin access to CI/CD needs to be only for admins, whereas apps are often public. There is a clear difference in needs for Concourse and the apps, so to minimize risks and attack surface, it makes sense to separate off Concourse from the live cluster(s).
