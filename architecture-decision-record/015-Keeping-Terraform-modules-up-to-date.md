# Keeping Terraform modules up to date

Date: 23/03/2020

## Status

✅ Accepted

## Context

We maintain a lot of [terraform modules] which teams use to manage AWS
resources for their namespaces.

In the past, lettings different namespaces use different versions of terraform
modules has caused problems because there has not always been a clear upgrade
path from a module that is several versions old, to the latest version.

In these cases, it would have been easier to have upgraded every namespace when
the significant change was made to the module, because at that point in time
the relevant knowledge and understanding are fresh. Leaving it until later
makes the whole process much more difficult.

## Decision

We decided to ensure that all namespaces are always using the latest version of
every module.

## Consequences

- We have to upgrade the module version in every affected namespace whenever we release a new version of a terraform module

[terraform modules]: https://github.com/ministryofjustice/cloud-platform#terraform-modules
