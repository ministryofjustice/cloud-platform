# 0. Record Architecture Decisions

Date: 2019-05-01

## Status

✅ Accepted

## Context

During our work on the Cloud Platform, we will have to make architectural
decisions about the platform and associated tools & processes.

When making decisions, we should record them somewhere for future reference, to
help us remember why we made them, and to help teams working in related areas
understand why we made them.

We should make our decisions public so that other teams can find them more
easily, and because [making things open makes things
better](https://www.gov.uk/guidance/government-design-principles#make-things-open-it-makes-things-better).

## Decision

We will use Architecture Decision Records, as described by Michael Nygard in
[this
article](http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions)

An architecture decision record is a short text file describing a single
decision.

We will keep ADRs in this public repository under decisions/[number]-[title].md

We should use a lightweight text formatting language like Markdown.

ADRs will be numbered sequentially and monotonically. Numbers will not be
reused.

If a decision is reversed, we will keep the old one around, but mark it as
superseded. (It's still relevant to know that it was the decision, but is no
longer the decision.)

We will use a format with just a few parts, so each document is easy to digest:

**Title** These documents have names that are short noun phrases. For example,
"ADR 1: Record architectural decisions" or "ADR 9: Use Docker for deployment"

**Status** A decision may be "proposed" if it's still under discussion, or
"accepted" once it is agreed. If a later ADR changes or reverses a decision, it
may be marked as "deprecated" or "superseded" with a reference to its
replacement.

**Context** This section describes the forces at play, including technological,
political, social, and local to the service. These forces are probably in
tension, and should be called out as such. The language in this section is
value-neutral. It is simply describing facts.

**Decision** This section describes our response to these forces. It is stated
in full sentences, with active voice. "We will ..."

**Consequences** This section describes the resulting context, after applying
the decision. All consequences should be listed here, not just the "positive"
ones. A particular decision may have positive, negative, and neutral
consequences, but all of them affect the team and service in the future.

The whole document should be one or two pages long. We will write each ADR as
if it is a conversation with a future person joining the team. This requires
good writing style, with full sentences organised into paragraphs. Bullets are
acceptable only for visual style, not as an excuse for writing sentence
fragments.

[adr-tools](https://github.com/npryce/adr-tools) can help us work with our ADRs
consistently.

We will link to these ADRs from other documentation where relevant.

## Consequences

One ADR describes one significant decision for the service. It should be
something that has an effect on how the rest of the service will run.

Developers and service stakeholders (and anyone else who's interested) can see
the ADRs, even as the team composition changes over time.

The motivation behind previous decisions is visible for everyone, present and
future. Nobody is left scratching their heads to understand, "What were they
thinking?" and the time to change old decisions will be clear from changes in
the service's context.

Having a central place to record decisions which affect all of our work will
make the sequence of decisions clear, and make it easier for us to refer back to
decisions later on.

This repo holds Architecture Decision Records, for the cloud platform team. We
will be using the architecture decision record to help keep a record of what
approaches we are currently taking to our infrastructure and to help our future
selves understand why those decisions were made. The approach is described by
Michael Nygard in [this
article](http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions).

Everyone can see the ADRs, even as the team composition changes over time. This
means motivation behind previous decisions is visible for everyone, present and
future. Nobody is left scratching their heads to understand, "What were they
thinking?" and the time to change old decisions will be clear from changes in
the project's context.
