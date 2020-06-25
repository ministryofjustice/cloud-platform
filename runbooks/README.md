# Runbooks

This directory contains the source files for the runbooks hosted at https://runbooks.cloud-platform.service.justice.gov.uk

A CircleCI pipeline compiles the source files from this directory to html, and updates the [target repository].

## Publisher docker image

A github action builds and pushes the [publisher docker image] to docker hub,
whenever a new [release] is created via the github user interface.

[target repository]: https://github.com/ministryofjustice/cloud-platform-runbooks
[publisher docker image]: https://hub.docker.com/repository/docker/ministryofjustice/cloud-platform-runbooks
[release]: https://github.com/ministryofjustice/cloud-platform/releases
