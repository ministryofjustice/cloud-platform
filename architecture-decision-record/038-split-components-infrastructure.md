# 38. Split Components Infrastructure

Date: 2024-03-04

## Status

🤔 Incomplete

## Context

Following the upgrade to Kubernetes 1.25, PodSecurityPolicy was removed, and the team chose to use Gatekeeper's mutating capabilities to replace the functionality that was previously provided by our PSP configurations. This meant a large change in the way our cluster operates and the modules at a components level. While upgrading it became clear that there are a large number of components dependent on each other that are not required. This makes the component level difficult to work with. We also found modules that used non standard methods of declaring dependency because of the version of terraform used at the time. Removing dependencies at the components level will allow us to focus on future features such as [multi cluster](036-multi-cluster.md) and blue green clusters.

We propose to add a new level to our build process called core. This splits our process into 3 distinct areas within our terraform:

- EKS - Everything required to bring up a basic cluster.
- Core - Core components that are required for our users to be able to safely and securely deploy their application to.
- Components - Optional components that provide additional functionality to a cluster

> Components are not optional in `live` (apart from the starter pack). When testing new components and features all components should be used to make sure there are no edge case issues.

### Considerations

- Developing this using test clusters will be easier than `live`, documentation and a focus on how we deploy this to `live`, and safely migrate existing terraform state into its new s3 location should be at the front of the work.
- Evaluate each component to determine whether they are core or optional components.
- Prometheus CRDs are a hard requirement on a lot of components, we can split these into `core` but would that increase complexity?
- Each stage must be a concrete deliverable.

## Decision

We have decided to:

- Start planning out the work required to split the components into `core`.
- Document the process of moving components.
- Test changes using test clusters.
- Apply the changes to `live`, `live-2` and `manager`

## Consequences

By implementing this solution we will be decreasing our component level complexity.

It opens up a lot of future roadmap work.
