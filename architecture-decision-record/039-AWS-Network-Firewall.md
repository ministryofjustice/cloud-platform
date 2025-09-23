# 38. AWS Network Firewall

Date: 2025-09-16

## Status

✅ Accepted

## Context

We need the ability to control and inspect traffic leaving and entering Cloud Platform EKS VPCs. This is important for:

- Enforcing security policies (e.g. blocking malicious destinations).
- Monitoring outbound connections for compliance and audit.
- Detecting suspicious behaviour, such as unexpected protocols or destinations.
- Protecting workloads by inspecting inbound traffic and blocking known malicious sources.

An earlier idea was to [create a centralised egress VPC] to achieve these goals. However, this ran into problems: test clusters use overlapping CIDRs, which prevents reliable connectivity to a shared egress VPC.

We need an approach that gives us inspection and control now, without large network redesigns.

## Decision

We will deploy [AWS Network Firewall] inside the Cloud Platform VPC.

- The Cloud Platform VPC will route all outbound traffic through firewall subnets before reaching the NAT Gateway and the internet.
- The Cloud Platform VPC will also route inbound traffic through the firewall for inspection before reaching workloads.
- Firewall rule groups (such as domain allowlists, port restrictions, or intrusion detection rules) will be managed centrally in code and applied consistently across the environment.
- Logs will be sent to CloudWatch for monitoring, alerting, and auditing.

We will also create firewall policies with the following threat detection rule groups enabled:

- [AttackInfrastructureStrictOrder] – detects and blocks communication with known attacker infrastructure.
- [MalwareDomainsStrictOrder] – blocks domains associated with malware distribution or activity.
- [BotNetCommandAndControlDomainsStrictOrder] – prevents communication with botnet command-and-control servers.

This ensures that the Cloud Platform VPC has consistent ingress and egress traffic inspection, control, and active threat detection at its edge.

## Consequences

**Pros**
- Provides the required control and inspection of both ingress and egress traffic in the Cloud Platform VPC.
- Adds built-in AWS threat detection for attacker infrastructure, malware, and botnets.

**Cons**
- Ingress and egress inspection may add latency for workloads, requiring performance testing.
- Some threat detection rules may cause false positives, requiring a process for tuning.


[AWS Network Firewall]: https://aws.amazon.com/network-firewall/
[create a centralised egress VPC]: https://github.com/ministryofjustice/cloud-platform/issues/7428
[AttackInfrastructureStrictOrder]: https://docs.aws.amazon.com/network-firewall/latest/developerguide/aws-managed-rule-groups-atd.html
[MalwareDomainsStrictOrder]: https://docs.aws.amazon.com/network-firewall/latest/developerguide/aws-managed-rule-groups-domain-list.html
[BotNetCommandAndControlDomainsStrictOrder]: https://docs.aws.amazon.com/network-firewall/latest/developerguide/aws-managed-rule-groups-domain-list.html


