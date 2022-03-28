# Kibana is open to all service teams

Date: 24/03/2020

## Status

✅ Accepted

## Context

We want users of the cloud platform to be able to access Kibana so that they can see the logs for their applications in a central place. AWS Kibana does not provide easy ways for users to authenticate. We need to put a proxy in front of Kibana so that users can authenticate with Github and then be redirected to the [Kibana dashboard][kibana-webconsole] to access their logs.

## Decision

It has been decided to use a combination of Auth0 and an OIDC proxy app. The application is managed in the [cloud-platform-terraform-monitoring repo][kibana-proxy] and configured ministryofjustice GitHub organization users to access Kibana.


## Consequences

- Kibana can be accessed by any member of the ministryofjustice GitHub organization, which will provide access to logs for all teams. Information written to log files can be of a sensitive nature, so it is important for users to keep sensitive data out of the logs.

[kibana-proxy]: https://github.com/ministryofjustice/cloud-platform-terraform-monitoring/blob/main/templates/oauth2-proxy.yaml.tpl
[kibana-webconsole]: https://kibana.cloud-platform.service.justice.gov.uk/_plugin/kibana
[kibana-access-task]: https://github.com/ministryofjustice/cloud-platform/issues/286
