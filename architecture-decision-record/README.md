# Architecture decision record

This repo holds Architecture Decision Records, for the cloud platform team. We will be using the architecture decision record to help keep a record of what approaches we are currently taking to our infrastructure and to help our future selves understand why those decisions were made. The approach is described by Michael Nygard in [this article](http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions).

Everyone can see the ADRs, even as the team composition changes over time. This means motivation behind previous decisions is visible for everyone, present and future. Nobody is left scratching their heads to understand, "What were they thinking?" and the time to change old decisions will be clear from changes in the project's context.

## How it works

An architecture decision record is a short text file describing a single decision. One ADR describes one significant decision for the project. It should be something that has an effect on how the rest of the project will run.

We keep architecture decision records in this repo using a naming structure of the form `NNN-description-of-decision.md`. We use [markdown](https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet) to format the records.

If a decision is reversed, we keep the old one around, but mark it as superseded. (It's still relevant to know that it was the decision, but is no longer the decision.)

## Format of a decision record

We use a format with just a few parts, so each document is easy to digest. The format has just a few parts.

**Title** These documents have names that are short noun phrases. For example, "ADR 1: Record architectural decisions" or "ADR 9: Use Docker for deployment"

**Status** A decision may be "proposed" if the project stakeholders haven't agreed with it yet, or "accepted" once it is agreed. If a later ADR changes or reverses a decision, it may be marked as "deprecated" or "superseded" with a reference to its replacement.

**Context** This section describes the forces at play, including technological, political, social, and project local. These forces are probably in tension, and should be called out as such. The language in this section is value-neutral. It is simply describing facts.

**Decision** This section describes our response to these forces. It is stated in full sentences, with active voice. "We will ..."

**Consequences** This section describes the resulting context, after applying the decision. All consequences should be listed here, not just the "positive" ones. A particular decision may have positive, negative, and neutral consequences, but all of them affect the team and project in the future.

The whole document should be one or two pages long. We will write each ADR as if it is a conversation with a future person joining the team. This requires good writing style, with full sentences organised into paragraphs. Bullets are acceptable only for visual style, not as an excuse for writing sentence fragments.

### Statuses:

* Proposed: 🤔
* Accepted: ✅
* Rejected: ❌
* Superseded: ⌛️
* Amended: ♻️
