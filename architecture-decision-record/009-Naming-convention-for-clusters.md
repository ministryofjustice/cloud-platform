# Naming convention for clusters

Date: 05/06/2018

## Status

✅ Accepted

## Context

As we are building our new platform on Kubernetes we have already found the need to build quite a few clusters. These have been for a range of purposes including users, testing new ideas ("sandbox"), testing new functionality ("test"), deploying apps to them ("non-production").

As we are still learning we are finding that:

1. we need to continue building new clusters for different purposes and
2. we often need to test the cluster creation process
3. we want to differentiate between clusters that have users on them and those that are for internal testing purposes
4. we do not want to differentiate cluster by function (e.g. "perf-test", "sandbox") or status ("non-production").

To make this easier we propose having a naming scheme that makes it easy to understand whether users are on that cluster but makes no other assumptions about what it is used for.

## Decision

We will name all clusters with the following naming scheme:

- `live-{n}` for any cluster that have users on them, for instance `live-1`.
- `test-{n}` for any cluster that do not have users on them and are used by the cloud platform team only, for instance `test-2`.

We will number the clusters sequentially.


## Consequences

1. Any new clusters that are built use this naming scheme (the first one using it is `live-1`)
2. We delete the current clusters: `sandbox`, `test` and `non-production` and transition users and functions to new clusters (`non-production` to `live-0`, `sandbox` and `test` to a new cluster called `test-0`)
3. We keep a central record of:
    1. The names of all clusters that we run
    2. If it has a particular function, what each cluster is for
