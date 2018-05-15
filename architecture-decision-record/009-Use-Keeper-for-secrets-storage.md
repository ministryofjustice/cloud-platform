# Use Keeper for secrets storage

Date: 15/05/2018

## Status

Proposed

## Context

The cloud platforms team is transitioning from an in-house deployed instance of Rattic (an open-source web app cloned from https://github.com/tildaslash/RatticWeb) to a managed service licensed from Keeper Security, Inc.
Once the DB migration and audit is done, all product teams will be able to use it, given that:

- The new service will be available without VPN at https://keepersecurity.com
- SSO is supported via Google's OAUTH and others
- MFA is supported and can be enforced
- RBAC and audit functions are much improved

In parallel, Rattic's functionality will be documented, and the data backed up before switching the app to read-only mode, followed by decommission once we've had positive feedback from all teams.

## Decision

We will be transitioning all the existing users and credentials from Rattic; progress tracked via https://waffle.io/ministryofjustice/cloud-platform?search=keeper

## Consequences

1. For the entire MoJ DS, secrets will be stored in Keeper
2. All the users will be mapped to Google accounts
3. Secrets will be consumed via a new web interface, mobile (Android, iOS) apps and a new API that will be documented and monitored by the Cloud Platform team
4. Platform maintenance and backups will be the responsibility of the service provider (Keeper Security, Inc) once the setup is validated by the following **add relevant parties here**
5. Secrets management will still be a shared responsibility for all teams
