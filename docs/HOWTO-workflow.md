# GitHub/Waffle.io workflow

## Purpose

This HOWTO is meant to be the simplest possible workflow for getting
things into a GitHub issue on one of Cloud Platform's repos and ensuring
that issue becomes a ticket in
[waffle.io](https://waffle.io/ministryofjustice/cloud-platform).

There are at least two different ways to do this.  I will cover the two
main methods in use in Cloud Platform at the time of writing
(2018-05-15).

## Adding a ticket using waffle.io

1) Sign in to your GitHub account.
2) Ensure your account is part of the organisation `ministryofjustice`.
3) Ensure you are on the `webops` team in that organisation.
4) Go to
[waffle.io](https://waffle.io/ministryofjustice/cloud-platform).
5) Enable your Waffle account for access to pubic repos.
6) Open the `ministryofjustice/cloud-platform` project board.
7) Click on the blue `Add Issue` button at the top of the page.
8) Add a brief, descriptive title.
5) Add some details about your issue.
9) Pick the repo where the issue/ticket should be created--if you aren't
sure which repo, ask the Product or Delivery Manager.
10) Select the column where the ticket should start:
    * `Inbox` is for as-yet-unestimated featrue work and should be **mostly** in priority order;
    * `Support To Do` for incoming support requests;
    * `To Do` for items to go into the current sprint or in one of the
      next two or three sprints.  These should be kept in priority order.
11) Associate the ticket with a label or milestone as required (both
are optional.
12) Save.

## Working on a ticket

1) Click on the person icon in the top-right-hand corner of the ticket.
2) Assign yourself.
3) Move the ticket to `To Do`.
4) Do the thing!
5) Create a Pull Request.
   * The ticket should have the main details of the job.
   * The Pull Request should contain an executive summary of this in its
    description.
6) Drop a link to the Pull Request into the `cloud-platform` slack
channel and request someone reviews it.
7) If no one picks it up within an hour or two, ping it again, or ask
someone to do it directly.
8) If changes are requested, either:
   1) Make those changes and ask the same person to review the PR again.
   1) Convince them the changes don't need to be made.
   1) Ask them to approve and add their requested changes as new
tickets/issues.
9) On GitHub, when your PR is approved and all QA steps have passed:
   1) `Squash and Merge` your branch to master.
   2) `Delete Branch`
10) Move your ticket to done.

## Getting a quick estimate online

If you added an unesitmated ticket to `To Do`, or you need to get an
estimate on a ticket do this:

1) Copy the ticket link into `#cloud-platform` slack
2) Copy this code into slack in a link immediately following the link: `/poll “Estimate” “1"“2”“3"“5”“8"`
3) After enough relevant people have answered the poll and any
differences in opinion have been discussed and resovled, add the
estimate to the ticket by clicking on the scales icon.

## Ticket workflow best practices

### Limit work in progress

Try to make sure you never have more than one active ticket and one
ticket waiting for a Pull Request to be approved in `To Do`.  If you
find yourself with more than one that needs active work (including one
where you have to do more work to address PR comments), move the others
back to `To Do` or `Support To Do` an unassign yourself.

### Limit people working on tickets

With a few exceptions-retro actions, for example-a ticket should never
have more than two assigness.  If there is a compelling reason to have
more than two, please check with one of the Team Leaders before making
the assignment(s).

### Push back on tickets if they lack necessary info

If you go to pick up a ticket and it doesn't have enough information,
the definition of done isn't clear, the circumstances have changed or
there is something else wrong with it, feel free to push back on it.
Label it as 'Needs info', bring it to the attention of the PM and DM,
and pick up another ticket instead.

