# Use RSpec for our cluster tests

Date: 04/07/2019

## Status

✅ Accepted

## Context

We need to be able to test that the production cluster(s) and any test clusters we create behave the way we expect them to.

So, we need automated tests that exercise a cluster and confirm that the effect we get is the one we expected.

## Decision

We chose the ruby testing framework [rspec] for this.

Although there are some go-based testing frameworks for kubernetes, there are some problems with them, and some benefits to rspec:

* The kubernetes testing frameworks seem quite immature, with limited documentation, tooling and other resources such as examples
* There is limited go expertise in the team
* RSpec is a very mature framework, with a lot of tooling, documentation and support
* There is a lot of ruby/rspec experience in the wider organisation
* Ruby is our scripting language of choice, so rspec fits with that

## Consequences

We have translated our existing bash-based smoke tests into Rspec [here][smoke-tests].

We have a structure for adding more tests as we need them.

We will be adding automated tests to our cluster build and monitoring pipelines, to ensure they stay up to date.

[rspec]: https://rspec.info/
[smoke-tests]: https://github.com/ministryofjustice/cloud-platform-smoke-tests
