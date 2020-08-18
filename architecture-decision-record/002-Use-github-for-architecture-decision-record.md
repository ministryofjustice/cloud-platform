# Use GitHub for architecture decision log

Date: 26/03/2018

## Status

✅ Accepted

## Context

The cloud platforms team has had a number of discussions about where to hold documentation. We have tried using confluence for technical documentation but it has largely gone stale through lack of updates.

For the development of the new platform we want to keep _technical_ documentation close to the code that implements that documentation. An example of this is the [kubernetes investigations](https://github.com/ministryofjustice/kubernetes-investigations) repo which holds our experiments into the use of kubernetes.

Putting technical documentation in GitHub has additional benefits:

* Using pull requests gives us a clear review and approval process
* It is part of the same workflow as other activities that we do on a day to day basis (e.g. writing code)
* The information can be held in the open and viewed by anyone

## Decision

1. Our architecture decision log for the new cloud platform will be held in GitHub

## Consequences

Our existing architecture decisions will need to be migrated from Confluence (see this [pull request](https://github.com/ministryofjustice/cloud-platform/pull/1)).
