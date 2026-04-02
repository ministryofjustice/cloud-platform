# 46. Use ECR for all User Images

Date: 2026-04-02

## Status

✅ Accepted

## Context

We currently allow users to host their Docker images anywhere, this has some implications:
 - We can't scan images until they are already deployed
 - We need to allow external deployment pipelines to access the cluster endpoint
 - We can't block images with vulnerabilities
 - We can't redeploy user images in a disaster scenario

## Decision

All user images will be hosted in ECR in the Container Platform

## Consequences

 - Users will need to change their deployment patterns
 - We will need to develop new deployment patterns for users (eg GitOps)
