<!--

Thank you for contributing changes to this document! Because we use a central repository
to synchronize this file across all our repositories, make sure to make your edits
in the correct file, which you can find here:

https://github.com/ory/meta/blob/master/templates/repository/common/CONTRIBUTING.md

-->

# Contributing to Ory Dockertest

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Introduction](#introduction)
- [FAQ](#faq)
- [How can I contribute?](#how-can-i-contribute)
- [Communication](#communication)
- [Contributing Code](#contributing-code)
- [Documentation](#documentation)
- [Disclosing vulnerabilities](#disclosing-vulnerabilities)
- [Code Style](#code-style)
  - [Working with Forks](#working-with-forks)
- [Conduct](#conduct)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Introduction

There are many ways in which you can contribute, beyond writing code. The goal
of this document is to provide a high-level overview of how you can get
involved.

_Please note_: We take Ory Dockertest's security and our users' trust very
seriously. If you believe you have found a security issue in Ory Dockertest,
please responsibly disclose by contacting us at security@ory.sh.

First: As a potential contributor, your changes and ideas are welcome at any
hour of the day or night, weekdays, weekends, and holidays. Please do not ever
hesitate to ask a question or send a pull request.

If you are unsure, just ask or submit the issue or pull request anyways. You
won't be yelled at for giving it your best effort. The worst that can happen is
that you'll be politely asked to change something. We appreciate any sort of
contributions, and don't want a wall of rules to get in the way of that.

That said, if you want to ensure that a pull request is likely to be merged,
talk to us! You can find out our thoughts and ensure that your contribution
won't clash or be obviated by Ory
Dockertest's normal direction. A great way to
do this is via
[Ory Dockertest Discussions](https://github.com/ory/meta/discussions)
or the [Ory Chat](https://www.ory.sh/chat).

## FAQ

- I am new to the community. Where can I find the
  [Ory Community Code of Conduct?](https://github.com/ory/dockertest/blob/master/CODE_OF_CONDUCT.md)

- I have a question. Where can I get
  [answers to questions regarding Ory Dockertest?](#communication)

- I would like to contribute but I am not sure how. Are there
  [easy ways to contribute?](#how-can-i-contribute)
  [Or good first issues?](https://github.com/search?l=&o=desc&q=label%3A%22help+wanted%22+label%3A%22good+first+issue%22+is%3Aopen+user%3Aory+user%3Aory-corp&s=updated&type=Issues)

- I want to talk to other Ory Dockertest users.
  [How can I become a part of the community?](#communication)

- I would like to know what I am agreeing to when I contribute to Ory
  Dockertest.
  Does Ory have
  [a Contributors License Agreement?](https://cla-assistant.io/ory/dockertest)

- I would like updates about new versions of Ory Dockertest.
  [How are new releases announced?](https://ory.us10.list-manage.com/subscribe?u=ffb1a878e4ec6c0ed312a3480&id=f605a41b53)

## How can I contribute?

If you want to start contributing code right away, we have a
[list of good first issues](https://github.com/ory/dockertest/labels/good%20first%20issue).

There are many other ways you can contribute without writing any code. Here are
a few things you can do to help out:

- **Give us a star.** It may not seem like much, but it really makes a
  difference. This is something that everyone can do to help out Ory Dockertest.
  Github stars help the project gain visibility and stand out.

- **Join the community.** Sometimes helping people can be as easy as listening
  to their problems and offering a different perspective. Join our Slack, have a
  look at discussions in the forum and take part in our weekly hangout. More
  info on this in [Communication](#communication).

- **Helping with open issues.** We have a lot of open issues for Ory Dockertest
  and some of them may lack necessary information, some are duplicates of older
  issues. You can help out by guiding people through the process of filling out
  the issue template, asking for clarifying information, or pointing them to
  existing issues that match their description of the problem.

- **Reviewing documentation changes.** Most documentation just needs a review
  for proper spelling and grammar. If you think a document can be improved in
  any way, feel free to hit the `edit` button at the top of the page. More info
  on contributing to documentation [here](#documentation).

- **Help with tests.** Some pull requests may lack proper tests or test plans.
  These are needed for the change to be implemented safely.

## Communication

We use [Slack](https://www.ory.sh/chat). You are welcome to drop in and ask
questions, discuss bugs and feature requests, talk to other users of Ory, etc.

Check out [Ory Dockertest Discussions](https://github.com/ory/meta/discussions). This is a great place for
in-depth discussions and lots of code examples, logs and similar data.

You can also join our community hangout, if you want to speak to the Ory team
directly or ask some questions. You can find more info on the hangouts in
[Slack](https://www.ory.sh/chat).

If you want to receive regular notifications about updates to Ory Dockertest,
consider joining the mailing list. We will _only_ send you vital information on
the projects that you are interested in.

Also [follow us on twitter](https://twitter.com/orycorp).

## Contributing Code

Unless you are fixing a known bug, we **strongly** recommend discussing it with
the core team via a GitHub issue or [in our chat](https://www.ory.sh/chat)
before getting started to ensure your work is consistent with Ory Dockertest's
roadmap and architecture.

All contributions are made via pull requests. To make a pull request, you will
need a GitHub account; if you are unclear on this process, see GitHub's
documentation on [forking](https://help.github.com/articles/fork-a-repo) and
[pull requests](https://help.github.com/articles/using-pull-requests). Pull
requests should be targeted at the `master` branch. Before creating a pull
request, go through this checklist:

1. Create a feature branch off of `master` so that changes do not get mixed up.
1. [Rebase](http://git-scm.com/book/en/Git-Branching-Rebasing) your local
   changes against the `master` branch.
1. Run the full project test suite with the `go test -tags sqlite ./...` (or
   equivalent) command and confirm that it passes.
1. Run `make format` if a `Makefile` is available, `gofmt -s` if the project is
   written in Go, `npm run format` if the project is written for NodeJS.
1. Ensure that each commit has a descriptive prefix. This ensures a uniform
   commit history and helps structure the changelog.  
   Please refer to this
   [list of prefixes for Dockertest](https://github.com/ory/dockertest/blob/master/.github/semantic.yml)
   for an overview.
1. Sign-up with CircleCI so that it has access to your repository with the
   branch containing your PR. Simply creating a CircleCI account is sufficient
   for the CI jobs to run, you do not need to setup a CircleCI project for the
   branch.

If a pull request is not ready to be reviewed yet
[it should be marked as a "Draft"](https://docs.github.com/en/github/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/changing-the-stage-of-a-pull-request).

Before your contributions can be reviewed you need to sign our
[Contributor License Agreement](https://cla-assistant.io/ory/dockertest).

This agreement defines the terms under which your code is contributed to Ory.
More specifically it declares that you have the right to, and actually do, grant
us the rights to use your contribution. You can see the Apache 2.0 license under
which our projects are published
[here](https://github.com/ory/meta/blob/master/LICENSE).

When pull requests fail testing, authors are expected to update their pull
requests to address the failures until the tests pass.

Pull requests eligible for review

1. follow the repository's code formatting conventions;
2. include tests which prove that the change works as intended and does not add
   regressions;
3. document the changes in the code and/or the project's documentation;
4. pass the CI pipeline;
5. have signed our
   [Contributor License Agreement](https://cla-assistant.io/ory/dockertest);
6. include a proper git commit message following the
   [Conventional Commit Specification](https://www.conventionalcommits.org/en/v1.0.0/).

If all of these items are checked, the pull request is ready to be reviewed and
you should change the status to "Ready for review" and
[request review from a maintainer](https://docs.github.com/en/github/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/requesting-a-pull-request-review).

Reviewers will approve the pull request once they are satisfied with the patch.

## Documentation

Please provide documentation when changing, removing, or adding features.
Documentation resides in the project's
[docs](https://github.com/ory/dockertest/tree/master/docs) folder. Generate API and
configuration reference documentation using `cd docs; npm run gen`.

For further instructions please head over to
[docs/README.md](https://github.com/ory/dockertest/blob/master/README.md).

## Disclosing vulnerabilities

Please disclose vulnerabilities exclusively to
[security@ory.sh](mailto:security@ory.sh). Do not use GitHub issues.

## Code Style

Please follow these guidelines when formatting source code:

- Go code should match the output of `gofmt -s` and pass `golangci-lint run`.
- NodeJS and JavaScript code should be prettified using `npm run format` where
  appropriate.

### Working with Forks

```
# First you clone the original repository
git clone git@github.com:ory/ory/dockertest.git

# Next you add a git remote that is your fork:
git remote add fork git@github.com:<YOUR-GITHUB-USERNAME-HERE>/ory/dockertest.git

# Next you fetch the latest changes from origin for master:
git fetch origin
git checkout master
git pull --rebase

# Next you create a new feature branch off of master:
git checkout my-feature-branch

# Now you do your work and commit your changes:
git add -A
git commit -a -m "fix: this is the subject line" -m "This is the body line. Closes #123"

# And the last step is pushing this to your fork
git push -u fork my-feature-branch
```

Now go to the project's GitHub Pull Request page and click "New pull request"

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
  member, please contact one of the channel ops or a member of the Ory Dockertest
  core team immediately.
- Likewise any spamming, trolling, flaming, baiting or other attention-stealing
  behaviour is not welcome.

We welcome discussion about creating a welcoming, safe, and productive
environment for the community. If you have any questions, feedback, or concerns
[please let us know](https://www.ory.sh/chat).
