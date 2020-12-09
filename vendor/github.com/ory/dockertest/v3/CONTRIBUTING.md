<!--

Thank you for contributing changes to this document! Because we use a central repository
to synchronize this file across all our repositories, make sure to make your edits
in the correct file, which you can find here:

https://github.com/ory/meta/blob/master/templates/repository/CONTRIBUTING.md

-->

# Contributing to ORY {{Project}}

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Introduction](#introduction)
- [Contributing Code](#contributing-code)
- [Disclosing vulnerabilities](#disclosing-vulnerabilities)
- [Code Style](#code-style)
- [Documentation](#documentation)
- [Pull request procedure](#pull-request-procedure)
- [Communication](#communication)
- [Conduct](#conduct)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Introduction

Please note: We take ORY {{Project}}'s security and our users' trust very
seriously. If you believe you have found a security issue in ORY {{Project}},
please responsibly disclose by contacting us at office@ory.sh.

First: if you're unsure or afraid of anything, just ask or submit the issue or
pull request anyways. You won't be yelled at for giving it your best effort. The
worst that can happen is that you'll be politely asked to change something. We
appreciate any sort of contributions, and don't want a wall of rules to get in
the way of that.

That said, if you want to ensure that a pull request is likely to be merged,
talk to us! You can find out our thoughts and ensure that your contribution
won't clash or be obviated by ORY {{Project}}'s normal direction. A great way to
do this is via the [ORY Community](https://community.ory.sh/) or join the
[ORY Chat](https://www.ory.sh/chat).

## Contributing Code

Unless you are fixing a known bug, we **strongly** recommend discussing it with
the core team via a GitHub issue or [in our chat](https://www.ory.sh/chat)
before getting started to ensure your work is consistent with ORY {{Project}}'s
roadmap and architecture.

All contributions are made via pull request. Note that **all patches from all
contributors get reviewed**. After a pull request is made other contributors
will offer feedback, and if the patch passes review a maintainer will accept it
with a comment. When pull requests fail testing, authors are expected to update
their pull requests to address the failures until the tests pass and the pull
request merges successfully.

At least one review from a maintainer is required for all patches (even patches
from maintainers).

Reviewers should leave a "LGTM" comment once they are satisfied with the patch.
If the patch was submitted by a maintainer with write access, the pull request
should be merged by the submitter after review.

## Disclosing vulnerabilities

Please disclose vulnerabilities exclusively to
[security@ory.sh](mailto:security@ory.sh). Do not use GitHub issues.

## Code Style

Please follow these guidelines when formatting source code:

- Go code should match the output of `gofmt -s` and pass `golangci-lint run`.
- NodeJS and JavaScript code should be prettified using `npm run format` where
  appropriate.

## Documentation

Please provide documentation when changing, removing, or adding features.
Documentation resides in the project's `docs` folder.

In cases where a project does not have a `docs` folder check the README for instructions.

The commands listed below work exclusively for projects with a `docs` folder

### Develop

To change the documentation locally, you need NodeJS installed and the project
checked out locally. Next, `cd` into `docs` and install the dependencies:

```shell script
$ cd docs
$ npm install
```

#### Start

To start a local development server with hot reloading, run:

```shell script
$ npm start
```

This command opens up a browser window. Please note that changes to the sidebar are not hot-reloaded
and require a restart of the command.

#### Build

The `npm build` generates static content into the `build` directory and can be
served using any static contents hosting service.

```shell script
$ npm build
```

## Pull request procedure

To make a pull request, you will need a GitHub account; if you are unclear on
this process, see GitHub's documentation on
[forking](https://help.github.com/articles/fork-a-repo) and
[pull requests](https://help.github.com/articles/using-pull-requests). Pull
requests should be targeted at the `master` branch. Before creating a pull
request, go through this checklist:

1. Create a feature branch off of `master` so that changes do not get mixed up.
1. [Rebase](http://git-scm.com/book/en/Git-Branching-Rebasing) your local
   changes against the `master` branch.
1. Run the full project test suite with the `go test ./...` (or equivalent)
   command and confirm that it passes.
1. Run `gofmt -s` (if the project is written in Go).
1. Ensure that each commit has a subsystem prefix (ex: `controller:`).

Pull requests will be treated as "review requests," and maintainers will give
feedback on the style and substance of the patch.

Normally, all pull requests must include tests that test your change.
Occasionally, a change will be very difficult to test for. In those cases,
please include a note in your commit message explaining why.

## How We Organize Our Work

All repositories in the [ORY organization](https://github.com/ory) have their issues and pull requests
monitored in the [ORY Monitoring Board](https://github.com/orgs/ory/projects/9). This allows
for a transparent backlog of unanswered issues and pull requests across the ecosystem from those
who are allowed to merge pull requests to the main branch.

The process is as follows:

1. _Cards_ represent open issues and pull requests and are automatically assigned to the **Triage** column if
   the author is not one of the maintainers and no maintainer has answered yet.
2. A maintainer assigns the issue or pull request to someone or adds the label _help wanted_
   which moves the card to **Requires Action**.
3. If a maintainer leaves a comment or review, the card moves to **Pending Reply**, implying that
   the original author needs to do something (e.g. implement a change, explain something in more detail, ...).
4. If a non-maintainer pushes changes to the pull request or leaves a comment, the card moves
   back to **Requires Action**.
5. If a card stays inactive for 60 days or more days, we assume that public interest in the issue
   or change has waned, **archiving** the card.
6. If the issue is closed or the pull request merged or closed, the card is **archived** as well.

We try our best to answer all issues and review all pull requests and hope that this transparent way
of keeping a backlog helps you better understand how heavy the workload is.

## Communication

We use [Slack](https://www.ory.sh/chat). You are welcome to drop in and ask
questions, discuss bugs, etc.

## Conduct

Whether you are a regular contributor or a newcomer, we care about making this
community a safe place for you and we've got your back.

- We are committed to providing a friendly, safe and welcoming environment for
  all, regardless of gender, sexual orientation, disability, ethnicity,
  religion, or similar personal characteristic.
- Please avoid using nicknames that might detract from a friendly, safe and
  welcoming environment for all.
- Be kind and courteous. There is no need to be mean or rude.
- We will exclude you from interaction if you insult, demean or harass anyone.
  In particular, we do not tolerate behavior that excludes people in socially
  marginalized groups.
- Private harassment is also unacceptable. No matter who you are, if you feel
  you have been or are being harassed or made uncomfortable by a community
  member, please contact one of the channel ops or a member of the ORY
  {{Project}} core team immediately.
- Likewise any spamming, trolling, flaming, baiting or other attention-stealing
  behaviour is not welcome.

We welcome discussion about creating a welcoming, safe, and productive
environment for the community. If you have any questions, feedback, or concerns
[please let us know](https://www.ory.sh/chat).
