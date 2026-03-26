# 44. New Platform Name

Date: 2026-03-26

## Status

✅ Accepted

## Context

We are rebuilding the cluster in the Modernisation Platform. For a period of time we will have two platforms running.  This creates naming issues with code, DNS, terraform modules and user guidance.

## Decision

We will call the new platform the "Container Platform".

## Consequences

 - We will refer to the new platform as the Container Platform
 - We will create a new set of user documentation for the container platform
 - We will register a new subdomain for `container-platform` and use that for the new platform DNS
 - Modules for the new platform will be named with the format `container-platform-terraform`
 - The team name will remain as "Cloud Platform", it is just the cluster platform product name that is changing.
 - This gives us the option to support other types of platforms with different technologies if the primary technology choice changes
