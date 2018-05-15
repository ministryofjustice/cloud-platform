# GitHub/Waffle.io workflow

## Purpose

This HOWTO is meant to be the simplest possible workflow for getting
things into a GitHub issue on one of Cloud Platform's repos and ensuring
that issue becomes a ticket in
[waffle.io](https://waffle.io/ministryofjustice/cloud-platform).

There are at least two different ways to do this.  I will touch on both,
but only cover the one we use in the Cloud Platform in any detail.

## Adding a ticket

1) Sign in to your GitHub account.
2) Ensure your account is part of the organisation `ministryofjustice`.
3) Ensure you are on the `webops` team in that organisation.
4) Go to
[waffle.io](https://waffle.io/ministryofjustice/cloud-platform).
5) Enable your Waffle account for access to pubic repos.
6) Open the `ministryofjustice/cloud-platform` project board.
7) Click on the blue `Add Issue` button at the top of the page.
8) Add a brief, descriptive title.
9) Pick the repo where the issue/ticket should be created--if you aren't
sure which repo, ask the Product or Delivery Manager.
10) Select the column where the ticket should start:
  * `Inbox` is for as-yet-unestimated featrue work and should be
    **mostly** in priority order;
  * `Support To Do` for incoming support requests;
  * `To Do` for items to go into the current sprint or in one of the
    next two or three sprints.  These should be kept in priority order.
11) Associate the ticket with a label or milestone as required (both
are optional.
12) Save.

## Getting a quick estimate online

If you added an unesitmated ticket to `To Do`, or you need to get an
estimate on a ticket do this:

1) Copy the ticket link into `#cloud-platform` slack
2) Copy this code into slack in a link immediately following the link: `/poll “Estimate” “1"“2”“3"“5”“8"`
3) After enough relevant people have answered the poll and any
differences in opinion have been discussed and resovled, add the
estimate to the ticket by clicking on the scales icon.
