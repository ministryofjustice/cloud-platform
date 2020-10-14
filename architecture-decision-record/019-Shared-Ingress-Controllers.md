# Shared Ingress Controllers for teams

Date: 09/09/2020

## Status

🤔 Proposed

## Context

The design proposed in [018](018-Dedicated-Ingress-Controllers.md) proved
impossible to implement. Every ingress controller requires an AWS Network Load
Balancer (NLB), and AWS have hard limits on the number of NLBs we can create
per VPC and availability zone (AZ).

## Decision

We will:

* Disable mod-security on the default ingress-controller

This should enable this ingress-controller to comfortably handle thousands of ingresses.

* Create and manage a set of ingress controllers with mod-security enabled

Most, if not all, production services will want the protection of a web
application firewall, and mod-security is the easiest to enable. We need to
ensure that each ingress controller only handles as many mod-security-enabled
ingresses as it can reliably cope with.

* Continue to have some dedicated ingress controllers

During our aborted migration to dedicated ingress controllers for every
namespace, several ingress controllers were created for specific services. We
will leave these in place to avoid additional disruption to these service
teams.

## Consequences

* We will have to manage a pool of "mod-security-enabled" ingress controllers
* There will be additional, dedicated ingress controllers which, ideally,
  should not exist.
* We need to find out what is the maximum number of mod-security-enabled
  ingresses a single ingress controller can safely handle, and ensure we never
  exceed this.
* We may need to prevent teams from defining ingresses outside of the
  environments repository, since our current setup would enable teams to
  'overload' a given ingress controller by using its ingress class in their
  ingress definitions. This could impact other services.
  
## UPDATE 2020-10-14

During the last kubernetes upgrade, we hit another AWS limit on the number of SecurityGroupRules (details in the [incident log]).

As part of the remediation, we removed all dedicated ingress controllers and migrated ingresses to the default controller, removing any modsecurity annotations.

[incident log]: https://runbooks.cloud-platform.service.justice.gov.uk/incident-log.html#incident-on-2020-10-06-09-07-intermittent-quot-micro-downtimes-quot-on-various-services-using-dedicated-ingress-controllers
