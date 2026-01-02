# 36. Decommission the Cloud Platform CLI Tool

Date: 2026-01-02

## Status

✅ Accepted

## Context

The Cloud Platform CLI Tool is used for various tasks, such as generating namespace files and running Terraform. Although useful, there are other ways in which these tasks can be achieved without the need for a separate tool which we need to maintain.

## Decision

We will decommission the CLI Tool.

## Consequences

- We will build new functionality to replace the existing tool
- We will follow existing patterns where possible, eg Modernisation Platform use GitHub Actions to generate user files
