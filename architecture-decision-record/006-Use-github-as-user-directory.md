# Use GitHub as our identity provider

Date 09/04/2018

## Status

✅ Accepted

## Context

As part of our [planning principles](https://docs.google.com/document/d/1kHaghp-68ooK-NwxozYkScGZThYJVrdOGWf4_K8Wo6s/edit) we highlighted "Building in access control" as a key principle for planning our building our new cloud platform.

Making this work for the new cloud platform means implementing ways that our users &mdash; mainly developers &mdash; can access the various bits of the new infrastructure. This is likely to include access to Kubernetes (CLI and API), AWS (things like S3, RDS), GitHub, and any tooling we put on top of Kubernetes that users will access as part of running their apps (e.g. ELK, [Prometheus](https://github.com/ministryofjustice/cloud-platform/blob/main/architecture-decision-record/026-Managed-Prometheus.md#choice-of-prometheus), [Concourse](https://github.com/ministryofjustice/cloud-platform/blob/main/architecture-decision-record/003-Use-Concourse-CI.md)).

At the current time there is no consistent access policy for tooling. We use a mixture of the Google domain, GitHub and AWS accounts to access and manage the various parts of our infrastructure. This makes it hard for users to make sure that they have the correct permissions to do what they need to do, resulting in lots of requests for permissions. It also makes it harder to manage the user lifecycle (adding, removing, updating user permissions) and to track exactly who has access to what.

We are proposing that we aim for a "single sign on" approach where users can use a single logon to access different resources. For this we will need a directory where we can store users and their permissions, including what teams they belong to and what roles they have.

The current most complete source of this information for people who will be the first users of the cloud platform is GitHub. So our proposal is to use GitHub as our initial user directory - authentication for the new services that we are building will be through GitHub.

## Decision

We will use GitHub as the identify provider for the cloud platform.

We will design and build the new cloud platform with the assumption that users will login to all components using a single GitHub id.

## Consequences

We will define users and groups in GitHub and use GitHub's integration tools to provide access to other tools that require authentication.

When adding new functionality or tooling to the platform we will design for using GitHub for authentication and authorisation. We will prefer tooling that integrates with GitHub directly, or has support for standards-compliant OIDC authentication (or SAML if that is not available) as Auth0 allows us to use Github identities as a generic authN/authZ source using those protocols. We will assess other solutions based on their ease of integration with GitHub.
