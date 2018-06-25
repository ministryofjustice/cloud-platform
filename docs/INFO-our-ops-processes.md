# Cloud Platform operational processes

This is a record of the operational processes that we will use to support our users and to maintain the cloud platform.

## Hours of support

Our hours of providing support are 10am - 5pm. During this time we will work on support requests from teams and make sure someone is available to answer questions in [`#ask-cloud-platform`](https://mojdt.slack.com/messages/C57UPMZLY/).

Outside of these hours we will respond to [high priority](#prioritising-incidents) incidents as per our [on call](#our-on-call-process) process.

## Working hours

A member of the support team will be available online from 9am and until 6pm (it doesn't have to be the same person for both). This team member can be working remotely. We are doing this so that we cover incident response for the time when the person on call has is coming into the office or going home at the end of the day.

This support team member should be ready to respond to any high priority incidents.

## The support sub team

The support sub team will consist of around 3 engineers. On average each cloud platform team member will be on the sub team for 2 weeks every 6 weeks.

The key activities for the day are:

_On starting the day (9AM)_
At least one engineer to:
* Check the [`#low-priority-alarms`](https://mojdt.slack.com/messages/C8QR5FQRX/) and [`#high-priority-alarms`](https://mojdt.slack.com/messages/C8PF51AT0/) slack channels for any issues to investigate
* Check the board for any user requests that have been raised and not assigned
* Get a handover from the on call engineer about any issues out of hours
* Read through [`#ask-cloud-platform`](https://mojdt.slack.com/messages/C57UPMZLY/) slack channel

_During support hours (10AM - 5PM)_
At least one engineer to:
* Actively participate in [`#ask-cloud-platform`](https://mojdt.slack.com/messages/C57UPMZLY/) to field support requests, triage, prioritise and fix
* Monitor [`#high-priority-alarms`](https://mojdt.slack.com/messages/C8PF51AT0/) and [`#low-priority-alarms`](https://mojdt.slack.com/messages/C8QR5FQRX/) for incidents, triage, prioritise and fix
* Monitor user support requests being added to the board and triage them for priority

The whole support team to:
* Manage incident process for incidents in progress
* Work on user requests for larger work
* Work on support backlog stories

_Before ending the day (6PM)_
At least one engineer to:
* Handover information about any planned work or in progress high priority incidents to on call engineer

## `#ask-cloud-platform` slack channel

The [`#ask-cloud-platform`](https://mojdt.slack.com/messages/C57UPMZLY/) channel is our main entry point for support. We encourage people to ask questions and report problems in this channel and we'll do the best we can to help them.

One engineer should be available to answer questions throughout the hours of support (10AM - 5PM). This will take a bit of co-ordination for lunch etc so the support team should make sure someone is nominated for each period.  

### Communication style

In our responses we should aim to be:
* friendly
* timely
* open &mdash; no question is stupid

When the solution is a quick one, it is good to have the whole conversation in the channel rather than in a direct message so that other people can search/see what the resolution was (it might help them or you in the future).

### What we communicate

The main purpose of [`#ask-cloud-platform`](https://mojdt.slack.com/messages/C57UPMZLY/) is to discuss the problems that people are having and help them to solve them.

Sometimes it is useful to make broadcast communications, examples include:

* when there is planned work on a service (e.g. we are completing maintenance on Postgres RDS instance in prod)
* when there is an incident in progress (e.g there is a problem with Jenkins, no one can access it right now)

In these situations:

* try not to use notifiers like `@here` and `@channel` unless absolutely necessary
* perhaps use visual clues in your message to show it is an announcement
* commit to a time when you will update on the situation and keep that commitment

So for example:

```
**ANNOUNCEMENT:**
As agreed with the relevant teams, the cloud platform team will be upgrading Postgres from version 9.2 (out of support) to version 9.6 on CLA, Moving People Safely and Prison Visits Booking RDS instances today.

We will be doing this work between 10AM and 3PM, we will update at 1PM and 3PM to say how it is going.
```

### User support tickets

If someone is asking for help that will be quick to do (less than 15 mins) or is mainly advice then keep the interaction in channel and get it done.

If someone needs something that takes longer, is more challenging to complete or you find you've spent longer than 15 mins on it, then continue to talk in the channel but also ask the person asking for help to raise a ticket using the GitHub issue link: [long version](https://github.com/ministryofjustice/cloud-platform/issues/new?template=cloud-platform-support-request.md&labels=support%20team) and [short version](https://goo.gl/msfGiS).

The purpose of the ticket is to keep a record of the work we are doing and how it is progressing. It is *not* the primary communication channel with the person who raised the problem &mdash; this should remain slack where you can provide a richer and more human interaction, answering questions if necessary. The engineer working on the ticket should update the ticket for our record as work progresses.

In many cases it can be more helpful for the engineer on support to create the ticket rather than the person reporting the problem as you will be able to provide more relevant context.

Once they have created the ticket it will appear in the `Support To Do` column of our [sprint board](https://waffle.io/ministryofjustice/cloud-platform). It can then be moved into the relevant columns as you work on it.

## Our incident process

An incident starts when a member of the support team says that it has. It is usually triggered by an alert that indicates a problem. The team member that calls it will send a message saying that there is an incident in progress to [`#cloud-platform`](https://mojdt.slack.com/messages/C514ETYJX/) and if user impacting, [`#ask-cloud-platform`](https://mojdt.slack.com/messages/C57UPMZLY/).

1. Support team chooses an incident lead. This person will be the main investigator of the incident. They start to work on the problem, calling on other team members (from the whole team) to help as required.

2. The rest of the support team communicates the incident out to those who are impacting, including giving updates at regular schedules. The people that they communicate to include:
    * Users who are affected by the problem &mdash; via [`#ask-cloud-platform`](https://mojdt.slack.com/messages/C57UPMZLY/) and the affected team's own slack channels.
    * Team members for awareness or as they might be able to help &mdash; via [`#cloud-platform`](https://mojdt.slack.com/messages/C514ETYJX/) and the `@cloud-platform-team` group
    * People in the team who manage communication with senior leadership in MoJ &mdash; Steve, Karen, Tony, Kalbir.

3. In a *high priority incident* (see below), the support team will gather every 1 hour to work out if additional people/skills are needed and update any external comms.

4. Once the incident is resolved, the support team will communicate that out and prepare for a postmortem.  

## Prioritising incidents

An incident is a system failure or degradation that has an impact on users of the cloud platform.

The size of that impact determines the priority of the incident. At the moment, we find it useful to categories incidents as high priority &mdash; we will begin work on any incident as soon as a member of the support team is available but high priority incidents require greater focus and communication.

#### High priority

* Service outages of user facing services
* Slowdown or degradation of service on a user facing service
* Failure of a system that stops one or more teams from working
* Failure of a system that creates a major vulnerability in the cloud platform

## Postmortems

In the cloud platform team we follow the practice of better understanding system failures through [blameless postmortems](https://codeascraft.com/2012/05/22/blameless-postmortems/).

The point of a postmortem is to understand the various factors that led up to a system failure and to try to identify things that can be improved to try to stop a similar event occurring in future.

Our approach to postmortems is to:

* For high priority incidents to hold a group session as soon as possible after the incident is closed
* Develop a timeline of the activities that led to the incident and its resolution
* Walkthrough that timeline and use it to identify opportunities for improvements
* Emerge with a prioritised list of improvements to be made, with owners
* Write up and publish the postmortem, giving a record of the timeline and any actions that have been agreed

We currently publish our postmortem reports on [pagerduty](https://moj-digital-tools.pagerduty.com/postmortems/). We have agreement to make these public and are in the process of working out how to do that.

## Our on call process

Team members who are on call manage an on call rota in pagerduty. The on call rota consists of a primary engineer and a secondary engineer.

Team members that are new to the rota should ensure that they have a secondary who has enough knowledge and experience of the process to support them.

The hours of on call are 7AM–10AM and 5PM–10PM weekdays, 7AM-10PM weekends and bank holidays. 9AM - 10AM and 5PM - 6PM are additionally covered by a team member from support who is online during those periods.

On call team members will respond to high priority incidents out of hours and will work on those incidents for up to 1 hour. If the engineer is not able to resolve the issue within that timeframe they will:

* put the service into maintenance mode
* inform the affected team
* put together a plan to resolve the issue in hours

## Our documentation

For our current systems based around the Template Deploy tooling, our documentation lives in [Confluence](https://dsdmoj.atlassian.net/wiki/spaces/PLAT/overview). This documentation includes architecture, runbooks and common issues.

During the development of the Cloud Platform we are putting our documentation into Git repos stored on Github. The starting point for this documentation is the [Cloud Platform repo](https://github.com/ministryofjustice/cloud-platform)

## Measuring how we are doing

We would like to measure:

* System availability
* Number of incidents and high priority incidents
* Mean time to recovery
* User happiness with support

We do not currently have a process in place for measuring these but we are working on it.

## Approved

The cloud platform team approved this document:

- Sablu Miah
-
