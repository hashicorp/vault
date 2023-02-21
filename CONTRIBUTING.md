# Contributing to Vault

**Please note:** We take Vault's security and our users' trust very seriously.
If you believe you have found a security issue in Vault, please responsibly
disclose by contacting us at security@hashicorp.com.

**First:** if you're unsure or afraid of _anything_, just ask or submit the
issue or pull request anyways. You won't be yelled at for giving it your best
effort. The worst that can happen is that you'll be politely asked to change
something. We appreciate any sort of contributions, and don't want a wall of
rules to get in the way of that. 

That said, if you want to ensure that a pull request is likely to be merged,
talk to us! You can find out our thoughts and ensure that your contribution
won't clash or be obviated by Vault's normal direction. A great way to do this
is via the [Vault Discussion Forum][2].

## Issues

This section will cover what we're looking for in terms of reporting issues.

By addressing all the points we're looking for, it raises the chances we can
quickly merge or address your contributions.

### Reporting an Issue

* Make sure you test against the latest released version. It is possible we
  already fixed the bug you're experiencing. Even better is if you can test
  against the `main` branch, as the bugs are regularly fixed but new versions
  are only released every few months.

* Provide steps to reproduce the issue, and if possible include the expected 
  results as well as the actual results. Please provide text, not screen shots!

* If you are seeing an internal Vault error (a status code of 5xx), please be
  sure to post relevant parts of (or the entire) Vault log, as often these
  errors are logged on the server but not reported to the user.

* If you experienced a panic, please create a [gist](https://gist.github.com)
  of the *entire* generated crash log for us to look at. Double check
  no sensitive items were in the log.

* Respond as promptly as possible to any questions made by the Vault
  team to your issue.

### Issue Lifecycle

1. The issue is reported.

2. The issue is verified and categorized by a Vault collaborator.
   Categorization is done via tags. For example, bugs are marked as "bugs".

3. Unless it is critical, the issue may be left for a period of time (sometimes
   many weeks), giving outside contributors -- maybe you!? -- a chance to
   address the issue.

4. The issue is addressed in a pull request or commit. The issue will be
   referenced in the commit message so that the code that fixes it is clearly
   linked.

5. The issue is closed.

6. Issues that are not reproducible and/or not gotten responses for a long time are
   stale issues. In order to provide faster responses and better engagement with
   the community, we strive to keep the issue tracker clean and the issue count
   low. In this regard, our current policy is to close stale issues after 30 days.
   Closed issues will still be indexed and available for future viewers. If users
   feel that the issue is still relevant, we encourage reopening them.

## Pull requests

When submitting a PR you should reference an existing issue. If no issue already exists, 
please create one. This can be skipped for trivial PRs like fixing typos.

Creating an issue in advance of working on the PR can help to avoid duplication of effort, 
e.g. maybe we know of existing related work. Or it may be that we can provide guidance 
that will help with your approach.

Your pull request should have a description of what it accomplishes, how it does so,
and why you chose the approach you did.  PRs should include unit tests that validate
correctness and the existing tests must pass.  Follow-up work to fix tests
does not need a fresh issue filed.

Someone will do a first pass review on your PR making sure it follows the guidelines 
in this document.  If it doesn't we'll mark the PR incomplete and ask you to follow
up on the missing requirements.

### Changelog Entries
Please include a file within your PR named `changelog/#.txt`, where `#` is your
pull request ID.  There are many examples under [changelog](changelog/), but
the general format is

````
```release-note:CATEGORY
COMPONENT: summary of change
```
````

CATEGORY is one of `security`, `change`, `feature`, `improvement`, or `bug`.
Your PR is almost certain to be one of `bug` or `improvement`, but don't
worry too much about getting it exactly right, we'll tell you if a change is 
needed.

To determine the relevant component, consult [CHANGELOG](CHANGELOG.md) and pick
whichever one you see that seems the closest match.

You do not need to include the link at the end of the summary that appears in
CHANGELOG.md, those are generated automatically by the changelog-building 
process.

### Vault UI

How you contribute to the UI depends on what you want to contribute. If that is 
a new feature, please submit an informational issue first.  That issue 
should include a short description of the proposed feature, the use case,
the approach you're taking, and the tests that would be written.  A mockup
is optional but encouraged.  

Bug fixes are welcome in PRs but existing tests must pass and updated logic 
should be handled in new tests.  You needn't submit an issue first to fix bugs.

Keep in mind that the UI should be consistent with other areas of Vault. 
The UI should be user-centered, informative, and include edge cases and errors— 
including accommodations for users who may not have permissions to view or 
interact with your feature. If you are not comfortable with UI design, a Vault
designer can take a look at your work— just be aware that this might mean 
it will add some time to the PR process. 

Finally, in your code, try to avoid logic-heavy templates (when possible, 
calculate values in the .js file instead of .hbs) and Ember anti-patterns. 
And most of all, if you have any questions, please ask!

## Setting up Go to work on Vault

If you have never worked with Go before, you will have to complete the
following steps listed in the README, under the section [Developing Vault][1].


[1]: https://github.com/hashicorp/vault#developing-vault
[2]: https://discuss.hashicorp.com/c/vault

## Contributor License Agreement

We require that all contributors sign our Contributor License Agreement ("CLA") before we can accept the contribution. 

[Learn more about why HashiCorp requires a CLA and what the CLA includes](https://www.hashicorp.com/cla)
