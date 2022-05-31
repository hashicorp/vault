# changelog

This folder holds changelog updates from commit 3bc7d15 onwards.

Release notes are text files with three lines:

 1. An opening code block with the `release-note:<MODE>` type annotation.

    For example:

        ```release-note:bug

    Valid modes are:

     - `bug`
     - `change`
     - `feature`
     - `improvement`
     - `deprecation`

 2. A component (for example, `secret/pki` or `sdk/framework` or), a colon and a space, and then a one-line description of the change.

 3. An ending code block.

This should be in a file named after the pull request number (e.g., `12345.txt`).

There are many examples in this folder; check one out if you're stuck!

See [hashicorp/go-changelog](https://github.com/hashicorp/go-changelog) for full documentation on the supported entries.
