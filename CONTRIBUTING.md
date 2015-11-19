# Contributing to Vault

**Please note:** We take Vault's security and our users' trust very seriously. 
If you believe you have found a security issue in Vault, please responsibly 
disclose by contacting us at security@hashicorp.com.

**First:** if you're unsure or afraid of _anything_, just ask
or submit the issue or pull request anyways. You won't be yelled at for
giving your best effort. The worst that can happen is that you'll be
politely asked to change something. We appreciate any sort of contributions,
and don't want a wall of rules to get in the way of that. 

You can also inquire on the [Vault Google Group][2], or in `#vault-tool` on Freenode.

This document will cover what we're looking for in terms of reporting issues. 
By addressing all the points we're looking for, it raises the chances we can 
quickly merge or address your contributions.

## Issues

### Reporting an Issue

* Make sure you test against the latest released version. It is possible
  we already fixed the bug you're experiencing.

* Provide steps to reproduce the issue, and if possible include the expected 
  results as well as the actual results

* If you experienced a panic, please create a [gist](https://gist.github.com)
  of the *entire* generated crash log for us to look at. Double check
  no sensitive items were in the log.

* Respond as promptly as possible to any questions made by the Vault
  team to your issue. Stale issues will be closed.

### Issue Lifecycle

1. The issue is reported.

2. The issue is verified and categorized by a Vault collaborator.
   Categorization is done via tags. For example, bugs are marked as "bugs".

3. Unless it is critical, the issue is left for a period of time (sometimes
   many weeks), giving outside contributors a chance to address the issue.

4. The issue is addressed in a pull request or commit. The issue will be
   referenced in the commit message so that the code that fixes it is clearly
   linked.

5. The issue is closed. Sometimes, valid issues will be closed to keep
   the issue tracker clean. The issue is still indexed and available for
   future viewers, or can be re-opened if necessary.

## Setting up Go to work on Vault

If you have never worked with Go before, you will have to complete the
following steps listed in the README, under the section [Developing Vault][1].


[1]: https://github.com/hashicorp/vault#developing-vault
[2]: https://groups.google.com/group/vault-tool