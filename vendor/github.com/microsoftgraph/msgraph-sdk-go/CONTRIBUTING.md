# Contributing to the Microsoft Graph SDK for Go

The Microsoft Graph SDK for Go is available for all manner of contribution. There are a couple of different recommended paths to get contributions into the released version of this SDK.

__NOTE__ A signed a contribution license agreement is required for all contributions, and is checked automatically on new pull requests. Please read and sign [the agreement](https://cla.microsoft.com/) before starting any work for this repository.

## File issues

The best way to get started with a contribution is to start a dialog with the owners of this repository. Sometimes features will be under development or out of scope for this SDK and it's best to check before starting work on contribution.

## Submit pull requests for trivial changes

If you are making a change that does not affect the interface components and does not affect other downstream callers, feel free to make a pull request against the __main__ branch.  The main branch will be updated frequently.

Revisions of this nature will result in a 0.0.X change of the version number.

## Submit pull requests for features

If major functionality is being added, or there will need to be gestation time for a change, it should be submitted against the __feature__ branch.

Revisions of this nature will result in a 0.X.X change of the version number.

## Commit message format

To support our automated release process, pull requests are required to follow the [Conventional Commit](https://www.conventionalcommits.org/en/v1.0.0/)
format.

Each commit message consists of a **header**, an optional **body** and an optional **footer**. The header is the first line of the commit and
MUST have a **type** (see below for a list of types) and a **description**. An optional **scope** can be added to the header to give extra context.

```
<type>[optional scope]: <short description>
<BLANK LINE>
<optional body>
<BLANK LINE>
<optional footer(s)>
```

The recommended commit types used are:

 - **feat** for feature updates (increments the _minor_ version)
 - **fix** for bug fixes (increments the _patch_ version)
 - **perf** for performance related changes e.g. optimizing an algorithm
 - **refactor** for code refactoring changes
 - **test** for test suite updates e.g. adding a test or fixing a test
 - **style** for changes that don't affect the meaning of code. e.g. formatting changes
 - **docs** for documentation updates e.g. ReadMe update or code documentation updates
 - **build** for build system changes (gradle updates, external dependency updates)
 - **ci** for CI configuration file changes e.g. updating a pipeline
 - **chore** for miscallaneous non-sdk changesin the repo e.g. removing an unused file

Adding a footer with the prefix **BREAKING CHANGE:** will cause an increment of the _major_ version.