# Cloud Platform operational processes

This is a record of the operational processes that we will use to support our users and to maintain the cloud platform.

## Hours of support

Our hours of providing support are 10am - 5pm. During this time we will work on support requests from teams and make sure someone is available to answer questions in #ask-cloud-platform.

Outside of these hours we will respond to high priority incidents as per our [on call]() process.

## Working hours

A member of the support team will be available online from 9am and until 6pm (it doesn't have to be the same person for both). This team member can be working remotely. We are doing this so that we cover incident response for the time when the person on call has is coming into the office or going home at the end of the day.

This support team member should be ready to respond to any high priority incidents.

## How does the support sub team work

The support sub team will consist of around 3 engineers. On average each cloud platform team member will be on the sub team for 2 weeks every 6 weeks.

The key activities for the day are:

_On starting the day (9AM)_
At least one engineer to:
* Check the #low-priority-alarms and #high-priority-alarms slack channels for any issues to investigate
* Check the board for any user requests that have been raised and not assigned
* Get a handover from the on call engineer about any issues out of hours
* Read through #ask-cloud-platform slack channel

_During support hours (10AM - 5PM)_
At least one engineer to:
* Actively participate in #ask-cloud-platform to field support requests, triage, prioritise and fix
* Monitor #high-priority-alarms and #low-priority alarms for incidents, triage, prioritise and fix
* Monitor user support requests being added to the board and triage them for priority

The whole support team to:
* Manage incident process for incidents in progress
* Work on user requests for larger work
* Work on support backlog stories

_Before ending the day (6PM)_
At least one engineer to:
* Handover information about any planned work or in progress high priority incidents to on call engineer

## #ask-cloud-platform

The #ask-cloud-platform channel is our main entry point for support. We encourage people to ask questions and report problems in this channel and we'll do the best we can to help them.

One engineer should be available to answer questions throughout the hours of support (10AM - 5PM). This will take a bit of co-ordination for lunch etc so the support team should make sure someone is nominated for each period.  

### Communication style

In our responses we should aim to be:
* friendly
* timely
* open &emdash; no question is stupid

When the solution is a quick one, it is good to have the whole conversation in the channel rather than in a direct message so that other people can search/see what the resolution was (it might help them or you in the future).

### What we communicate

The main purpose of #ask-cloud-platform is to discuss the problems that people are having and help them to solve them.

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

If someone needs something that takes longer, is more challenging to complete or you find you've spent longer than 15 mins on it, then continue to talk in the channel but also ask the person asking for help to raise a ticket using the GitHub issue link: [long version](https://github.com/ministryofjustice/cloud-platform/issues/new?template=cloud-platform-support-request.md&labels=support%20team) and [short version](goo.gl/msfGiS).

The purpose of the ticket is to keep a record of the work we are doing and how it is progressing. It is *not* the primary communication channel with the person who raised the problem &emdash; this should remain slack where you can provide a richer and more human interaction, answering questions if necessary.

In many cases it can be more helpful for the you (the engineer on support) to create the ticket rather than the person reporting the problem as you will be able to provide more relevant context.

Once they have created the ticket it will appear in the `Support To Do` column of our [sprint board](https://waffle.io/ministryofjustice/cloud-platform). It can then be moved into the relevant columns as you work on it.

## What is our incident process

## How do we prioritise incidents

## How we do postmortems

## How we do on call

## How do we measure how we are doing
