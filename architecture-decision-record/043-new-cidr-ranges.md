# 43. New Cloud Platform CIDR ranges

Date: 2026-02-13

## Status

✅ Accepted

## Context

To build the new live and nonlive clusters we need to use an RFC 1918 CIDR range that is unique so that we can route traffic from and to other areas of the MoJ via the MoJ Transit Gateway.

## Decision

We have been allocated the following CIDR ranges from the MoJ Networking team:

```
10.41.0.0/16
10.195.0.0/16
```

These are recoreded in the MoJ IPAM solution.  More details can be found [here](https://justiceuk.sharepoint.com/:f:/r/sites/CloudPlatform/Shared%20Documents/Cloud%20Platform/CIDR%20Range%20allocations?csf=1&web=1&e=FOyrbd)

## Consequences

- We will use the above CIDR ranges for the new cluster
