# 36. Use GitHub Actions

Date: 2026-01-02

## Status

✅ Accepted

## Context

Cloud Platform uses a self managed Concourse instance hosted in a management EKS cluster. This runs all of the Cloud Platform deployment pipelines.

The Cloud Platform will be moving to the Modernisation Platform (see ADR 39), the Modernisation Platform uses GitHub Actions for deployment pipelines.

## Decision

We will switch to using GitHub Actions for our deployment pipelines.

## Consequences

- We will decomission our Concourse pipelines once the move to Modernisation Platform has been completed
- We will decomission our management cluster once the move to Modernisation Platform has been completed
- We will use the Modernisation Platform pipelines where possible
- We will create new GitHub Actions in a our repositories where appropriate
