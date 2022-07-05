# Deprecated TLS Versions

Date: 04/07/2022

## Status

✅ Accepted

## Context

The current Cloud Platform Ingress Controllers support TLS 1.0, 1.1, and 1.2. This is due to end users of services being hosted on the platform having to rely on unsupported browser versions. We have had to retain for older TLS versions to support this ongoing dependancy.

Since the original decision to support older TLS versions a number of things have happended:

1. TLS version 1.3 needs to be supported
2. The user estate has undergone a number of end user device refreshes which have introduced modern browsers
3. There is a recurring issue of legacy TLS versions being identified in Service Team ITHCs
4. [NCSC Guidance on legacy TLS versions](https://www.ncsc.gov.uk/guidance/using-tls-to-protect-data#section_5)

We are now offering users new Ingress Controllers that offer support for 1.2 and 1.3, and we need to plan how we will manage the older versions of TLS.

## Decision

Remove support for TLS versions 1.0 and 1.1 on the new Ingress Controller.

Add support for TLS version 1.3

## Consequences

There is a small risk that a user on the estate may be using an unsupported browser that uses a deprecated version of TLS. Following a number end user device refresh projects we believe that users have been moved to browsers that support modern TLS versions.

To mitigate this risk we will advise our users to move non-production workloads first, test that moving to the new ingress controller doesn't cause an issuses, and then move production workloads.

Removing unsupported TLS versions will improve security on the platform.
