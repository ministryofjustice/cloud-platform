## Notes on PRs and managing merges

### Kill all your darlings (paraphrasing William Faulkner for the digital age) 

**You own your issues, ticket and PRs.** It's your job to:

  * Get your PR reviewed and merged.  It's OK to have a finished issued sitting in `In progress` for a little while.  However, it shouldn't be the norm; most of the time (read '80% of your issues or more') should be reviewed, merged and moved to `Done` before you start a new ticket/issue.
  * Make sure all the acceptance criteria are complete before you ask someone for a formal review.  If you aren't sure, raise a `WIP` PR and ask someone to look at it.
  * Make sure the ticket/issue has acceptance criteria before starting it.  If it doesn't, then don't start it–push it back to your DM/PM.  
    * You can write your own criteria.  This is particularly relevant if you have started the ticket/issue, yourself.  
    * However, if you didn't raise the ticket and it is missing some or all of its criteria and you know what these should be, don't hesitate to add them in.  If you do, make sure the original owner of the ticket is aware and agrees with your assessment before you start.
  * Connect your PRs to their associated issue(s) so waffle puts them together on the board - put `connects to <issue number or URI>` at the **TOP** of the PR description.  You can also put `closes <issue number or URI>` if you want github to automatically mark the issue as done when the PR is merged.  Use `closes` cautiously, not all PRs **close** the issue they were raised against. 
  * **Leave whoever is reviewing your PR alone to make their own assessment.** A PR review is a chance to get an objective assessment based on the strengths and weaknesses of your work. If you pair with the reviewer there is a risk your perspective will colour theirs.  
    * Feel free to pair/discuss comments and observations with them **after** they've made their initial review your PR.  Leave them to it until then.
 * Merge your branch ASAP after approval.
 * Delete your branch after merge. 
 * Move your ticket to `Done`.

## Dealing with the review

You can either fix the issues as requested or try to convince the reviewer to change their position.  It is best to do this via comments on the PR so we have history for the decision in the future.  However, if the issue is complex, it is probably best to discuss it with them and only note the salient points in the history.  

Please be mindful of the possibility that your reviewer might not be on the same level as yourself - don't accidentally browbeat someone into accepting your view of things. 

## Advisory PRs

You can issue a PR against your work any time you like.  PRs are a good tool for getting input on your work before you are ready to merge.  If you are doing this, just make sure your PR is marked `WIP` in the title: `[WIP] the thing I'm working on`.  You can also label the PR with `WIP` (make the label if it doesn't already exist on the repo). 

## Reviewing

* Make sure all the acceptance criteria are complete. 
* Don't approve a PR if there aren't any obvious acceptance criteria 
  * These can take many forms–checklists/bullet points are great, but it could also be a set of Given/When/Then statements (gherkin), or a narrative description. 
* Feel free to ask for help in the initial review.  Don't ask the originator, though (see above). 