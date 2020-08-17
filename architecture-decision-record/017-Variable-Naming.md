# Variable names in YAML & Terraform files are in snake case

Date: 17/08/2020

## Status

🤔 Proposed

## Context

We have a lot of scripts, pipeline definitions, terraform files, yaml files and
templates which need to define and use variables. We want a consistent
convention for naming these so that, as we write code in multiple,
inter-dependent repositories, we can be confident that the names we are using
are correct.

## Decision

We will always use snake case (e.g. `foo_bar`) for variable names which appear
in terraform/yaml files and templates.

## Consequences

We should end up with fewer errors arising from scripts in one repo expecting,
e.g. `repo_name`, but a template in another repo defining `repo-name` instead.
