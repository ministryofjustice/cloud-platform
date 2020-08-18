# Introduce Open Policy Agent

Date: 07/06/2019

## Status

✅ Accepted

## Context

On the Cloud Platform, there is a need to implement various policies to safeguard our tenant applications and to enforce best practices.

Kubernetes offers various mechanisms that cover some of our needs (eg.: `ResourceQuotas` to prevent resource exhaustion and `PodSecurityPolicies` to enforce non-root containers) but there are other areas for which there is no builtin solution. However, kubernetes implements a Dynamic Admission Control API which introduces [admission webhooks][admission-control]. This API provides an easy way with which to expand on the existing admission controllers (built in the apiserver).

Our immediate need was to prevent users from reusing hostnames in `Ingresses`. Although our ingress controller prevents hijacking of hostnames, it does so silently and furthermore, this is not a documented behaviour. Therefore, we decided that the user should not be allowed to reuse hostnames already defined in other `Ingresses` and receive a useful error message if they try to do that.

## Decision

We explored a number of existing solutions in the open source community, as well as the possibility of implementing our own and we also discussed the issue with other organisations that use kubernetes before reaching a conclusion.

Eventually we decided to introduce the [Open Policy Agent][open-policy-agent]:
- It is a generic framework for building and enforcing policies (whereas most other existing implementations were designed around specific problems)
- The policies are defined in a declarative, high-level language
- It is designed for cloud-native environments
- It provides a kubernetes integration
- It provides a way by which to unit test the policies
- The project is adopted by CNCF

Although the project is still in alpha and very likely to change in the near future, we decided that it is stable enough for our needs and worth adopting even at these early stages, since the benefits outweigh the cost.

## Consequences

By building the foundation, we have introduced a solid framework on which to build policies for our multi-tenant platform and have solved our immediate need of protecting application hostnames.

The project seems to have gained significant velocity and as it evolves we will need to follow the development closely and adapt our implementation.

## Notes

The first policy we introduced solves the aforementioned issue with hostnames in `Ingresses` and can be found [here](https://github.com/ministryofjustice/cloud-platform-infrastructure/tree/907918c3022f581d4ebbae828e3b8200b6f067cd/terraform/cloud-platform-components/resources/opa/policies).

[admission-control]: https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/
[open-policy-agent]: https://openpolicyagent.org
